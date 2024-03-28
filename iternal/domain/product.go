package domain

type Product struct {
	Base
	Name     string
	Orders   []Order   `gorm:"many2many:orders_products"`
	Palettes []Palette `gorm:"many2many:palettes_products"`
}
