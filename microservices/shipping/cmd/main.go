package main

import (
	"log"
	"net"

	"github.com/ruandg/microservices-proto/golang/shipping"
	"github.com/ruandg/microservices/shipping/internal/adapters/grpc"
	"github.com/ruandg/microservices/shipping/internal/application/core/api"
	g "google.golang.org/grpc"
)

func main() {
	// 1. Inicializa o Coração (Service) que criamos no Passo 1
	shippingService := api.NewService()

	// 2. Inicializa o Adaptador gRPC que criamos no Passo 2
	grpcAdapter := grpc.NewAdapter(shippingService)

	// 3. Define que o serviço vai ouvir na porta 3002
	listen, err := net.Listen("tcp", ":3002")
	if err != nil {
		log.Fatalf("Falha ao abrir a porta 3002: %v", err)
	}

	// 4. Cria o servidor gRPC e registra o nosso serviço nele
	server := g.NewServer()
	shipping.RegisterShippingServer(server, grpcAdapter)

	log.Println("✅ Serviço Shipping rodando na porta 3002...")
	
	// 5. Começa a servir as requisições
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Falha ao rodar o servidor gRPC: %v", err)
	}
}