package domain

type Palette struct {
	Base
	Name           string
	Products       []Product `gorm:"many2many:palettes_products"`
	OrdersProducts OrdersProducts
}
