package tools

import (
	"Arkadiy_2Service/iternal/model"
	"fmt"
	"sort"
)

func Sort(orders []model.Assemblings) {
	count := ""

	sort.Slice(orders,
		func(i, j int) bool {
			return orders[i].Palette < orders[j].Palette
		})

	for _, order := range orders {
		if order.Palette != count {
			fmt.Println("___________", order.Palette, "_____________")
			count = order.Palette
		}
		fmt.Println(fmt.Sprintf("id заказа: %v | id товара: %v | наименование: %v | колличество: %v | стеллаж: %v", order.OrderID, order.ProductID, order.ProductName, order.Count, order.Palette))
	}

}
