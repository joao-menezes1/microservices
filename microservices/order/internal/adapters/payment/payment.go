package payment_adapter

import (
    "context"
    "time"

    grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
    "github.com/ruandg/microservices-proto/golang/payment"
    "github.com/ruandg/microservices/order/internal/application/core/domain"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
    payment payment.PaymentClient 
}

// func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
//     var opts []grpc.DialOption
//     opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
//     conn, err := grpc.Dial(paymentServiceUrl, opts...)
//     if err != nil {
//         return nil, err
//     }
//     client := payment.NewPaymentClient(conn)
//     return &Adapter{payment: client}, nil
// }

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
    var opts []grpc.DialOption

    // Configurando o Interceptor de Retentativas
    opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
        grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted, codes.DeadlineExceeded),
        grpc_retry.WithMax(5),
        grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
    )))

    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

    // Estabelecendo a conexão com o Interceptor
    conn, err := grpc.Dial(paymentServiceUrl, opts...)
    if err != nil {
        return nil, err
    }
    client := payment.NewPaymentClient(conn)
    return &Adapter{payment: client}, nil
}

// func (a *Adapter) Charge(order *domain.Order) error {
//     // 1. Criar um contexto com timeout de 2 segundos
//     ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//     defer cancel() // Importante para liberar recursos após o término da chamada

//     // 2. Passar o 'ctx' criado em vez de context.Background()
//     _, err := a.payment.Create(ctx, &payment.CreatePaymentRequest{
//         UserId:     order.CustomerID,
//         OrderId:    order.ID,
//         TotalPrice: order.TotalPrice(),
//     })

//     // 3. Verificar se ocorreu um erro de Timeout
//     if err != nil {
//         // Converte o erro genérico para um status gRPC
//         if st, ok := status.FromError(err); ok {
//             if st.Code() == codes.DeadlineExceeded {
//                 log.Printf("Erro: O limite de tempo (timeout) foi atingido ao chamar o serviço de Payment para o pedido %d", order.ID)
//             }
//         }
//         return err
//     }

//     return nil
// }

func (a *Adapter) Charge(order *domain.Order) error {
    // O interceptor vai tentar 5 vezes, e cada tentativa terá seu próprio timeout de 2s
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    _, err := a.payment.Create(ctx, &payment.CreatePaymentRequest{
        UserId:     order.CustomerID,
        OrderId:    order.ID,
        TotalPrice: order.TotalPrice(),
    })

    return err
}