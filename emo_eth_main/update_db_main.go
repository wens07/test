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

func main() {

	mysql_conn_str_tmp := "wq:123456@tcp(192.168.1.123:3306)/ethAddrBalance?charset=utf8"
	db := eth_lib.Connect_db(mysql_conn_str_tmp)
	defer db.Close()

	select_str := "SELECT addr from air_drop where LENGTH(balance) > 18"

	fmt.Println(select_str)

	rows, err := db.Query(select_str)
	eth_lib.CheckErr(err)

	var addr string
	var index int = 0

	for rows.Next() {

		if err := rows.Scan(&addr); err != nil {
			eth_lib.CheckErr(err)
		}

		if err != nil {
			eth_lib.CheckErr(err)
		}

		//to do should get by blocknum
		var block_num uint64 = 4730666
		balance_int_string := eth_lib.Remote_eth_getBalance(addr, block_num)
		precision := 19

		//fmt.Println(balance_int_string)

		//update_str := "update test_addr_balance set balance = " + `"` + balance + `"` + " where addr = " + `"` + addr + `"`
		// the balance should larger than "1000000000000000000"
		if len(balance_int_string) >= precision {

			newbalance := balance_int_string[:len(balance_int_string)-(precision-1)]

			update_str := "update air_drop set balance = " + `"` + newbalance + `"` + " where addr = " + `"` + addr + `"`

			fmt.Println(update_str)

			_, err := db.Exec(update_str)
			eth_lib.CheckErr(err)

		}

		index++
	}

	fmt.Println(index)

}
