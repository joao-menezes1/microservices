package api

import (
	"google.golang.org/grpc/codes"   
    "google.golang.org/grpc/status"  
	"github.com/ruandg/microservices/order/internal/application/core/domain"
	"github.com/ruandg/microservices/order/internal/ports"
)

type Application struct {
	db ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db: db,
		payment: payment,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {

	totalItems := 0
	for _, item := range order.OrderItems {
		totalItems += int(item.Quantity)
	}

	if totalItems > 50 {
		return domain.Order{}, status.Errorf(codes.InvalidArgument, "total items exceeds the limit of 50")
	}


	err := a.db.Save(&order)
	if err != nil {
		order.Status = "Canceled"
		_ = a.db.Update(&order)
		return domain.Order{}, err
	}

	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		order.Status = "Canceled"
		_ = a.db.Update(&order)
		return domain.Order{}, paymentErr
	}

	order.Status = "Paid"
	updateErr := a.db.Update(&order)
	if updateErr != nil {
		return domain.Order{}, updateErr
	}
	return order, nil
}