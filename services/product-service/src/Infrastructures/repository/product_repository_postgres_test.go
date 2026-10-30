package repository

import (
	"log"
	"os"
	"product-service/config"
	"product-service/src/Domains/entities"
	"product-service/src/Infrastructures/database"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var productRepository *ProductRepositoryPostgres

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("Peringatan: Gagal memuat .env dari root, mencoba path alternatif...")
	}

	dsn := config.DSN()
	database.ConnectDB(dsn)

	productRepository = &ProductRepositoryPostgres{
		DB: database.DB,
	}

	code := m.Run()

	os.Exit(code)
}

func Test_CreateProduct(t *testing.T) {
	// 1. Siapkan data dummy produk untuk pengujian
	// Gunakan ID unik agar tidak berbenturan dengan data lain di database
	product := entities.Product{
		ID:          "test-prod-uuid-001",
		StoreID:     "store-test-999",
		Name:        "Produk Testing Unit",
		Description: "Deskripsi untuk keperluan testing database postgres",
		Price:       75000.0,
		Stock:       15,
	}

	// 2. Gunakan t.Cleanup untuk memastikan data dihapus setelah test selesai
	// Ini akan berjalan otomatis meskipun test sukses ataupun gagal
	t.Cleanup(func() {
		// Menghapus data dari database berdasarkan ID menggunakan Gorm
		// Unscoped() digunakan untuk menghapus secara permanen jika menggunakan soft delete
		err := productRepository.DB.Unscoped().Delete(&entities.Product{}, "id = ?", product.ID).Error
		if err != nil {
			t.Logf("Gagal membersihkan data testing: %v", err)
		}
	})

	// 3. Jalankan fungsi CreateProduct yang ingin diuji
	err := productRepository.CreateProduct(&product)

	// 4. Lakukan assertion untuk memastikan tidak ada error
	assert.NoError(t, err)

	// 5. (Opsional) Verifikasi tambahan: pastikan data benar-benar masuk ke database
	var savedProduct entities.Product
	errFind := productRepository.DB.First(&savedProduct, "id = ?", product.ID).Error

	assert.NoError(t, errFind)
	assert.Equal(t, product.ID, savedProduct.ID)
	assert.Equal(t, product.Name, savedProduct.Name)
	assert.Equal(t, product.Price, savedProduct.Price)
}

func Test_DecreaseStock(t *testing.T) {
	t.Run("if stock < quantity", func(t *testing.T) {
		productID := "test-prod-insufficient-01"
		dummyProduct := entities.Product{
			ID:          productID,
			StoreID:     "store-test-999",
			Name:        "Test Mouse",
			Description: "Wireless Mouse",
			Price:       150000.0,
			Stock:       3, // Stok sedikit
		}

		err := database.DB.Create(&dummyProduct).Error
		assert.NoError(t, err)

		t.Cleanup(func() {
			database.DB.Unscoped().Delete(&entities.Product{}, "id = ?", productID)
		})

		err = productRepository.DecreaseStock(productID, 5)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrInsufficientStock, err)
	})

	t.Run("berhasil", func(t *testing.T) {
		productID := "test-prod-success-01"
		dummyProduct := entities.Product{
			ID:          productID,
			StoreID:     "store-test-999",
			Name:        "Test Monitor",
			Description: "LED Monitor",
			Price:       2000000.0,
			Stock:       10, // Stok cukup
		}

		err := database.DB.Create(&dummyProduct).Error
		assert.NoError(t, err)

		t.Cleanup(func() {
			database.DB.Unscoped().Delete(&entities.Product{}, "id = ?", productID)
		})

		err = productRepository.DecreaseStock(productID, 4)

		assert.NoError(t, err)

		var updatedProduct entities.Product
		database.DB.First(&updatedProduct, "id = ?", productID)
		assert.Equal(t, 6, updatedProduct.Stock)
	})
}

func TestProductRepositoryPostgres_GetProductById(t *testing.T) {
	t.Run("should return error when product not found", func(t *testing.T) {
		nonExistentID := "non-existent-id-999"
		product, err := productRepository.GetProductById(nonExistentID)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrProductNotFound, err)
		assert.Equal(t, entities.Product{}, product)
	})

	t.Run("should return product successfully when product exists", func(t *testing.T) {
		productID := "test-product-id-123"
		dummyProduct := entities.Product{
			ID:          productID,
			StoreID:     "store-test-999",
			Name:        "Mechanical Keyboard",
			Description: "Wireless RGB Mechanical Keyboard",
			Price:       1250000.0,
			Stock:       25,
		}

		err := database.DB.Create(&dummyProduct).Error
		assert.NoError(t, err)

		t.Cleanup(func() {
			database.DB.Unscoped().Delete(&entities.Product{}, "id = ?", productID)
		})

		product, err := productRepository.GetProductById(productID)

		assert.NoError(t, err)
		assert.Equal(t, productID, product.ID)
		assert.Equal(t, dummyProduct.Name, product.Name)
		assert.Equal(t, dummyProduct.StoreID, product.StoreID)
		assert.Equal(t, dummyProduct.Stock, product.Stock)
	})
}

func TestProductRepositoryPostgres_UpdateStock(t *testing.T) {
	t.Run("should update product stock successfully", func(t *testing.T) {
		// 1. Arrange: Buat data dummy produk dengan stok awal (misal: 10)
		productID := "test-update-stock-01"
		dummyProduct := entities.Product{
			ID:          productID,
			StoreID:     "store-test-999",
			Name:        "Test Product Update Stock",
			Description: "Description",
			Price:       100000.0,
			Stock:       10, // Stok awal
		}

		err := database.DB.Create(&dummyProduct).Error
		assert.NoError(t, err)

		// Bersihkan database setelah test selesai menggunakan t.Cleanup
		t.Cleanup(func() {
			database.DB.Unscoped().Delete(&entities.Product{}, "id = ?", productID)
		})

		// 2. Act: Panggil fungsi UpdateStock dengan nilai stok baru (misal: 50)
		newStock := 50
		err = productRepository.UpdateStock(productID, newStock)

		// 3. Assert: Pastikan tidak ada error
		assert.NoError(t, err)

		// Verifikasi langsung ke database bahwa stok benar-benar berubah menjadi 50
		var updatedProduct entities.Product
		database.DB.First(&updatedProduct, "id = ?", productID)
		assert.Equal(t, newStock, updatedProduct.Stock)
	})
}
