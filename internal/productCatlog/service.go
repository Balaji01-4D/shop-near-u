package productcatlog

import "shop-near-u/internal/models"

type Service struct {
	repository *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repository: r}
}

func (s *Service) CreateCatalogProduct(product *CreateCatalogProductDTO) error {
	catalogProduct := &models.CatalogProduct{
		Name:       product.Name,
		Brand:      product.Brand,
		Category:   product.Category,
		Desciption: product.Description,
		ImageURL:   product.ImageURL,
	}
	return s.repository.CreateCatalogProduct(catalogProduct)
}

func (s *Service) SuggestCatalogProducts(keyword string, limit int) (*[]models.CatalogProduct, error) {
	return s.repository.Suggest(keyword, limit)
}
