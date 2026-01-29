package api

import (
	"github.com/ruandg/microservices/shipping/internal/application/core/domain"
)

type Service struct {}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(shipping domain.Shipping) (domain.Shipping, error) {
	totalUnits := int32(0)
	for _, item := range shipping.Items {
		totalUnits += item.Quantity
	}

	// Regra: 1 dia fixo + (unidades / 5)
	// Exemplo: 4 un / 5 = 0 -> 1 dia total
	// Exemplo: 5 un / 5 = 1 -> 2 dias total
	shipping.DeliveryDays = 1 + (totalUnits / 5)

	return shipping, nil
}