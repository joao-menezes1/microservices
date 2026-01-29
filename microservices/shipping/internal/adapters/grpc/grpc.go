package grpc

import (
	"github.com/ruandg/microservices-proto/golang/shipping"
	"github.com/ruandg/microservices/shipping/internal/ports"
	//"google.golang.org/grpc"
)

type Adapter struct {
	api ports.ShippingPort // Porta que liga o adaptador ao serviço
	shipping.UnimplementedShippingServer
}

func NewAdapter(api ports.ShippingPort) *Adapter {
	return &Adapter{api: api}
}