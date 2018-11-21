/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2017 . All rights reserved.
 */

package emo_eth_main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/wens07/eth_lib"
)

func main() {

	//eth_lib.ETH_getTrxHashsByBlockNumber(2257656)
	//hexstr := eth_lib.Int2hexstr(0)
	//fmt.Println(hexstr)
	//eth_lib.ETH_getTrxHashsByBlockNumber_debug(0)
	//
	//fmt.Println(len(data))
	//fmt.Println(reflect.TypeOf(data))
	//addr_from, addr_to := eth_lib.ETH_getTransactionByHash("0x3540984409284cc903d99f1ddb9f6e51534a55c292d6f01c05a43fbb3ffdce5c")
	//eth_lib.ETH_getBalance(addr_from)
	//eth_lib.ETH_getBalance("0x1c19839c88c0cdb4D0b10208833101D75454d7C1")

	//for _, v := range data {
	//	fmt.Println(v)
	//}

	//fmt.Println(eth_lib.ETH_blockNumber())
	//fmt.Println(eth_lib.WEB3_clientVersion())
	//fmt.Println(eth_lib.NET_peerCount())

	//db := eth_lib.Connect_db(mysql_conn_str)
	//eth_lib.Insert_address_byscanningblock_(792629, db)
	//addr := "0x63a9975ba31b0b9626b34300f7f627147df1f526"
	//insert_str := "insert into addr_balance (addr) VALUES (" + `"` + addr + `"` + ")"
	//
	//fmt.Println(insert_str)
	//
	//_, err := db.Exec(insert_str)
	//if err != nil {
	//	fmt.Println("insert error!")
	//}
	//
	//rows, err := db.Query("select * from addr_balance")
	//eth_lib.CheckErr(err)
	//
	//data := eth_lib.Get_result(rows)
	//
	//for k, v := range data {
	//	fmt.Println(k, v)
	//}

	runtime.GOMAXPROCS(4)

	var startblocknum uint64 = 4080000
	var endblocknum uint64 = 4102000
	var step uint64 = 1000

	start := time.Now()

	for startblocknum < endblocknum {

		var wg sync.WaitGroup
		wg.Add(4)

		go eth_lib.Multigoroutine_scan_block(startblocknum, startblocknum+step, &wg)
		go eth_lib.Multigoroutine_scan_block(startblocknum+step+1, startblocknum+2*step, &wg)
		go eth_lib.Multigoroutine_scan_block(startblocknum+2*step+1, startblocknum+3*step, &wg)
		go eth_lib.Multigoroutine_scan_block(startblocknum+3*step+1, startblocknum+4*step, &wg)

		wg.Wait()

		startblocknum = startblocknum + 4*step + 1
	}

	duration := time.Since(start)

	fmt.Println("executation time: ", duration)

}
