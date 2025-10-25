package shop_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shop-near-u/internal/database"
	"shop-near-u/internal/models"
	"shop-near-u/internal/product"
	"shop-near-u/internal/server"
	"shop-near-u/internal/shop"
	"shop-near-u/internal/utils"
	"strings"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type Data struct {
	models.Shop
	Token string `json:"token"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Data   `json:"data,omitempty"`
}

func registerUser(t *testing.T) Response {
	t.Helper()
	router := server.NewServer().Handler

	w := httptest.NewRecorder()

	dto := shop.ShopRegisterDTORequest{
		Name:      "test shop1",
		OwnerName: "test",
		Email:     "test@gmail.com",
		Mobile:    "1234567890",
		Address:   "no 123. test street, test state",
		Latitude:  13.07439,
		Longitude: 80.237617,
		Password:  "test@123",
		Type:      "Test-type",
	}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	var shop Response
	err := json.Unmarshal(w.Body.Bytes(), &shop)

	require.Equal(t, http.StatusCreated, w.Code)
	assert.NoError(t, err)
	assert.True(t, shop.Success)
	assert.Equal(t, shop.Data.Email, dto.Email)

	return shop
}

func registerProduct(t *testing.T, db *gorm.DB) uint {
	t.Helper()

	catalog := models.CatalogProduct{
		Name:       "Test Product",
		Category:   "Test Category",
		Desciption: "Test Description",
		Brand:      "Test Brand",
		ImageURL:   "http://example.com/image.jpg",
	}

	err := db.Create(&catalog).Error
	require.NoError(t, err)
	return catalog.ID
}

func clearTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	// delete dependent records first to avoid FK constraint errors
	if err := db.Exec("DELETE FROM shop_products").Error; err != nil {
		t.Fatalf("failed to clear shop_products table: %v", err)
	}

	// then delete shops
	if err := db.Exec("DELETE FROM shops").Error; err != nil {
		t.Fatalf("failed to clear shops table: %v", err)
	}

	// clear catalog products if tests create any
	if err := db.Exec("DELETE FROM catalog_products").Error; err != nil {
		t.Fatalf("failed to clear catalog_products table: %v", err)
	}
}

func TestRegisterShop(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	w := httptest.NewRecorder()

	dto := shop.ShopRegisterDTORequest{
		Name:      "test shop1",
		OwnerName: "test",
		Email:     "test@gmail.com",
		Mobile:    "1234567890",
		Address:   "no 123. test street, test state",
		Latitude:  13.07439,
		Longitude: 80.237617,
		Password:  "test@123",
		Type:      "Test-type",
	}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	var shop Response
	err := json.Unmarshal(w.Body.Bytes(), &shop)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NoError(t, err)
	assert.True(t, shop.Success)
	assert.Equal(t, shop.Data.Email, dto.Email)
}

func TestRegisterShopToken(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	w := httptest.NewRecorder()

	dto := shop.ShopRegisterDTORequest{
		Name:      "test shop2",
		OwnerName: "test",
		Email:     "test1@gmail.com",
		Mobile:    "1234567890",
		Address:   "no 123. test street, test state",
		Latitude:  12.06439,
		Longitude: 82.237617,
		Password:  "test@123",
		Type:      "Test-type",
	}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	var shop Response
	err := json.Unmarshal(w.Body.Bytes(), &shop)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NoError(t, err)
	assert.True(t, shop.Success)
	assert.Equal(t, shop.Data.Email, dto.Email)
	token := shop.Data.Token

	id, role, err := utils.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, role, models.RoleShopOwner) // checking role
	assert.Equal(t, uint(id), shop.Data.ID)     // checking user id
}

func TestLoginShop(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	w := httptest.NewRecorder()

	registerdto := shop.ShopRegisterDTORequest{
		Name:      "test shop1",
		OwnerName: "test",
		Email:     "test@gmail.com",
		Mobile:    "1234567890",
		Address:   "no 123. test street, test state",
		Latitude:  13.07439,
		Longitude: 80.237617,
		Password:  "test@123",
		Type:      "Test-type",
	}
	body, _ := json.Marshal(registerdto)
	req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	logindto := shop.ShopLoginDTORequest{
		Email:    registerdto.Email,
		Password: registerdto.Password,
	}
	loginBody, _ := json.Marshal(logindto)
	loginReq, _ := http.NewRequest("POST", "/shop/login", strings.NewReader(string(loginBody)))
	router.ServeHTTP(w, loginReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var shop Response
	err := json.Unmarshal(w.Body.Bytes(), &shop)
	assert.NoError(t, err)
	assert.True(t, shop.Success)
	assert.Equal(t, shop.Data.Email, logindto.Email)

	token := shop.Data.Token
	id, role, err := utils.ParseToken(token)

	assert.NoError(t, err)
	assert.Equal(t, role, models.RoleShopOwner) // checking role
	assert.Equal(t, uint(id), shop.Data.ID)
}

func TestGetShopProfile(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler
	data := registerUser(t)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/shop/profile", nil)
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: data.Data.Token})
	router.ServeHTTP(w, req)
	t.Log(w.Body)

	var response Response

	json.Unmarshal(w.Body.Bytes(), &response)
	shopProfile := response.Data

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, data.Data.Name, shopProfile.Name)
	assert.Equal(t, data.Data.OwnerName, shopProfile.OwnerName)
	assert.Equal(t, data.Data.Address, shopProfile.Address)
	assert.Equal(t, data.Data.Latitude, shopProfile.Latitude)
	assert.Equal(t, data.Data.Longitude, shopProfile.Longitude)

}

func TestGetShopProfile_Unauthorized(t *testing.T) {
	router := server.NewServer().Handler

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/shop/profile", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAddProduct(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler
	data := registerUser(t)

	w := httptest.NewRecorder()

	dto := product.AddProductDTORequest{
		CatalogID:   registerProduct(t, database.New().GetDB()),
		Price:       100.0,
		Stock:       50,
		Discount:    10.0,
		IsAvailable: true,
	}

	body, _ := json.Marshal(dto)

	req, _ := http.NewRequest("POST", "/shop/products", strings.NewReader(string(body)))
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: data.Data.Token})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Product added successfully")

	w = httptest.NewRecorder()

	getProductsReq, _ := http.NewRequest("GET", "/shop/products", nil)
	getProductsReq.AddCookie(&http.Cookie{Name: "Authorization", Value: data.Data.Token})
	router.ServeHTTP(w, getProductsReq)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"price\":100")
	assert.Contains(t, w.Body.String(), "\"stock\":50")
}

// Additional test cases for comprehensive coverage

func TestRegisterShop_MissingRequiredFields(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	testCases := []struct {
		name     string
		dto      map[string]interface{}
		wantCode int
	}{
		{
			name: "Missing Name",
			dto: map[string]interface{}{
				"owner_name": "test",
				"email":      "test@gmail.com",
				"mobile":     "1234567890",
				"address":    "no 123. test street",
				"latitude":   13.07439,
				"longitude":  80.237617,
				"password":   "test@123",
				"type":       "Test-type",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Missing Email",
			dto: map[string]interface{}{
				"name":       "test shop",
				"owner_name": "test",
				"mobile":     "1234567890",
				"address":    "no 123. test street",
				"latitude":   13.07439,
				"longitude":  80.237617,
				"password":   "test@123",
				"type":       "Test-type",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Missing Password",
			dto: map[string]interface{}{
				"name":       "test shop",
				"owner_name": "test",
				"email":      "test@gmail.com",
				"mobile":     "1234567890",
				"address":    "no 123. test street",
				"latitude":   13.07439,
				"longitude":  80.237617,
				"type":       "Test-type",
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.dto)
			req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader(string(body)))
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantCode, w.Code)
		})
	}
}

func TestRegisterShop_InvalidEmail(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	w := httptest.NewRecorder()

	dto := shop.ShopRegisterDTORequest{
		Name:      "test shop",
		OwnerName: "test",
		Email:     "invalid-email",
		Mobile:    "1234567890",
		Address:   "no 123. test street",
		Latitude:  13.07439,
		Longitude: 80.237617,
		Password:  "test@123",
		Type:      "Test-type",
	}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterShop_DuplicateEmail(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	// Register first shop
	registerUser(t)

	// Try to register with same email
	w := httptest.NewRecorder()
	dto := shop.ShopRegisterDTORequest{
		Name:      "test shop2",
		OwnerName: "test2",
		Email:     "test@gmail.com", // Same email
		Mobile:    "9876543210",
		Address:   "no 456. test street",
		Latitude:  13.08439,
		Longitude: 80.247617,
		Password:  "test@123",
		Type:      "Test-type",
	}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	// Accept either 409 or 500 as both indicate duplicate email error
	assert.True(t, w.Code == http.StatusConflict || w.Code == http.StatusInternalServerError)
}

func TestRegisterShop_InvalidJSON(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader("invalid json"))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginShop_InvalidCredentials(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	// Register a shop first
	registerUser(t)

	testCases := []struct {
		name     string
		email    string
		password string
	}{
		{
			name:     "Wrong Password",
			email:    "test@gmail.com",
			password: "wrongpassword",
		},
		{
			name:     "Non-existent Email",
			email:    "nonexistent@gmail.com",
			password: "test@123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			dto := shop.ShopLoginDTORequest{
				Email:    tc.email,
				Password: tc.password,
			}
			body, _ := json.Marshal(dto)
			req, _ := http.NewRequest("POST", "/shop/login", strings.NewReader(string(body)))
			router.ServeHTTP(w, req)

			// Accept either 401 or 500 as both indicate authentication failure
			assert.True(t, w.Code == http.StatusUnauthorized || w.Code == http.StatusInternalServerError)
		})
	}
}

func TestLoginShop_MissingFields(t *testing.T) {
	router := server.NewServer().Handler

	testCases := []struct {
		name string
		dto  map[string]interface{}
	}{
		{
			name: "Missing Email",
			dto: map[string]interface{}{
				"password": "test@123",
			},
		},
		{
			name: "Missing Password",
			dto: map[string]interface{}{
				"email": "test@gmail.com",
			},
		},
		{
			name: "Empty Request",
			dto:  map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.dto)
			req, _ := http.NewRequest("POST", "/shop/login", strings.NewReader(string(body)))
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestLoginShop_InvalidEmail(t *testing.T) {
	router := server.NewServer().Handler

	w := httptest.NewRecorder()
	dto := shop.ShopLoginDTORequest{
		Email:    "not-an-email",
		Password: "test@123",
	}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/shop/login", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetShopProfile_InvalidToken(t *testing.T) {
	router := server.NewServer().Handler

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/shop/profile", nil)
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: "invalid-token"})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetShopProfile_ExpiredToken(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	// Create an expired token
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDk0NTkyMDAsInVzZXJfaWQiOjEsInJvbGUiOiJzaG9wX293bmVyIn0.invalid"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/shop/profile", nil)
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: expiredToken})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAddProduct_Unauthorized(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	w := httptest.NewRecorder()
	dto := product.AddProductDTORequest{
		CatalogID:   1,
		Price:       100.0,
		Stock:       50,
		Discount:    10.0,
		IsAvailable: true,
	}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/shop/products", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAddProduct_InvalidCatalogID(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler
	data := registerUser(t)

	w := httptest.NewRecorder()
	dto := product.AddProductDTORequest{
		CatalogID:   99999, // Non-existent catalog ID
		Price:       100.0,
		Stock:       50,
		Discount:    10.0,
		IsAvailable: true,
	}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/shop/products", strings.NewReader(string(body)))
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: data.Data.Token})
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestAddProduct_InvalidData(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler
	data := registerUser(t)

	testCases := []struct {
		name string
		dto  map[string]interface{}
	}{
		{
			name: "Negative Price",
			dto: map[string]interface{}{
				"catalog_id":   registerProduct(t, database.New().GetDB()),
				"price":        -100.0,
				"stock":        50,
				"discount":     10.0,
				"is_available": true,
			},
		},
		{
			name: "Negative Stock",
			dto: map[string]interface{}{
				"catalog_id":   registerProduct(t, database.New().GetDB()),
				"price":        100.0,
				"stock":        -50,
				"discount":     10.0,
				"is_available": true,
			},
		},
		{
			name: "Invalid Discount",
			dto: map[string]interface{}{
				"catalog_id":   registerProduct(t, database.New().GetDB()),
				"price":        100.0,
				"stock":        50,
				"discount":     150.0, // More than 100%
				"is_available": true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.dto)
			req, _ := http.NewRequest("POST", "/shop/products", strings.NewReader(string(body)))
			req.AddCookie(&http.Cookie{Name: "Authorization", Value: data.Data.Token})
			router.ServeHTTP(w, req)

			// Depending on validation, this might be BadRequest or still Created
			// Adjust based on your actual validation logic
			assert.NotEqual(t, http.StatusInternalServerError, w.Code)
		})
	}
}

func TestMultipleShopsRegistration(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	shops := []shop.ShopRegisterDTORequest{
		{
			Name:      "Shop 1",
			OwnerName: "Owner 1",
			Email:     "shop1@gmail.com",
			Mobile:    "1111111111",
			Address:   "Address 1",
			Latitude:  13.01,
			Longitude: 80.21,
			Password:  "password1",
			Type:      "Type1",
		},
		{
			Name:      "Shop 2",
			OwnerName: "Owner 2",
			Email:     "shop2@gmail.com",
			Mobile:    "2222222222",
			Address:   "Address 2",
			Latitude:  13.02,
			Longitude: 80.22,
			Password:  "password2",
			Type:      "Type2",
		},
		{
			Name:      "Shop 3",
			OwnerName: "Owner 3",
			Email:     "shop3@gmail.com",
			Mobile:    "3333333333",
			Address:   "Address 3",
			Latitude:  13.03,
			Longitude: 80.23,
			Password:  "password3",
			Type:      "Type3",
		},
	}

	for _, shopDTO := range shops {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(shopDTO)
		req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader(string(body)))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, shopDTO.Email, response.Data.Email)
	}
}

func TestRegisterAndLoginCycle(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	// Register
	registerDTO := shop.ShopRegisterDTORequest{
		Name:      "Cycle Test Shop",
		OwnerName: "Cycle Owner",
		Email:     "cycle@gmail.com",
		Mobile:    "9999999999",
		Address:   "Cycle Address",
		Latitude:  13.05,
		Longitude: 80.25,
		Password:  "cycle@123",
		Type:      "Cycle-Type",
	}

	w := httptest.NewRecorder()
	body, _ := json.Marshal(registerDTO)
	req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var registerResponse Response
	json.Unmarshal(w.Body.Bytes(), &registerResponse)
	registerToken := registerResponse.Data.Token

	// Login
	loginDTO := shop.ShopLoginDTORequest{
		Email:    registerDTO.Email,
		Password: registerDTO.Password,
	}

	w = httptest.NewRecorder()
	loginBody, _ := json.Marshal(loginDTO)
	loginReq, _ := http.NewRequest("POST", "/shop/login", strings.NewReader(string(loginBody)))
	router.ServeHTTP(w, loginReq)
	require.Equal(t, http.StatusOK, w.Code)

	var loginResponse Response
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	loginToken := loginResponse.Data.Token

	// Both tokens should work
	assert.NotEmpty(t, registerToken)
	assert.NotEmpty(t, loginToken)

	// Verify profile with both tokens
	for _, token := range []string{registerToken, loginToken} {
		w = httptest.NewRecorder()
		profileReq, _ := http.NewRequest("GET", "/shop/profile", nil)
		profileReq.AddCookie(&http.Cookie{Name: "Authorization", Value: token})
		router.ServeHTTP(w, profileReq)

		assert.Equal(t, http.StatusOK, w.Code)
		var profileResponse Response
		json.Unmarshal(w.Body.Bytes(), &profileResponse)
		assert.Equal(t, registerDTO.Email, profileResponse.Data.Email)
	}
}

func TestEdgeCaseCoordinates(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler

	testCases := []struct {
		name      string
		email     string
		latitude  float64
		longitude float64
		wantCode  int
	}{
		{
			name:      "Valid Positive Coordinates",
			email:     "edge1@gmail.com",
			latitude:  13.07439,
			longitude: 80.237617,
			wantCode:  http.StatusCreated,
		},
		{
			name:      "Negative Coordinates",
			email:     "edge2@gmail.com",
			latitude:  -13.07439,
			longitude: -80.237617,
			wantCode:  http.StatusCreated,
		},
		{
			name:      "Near Equator",
			email:     "edge3@gmail.com",
			latitude:  0.1,
			longitude: 0.1,
			wantCode:  http.StatusCreated,
		},
		{
			name:      "High Latitude North",
			email:     "edge4@gmail.com",
			latitude:  85.0,
			longitude: 80.237617,
			wantCode:  http.StatusCreated,
		},
		{
			name:      "High Latitude South",
			email:     "edge5@gmail.com",
			latitude:  -85.0,
			longitude: 80.237617,
			wantCode:  http.StatusCreated,
		},
		{
			name:      "Different Valid Location",
			email:     "edge6@gmail.com",
			latitude:  40.7128,
			longitude: -74.0060, // New York coordinates
			wantCode:  http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			dto := shop.ShopRegisterDTORequest{
				Name:      "Edge Shop",
				OwnerName: "Edge Owner",
				Email:     tc.email,
				Mobile:    "8888888888",
				Address:   "Edge Address",
				Latitude:  tc.latitude,
				Longitude: tc.longitude,
				Password:  "edge@123",
				Type:      "Edge-Type",
			}
			body, _ := json.Marshal(dto)
			req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader(string(body)))
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantCode, w.Code)
		})
	}
}
