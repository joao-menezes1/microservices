package db

import (

	"fmt"
	"github.com/ruandg/microservices/order/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

)


type Order struct {
	gorm.Model
	CustomerID int64
	Status    string
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderID     uint
	ProductCode string
	UnitPrice   float32
	Quantity    int32
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}
	err := db.AutoMigrate(&Order{}, &OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}

func (a Adapter) Get(id string) (domain.Order, error) {
	var orderEntity Order
	res := a.db.First(&orderEntity, id)
	var orderItems []domain.OrderItem

	for _, orderItem := range orderEntity.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	order := domain.Order{
		ID:         int64(orderEntity.ID),
		CustomerID: orderEntity.CustomerID,
		Status:     orderEntity.Status,
		OrderItems: orderItems,
		CreatedAt:  orderEntity.CreatedAt.Unix(),
	}

	return order, res.Error
}


func (a Adapter) Save(order *domain.Order) error {
	var orderItems []OrderItem
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	orderModel := Order{
		CustomerID: order.CustomerID,
		Status:     order.Status,
		OrderItems: orderItems,	
	}
	res := a.db.Create(&orderModel)
	if res.Error == nil {
		order.ID = int64(orderModel.ID)
	}
	return res.Error
}

// GetProduct busca um produto na tabela 'products' pelo código
func (a Adapter) GetProduct(code string) (domain.Product, error) {
	var productEntity domain.Product
	
	// O GORM vai rodar: SELECT * FROM products WHERE code = 'P1' LIMIT 1;
	res := a.db.First(&productEntity, "code = ?", code)
	
	if res.Error != nil {
		return domain.Product{}, res.Error
	}
	
	return productEntity, nil
}


func (a Adapter) Update(order *domain.Order) error {
    // Usamos &Order{} (a struct local do banco) para que o GORM saiba o esquema.
    // Usamos Where para garantir que estamos atualizando o ID correto.
    return a.db.Model(&Order{}).
        Where("id = ?", order.ID).
        Update("status", order.Status).Error
}