# Создание и Выдача сборщику заказа.
____
**Необходимые пакеты**
 + Gorm
 + Logrus

 **МОЖНО ИСПОЛЬЗОВАТЬ КОМАНДУ** : тогда мы докачаем зависимости которые описаны в go.mod

```go
go get ./..
```
____
# Подключение к DB и настройка миграций.

**Подключение к DB**: 
```go

var DB *gorm.DB

func InitPgSQL() error {
	var db *gorm.DB

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		Env.DbUser,
		Env.DbPassword,
		Env.DbHost,
		Env.DbPort,
		Env.DbName,
	)

	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		return err
	}

	DB = db

	return nil
}
```

**Миграция**:

пример:
```go
m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
	    {
			ID: userID.String(),
			// передаем структуру на основании которой будет создана таблица
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&domain.User{}).Error
				if err != nil {
					return err
				}
				return nil
			},
			// это метод отмены миграции ни разу не использовал
			Rollback: func(tx *gorm.DB) error {
				err := tx.DropTable("users").Error
				if err != nil {
					return err
				}
				return nil
			},
		}
        })
```
_____
# Функционал

**Создание заказа**
>За это отвечает функция репозитория CheckAndUpdateOPalette, ей необходимо передать id продукта, колличество и id заказа.

*Код :*

```go
func (r *Repo) CheckAndUpdateOPalette(productID, count, orderID int) error {
	var pp domain.PalettesProducts
	var countProducts []domain.PalettesProducts
	var paMass []model.Order
	sumProduct := 0

	err := config.DB.
		Table("palettes_products").
		Where("product_id =?", productID).
		Find(&countProducts).
		Error

	if err != nil {
		return err
	}

	for _, ppPart := range countProducts {
		sumProduct += ppPart.Count
	}

	if sumProduct < count {
		return errors.New("не хватает товара")
	}

	err = config.DB.
		Table("palettes_products").
		Where("product_id =? AND status = true", productID).
		Find(&pp).
		Error
	if err != nil {
		return err
	}

	if pp.Count < count {

		if pp.Count != 0 {
			err = config.DB.
				Table("palettes_products").
				Where("product_id =? AND status = true", productID).
				UpdateColumn("count", gorm.Expr("count - ?", pp.Count)).
				Error
			if err != nil {
				return err
			}
			count = count - pp.Count

			paMass = append(paMass, model.Order{
				Count:     count,
				ProductID: productID,
				PaletteID: pp.PaletteID,
				OrderID:   orderID,
			})
		}

		for count != 0 {
			err = config.DB.
				Table("palettes_products").
				Where("product_id =? AND count > 0", productID).
				Find(&pp).
				Error
			if err != nil {
				return err
			}

			if pp.Count >= count {
				err = config.DB.
					Table("palettes_products").
					Where("product_id =? AND palette_id =?", productID, pp.PaletteID).
					UpdateColumn("count", gorm.Expr("count -?", count)).
					Error
				if err != nil {
					return err
				}
				count -= count
				paMass = append(paMass, model.Order{
					Count:     count,
					ProductID: productID,
					PaletteID: pp.PaletteID,
					OrderID:   orderID,
				})

				err = r.PostOrdersProducts(paMass)
				if err != nil {
					return err
				}

				fmt.Println("заказ ", orderID, " занесен в базу.")
				return nil
			}

			err = config.DB.
				Table("palettes_products").
				Where("product_id =? AND palette_id =?", productID, pp.PaletteID).
				UpdateColumn("count", gorm.Expr("count -?", pp.Count)).
				Error
			if err != nil {
				return err
			}
			count -= pp.Count
			paMass = append(paMass, model.Order{
				Count:     count,
				ProductID: productID,
				PaletteID: pp.PaletteID,
				OrderID:   orderID,
			})
		}
	}

	err = config.DB.
		Table("palettes_products").
		Where("product_id =? AND status = true", productID).
		UpdateColumn("count", gorm.Expr("count -?", count)).
		Error
	if err != nil {
		return err
	}
	paMass = append(paMass, model.Order{
		Count:     count,
		ProductID: productID,
		PaletteID: pp.PaletteID,
		OrderID:   orderID,
	})

	err = r.PostOrdersProducts(paMass)
	if err != nil {
		return err
	}

	fmt.Println("заказ ", orderID, " занесен в базу.")
	return nil
}
```

**Создание списка заказов для сборщика.**
>Необходимо передать номера заказов в консоле при вызове main функции.

*Код:*
```go

func (r *Repo) TakeOrders(orders []int) error {
	var collectedOrder []model.Assemblings
	var collectedOrders [][]model.Assemblings

	tx := config.DB.Begin()

	for _, orderID := range orders {
		err := tx.Select("op.product_id as product_id, o.id as order_id, p.name as product_name, op.count as count, pa.name as palette").
			Table("orders o").
			Joins("JOIN orders_products op ON o.id = op.order_id AND(o.id =?)", orderID).
			Joins("JOIN products p ON op.product_id = p.id").
			Joins("JOIN palettes pa ON op.palette_id = pa.id").
			Order("p.name DESC").
			Scan(&collectedOrder).
			Error
		if err != nil {
			tx.Rollback()
			return err
		}
		collectedOrders = append(collectedOrders, collectedOrder)
	}
	tx.Commit()

	tools.Sort(collectedOrders)

	return nil
}
```
> В конце была вызвана функция сортировки заказов по стеллажам ее код приведен ниже.

```go
func Sort(orders [][]model.Assemblings) {
	response := make(map[string][]model.Assemblings)

	for _, g := range orders {
		for _, gg := range g {
			if _, ok := response[gg.Palette]; !ok {
				response[gg.Palette] = append(response[gg.Palette])
			}

			response[gg.Palette] = append(response[gg.Palette], gg)
		}
	}

	keys := make([]string, 0, len(response))
	for k := range response {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, palette := range keys {
		fmt.Println("___________", palette, "_____________")
		for _, v := range response[palette] {
			fmt.Println(fmt.Sprintf("id заказа: %v | id товара: %v | наименование: %v | колличество: %v | стеллаж: %v", v.OrderID, v.ProductID, v.ProductName, v.Count, v.Palette))
		}
	}

}
```
# **Схема DB**

![схема бв.png]()