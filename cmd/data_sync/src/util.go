package src

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"fmt"
	"io/ioutil"
	"sort"
	"database/sql"
	"os"
	"encoding/json"
	"strconv"
)

const (
	SECRET_KEY = "5f5a740a75ed31c67c5d16eded29d30d"
	HOST = "http://www.kuaifazs.com"
	DB_USER = "kftest"
	DB_PASS = "123456"
	DB_HOST = "10.8.230.17"
	DB_PORT = "3308"
	DB_NAME = "work_together"
)

func dbConnect() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME))
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}

// @Param query  string "e.g. SELECT MAX(update_time) FROM table_name;"
func DBGetLastUpdateTime(query string) int64 {
	db := dbConnect()
	defer db.Close()
	rows, err := db.Query(query)
	checkError(err)
	var update_time int64
	for rows.Next() {
		err := rows.Scan(&update_time)
		checkError(err)
	}
	return update_time
}

func DBQuery(query string) {
	db := dbConnect()
	defer db.Close()
	rows, err := db.Query(query)
	checkError(err)
	var id int64
	var game_id int64
	var amount float64
	for rows.Next() {
		err := rows.Scan(&id, &game_id, &amount)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(id, game_id, amount)
	}
}

func DBDelete(query string) {
	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
}

// @Param query string "e.g."REPLACE INTO `order`(id, game_id,amount) VALUES(1, 1,111111.1),(2, 2,222222.2)""
func DBUpdate(query string) {
	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	res, err := stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	num, _ := res.RowsAffected()
	fmt.Println("Rows Affected:", num)
	defer db.Close()
}

func ParseURL(path string, parms map[string]string) string {
	keys := make([]string, len(parms))
	i := 0
	for k, _ := range parms {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	ret := ""
	for _, key := range keys {
		ret = ret + key + "=" + parms[key] + "&"
	}
	sign := GetSign(ret[:len(ret) - 1])
	return HOST + path + "?" + ret + "sign=" + sign
}

func GetSign(parms string) string {
	m := md5.New()
	m.Write([]byte(parms))
	tmp_md5 := hex.EncodeToString(m.Sum(nil))
	mm := md5.New()
	mm.Write([]byte(tmp_md5 + SECRET_KEY))
	return hex.EncodeToString(mm.Sum(nil))
}

func GetData(path string, params map[string]string) string {
	addr := ParseURL(path, params)
	fmt.Println(addr)
	rsp, err := http.Get(addr)
	checkError(err)
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	checkError(err)
	//fmt.Println("data type:", reflect.TypeOf(body))
	return string(body)
}

func UpdateGame(path string, params map[string]string) {
	data := GetData(path, params)
	//fmt.Println(data)
	var game GameResponse
	json.Unmarshal([]byte(data), &game)
	fmt.Println("game num:", len(game.Data))
	if game.Status == "0" && len(game.Data) > 0 {
		query_sql := "REPLACE INTO `game_all` (game_id, game_name, icon, show_icon) VALUES"
		for _, row := range game.Data {
			query_sql = query_sql + fmt.Sprintf("(%d,'%s','%s','%s'),", row.Game_id, row.Name, row.Icon, row.Show_icon)
		}
		query_sql = query_sql[:len(query_sql) - 1]
		DBUpdate(query_sql)
	} else {
		fmt.Println("Info: no data for update.")
	}
}

func UpdateChannel(path string, params map[string]string) {
	data := GetData(path, params)
	//fmt.Println(data)
	var channel ChannelResponse
	json.Unmarshal([]byte(data), &channel)
	if channel.Status == "0" && len(channel.Data) > 0 {
		query_sql := "REPLACE INTO `channel` (channel_id, name, cp, platform) VALUES"
		for _, row := range channel.Data {
			query_sql = query_sql + fmt.Sprintf("(%d,'%s','%s','%s'),", row.Id, row.Channel_name, row.Cp, row.Platform)
		}
		query_sql = query_sql[:len(query_sql) - 1]
		DBUpdate(query_sql)
	} else {
		fmt.Println("Info: no data for update.")
	}
}

func UpdateOrder(path string, params map[string]string) {
	data := GetData(path, params)
	var order OrderResponse
	json.Unmarshal([]byte(data), &order)
	if order.Status == "0" && len(order.Data) > 0 {
		DBDelete("DELETE FROM `order` WHERE status=4 and `date`='"+params["date"]+"'")
		query_sql := "REPLACE INTO `order` (game_id, cp, date, amount, status, update_time) VALUES"
		for _, row := range order.Data {
			amount, _ := strconv.ParseFloat(row.Amount, 64)
			query_sql = query_sql + fmt.Sprintf("(%d,'%s','%s',%e,%d,%d),", row.Game_id, row.Cp, row.Date, amount, row.Order_status, row.Update_time)
		}
		query_sql = query_sql[:len(query_sql) - 1]
		DBUpdate(query_sql)
	} else {
		fmt.Println("Info: no data for update.")
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}


