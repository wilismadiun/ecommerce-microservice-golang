package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONMap map[string]interface{}

// Value mengonversi map ke format JSON []byte agar bisa dibaca oleh driver database (PostgreSQL, MySQL, dll)
func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan mengonversi data JSON/JSONB dari database kembali menjadi map[string]interface{}
func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &m)
}

type Product struct {
	ID          string
	StoreID     string  `json:"store_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Attributes  JSONMap `gorm:"type:jsonb" json:"attributes"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func VerifyProdyct(product Product) error {
	if product.StoreID == "" {
		return errors.New("store id is required")
	}
	if product.Name == "" {
		return errors.New("name is required")
	}
	if product.Description == "" {
		return errors.New("description is required")
	}
	if product.Price <= 0 {
		return errors.New("price is required and must be greater than 0")
	}
	if product.Stock <= 0 {
		return errors.New("stock is required and must be greater than or equal to 0") // Jika stock 0 dianggap kosong/tidak ada
	}
	return nil
}
