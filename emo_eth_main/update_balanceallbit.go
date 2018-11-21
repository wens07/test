/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2017 . All rights reserved.
 */

package emo_eth_main

import (
	"fmt"
	"time"

	"github.com/wens07/eth_lib"
)

func update_balance_allbit(num uint64) int {

	mysql_conn_str_tmp := "wq:123456@tcp(192.168.1.123:3306)/ethAddrBalance?charset=utf8"
	db := eth_lib.Connect_db(mysql_conn_str_tmp)
	defer db.Close()

	select_str := "select addr from air_drop_allbalancebit where `check` is null  limit " + fmt.Sprintf("%d", num)

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
		var block_num uint64 = 4730666
		balance := eth_lib.ETH_getBalance_by_block(addr, block_num)

		check := "1"
		update_str := "update air_drop_allbalancebit set `check` = " + `"` + check + `"` + " , balance = " + `"` + balance + `"` + " where addr = " + `"` + addr + `"`
		fmt.Println(update_str)

		_, err := db.Exec(update_str)
		eth_lib.CheckErr(err)

	}

	return index

}

func main() {

	//update_balance_allbit(3)

	for {

		res := update_balance_allbit(1000)

		time.Sleep(5 * time.Second)

		if res < 1000 {
			break
		}
	}

}
