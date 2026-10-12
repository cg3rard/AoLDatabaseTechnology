package middleware

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Peran pengguna (sama dengan ENUM `role` di tabel users).
const (
	RoleAdmin      = "admin"
	RoleResident   = "resident"
	RoleLaundry    = "laundry_staff"
	RoleSecurity   = "security"
	RoleTechnician = "technician"
)

// Claims adalah identitas pengguna yang sedang login (diambil dari token).
type Claims struct {
	UserID uint
	Role   string
	UnitID uint // 0 jika pengguna tidak terhubung ke unit
}

type tokenClaims struct {
	Role   string `json:"role"`
	UnitID uint   `json:"unit_id"`
	jwt.RegisteredClaims
}

const claimsKey = "claims"

// IssueToken membuat JWT (HS256) yang berlaku selama ttl.
func IssueToken(secret []byte, ttl time.Duration, userID uint, role string, unitID uint) (string, error) {
	now := time.Now()
	claims := tokenClaims{
		Role:   role,
		UnitID: unitID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(userID), 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func unauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "login diperlukan atau token tidak valid"})
}

// Auth memeriksa header "Authorization: Bearer <token>" dan menyimpan Claims di request.
func Auth(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		const prefix = "Bearer "
		h := c.Get("Authorization")
		if len(h) <= len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
			return unauthorized(c)
		}

		var tc tokenClaims
		tok, err := jwt.ParseWithClaims(
			h[len(prefix):], &tc,
			func(t *jwt.Token) (any, error) { return secret, nil },
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
		)
		if err != nil || !tok.Valid {
			return unauthorized(c)
		}
		uid, err := strconv.ParseUint(tc.Subject, 10, 64)
		if err != nil || uid == 0 {
			return unauthorized(c)
		}

		c.Locals(claimsKey, &Claims{UserID: uint(uid), Role: tc.Role, UnitID: tc.UnitID})
		return c.Next()
	}
}

// GetClaims mengambil identitas pengguna. Hanya valid setelah middleware Auth.
func GetClaims(c *fiber.Ctx) *Claims {
	cl, _ := c.Locals(claimsKey).(*Claims)
	return cl
}

// RequireRoles hanya meloloskan pengguna dengan salah satu peran yang disebut.
// Tanpa argumen, semua pengguna yang sudah login diloloskan.
func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl := GetClaims(c)
		if cl == nil {
			return unauthorized(c)
		}
		if len(roles) == 0 || slices.Contains(roles, cl.Role) {
			return c.Next()
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "peran Anda tidak boleh mengakses ini"})
	}
}
