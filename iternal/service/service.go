package service

import (
	"Arkadiy_2Service/iternal/repository"
)

var repo = repository.NewRepo()

//func PostOrder(orderID, productID, productsCount int) error {
//	result, err := repo.CheckAndUpdateOPalette(productID, productsCount)
//	if err != nil {
//		return err
//	}
//}
