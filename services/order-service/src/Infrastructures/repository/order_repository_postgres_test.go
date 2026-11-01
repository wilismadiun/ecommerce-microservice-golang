package repository

import (
	"log"
	"order-service/config"
	"order-service/src/Domains/entities"
	"order-service/src/Infrastructures/database"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var orderRepository *OrderRepositoryPostgres

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("Peringatan: Gagal memuat .env dari root, mencoba path alternatif...")
	}

	dsn := config.DSN()
	database.ConnectDB(dsn)

	orderRepository = &OrderRepositoryPostgres{
		DB: database.DB,
	}

	code := m.Run()

	os.Exit(code)
}

func Test_OrderRepository(t *testing.T) {
	orderID := "order-test-uuid-001"
	item1ID := "item-test-uuid-001"
	item2ID := "item-test-uuid-002"

	dummyOrder := entities.Order{
		ID:              orderID,
		UserID:          "user-test-uuid-123",
		StoreID:         "store-test-999",
		TotalAmount:     2950000.0,
		Status:          "PENDING",
		ShippingAddress: "Jl. Merdeka No. 45, Jakarta",
		Items: []entities.OrderItem{
			{
				ID:          item1ID,
				OrderID:     orderID,
				ProductID:   "prod-test-001",
				ProductName: "Mechanical Keyboard",
				Quantity:    2,
				Price:       1250000.0,
			},
			{
				ID:          item2ID,
				OrderID:     orderID,
				ProductID:   "prod-test-002",
				ProductName: "Gaming Mouse",
				Quantity:    1,
				Price:       450000.0,
			},
		},
	}

	// Bersihkan database setelah test selesai menggunakan t.Cleanup
	// Karena menggunakan ON DELETE CASCADE pada foreign key, menghapus order akan otomatis menghapus item-nya.
	t.Cleanup(func() {
		database.DB.Unscoped().Delete(&entities.Order{}, "id = ?", orderID)
	})

	// 2. Act: Panggil fungsi CreateOrder (sesuai fungsi repository manual Anda)
	err := orderRepository.CreateOrder(dummyOrder)

	// 3. Assert: Pastikan tidak ada error
	assert.NoError(t, err)

	// Verifikasi: Pastikan data order tersimpan di database beserta relasi item-nya
	var savedOrder entities.Order
	err = database.DB.Preload("Items").First(&savedOrder, "id = ?", orderID).Error
	assert.NoError(t, err)
	assert.Equal(t, dummyOrder.UserID, savedOrder.UserID)
	assert.Equal(t, dummyOrder.StoreID, savedOrder.StoreID)
	assert.Equal(t, dummyOrder.TotalAmount, savedOrder.TotalAmount)
	assert.Equal(t, dummyOrder.Status, savedOrder.Status)
	assert.Equal(t, dummyOrder.ShippingAddress, savedOrder.ShippingAddress)

	// Verifikasi jumlah dan isi items
	assert.Len(t, savedOrder.Items, 2)
	assert.Equal(t, "Mechanical Keyboard", savedOrder.Items[0].ProductName)
	assert.Equal(t, "Gaming Mouse", savedOrder.Items[1].ProductName)
}
