package handlers

import (
	"errors"
	"log"
	"strings"

	mysqldrv "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/cg3rard/AoLDatabaseTechnology/backend/middleware"
)

// Resource adalah CRUD generik untuk satu tabel (T = struct model).
// Cukup isi opsinya; List/Get/Create/Update/Delete sudah disediakan.
type Resource[T any] struct {
	// Kolom yang boleh dipakai sebagai filter query, mis. ?unit_id=1&status=open.
	Filters []string
	// Urutan hasil, default "id".
	Order string
	// Membatasi baris yang boleh dilihat pengguna (mis. penghuni hanya unitnya sendiri).
	// Mengembalikan kondisi SQL + argumen, atau "" jika tanpa batasan.
	Scope func(cl *middleware.Claims) (string, []any)
	// Dijalankan sebelum INSERT (dalam transaksi): isi field otomatis dan validasi.
	OnCreate func(tx *gorm.DB, cl *middleware.Claims, item *T) error
	// Dijalankan sebelum UPDATE: old = data lama, item = data baru.
	OnUpdate func(tx *gorm.DB, cl *middleware.Claims, old, item *T) error
}

func (r Resource[T]) scoped(q *gorm.DB, cl *middleware.Claims) *gorm.DB {
	if r.Scope != nil {
		if cond, args := r.Scope(cl); cond != "" {
			q = q.Where(cond, args...)
		}
	}
	return q
}

// writeError menerjemahkan error database/hook menjadi respons HTTP yang ramah.
func writeError(c *fiber.Ctx, err error, fallback string) error {
	var fe *fiber.Error
	if errors.As(err, &fe) {
		return errJSON(c, fe.Code, fe.Message)
	}
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errJSON(c, 404, "data tidak ditemukan")
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return errJSON(c, 409, "nilai yang harus unik sudah dipakai")
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return errJSON(c, 409, "data terkait tidak ditemukan, atau data ini masih dipakai data lain")
	}
	var me *mysqldrv.MySQLError
	if errors.As(err, &me) {
		switch me.Number {
		case 1451, 1452: // pelanggaran foreign key
			return errJSON(c, 409, "data terkait tidak ditemukan, atau data ini masih dipakai data lain")
		case 1062:
			return errJSON(c, 409, "nilai yang harus unik sudah dipakai")
		case 1048, 1264, 1265, 1292, 1366, 1406, 3819: // NULL, di luar rentang, ENUM/tipe salah, terlalu panjang
			log.Println("input ditolak database:", err)
			return errJSON(c, 400, "data tidak valid atau ada field wajib yang kosong")
		}
	}
	log.Println(fallback+":", err)
	return errJSON(c, 500, fallback)
}

func filterValue(v string) any {
	switch strings.ToLower(v) {
	case "true":
		return 1
	case "false":
		return 0
	}
	return v
}

func clamp(v, lo, hi int) int {
	return max(lo, min(v, hi))
}

// GET /api/<resource>?limit=&offset=&<filter>=
func (r Resource[T]) List(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl := middleware.GetClaims(c)
		q := r.scoped(db.Model(new(T)), cl)
		for _, f := range r.Filters {
			if v := c.Query(f); v != "" {
				q = q.Where("`"+f+"` = ?", filterValue(v))
			}
		}
		order := r.Order
		if order == "" {
			order = "id"
		}
		limit := clamp(c.QueryInt("limit", 100), 1, 500)
		offset := max(c.QueryInt("offset", 0), 0)

		items := make([]T, 0)
		if err := q.Order(order).Limit(limit).Offset(offset).Find(&items).Error; err != nil {
			return writeError(c, err, "gagal mengambil data")
		}
		return c.JSON(items)
	}
}

// GET /api/<resource>/:id
func (r Resource[T]) Get(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl := middleware.GetClaims(c)
		id, err := c.ParamsInt("id")
		if err != nil || id <= 0 {
			return errJSON(c, 400, "id harus berupa angka positif")
		}
		var item T
		if err := r.scoped(db, cl).First(&item, id).Error; err != nil {
			return writeError(c, err, "gagal mengambil data")
		}
		return c.JSON(item)
	}
}

// POST /api/<resource>
func (r Resource[T]) Create(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl := middleware.GetClaims(c)
		var item T
		if err := c.BodyParser(&item); err != nil {
			return errJSON(c, 400, "format JSON tidak valid")
		}
		err := db.Transaction(func(tx *gorm.DB) error {
			if r.OnCreate != nil {
				if err := r.OnCreate(tx, cl, &item); err != nil {
					return err
				}
			}
			// Omit("id"): id selalu dibuat database, abaikan kiriman klien.
			return tx.Omit("id").Create(&item).Error
		})
		if err != nil {
			return writeError(c, err, "gagal menyimpan data")
		}
		// Baca ulang agar nilai default database (status, waktu, dll.) ikut tampil.
		_ = db.First(&item).Error
		return c.Status(201).JSON(item)
	}
}

// PUT /api/<resource>/:id  (field yang tidak dikirim tetap memakai nilai lama)
func (r Resource[T]) Update(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl := middleware.GetClaims(c)
		id, err := c.ParamsInt("id")
		if err != nil || id <= 0 {
			return errJSON(c, 400, "id harus berupa angka positif")
		}
		var old T
		if err := r.scoped(db, cl).First(&old, id).Error; err != nil {
			return writeError(c, err, "gagal mengambil data")
		}

		item := old // mulai dari data lama, lalu timpa dengan kiriman klien
		if err := c.BodyParser(&item); err != nil {
			return errJSON(c, 400, "format JSON tidak valid")
		}
		if r.OnUpdate != nil {
			if err := r.OnUpdate(db, cl, &old, &item); err != nil {
				return writeError(c, err, "gagal mengubah data")
			}
		}
		// Select("*") menulis semua kolom (termasuk nilai nol); id dan created_at tidak boleh berubah.
		err = db.Model(&old).Select("*").Omit("id", "created_at").Updates(&item).Error
		if err != nil {
			return writeError(c, err, "gagal mengubah data")
		}

		var out T
		if err := db.First(&out, id).Error; err != nil {
			return writeError(c, err, "gagal mengambil data")
		}
		return c.JSON(out)
	}
}

// DELETE /api/<resource>/:id
func (r Resource[T]) Delete(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil || id <= 0 {
			return errJSON(c, 400, "id harus berupa angka positif")
		}
		result := db.Delete(new(T), id)
		if result.Error != nil {
			return writeError(c, result.Error, "gagal menghapus data")
		}
		if result.RowsAffected == 0 {
			return errJSON(c, 404, "data tidak ditemukan")
		}
		return c.SendStatus(204)
	}
}
