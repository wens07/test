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

func get_all_emo_of_biyinaddr(num int) int {

	biyin_mysql_conn_str := "wq:123456@tcp(192.168.1.123:3306)/ethAddrBalance?charset=utf8"
	db := eth_lib.Connect_db(biyin_mysql_conn_str)
	defer db.Close()

	select_str := "select addr from biyin_address_currentblock where balance is null  limit " + fmt.Sprintf("%d", num)

	fmt.Println(select_str)

	rows, err := db.Query(select_str)
	eth_lib.CheckErr(err)

	var addr string
	var index int = 0

	for rows.Next() {
		index++

		if err := rows.Scan(&addr); err != nil {
			eth_lib.CheckErr(err)
		}

		fmt.Println(addr)

		//air drop state block num
		//var block_num uint64 = 4730666
		balance := eth_lib.ETH_getBalance(addr)

		check := "1"
		update_str := "update biyin_address_currentblock set `check` = " + `"` + check + `"` + " , balance = " + `"` + balance + `"` + " where addr = " + `"` + addr + `"`
		fmt.Println(update_str)

		_, err := db.Exec(update_str)
		eth_lib.CheckErr(err)

	}

	return index
}

func update_biyinaddr_airdroped(num int) int {

	biyin_mysql_conn_str := "wq:123456@tcp(192.168.1.123:3306)/ethAddrBalance?charset=utf8"
	db := eth_lib.Connect_db(biyin_mysql_conn_str)
	defer db.Close()

	select_str := "select addr from biyin_address where `check` is null limit " + fmt.Sprintf("%d", num)
	//fmt.Println(select_str)

	rows, err := db.Query(select_str)
	eth_lib.CheckErr(err)

	var addr string
	var index int = 0

	for rows.Next() {
		index++

		if err := rows.Scan(&addr); err != nil {
			eth_lib.CheckErr(err)
		}

		fmt.Println(addr)

		select_str = "select `check` from air_drop_check where addr = " + `"` + addr + `"` + " and ( `check` = '1' or `check` = '-1' )"

		result_row, _ := db.Query(select_str)
		var check string
		var update_str string

		if result_row.Next() {

			if result_err := result_row.Scan(&check); result_err != nil {
				eth_lib.CheckErr(result_err)
			}

			if check == "1" {

				update_check := "2"
				update_str = "update biyin_address set `check` = " + `"` + update_check + `"` + " where addr = " + `"` + addr + `"`

			} else if check == "-1" {

				update_check := "-1"
				update_str = "update biyin_address set `check` = " + `"` + update_check + `"` + " where addr = " + `"` + addr + `"`

			}

		} else {

			update_check := "1"
			update_str = "update biyin_address set `check` = " + `"` + update_check + `"` + " where addr = " + `"` + addr + `"`

		}

		fmt.Println(check, update_str)

		_, err := db.Exec(update_str)
		eth_lib.CheckErr(err)

	}

	return index

}

func main() {

	total := 0

	for {

		fmt.Println("total num: ", total)
		res := get_all_emo_of_biyinaddr(1000)
		total += res

		if res < 1000 {
			break
		}

	}

	//
	//for {
	//
	//	fmt.Println("select num: ", total)
	//	res := update_biyinaddr_airdroped(30)
	//	total += res
	//
	//	if res < 30 {
	//		break
	//	}
	//
	//}

}
