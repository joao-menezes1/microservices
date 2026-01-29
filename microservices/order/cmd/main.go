package main

import (
    "log"
    "github.com/ruandg/microservices/order/config"
    "github.com/ruandg/microservices/order/internal/adapters/db"
    "github.com/ruandg/microservices/order/internal/adapters/payment"
    "github.com/ruandg/microservices/order/internal/adapters/shipping" // 1. Importe o novo adapter
    "github.com/ruandg/microservices/order/internal/adapters/grpc"
    "github.com/ruandg/microservices/order/internal/application/core/api"
)

func main() {
    // Inicialização do Banco de Dados
    dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
    if err != nil {
        log.Fatalf("failed to connect to database. Error: %v", err)
    }

    // Inicialização do Adapter de Pagamento (Porta 3001)
    paymentAdapter, err := payment_adapter.NewAdapter(config.GetPaymentServiceURL())
    if err != nil {
        log.Fatalf("falha ao iniciar o payment stub. Error: %v", err)
    }

    // 2. Inicialização do Adapter de Shipping (Porta 3002)
    // Certifique-se que config.GetShippingServiceURL() retorne "localhost:3002"
    shippingAdapter, err := shipping.NewAdapter(config.GetShippingServiceURL())
    if err != nil {
        log.Fatalf("falha ao iniciar o shipping stub. Error: %v", err)
    }

    // 3. Atualização da Application: agora passamos db, payment E shipping
    application := api.NewApplication(dbAdapter, paymentAdapter, shippingAdapter)

    // Inicialização do Servidor gRPC do Order (Porta 3000)
    grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
    grpcAdapter.Run()
}