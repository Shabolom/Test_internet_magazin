package repository

//
//import (
//	"Arkadiy_2Service/config"
//	"Arkadiy_2Service/iternal/domain"
//	"Arkadiy_2Service/iternal/model"
//	"Arkadiy_2Service/tools"
//	"errors"
//	"fmt"
//	"github.com/jinzhu/gorm"
//)
//
//type Repo struct {
//}
//
//func NewRepo() *Repo {
//	return &Repo{}
//}
//
//func (r *Repo) CheckAndUpdateOPalette(productID, count, orderID int) error {
//	var pp domain.PalettesProducts
//	var countProducts []domain.PalettesProducts
//	var paMass []model.Order
//	sumProduct := 0
//
//	err := config.DB.
//		Table("palettes_products").
//		Where("product_id =?", productID).
//		Find(&countProducts).
//		Error
//
//	if err != nil {
//		return err
//	}
//
//	for _, ppPart := range countProducts {
//		sumProduct += ppPart.Count
//	}
//
//	if sumProduct < count {
//		return errors.New("не хватает товара")
//	}
//
//	err = config.DB.
//		Table("palettes_products").
//		Where("product_id =? AND status = true", productID).
//		Find(&pp).
//		Error
//	if err != nil {
//		return err
//	}
//
//	if pp.Count < count {
//
//		if pp.Count != 0 {
//			err = config.DB.
//				Table("palettes_products").
//				Where("product_id =? AND status = true", productID).
//				UpdateColumn("count", gorm.Expr("count - ?", pp.Count)).
//				Error
//			if err != nil {
//				return err
//			}
//			count = count - pp.Count
//
//			paMass = append(paMass, model.Order{
//				Count:     count,
//				ProductID: productID,
//				PaletteID: pp.PaletteID,
//				OrderID:   orderID,
//			})
//		}
//
//		for count != 0 {
//			err = config.DB.
//				Table("palettes_products").
//				Where("product_id =? AND count > 0", productID).
//				Find(&pp).
//				Error
//			if err != nil {
//				return err
//			}
//
//			if pp.Count >= count {
//				err = config.DB.
//					Table("palettes_products").
//					Where("product_id =? AND palette_id =?", productID, pp.PaletteID).
//					UpdateColumn("count", gorm.Expr("count -?", count)).
//					Error
//				if err != nil {
//					return err
//				}
//				count -= count
//				paMass = append(paMass, model.Order{
//					Count:     count,
//					ProductID: productID,
//					PaletteID: pp.PaletteID,
//					OrderID:   orderID,
//				})
//
//				err = r.PostOrdersProducts(paMass)
//				if err != nil {
//					return err
//				}
//
//				fmt.Println("заказ ", orderID, " занесен в базу.")
//				return nil
//			}
//
//			err = config.DB.
//				Table("palettes_products").
//				Where("product_id =? AND palette_id =?", productID, pp.PaletteID).
//				UpdateColumn("count", gorm.Expr("count -?", pp.Count)).
//				Error
//			if err != nil {
//				return err
//			}
//			count -= pp.Count
//			paMass = append(paMass, model.Order{
//				Count:     count,
//				ProductID: productID,
//				PaletteID: pp.PaletteID,
//				OrderID:   orderID,
//			})
//		}
//	}
//
//	err = config.DB.
//		Table("palettes_products").
//		Where("product_id =? AND status = true", productID).
//		UpdateColumn("count", gorm.Expr("count -?", count)).
//		Error
//	if err != nil {
//		return err
//	}
//	paMass = append(paMass, model.Order{
//		Count:     count,
//		ProductID: productID,
//		PaletteID: pp.PaletteID,
//		OrderID:   orderID,
//	})
//
//	err = r.PostOrdersProducts(paMass)
//	if err != nil {
//		return err
//	}
//
//	fmt.Println("заказ ", orderID, " занесен в базу.")
//	return nil
//}
//
//func (r *Repo) PostOrdersProducts(orders []model.Order) error {
//	tx := config.DB.Begin()
//	orderPattern := orders[0]
//
//	err := r.FindOrder(orderPattern.OrderID)
//
//	if err != nil {
//		err = r.CreateOrder(orderPattern.OrderID)
//
//		if err != nil {
//			return err
//		}
//	}
//
//	for _, order := range orders {
//		err := tx.
//			Table("orders_products").
//			Create(&domain.OrdersProducts{
//				ProductID: order.ProductID,
//				OrderID:   order.OrderID,
//				Count:     order.Count,
//				PaletteID: order.PaletteID,
//			}).
//			Error
//		if err != nil {
//			tx.Rollback()
//			return err
//		}
//	}
//	tx.Commit()
//
//	return nil
//}
//
//func (r *Repo) TakeOrders(orders []int) error {
//	var collectedOrder []model.Assemblings
//	var collectedOrders [][]model.Assemblings
//
//	tx := config.DB.Begin()
//
//	for _, orderID := range orders {
//		err := tx.Select("op.product_id as product_id, o.id as order_id, p.name as product_name, op.count as count, pa.name as palette").
//			Table("orders o").
//			Joins("JOIN orders_products op ON o.id = op.order_id AND(o.id =?)", orderID).
//			Joins("JOIN products p ON op.product_id = p.id").
//			Joins("JOIN palettes pa ON op.palette_id = pa.id").
//			Order("p.name DESC").
//			Scan(&collectedOrder).
//			Error
//		if err != nil {
//			tx.Rollback()
//			return err
//		}
//		collectedOrders = append(collectedOrders, collectedOrder)
//	}
//	tx.Commit()
//
//	tools.Sort(collectedOrders)
//
//	return nil
//}
//
//func (r *Repo) FindOrder(orderID int) error {
//	var order domain.Order
//
//	err := config.DB.
//		Where("id =?", orderID).
//		First(&order).
//		Error
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (r *Repo) CreateOrder(orderID int) error {
//	err := config.DB.
//		Create(&domain.Order{
//			Base: domain.Base{ID: orderID},
//		}).
//		Error
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
