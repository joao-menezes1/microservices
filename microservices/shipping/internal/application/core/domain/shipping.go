package domain

type ShippingItem struct {
	ProductCode string  `json:"product_code"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int32   `json:"quantity"`
}

type Shipping struct {
	ID           int64          `json:"id"`
	OrderID      int64          `json:"order_id"`
	CustomerID   int64          `json:"customer_id"`
	Items        []ShippingItem `json:"items"`
	DeliveryDays int32          `json:"delivery_days"`
}