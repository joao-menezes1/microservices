package shipping

import (
	"context"
	"github.com/ruandg/microservices-proto/golang/shipping"
	"github.com/ruandg/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	shipping shipping.ShippingClient
}

func NewAdapter(address string) (*Adapter, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := shipping.NewShippingClient(conn)
	return &Adapter{shipping: client}, nil
}

func (a *Adapter) Create(order *domain.Order) (int32, error) {
	// 1. Criar a lista de itens no novo formato ShippingItem definido no proto
	var shippingItems []*shipping.ShippingItem

	for _, item := range order.OrderItems {
		shippingItems = append(shippingItems, &shipping.ShippingItem{
			ProductCode: item.ProductCode,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
		})
	}

	// 2. Enviar a requisição completa com os itens mapeados
	res, err := a.shipping.Create(context.Background(), &shipping.CreateShippingRequest{
		OrderId:    int64(order.ID),
		CustomerId: order.CustomerID,
		Items:      shippingItems, // Agora passamos a lista aqui
	})

	if err != nil {
		return 0, err
	}
	
	return res.DeliveryDays, nil
}