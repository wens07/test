/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2017 . All rights reserved.
 */

package emo_eth_main

import (
	"fmt"

	"github.com/wens07/eth_lib"
)

func check_trx_whether_success(num uint64) int {

	mysql_conn_str_tmp := "wq:123456@tcp(192.168.1.123:3306)/ethAddrBalance?charset=utf8"
	db := eth_lib.Connect_db(mysql_conn_str_tmp)
	defer db.Close()

	select_str := "select addr, trx_id from air_drop_check where trx_id is not null and `check` is null limit " + fmt.Sprintf("%d", num)

	fmt.Println(select_str)

	rows, err := db.Query(select_str)
	eth_lib.CheckErr(err)

	var addr, trx_id string
	var index int = 0

	for rows.Next() {
		index++

		if err := rows.Scan(&addr, &trx_id); err != nil {
			eth_lib.CheckErr(err)
		}

		fmt.Println(addr, "\t", trx_id)
		successful := eth_lib.ETH_check_transaction_successful(trx_id)
		var check, update_str string

		if successful {

			check = "1"
			update_str = "update air_drop_check set `check` = " + `"` + check + `"` + " where addr = " + `"` + addr + `"`
			fmt.Println(update_str)

		} else {

			check = "0"
			update_str = "update air_drop_check set `check` = " + `"` + check + `"` + " where addr = " + `"` + addr + `"`
			fmt.Println(update_str)

		}

		_, err := db.Exec(update_str)
		eth_lib.CheckErr(err)

	}

	return index

}

func main() {

	//from, to := eth_lib.ETH_getTransactionByHash("0x1d03089418abf4dab380631be855da81053a489111240c31fadb171f3a77318c")
	//fmt.Println(from, to)
	//
	//res := eth_lib.ETH_check_transaction_successful("0x1d03089418abf4dab380631be855da81053a489111240c31fadb171f3a77318b")
	//fmt.Println(res)

	total := 0

	for {

		fmt.Println("total num: ", total)

		resnum := check_trx_whether_success(10000)
		total += resnum

		if resnum < 10000 {
			break
		}

	}

}
