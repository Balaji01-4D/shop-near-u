package product

import "shop-near-u/internal/models"

type Service struct {
	repository *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repository: r}
}

func (s *Service) AddProduct(dto *AddProductDTORequest, shopID uint) error {
	product := &models.ShopProduct{
		ShopID:      shopID,
		CatalogID:   dto.CatalogID,
		Price:       dto.Price,
		Stock:       dto.Stock,
		Discount:    dto.Discount,
		IsAvailable: dto.IsAvailable,
	}

	return s.repository.AddProduct(product)
}

func (s *Service) GetProductsByShopID(shopID uint) ([]models.ShopProduct, error) {
	return s.repository.GetProductsByShopID(shopID)
}

func (s *Service) GetProductByID(productID uint) (*models.ShopProduct, error) {
	return s.repository.GetProductByID(productID)
}

func (s *Service) UpdateProduct(productDTO *ProductUpdateDTORequest) error {
	product := &models.ShopProduct{
		ID:          productDTO.ID,
		Price:       productDTO.Price,
		Stock:       productDTO.Stock,
		Discount:    productDTO.Discount,
		IsAvailable: productDTO.IsAvailable,
	}
	return s.repository.UpdateProduct(product)
}

func (s *Service) DeleteProduct(productID uint) error {
	return s.repository.DeleteProduct(productID)
}
