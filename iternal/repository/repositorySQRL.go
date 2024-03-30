package repository

import (
	"Arkadiy_2Service/config"
	"Arkadiy_2Service/iternal/model"
	"Arkadiy_2Service/tools"
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"math/rand"
	"time"
)

type Repository2 struct {
}

func NewRepo2() *Repository2 {
	return &Repository2{}
}

var ctx, dbQuerySpan = otel.Tracer("").Start(context.Background(), "db_query: DeleteDevice")

func (r *Repository2) AssemblingOrder(ordersID []int) error {
	var orderPart model.Assemblings
	var order []model.Assemblings

	for _, orderID := range ordersID {

		sql, args, err := config.Sq.
			Select("op.product_id", "o.id", "p.name", "op.product_count", "pa.name").
			From("orders o").
			Where("o.id = $1", orderID).
			Join("orders_products op ON op.order_id = o.id").
			Join("palettes pa ON op.palette_id = pa.id").
			Join("products p ON op.product_id = p.id").
			GroupBy("op.product_id, o.id, p.name, pa.name, op.product_count").
			ToSql()
		if err != nil {
			return err
		}

		dbQuerySpan.SetAttributes(attribute.String("sql.query", sql))
		rows, err := config.Pool.Query(ctx, sql, args...)
		if err != nil {
			return err
		}

		for rows.Next() {
			err = rows.Scan(&orderPart.ProductID, &orderPart.OrderID, &orderPart.ProductName, &orderPart.Count, &orderPart.Palette)
			if err != nil {
				return err
			}
			fmt.Println(rows.Values())

			order = append(order, orderPart)
		}
		defer rows.Close()
	}

	tools.Sort(order)
	defer dbQuerySpan.End()
	return nil
}

func (r *Repository2) CheckAndUpdateOPalette(orderID, productID, count int) error {
	var assembledOrder []model.Order

	answer, err := r.CheckSum(productID, count)
	if err != nil {
		return err
	}
	if answer != true {
		return errors.New("не хватает товара")
	}

	paletteCount, paletteID, answer, err := r.CheckPaletteCount(true, productID, count)
	if err != nil {
		fmt.Println(1)
		return err
	}

	if answer != true {

		if paletteCount > 0 {
			err2 := r.UpdatePalette(paletteID, productID, paletteCount, paletteCount)
			if err2 != nil {
				fmt.Println(2)
				return err2
			}

			assembledOrder = append(assembledOrder, model.Order{
				OrderID:   orderID,
				Count:     paletteCount,
				ProductID: productID,
				PaletteID: paletteID,
			})

			count -= paletteCount
		}

		for count != 0 {

			paletteCount, paletteID, answer, err = r.CheckPaletteCount(false, productID, count)

			if answer == true {
				err = r.UpdatePalette(paletteID, productID, paletteCount, count)
				if err != nil {
					fmt.Println(2)
					return err
				}

				assembledOrder = append(assembledOrder, model.Order{
					OrderID:   orderID,
					Count:     count,
					ProductID: productID,
					PaletteID: paletteID,
				})

				count = 0

				err = r.InsertOrder(orderID)
				if err != nil {
					fmt.Println(3)
					return err
				}

				err = r.UpdateOrdersProducts(assembledOrder)
				if err != nil {
					fmt.Println(4)
					return err
				}

				return nil
			}

			err = r.UpdatePalette(paletteID, productID, paletteCount, paletteCount)
			if err != nil {
				fmt.Println(2)
				return err
			}

			assembledOrder = append(assembledOrder, model.Order{
				OrderID:   orderID,
				Count:     paletteCount,
				ProductID: productID,
				PaletteID: paletteID,
			})

			count -= paletteCount
		}
	}

	err2 := r.UpdatePalette(paletteID, productID, paletteCount, count)
	if err2 != nil {
		fmt.Println(5)
		return err2
	}
	assembledOrder = append(assembledOrder, model.Order{
		OrderID:   orderID,
		Count:     count,
		ProductID: productID,
		PaletteID: paletteID,
	})

	err2 = r.InsertOrder(orderID)
	if err2 != nil {
		fmt.Println(6)
		return err2
	}

	err = r.UpdateOrdersProducts(assembledOrder)
	if err != nil {
		fmt.Println(7)
		return err
	}

	return nil
}

// CheckPaletteCount возвращает paletteCount, paletteID и bool которые показывают хватает на нем товара или нет.
func (r *Repository2) CheckPaletteCount(status bool, productID, orderCount int) (int, int, bool, error) {
	var count int
	var paletteID int

	sql, args, err := config.Sq.
		Select("pp.product_count", "pp.palette_id").
		From("palettes_products pp").
		Where("pp.product_id = $1", productID).
		Where("pp.palette_status = $2", status).
		Where("pp.product_count > 0").
		GroupBy("pp.palette_id", "pp.product_count").
		ToSql()
	if err != nil {
		return 0, 0, false, err
	}

	rows, err := config.Pool.Query(ctx, sql, args...)
	if err != nil {
		return 0, 0, false, err
	}
	for rows.Next() {
		err = rows.Scan(&count, &paletteID)
		if err != nil {
			return 0, 0, false, err
		}
	}
	defer rows.Close()

	if count < orderCount {
		return count, paletteID, false, nil
	}

	return count, paletteID, true, nil
}

// UpdatePalette тут происходит обнавление суммы товаров на палетах.
func (r *Repository2) UpdatePalette(paletteID, productID, paletteCount, sumReduce int) error {

	paletteCount -= sumReduce

	sql, args, err := config.Sq.
		Update("palettes_products").
		Set("product_count", paletteCount).
		Where("product_id = $2", productID).
		Where("palette_id = $3", paletteID).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := config.Pool.Query(ctx, sql, args...)

	for rows.Next() {
		err = rows.Scan()
		if err != nil {
			return err
		}
	}
	defer rows.Close()

	return nil
}

func (r *Repository2) InsertOrdersProducts(order model.Order) error {
	id := rand.New(rand.NewSource(time.Now().Unix())).Int31()

	sql, args, err := config.Sq.
		Insert("orders_products").
		Columns("id", "order_id", "product_id", "palette_id", "product_count").
		Values(id, order.OrderID, order.ProductID, order.PaletteID, order.Count).
		ToSql()
	if err != nil {
		return err
	}
	_, err = config.Pool.Exec(ctx, sql, args...)

	return nil
}

func (r *Repository2) CheckSum(productID, count int) (bool, error) {
	var dbCount int

	sql, args, err := config.Sq.
		Select("SUM(pp.product_count)").
		From("palettes_products pp").
		Where("pp.product_id = $1", productID).
		ToSql()
	if err != nil {
		return false, err
	}

	row := config.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&dbCount)
	if err != nil {
		return false, err
	}

	if dbCount < count {
		return false, nil
	}

	return true, nil
}

func (r *Repository2) InsertOrder(orderID int) error {
	err := r.CheckOrder(orderID)
	if err == nil {
		return nil
	}

	sql, args, err := config.Sq.
		Insert("orders").
		Columns("id").
		Values(orderID).
		ToSql()
	if err != nil {
		return err
	}

	_, err = config.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository2) UpdateOrdersProducts(order []model.Order) error {

	for _, value := range order {
		orderCount, err := r.CheckOrdersProducts(value.OrderID, value.ProductID, value.PaletteID)

		if err != nil {
			err2 := r.InsertOrdersProducts(value)
			if err2 != nil {
				return err2
			}
		} else {
			orderCount += value.Count

			sql, args, err2 := config.Sq.
				Update("orders_products").
				Set("product_count", orderCount).
				Where("product_id = $2", value.ProductID).
				Where("order_id = $3", value.OrderID).
				Where("palette_id = $4", value.PaletteID).
				ToSql()
			if err2 != nil {
				return err2
			}
			fmt.Println(sql, args)
			_, err2 = config.Pool.Exec(ctx, sql, args...)
			if err2 != nil {
				return err2
			}
		}
	}

	fmt.Println("товар ID:", order[0].OrderID, " успешно занесен в базу.")
	return nil
}

func (r *Repository2) CheckOrder(orderID int) error {
	var id int
	sql, args, err := config.Sq.
		Select("o.id").
		From("orders o").
		Where("o.id = $1", orderID).
		ToSql()

	if err != nil {
		return err
	}

	rows := config.Pool.QueryRow(ctx, sql, args...)
	err = rows.Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository2) CheckOrdersProducts(orderID, productID, paletteID int) (int, error) {
	var count int
	sql, args, err := config.Sq.
		Select("op.product_count").
		From("orders_products op").
		Where("op.order_id = $1", orderID).
		Where("op.product_id = $2", productID).
		Where("op.palette_id = $3", paletteID).
		ToSql()
	if err != nil {
		return 0, err
	}

	row := config.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
