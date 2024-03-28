package model

type Place struct {
	Place      string
	Assemblies []Assembling
}

type Assembling struct {
	ProductID int
	OrderID   int
	Count     int
}
