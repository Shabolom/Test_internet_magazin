package tools

import (
	"Arkadiy_2Service/iternal/model"
	"fmt"
)

func Sort(orders [][]model.Sborka) {
	ans := make(map[string][]model.Sborka)
	response := make(map[string][]model.Sborka)

	for _, g := range orders {
		for _, gg := range g {
			if _, ok := ans[gg.Palette]; !ok {
				ans[gg.Palette] = append(ans[gg.Palette])
			}

			ans[gg.Palette] = append(ans[gg.Palette], gg)
		}
	}

	for _, char := range "ABCDEFGH" {
		if v, ok := ans[string(char)]; ok {
			for _, va := range v {
				response[string(char)] = append(response[string(char)], va)
			}
		}
	}

	for palette, data := range response {
		fmt.Println("___________", palette, "_____________")
		for _, v := range data {
			fmt.Println(fmt.Sprintf("id заказа %v | id товара %v | наименование %v | колличество %v | стеллаж %v", v.OrderID, v.ProductID, v.ProductName, v.Count, v.Palette))
		}
	}

}
