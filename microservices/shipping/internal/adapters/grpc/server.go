package grpc

import (
	"context"
	"github.com/ruandg/microservices-proto/golang/shipping"
	"github.com/ruandg/microservices/shipping/internal/application/core/domain"
)

// Create é a função que o gRPC chama quando chega uma requisição
func (a *Adapter) Create(ctx context.Context, request *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error) {
	var items []domain.ShippingItem
	
	// Convertendo os itens que vieram do gRPC para o nosso Domínio
	for _, item := range request.Items {
		items = append(items, domain.ShippingItem{
			ProductCode: item.ProductCode,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
		})
	}

	newShipping := domain.Shipping{
		OrderID:    request.OrderId,
		CustomerID: request.CustomerId,
		Items:      items,
	}

	// Chama a lógica de negócio (o service do Passo 1)
	result, err := a.api.Create(newShipping)
	if err != nil {
		return nil, err
	}

	// Retorna a resposta no formato gRPC
	return &shipping.CreateShippingResponse{
		ShippingId:   result.ID,
		DeliveryDays: result.DeliveryDays,
	}, nil
}