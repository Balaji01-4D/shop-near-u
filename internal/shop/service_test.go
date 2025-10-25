package shop

import (
	"reflect"
	"shop-near-u/internal/models"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(postgres.Open("host=localhost user=balaji password=balaji2005 dbname=testdb sslmode=disable"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&models.Shop{}); err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}
	return db
}

func clearTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Exec("DELETE FROM shops").Error; err != nil {
		t.Fatalf("failed to clear shops table: %v", err)
	}
}

func createTestRepository(t *testing.T, db *gorm.DB) *Repository {
	t.Helper()
	return NewRepository(db)
}

func TestService_RegisterShop(t *testing.T) {
	db := setupTestDB(t)
	defer clearTestDB(t, db)

	shop := &models.Shop{
		Name:    "Test Shop",
		OwnerName: "John Doe",
		Type:   "Grocery",
		Email:    "test@123.com",
		Password: "password123",
		Address:  "123 Test St",
		Latitude:   37.7749,
		Longitude: -122.4194,
		SupportsDelivery: true,
	}

}

func TestService_AuthenticateShop(t *testing.T) {
	type fields struct {
		repository *Repository
	}
	type args struct {
		request *ShopLoginDTORequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Shop
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repository: tt.fields.repository,
			}
			got, err := s.AuthenticateShop(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AuthenticateShop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.AuthenticateShop() = %v, want %v", got, tt.want)
			}
		})
	}
}
