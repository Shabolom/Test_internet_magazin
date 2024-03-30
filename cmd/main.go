package main

import (
	"Arkadiy_2Service/config"
	"Arkadiy_2Service/iternal/repository"
	"Arkadiy_2Service/tools"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type order struct {
	productID int
	count     int
	orderID   int
}

var argsOrders []int
var orders []order

func main() {
	config.CheckFlagEnv()
	tools.InitLogger()

	err := config.InitPgSQL()
	if err != nil {
		log.WithField("component", "initialization").Panic(err)
		defer config.Pool.Close()
	}
	//migrate.Migrate()

	args := os.Args
	for _, v := range args {
		orderID, _ := strconv.Atoi(v)
		argsOrders = append(argsOrders, orderID)
	}

	//err = repository.NewRepo2().CheckAndUpdateOPalette(3, 4, 1)
	//if err != nil {
	//	fmt.Println(err)
	//}

	err = repository.NewRepo2().AssemblingOrder([]int{1, 2, 3})
	if err != nil {
		println(err.Error())
	}

	//or := [][]int{{1, 2, 10}, {3, 1, 10}, {6, 1, 10}, {1, 3, 11}, {1, 3, 14}, {4, 4, 14}, {5, 1, 15}}
	//for _, ord := range or {
	//	orderEntity := &order{
	//		productID: ord[0],
	//		count:     ord[1],
	//		orderID:   ord[2],
	//	}
	//	orders = append(orders, *orderEntity)
	//}
	//
	//for _, value := range orders {
	//	err = repository.NewRepo().CheckAndUpdateOPalette(value.productID, value.count, value.orderID)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//
	//repository.NewRepo().TakeOrders(argsOrders)
}
