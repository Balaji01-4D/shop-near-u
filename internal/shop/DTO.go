package shop

type ShopRegisterDTORequest struct {
	Name      string `json:"name" binding:"required"`
	OwnerName string `json:"owner_name" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Mobile    string `json:"mobile" binding:"required"`
	Password  string `json:"password" binding:"required"`

	Address   string  `json:"address" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type ShopRegisterDTOResponse struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	OwnerName string  `json:"owner_name"`
	Type      string  `json:"type"`
	Email     string  `json:"email"`
	Mobile    string  `json:"mobile"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
