package ports

import "github.com/ruandg/microservices/shipping/internal/application/core/domain"

type ShippingPort interface {
	Create(shipping domain.Shipping) (domain.Shipping, error)
}