package tools

import (
	"Arkadiy_2Service/iternal/model"
	"fmt"
	"sort"
)

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
