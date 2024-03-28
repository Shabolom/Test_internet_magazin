package domain

type OrdersProducts struct {
	ProductID int `gorm:"primaryKey"`
	OrderID   int `gorm:"primaryKey"`
	PaletteID int
	Count     int
}
