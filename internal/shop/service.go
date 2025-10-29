package shop

import (
	"shop-near-u/internal/models"
	"shop-near-u/internal/utils"

	"github.com/restayway/gogis"
)

type Service struct {
	repository *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repository: r}
}

func (s *Service) RegisterShop(registerDTO *ShopRegisterDTORequest) (*models.Shop, error) {

	password, err := utils.HashPassword(registerDTO.Password)
	if err != nil {
		return nil, err
	}
	shop := &models.Shop{
		Name:      registerDTO.Name,
		OwnerName: registerDTO.OwnerName,
		Type:      registerDTO.Type,
		Password:  password,
		Email:     registerDTO.Email,
		Mobile:    registerDTO.Mobile,
		Address:   registerDTO.Address,
		Latitude:  registerDTO.Latitude,
		Longitude: registerDTO.Longitude,
		Location: gogis.Point{
			Lng: registerDTO.Longitude,
			Lat: registerDTO.Latitude,
		},
	}

	if err := s.repository.Create(shop); err != nil {
		return nil, err
	}

	return shop, nil
}

func (s *Service) AuthenticateShop(request *ShopLoginDTORequest) (*models.Shop, error) {
	shop, err := s.repository.FindByEmail(request.Email)
	if err != nil {
		return nil, err
	}

	if shop == nil {
		return nil, nil
	}

	if err := utils.CheckPasswordHash(request.Password, shop.Password); err != nil {
		return nil, nil
	}

	return shop, nil
}

func (s *Service) GetShopByID(shopID uint) (*models.Shop, error) {
	return  s.repository.FindByID(shopID)

}

func (s *Service) GetNearbyShops(lat float64, lon float64, radius float64, limit int) ([]NearByShopsDTORespone, error) {
	shops, err := s.repository.FindNearbyShops(lat, lon, radius, limit)
	if err != nil {
		return nil, err
	}

	return shops, nil
}

func (s *Service) SubscribeShop(userID uint, shopID uint) (uint, error) {
	subscriberCount, err := s.repository.SubscribeShop(shopID, userID)
	if err != nil {
		return 0, err
	}
	return subscriberCount, nil
}

func (s *Service) UpdateShopStatus(shopID uint, status bool) error {
	return s.repository.UpdateShopStatus(shopID, status)
}

func (s *Service) UnsubscribeShop(userID uint, shopID uint) (uint, error) {
	subscriberCount, err := s.repository.UnsubscribeShop(shopID, userID)
	if err != nil {
		return 0, err
	}
	return subscriberCount, nil
}

func (s *Service) GetShopDetails(shopID uint, userID uint) (*models.Shop, bool, error) {
	shop, isSubscribed, err := s.repository.GetShopDetails(shopID, userID)
	if err != nil {
		return nil, false, err
	}
	return shop, isSubscribed, nil
}
