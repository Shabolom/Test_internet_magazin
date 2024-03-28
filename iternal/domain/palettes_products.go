package domain

type PalettesProducts struct {
	PaletteID int  `gorm:"primaryKey"`
	ProductID int  `gorm:"primaryKey"`
	Count     int  `gorm:"colum:count; type:int"`
	Status    bool `gorm:"colum:status; type:bool"`
}
