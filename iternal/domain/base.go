package domain

type Base struct {
	ID int `gorm:"type:int;" gorm:"primaryKey" json:"id" gorm:"index:id"`
}
