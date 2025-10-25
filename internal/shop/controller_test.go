package shop_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shop-near-u/internal/database"
	"shop-near-u/internal/models"
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
	Token	 string `json:"token"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    Data `json:"data,omitempty"`
}

func registerUser(t *testing.T) Response{
	t.Helper()
	router := server.NewServer().Handler

	w := httptest.NewRecorder()

	dto := shop.ShopRegisterDTORequest {
		Name: "test shop1",
		OwnerName: "test",
		Email: "test@gmail.com",
		Mobile: "1234567890",
		Address: "no 123. test street, test state",
		Latitude: 13.07439,
		Longitude: 80.237617,
		Password: "test@123",
		Type: "Test-type",

	}
	body, _  := json.Marshal(dto)
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


func clearTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Exec("DELETE FROM shops").Error; err != nil {
		t.Fatalf("failed to clear shops table: %v", err)
	}
}

func TestRegisterShop(t *testing.T) {
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler


	w := httptest.NewRecorder()

	dto := shop.ShopRegisterDTORequest {
		Name: "test shop1",
		OwnerName: "test",
		Email: "test@gmail.com",
		Mobile: "1234567890",
		Address: "no 123. test street, test state",
		Latitude: 13.07439,
		Longitude: 80.237617,
		Password: "test@123",
		Type: "Test-type",

	}
	body, _  := json.Marshal(dto)
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

	dto := shop.ShopRegisterDTORequest {
		Name: "test shop2",
		OwnerName: "test",
		Email: "test1@gmail.com",
		Mobile: "1234567890",
		Address: "no 123. test street, test state",
		Latitude: 12.06439,
		Longitude: 82.237617,
		Password: "test@123",
		Type: "Test-type",

	}
	body, _  := json.Marshal(dto)
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
	assert.Equal(t, role, models.RoleShopOwner)		// checking role
	assert.Equal(t, uint(id), shop.Data.ID)			// checking user id
}

func TestLoginShop(t *testing.T){
	defer clearTestDB(t, database.New().GetDB())
	router := server.NewServer().Handler


	w := httptest.NewRecorder()

	registerdto := shop.ShopRegisterDTORequest {
		Name: "test shop1",
		OwnerName: "test",
		Email: "test@gmail.com",
		Mobile: "1234567890",
		Address: "no 123. test street, test state",
		Latitude: 13.07439,
		Longitude: 80.237617,
		Password: "test@123",
		Type: "Test-type",

	}
	body, _  := json.Marshal(registerdto)
	req, _ := http.NewRequest("POST", "/shop/register", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)


	assert.Equal(t, http.StatusCreated, w.Code)

	
	w = httptest.NewRecorder()
	logindto := shop.ShopLoginDTORequest{
		Email: registerdto.Email,
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
	assert.Equal(t, role, models.RoleShopOwner)		// checking role
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