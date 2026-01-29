package ports

import "github.com/ruandg/microservices/order/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(*domain.Order) error
	Update(*domain.Order) error
	GetProduct(code string) (domain.Product, error)
}