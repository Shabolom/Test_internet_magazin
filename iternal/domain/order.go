package domain

type Order struct {
	Base
	Products []Product `gorm:"many2many:orders_products"`
}
