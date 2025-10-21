package product

type AddProductDTORequest struct {
	CatalogID   uint    `json:"catalog_id" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	Discount    float64 `json:"discount" binding:"gte=0"`
	IsAvailable bool    `json:"is_available"`
}

type ProductUpdateDTORequest struct {
	ID          uint    `json:"id" binding:"required"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
	Discount    float64 `json:"discount" binding:"omitempty,gte=0"`
	IsAvailable bool    `json:"is_available"`
}
