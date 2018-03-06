package main

import (
	"strconv"
	"time"
	"kuaifa.com/kuaifa/work-together/cmd/data_sync/src"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"fmt"
)

const (
	GAME_PATH    = "/forgoapi/game/list"
	CHANNEL_PATH = "/forgoapi/channel/list"
	ORDER_PATH   = "/forgoapi/order/daydata"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "command", "[arguments]")
		fmt.Println("Avaliable Commands:")
		fmt.Println("\tgame     game_id          update the game data")
		fmt.Println("\tchannel  channel_id       update the channel data")
		fmt.Println("\torder    date             update the order data")
		os.Exit(0)
	}
	switch os.Args[1] {
	case "game":
		game_params := map[string]string{
			"game_id":   os.Args[2],
			"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		}
		src.UpdateGame(GAME_PATH, game_params)
	case "channel":
		channel_params := map[string]string{
			"channel_id": os.Args[2],
			"timestamp":  strconv.FormatInt(time.Now().Unix(), 10),
		}
		src.UpdateChannel(CHANNEL_PATH, channel_params)
	case "order":
		order_params := map[string]string{
			"date":      os.Args[2],
			"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		}
		src.UpdateOrder(ORDER_PATH, order_params)
	default:
		fmt.Println("ERROR: wrong commond")
	}

	//game_params := map[string]string{
	//	"game_id": "0",
	//	"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	//}
	//game_addr := "/forgoapi/game/list"
	//src.UpdateGame(game_addr, game_params)

	//channel_params := map[string]string{
	//	"channel_id": "100",
	//	"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	//}
	//channel_addr := "/forgoapi/channel/list"
	//src.UpdateChannel(channel_addr, channel_params)

	//order_params := map[string]string{
	//	"date": "2016-12-31",
	//	"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	//}
	//order_addr := "/forgoapi/order/daydata"
	//data := src.GetData(order_addr, order_params)
	//fmt.Println(data)
	//parms := map[string]string{
	//	"start": "2016-12-01",
	//	"end": "2016-12-02",
	//	"game_id": "732",
	//	"cp": "",
	//	"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	//}
	//data := src.GetData("/forgoapi/order/export", parms)
	//fmt.Println(data)
	//fmt.Println(src.GetData(addr))

	//src.DBQuery("select id,game_id,amount from `order`")
	//query := "REPLACE INTO `order`(id, game_id,amount) VALUES(1, 1,111111.1),(2, 2,222222.2)"
	//src.DBUpdate(query)
}
