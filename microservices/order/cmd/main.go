package main

import (
	"log"
	"github.com/ruandg/microservices/order/config"
	"github.com/ruandg/microservices/order/internal/adapters/db"

	"github.com/ruandg/microservices/order/internal/adapters/payment"

	"github.com/ruandg/microservices/order/internal/adapters/grpc"

	"github.com/ruandg/microservices/order/internal/application/core/api"

)

func main() {
	dbAdapter , err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to connect to database. Error: %v", err)
	}

	paymentAdapter, err := payment_adapter.NewAdapter(config.GetPaymentServiceURL())
	if err != nil {
		log.Fatalf("falha ao iniciar o payment stub. Error: %v", err)
	}


	application := api.NewApplication(dbAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
