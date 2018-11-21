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

func check_address_whether_update(num uint64) int {

	mysql_conn_str_tmp := "wq:123456@tcp(192.168.1.123:3306)/ethAddrBalance?charset=utf8"
	db := eth_lib.Connect_db(mysql_conn_str_tmp)
	defer db.Close()

	select_str := "select addr from addr_balance where `check` = '0' limit " + fmt.Sprintf("%d", num)

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

		//to do should get by blocknum
		var block_num uint64 = 4730666
		balance_int_string := eth_lib.Remote_eth_getBalance(addr, block_num)
		precision := 19

		//fmt.Println(balance_int_string)

		//update_str := "update test_addr_balance set balance = " + `"` + balance + `"` + " where addr = " + `"` + addr + `"`
		// the balance should larger than "1000000000000000000"
		if len(balance_int_string) >= precision {

			newbalance := balance_int_string[:len(balance_int_string)-(precision-1)]

			insert_str := "insert ignore into air_drop (addr, balance) VALUES (" + `"` + addr + `"` + `, "` + newbalance + `"` + ")"

			fmt.Println(insert_str)

			_, err := db.Exec(insert_str)
			eth_lib.CheckErr(err)

		}

		check := "1"
		update_str := "update addr_balance set `check` = " + `"` + check + `"` + " where addr = " + `"` + addr + `"`
		fmt.Println(update_str)

		_, err := db.Exec(update_str)
		eth_lib.CheckErr(err)

	}

	return index

}

func main() {

	for {
		numchecked := check_address_whether_update(1000)

		if numchecked == 0 {
			break
		}
	}

}
