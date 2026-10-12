package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

// Date mewakili kolom MySQL bertipe DATE.
// Di JSON ditulis "YYYY-MM-DD" (bukan format panjang RFC3339 seperti time.Time biasa).
type Date struct {
	time.Time
}

// Value dipanggil GORM saat menyimpan ke database.
func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

// Scan dipanggil GORM saat membaca dari database.
func (d *Date) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		d.Time = time.Time{}
	case time.Time:
		d.Time = v
	case []byte:
		return d.parse(string(v))
	case string:
		return d.parse(v)
	default:
		return fmt.Errorf("Date.Scan: tipe %T tidak didukung", src)
	}
	return nil
}

func (d *Date) parse(s string) error {
	t, err := time.ParseInLocation(dateLayout, s, time.Local)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Format(dateLayout) + `"`), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}
	if t, err := time.ParseInLocation(dateLayout, s, time.Local); err == nil {
		d.Time = t
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("tanggal harus berformat YYYY-MM-DD")
	}
	d.Time = t
	return nil
}
