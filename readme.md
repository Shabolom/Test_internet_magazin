# Создание и Выдача сборщику заказа.
____
**Необходимые пакеты**
 + Gorm
 + Logrus

 **МОЖНО ИСПОЛЬЗОВАТЬ КОМАНДУ** : тогда мы докачаем зависимости которые описаны в go.mod

```go
go get ./..
```
___
# Для запуска миграции с помощью migrate go.

**Миграция:**
```go
migrate -source file://migrations -database postgresql://postgres:1234@localhost:5432/Uchoba?sslmode=disable 
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

---
# Заполнение BD при помощи squirrel.go и pgxpool.go

**Подключение к базе данных при помощи pgxpool**
```go
// подключение к бд
pool, err := pgxpool.Connect(context.Background(), connectionString)
if err != nil {
return err
}
```
>connectionString представляет из себя строку подключения к базе данных.

**Передаем настройки плейсхолдера для squirrel**
>Тут мы определяем какой знак будет воспринимать генератор sql запороса для замены передоваемыми данными.

```go
// тут представлен плейсхолдер $1,$2...
sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
```

***Генирация sql запроса при помощи squirrel***

Сам запрос строится согласно правилам sql.
>Пример Select запроса с Join.

```go
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
```
При выполнении функции мы получаем строку sql запроса которую нужно будет передать в обработчик sql запроса pgxpool, но тк у нас объявленны плейсхолдеры передаваемые значения будут "$1" в связи с чем обработчику необходимо передать так-же и аргументы.
>пример sql строки ответа функции
```go
SELECT op.product_id, o.id, p.name, op.product_count, pa.name FROM orders o JOIN orders_products op ON op.order_id = o.id JOIN palettes pa ON op.palette_id = pa.id JOIN products p ON op.product_id = p.id WHERE o.id = $1 GROUP BY
 op.product_id, o.id, p.name, pa.name, op.product_count
```
Пример передачи строки и аргуметов обработчику.
```go
// Pool.Query используется когда ожидается получение нескольких чтрок ответа для дальнейшей их обработки
rows, err := config.Pool.Query(ctx, sql, args...)
		if err != nil {
			return err
		}
        // обработка по строчно ответов
		for rows.Next() {
			err = rows.Scan(&orderPart.ProductID, &orderPart.OrderID, &orderPart.ProductName, &orderPart.Count, &orderPart.Palette)
			if err != nil {
				return err
			}
			order = append(order, orderPart)
		}
		defer rows.Close()
	}
```

----
# Основные функции
 
***При работе с squirrel***
+ Создание заказа
> За создания заказа отвечает функция CheckAndUpdateOPalette она проверяет создан ли был ранее и обновляет его или создает. 
```go
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
```
+ Получение заказов
>Получение заказов происходить при помощи AssemblingOrder 
```go
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
			order = append(order, orderPart)
		}
		defer rows.Close()
	}

	tools.Sort(order)
	defer dbQuerySpan.End()
	return nil
}
```
Пример ответа консоли.

![схема бв.png]()

----
# **Схема DB Gorm**

![схема бв.png]()