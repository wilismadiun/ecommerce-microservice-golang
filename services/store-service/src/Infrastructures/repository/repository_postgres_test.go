package repository

import (
	"errors"
	"log"
	"os"
	"store-service/config"
	"store-service/src/Domains/entities"
	"store-service/src/Infrastructures/database"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var storeRepo *StoreRepositoryPostgres
var userStoreRepo *UserStoreRepositoryPostgres

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("Peringatan: Gagal memuat .env dari root, mencoba path alternatif...")
	}

	dsn := config.DSN()
	database.ConnectDB(dsn)

	storeRepo = &StoreRepositoryPostgres{
		DB: database.DB,
	}

	userStoreRepo = &UserStoreRepositoryPostgres{
		DB: database.DB,
	}

	code := m.Run()

	os.Exit(code)
}

func Test_FindStoreByName(t *testing.T) {
	// Skenario Berhasil: Cari toko berdasarkan nama yang ada di database

	// 1. Siapkan data dummy unik untuk dites
	dummyStore := entities.Store{
		Id:          "store-test-uuid-999",
		Name:        "Toko Unit Test Berhasil",
		Description: "Ini adalah toko untuk pengujian find by name",
		Phone:       "081234567890",
		Email:       "unittest@example.com",
		Address:     "Jl. Pengujian No. 123",
	}

	// 2. Masukkan data dummy ke database sebelum pengujian
	err := storeRepo.DB.Create(&dummyStore).Error
	assert.NoError(t, err, "Gagal memasukkan data dummy ke database")

	// 3. Pastikan data dihapus kembali di akhir test (cleanup) baik sukses maupun gagal
	t.Cleanup(func() {
		err := storeRepo.DB.Where("id = ?", dummyStore.Id).Delete(&entities.Store{}).Error
		if err != nil {
			t.Logf("Gagal membersihkan data test: %v", err)
		}
	})

	// 4. Jalankan fungsi yang diuji
	result, err := storeRepo.FindStoreByName(dummyStore.Name)

	// 5. Validasi hasil menggunakan testify/assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Id)
	assert.Equal(t, dummyStore.Id, result.Id)
	assert.Equal(t, dummyStore.Name, result.Name)
	assert.Equal(t, dummyStore.Email, result.Email)
}

func Test_AddStore(t *testing.T) {
	// Skenario Berhasil: Menambahkan store baru ke database

	// 1. Siapkan data store baru
	newStore := entities.Store{
		Id:          "store-add-uuid-999",
		Name:        "Toko Baru Test",
		Description: "Toko untuk pengujian add store",
		Phone:       "081298765432",
		Email:       "addstore@example.com",
		Address:     "Jl. Tambah No. 45",
	}

	// 2. Pastikan data dihapus di akhir test (cleanup) jika berhasil masuk,
	// atau untuk berjaga-jaga jika test gagal di tengah jalan.
	t.Cleanup(func() {
		err := storeRepo.DB.Where("id = ?", newStore.Id).Delete(&entities.Store{}).Error
		if err != nil {
			t.Logf("Gagal membersihkan data test AddStore: %v", err)
		}
	})

	// 3. Jalankan fungsi AddStore
	err := storeRepo.AddStore(&newStore)

	// 4. Validasi hasil menggunakan testify/assert
	assert.NoError(t, err)

	// 5. Verifikasi tambahan: cek apakah data benar-benar masuk ke database
	var foundStore entities.Store
	err = storeRepo.DB.Where("id = ?", newStore.Id).First(&foundStore).Error
	assert.NoError(t, err)
	assert.Equal(t, newStore.Id, foundStore.Id)
	assert.Equal(t, newStore.Name, foundStore.Name)
	assert.Equal(t, newStore.Email, foundStore.Email)
}

func Test_DeleteStoreById(t *testing.T) {
	// Skenario Berhasil: Menghapus store berdasarkan ID

	// 1. Siapkan data dummy yang akan dimasukkan dulu ke database
	dummyStore := entities.Store{
		Id:          "store-delete-uuid-999",
		Name:        "Toko Mau Dihapus",
		Description: "Toko untuk pengujian delete",
		Phone:       "081233334444",
		Email:       "deletestore@example.com",
		Address:     "Jl. Hapus No. 99",
	}

	// Masukkan data ke database
	err := storeRepo.DB.Create(&dummyStore).Error
	assert.NoError(t, err, "Gagal menyiapkan data dummy untuk test delete")

	// 2. Gunakan t.Cleanup sebagai antisipasi jika test gagal sebelum proses delete,
	// agar data dummy tetap terhapus dari database (membersihkan sampah test).
	t.Cleanup(func() {
		_ = storeRepo.DB.Where("id = ?", dummyStore.Id).Delete(&entities.Store{}).Error
	})

	// 3. Jalankan fungsi DeleteStoreById yang diuji
	err = storeRepo.DeleteStoreById(dummyStore.Id)

	// 4. Validasi bahwa fungsi berjalan tanpa error
	assert.NoError(t, err)

	// 5. Verifikasi tambahan: pastikan data benar-benar sudah hilang dari database
	var foundStore entities.Store
	err = storeRepo.DB.Where("id = ?", dummyStore.Id).First(&foundStore).Error

	// GORM akan mengembalikan error record not found jika data berhasil dihapus
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Data seharusnya sudah tidak ada di database")
}

func Test_FindUserStoreByUser(t *testing.T) {
	// Skenario Berhasil: Menemukan daftar user store berdasarkan user_id

	// 1. Buat data dummy Store terlebih dahulu (karena ada foreign key constraint)
	dummyStore := entities.Store{
		Id:          "store-uuid-999",
		Name:        "Toko Relasi Test",
		Description: "Deskripsi toko test",
		Phone:       "081234567890",
		Email:       "storetest@example.com",
		Address:     "Jl. Test No. 1",
	}

	err := storeRepo.DB.Create(&dummyStore).Error
	assert.NoError(t, err, "Gagal menyiapkan data dummy Store")

	// 2. Siapkan data dummy UserStore yang terhubung ke Store di atas
	dummyUserStore := entities.UserStore{
		Id:      "user-store-uuid-999",
		StoreId: dummyStore.Id, // Mengacu ke Store yang baru dibuat
		UserId:  "user-test-uuid-123",
		RoleId:  "1",
	}

	err = userStoreRepo.DB.Create(&dummyUserStore).Error
	assert.NoError(t, err, "Gagal menyiapkan data dummy UserStore")

	// 3. Pastikan semua data dummy dihapus di akhir test (cleanup)
	// Urutan penghapusan penting: hapus relasinya dulu (UserStore), baru data induknya (Store)
	t.Cleanup(func() {
		_ = userStoreRepo.DB.Where("id = ?", dummyUserStore.Id).Delete(&entities.UserStore{}).Error
		_ = storeRepo.DB.Where("id = ?", dummyStore.Id).Delete(&entities.Store{}).Error
	})

	// 4. Jalankan fungsi FindUserStoreByUser yang diuji
	results, err := userStoreRepo.FindUserStoreByUser(dummyUserStore.UserId)

	// 5. Validasi hasil menggunakan testify/assert
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	assert.Len(t, results, 1)
	assert.Equal(t, dummyUserStore.Id, results[0].Id)
	assert.Equal(t, dummyUserStore.StoreId, results[0].StoreId)
	assert.Equal(t, dummyUserStore.UserId, results[0].UserId)
	assert.Equal(t, dummyUserStore.RoleId, results[0].RoleId)
}

func Test_AddUserStore(t *testing.T) {
	// Skenario Berhasil: Menambahkan relasi user store ke database

	// 1. Buat data dummy Store terlebih dahulu (karena foreign key constraint ke stores)
	dummyStore := entities.Store{
		Id:          "store-add-user-uuid-999",
		Name:        "Toko Untuk Add User Store",
		Description: "Deskripsi toko",
		Phone:       "081234567890",
		Email:       "storeadduser@example.com",
		Address:     "Jl. Add Store No. 1",
	}

	err := storeRepo.DB.Create(&dummyStore).Error
	assert.NoError(t, err, "Gagal menyiapkan data dummy Store untuk Test_AddUserStore")

	// Catatan: Jika tabel roles Anda wajib diisi, pastikan role dengan nama "1" sudah ada di database,
	// atau sesuaikan nilai RoleId dengan role yang valid di database Anda.

	// 2. Siapkan data dummy UserStore
	newUserStore := entities.UserStore{
		Id:      "user-store-add-uuid-999",
		StoreId: dummyStore.Id,
		UserId:  "user-test-uuid-456",
		RoleId:  "1",
	}

	// 3. Pastikan data dibersihkan di akhir test (cleanup) dengan urutan yang benar
	t.Cleanup(func() {
		_ = userStoreRepo.DB.Where("id = ?", newUserStore.Id).Delete(&entities.UserStore{}).Error
		_ = storeRepo.DB.Where("id = ?", dummyStore.Id).Delete(&entities.Store{}).Error
	})

	// 4. Jalankan fungsi AddUserStore yang diuji
	err = userStoreRepo.AddUserStore(newUserStore)

	// 5. Validasi bahwa fungsi berjalan tanpa error
	assert.NoError(t, err)

	// 6. Verifikasi tambahan: cek apakah data benar-benar tersimpan di database
	var foundUserStore entities.UserStore
	err = userStoreRepo.DB.Where("id = ?", newUserStore.Id).First(&foundUserStore).Error
	assert.NoError(t, err)
	assert.Equal(t, newUserStore.Id, foundUserStore.Id)
	assert.Equal(t, newUserStore.StoreId, foundUserStore.StoreId)
	assert.Equal(t, newUserStore.UserId, foundUserStore.UserId)
	assert.Equal(t, newUserStore.RoleId, foundUserStore.RoleId)
}
