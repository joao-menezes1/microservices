package domain

// Product representa a tabela que você acabou de criar no MySQL
type Product struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}