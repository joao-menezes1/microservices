package ports

import "github.com/ruandg/microservices/order/internal/application/core/domain"

type ShippingPort interface {
    // Esta função recebe um ponteiro para o pedido e retorna o prazo em dias
    Create(order *domain.Order) (int32, error)
}