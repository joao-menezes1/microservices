package api

import (
    //"fmt"
    "log"
    "google.golang.org/grpc/codes"   
    "google.golang.org/grpc/status"  
    "github.com/ruandg/microservices/order/internal/application/core/domain"
    "github.com/ruandg/microservices/order/internal/ports"
)

type Application struct {
    db      ports.DBPort
    payment ports.PaymentPort
    shipping ports.ShippingPort // Adicionamos a porta de Shipping aqui
}

// Atualizamos o construtor para receber o shipping
func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
    return &Application{
        db:       db,
        payment:  payment,
        shipping: shipping,
    }
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {

    // --- 1. VALIDAÇÃO DE ESTOQUE (Requisito 1.2) ---
    // Verificamos se cada item do pedido existe na tabela 'products'
    for _, item := range order.OrderItems {
        _, err := a.db.GetProduct(item.ProductCode)
        if err != nil {
            // Se o produto não existe, retornamos erro imediatamente
            return domain.Order{}, status.Errorf(codes.NotFound, "produto '%s' não encontrado no estoque", item.ProductCode)
        }
    }

    // --- 2. VALIDAÇÃO DE QUANTIDADE TOTAL ---
    totalItems := 0
    for _, item := range order.OrderItems {
        totalItems += int(item.Quantity)
    }

    if totalItems > 50 {
        return domain.Order{}, status.Errorf(codes.InvalidArgument, "total items exceeds the limit of 50")
    }

    // --- 3. PERSISTÊNCIA INICIAL ---
    err := a.db.Save(&order)
    if err != nil {
        return domain.Order{}, err
    }

    // --- 4. PAGAMENTO ---
    paymentErr := a.payment.Charge(&order)
    if paymentErr != nil {
        order.Status = "Canceled"
        _ = a.db.Update(&order)
        return domain.Order{}, paymentErr
    }

    // --- 5. SHIPPING (Requisito: Só chama se o pagamento der certo) ---
    // Chamamos o microsserviço que você criou na porta 3002
    deliveryDays, shippingErr := a.shipping.Create(&order)
    if shippingErr != nil {
        // Aqui decidimos: se o frete falhar, cancelamos ou apenas logamos? 
        // Geralmente, se o pagamento foi feito, tentamos tratar o erro ou logar.
        log.Printf("Pagamento aprovado, mas erro ao solicitar entrega: %v", shippingErr)
    } else {
        log.Printf("Entrega agendada com sucesso! Prazo: %d dias", deliveryDays)
    }

    // --- 6. FINALIZAÇÃO ---
    order.Status = "Paid"
    updateErr := a.db.Update(&order)
    if updateErr != nil {
        return domain.Order{}, updateErr
    }
    
    return order, nil
}