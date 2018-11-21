/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2017 . All rights reserved.
 */

package emo_eth_main

import (
	"fmt"
	"log"
	"os"

	"github.com/bitly/go-simplejson"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/parnurzeal/gorequest"
	"github.com/wens07/eth_lib"
)

const test_remote_api_url string = "http://192.168.1.178:8546"

//const test_remote_api_url string = "http://192.168.1.124:8588"

func airdrop_rpc_call(args ...string) *simplejson.Json {

	id := "1"
	method := args[0]
	params := "[]"
	if len(args) > 1 {
		params = args[1]
	}

	postBody := `{"jsonrpc":"2.0","method":"` + method + `","params":` + params + `,"id":` + id + `}`
	fmt.Println("postBody: " + postBody)

	_, body, errs := gorequest.New().Post(test_remote_api_url).
		Send(postBody).
		End()

	if errs != nil {
		fmt.Println(errs)
		panic(errs)
	}

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		log.Fatalln(err)
	}

	return js
}

func airdrop_eth_getBalance(addr string, blocknum uint64) string {

	hexblocknum := eth_lib.Int2hexstr(blocknum)
	fmt.Println(hexblocknum)

	callArgs := `["` + addr + `",` + `"` + hexblocknum + `"` + "]"

	res := airdrop_rpc_call("eth_getBalance", callArgs).Get("result").MustString()

	//fmt.Println(res)

	balance, err := hexutil.DecodeBig(res)

	eth_lib.CheckErr(err)

	res = fmt.Sprintf("%s", balance)

	//fmt.Println(res)

	return res

}

func airdrop_eth_sendTransaction(sender_addr string, contract_addr string, gas_param uint64, data string) string {

	gas := eth_lib.Int2hexstr(gas_param)
	gasPrice := eth_lib.Int2hexstr(21000000000)

	callArgs := `[{` + `"from": ` + `"` + sender_addr + `", ` + `"to": ` + `"` + contract_addr + `", ` + `"gas": ` + `"` + gas + `", `
	callArgs += `"gasPrice": ` + `"` + gasPrice + `", ` + `"data": ` + `"` + data + `"` + "}]"

	res := airdrop_rpc_call("eth_sendTransaction", callArgs)

	fmt.Println(res)

	trxid := res.Get("result").MustString()

	fmt.Println(trxid)
	return trxid

}

func airdrop_eth_getTransactionCount(addr string) uint64 {

	callArgs := `["` + addr + `", "latest" ]`

	res := airdrop_rpc_call("eth_getTransactionCount", callArgs).Get("result").MustString()

	fmt.Println(res)
	num := eth_lib.Hexstr2int(res)

	return num
}

func write_addr_balance_script(start uint64, step uint64) {

	mysql_conn_str_tmp := "wq:123456@tcp(192.168.1.123:3306)/test_wq?charset=utf8"
	db := eth_lib.Connect_db(mysql_conn_str_tmp)
	defer db.Close()

	fmt.Println(start, step)
	select_str := "select addr, balance from test_addr_balance limit " + fmt.Sprintf("%d", start) + `, ` + fmt.Sprintf("%d", step)

	fmt.Println(select_str)

	rows, err := db.Query(select_str)
	eth_lib.CheckErr(err)

	var addr string
	var balance string

	fh_addr, err_addr := os.OpenFile("C:\\Users\\wensw\\Desktop\\addr.txt", os.O_RDWR|os.O_APPEND, 0755)
	if err_addr != nil {
		eth_lib.CheckErr(err_addr)
	}
	fh_addr.WriteString("[")

	fh_balance, err_balance := os.OpenFile("C:\\Users\\wensw\\Desktop\\balance.txt", os.O_RDWR|os.O_APPEND, 0755)
	if err_balance != nil {
		eth_lib.CheckErr(err_balance)
	}
	fh_balance.WriteString("[")

	for rows.Next() {

		if err := rows.Scan(&addr, &balance); err != nil {
			eth_lib.CheckErr(err)
		}

		fmt.Println(addr)
		fmt.Println(balance)

		fh_addr.WriteString(`"` + addr + `"`)

		fh_addr.WriteString(", ")

		fh_balance.WriteString(`"` + balance + `"`)

		fh_balance.WriteString(", ")

	}

	fh_addr.WriteString("]\n")
	fh_balance.WriteString("]\n")

}

func airdrop_balance(offset uint64) int {

	mysql_conn_str_tmp := "wq:123456@tcp(192.168.1.123:3306)/ethAddrBalance?charset=utf8"
	db := eth_lib.Connect_db(mysql_conn_str_tmp)
	defer db.Close()

	//select_str := "select addr, balance from air_drop where trx_id is null ORDER BY CONVERT(balance, DECIMAL) DESC limit " + fmt.Sprintf("%d", offset)
	//select_str := "select addr, balance from air_drop where trx_id is null ORDER BY CONVERT(balance, DECIMAL) ASC limit " + fmt.Sprintf("%d", offset)
	trx_id_str := "0xbc3a00371cb8d351896a814cf0f0abf1ec3b3e173a165c99b74a3d1fd85ec4ef"
	select_str := "select addr, balance from one_transaction_resend where trx_id = " + `"` + trx_id_str + `"` + " limit " + fmt.Sprintf("%d", offset)

	fmt.Println(select_str)

	rows, err := db.Query(select_str)
	eth_lib.CheckErr(err)

	var addr string
	var balance string

	param_addr_arr := make([]string, offset)
	param_balance_arr := make([]string, offset)
	var index = 0

	for rows.Next() {

		if err := rows.Scan(&addr, &balance); err != nil {
			eth_lib.CheckErr(err)
		}

		param_addr_arr[index] = addr

		if err != nil {
			eth_lib.CheckErr(err)
		}

		balance = eth_lib.Remote_eth_getBalance(addr, 4730666)
		param_balance_arr[index] = balance

		//update_str := "update air_drop set balance = " + `"` + balance + `"` + " where addr = " + `"` + addr + `"`
		////fmt.Println(update_str)
		//
		//_, err := db.Exec(update_str)
		//eth_lib.CheckErr(err)

		index++
	}

	fmt.Println(param_balance_arr[:])
	fmt.Println(param_addr_arr[:])
	//addr_arr := fmt.Sprintf("%s", param_addr_arr[:])
	//balance_arr := fmt.Sprintf("%s", param_balance_arr[:])
	//
	//fh, err_addr := os.OpenFile("C:\\Users\\wensw\\Desktop\\final_air_drop.txt", os.O_RDWR|os.O_APPEND, 0755)
	//if err_addr != nil {
	//	eth_lib.CheckErr(err_addr)
	//}
	//
	//fh.WriteString(addr_arr)
	//fh.WriteString("\n")
	//fh.WriteString(balance_arr)
	//fh.WriteString("\n")
	//fh.WriteString(trxid + "\n")

	data := eth_lib.Get_parameter_for_call_contract(param_addr_arr[:], param_balance_arr[:])

	gas := eth_lib.ETH_estimateGas("0x1c19839c88c0cdb4d0b10208833101d75454d7c1", "0xd037a81b22e7f814bc6f87d50e5bd67d8c329fa2", data) + 100000
	fmt.Println(gas)
	if gas > 8000000 {
		panic("gas exceeds block gas limit!")
	}

	trxid := airdrop_eth_sendTransaction("0x1c19839c88c0cdb4d0b10208833101d75454d7c1", "0xd037a81b22e7f814bc6f87d50e5bd67d8c329fa2", gas, data)

	if trxid == "" {
		return 0
	} else {

		for _, v := range param_addr_arr {

			update_str := "update one_transaction_resend set trx_id = " + `"` + trxid + `"` + " where addr = " + `"` + v + `"`
			//fmt.Println(update_str)

			_, err := db.Exec(update_str)
			eth_lib.CheckErr(err)

		}

		return index

	}

}

func test_airdrop(offset uint64) int {

	mysql_conn_str_tmp := "wq:123456@tcp(192.168.1.123:3306)/ethAddrBalance?charset=utf8"
	db := eth_lib.Connect_db(mysql_conn_str_tmp)
	defer db.Close()

	//select_str := "select addr, balance from air_drop where trx_id is null ORDER BY CONVERT(balance, DECIMAL) DESC limit " + fmt.Sprintf("%d", offset)
	select_str := "select addr, balance from air_drop where trx_id is null ORDER BY CONVERT(balance, DECIMAL) ASC limit " + fmt.Sprintf("%d", offset)

	fmt.Println(select_str)

	rows, err := db.Query(select_str)
	eth_lib.CheckErr(err)

	var addr string
	var balance string

	param_addr_arr := make([]string, offset)
	param_balance_arr := make([]string, offset)
	var index = 0

	for rows.Next() {

		if err := rows.Scan(&addr, &balance); err != nil {
			eth_lib.CheckErr(err)
		}

		param_addr_arr[index] = addr

		if err != nil {
			eth_lib.CheckErr(err)
		}

		balance = eth_lib.Remote_eth_getBalance(addr, 4730666)
		param_balance_arr[index] = balance

		update_str := "update air_drop set balance = " + `"` + balance + `"` + " where addr = " + `"` + addr + `"`
		//fmt.Println(update_str)

		_, err := db.Exec(update_str)
		eth_lib.CheckErr(err)

		index++
	}

	fmt.Println(param_balance_arr[:])
	fmt.Println(param_addr_arr[:])
	//addr_arr := fmt.Sprintf("%s", param_addr_arr[:])
	//balance_arr := fmt.Sprintf("%s", param_balance_arr[:])
	//
	//fh, err_addr := os.OpenFile("C:\\Users\\wensw\\Desktop\\final_air_drop.txt", os.O_RDWR|os.O_APPEND, 0755)
	//if err_addr != nil {
	//	eth_lib.CheckErr(err_addr)
	//}
	//
	//fh.WriteString(addr_arr)
	//fh.WriteString("\n")
	//fh.WriteString(balance_arr)
	//fh.WriteString("\n")
	//fh.WriteString(trxid + "\n")

	data := eth_lib.Get_parameter_for_call_contract(param_addr_arr[:], param_balance_arr[:])

	gas := eth_lib.ETH_estimateGas("0x1e470f8a0f46e62e64a7911e7ec51000acca4b23", "0x52be9e526e18937eb9f842a7a4a7104b2d2543ec", data) + 1000000

	trxid := airdrop_eth_sendTransaction("0x1e470f8a0f46e62e64a7911e7ec51000acca4b23", "0x52be9e526e18937eb9f842a7a4a7104b2d2543ec", gas, data)
	fmt.Println(trxid)
	return index

	//
	//if trxid == "" {
	//	return 0
	//} else {
	//
	//	for _, v := range param_addr_arr {
	//
	//		update_str := "update air_drop set trx_id = " + `"` + trxid + `"` + " where addr = " + `"` + v + `"`
	//		//fmt.Println(update_str)
	//
	//		_, err := db.Exec(update_str)
	//		eth_lib.CheckErr(err)
	//
	//	}
	//
	//	return index
	//
	//}

}

func main() {
	//for i := 0; i < 10; i++ {
	//
	//	fmt.Println(airdrop_balance(100))
	//}
	//
	////fmt.Println(test_airdrop(3))
	fmt.Println(airdrop_balance(50))
	//fmt.Println(time.Now())

	//total := 0
	//
	//for {
	//
	//	time.Sleep(5 * time.Second)
	//	if total >= 758 {
	//		break
	//	}
	//
	//	pending_trx_num := eth_lib.ETH_pendingTransactions()
	//
	//	fmt.Println("current transaction num ", total)
	//
	//	if pending_trx_num < 1 {
	//
	//		res := airdrop_balance(100)
	//		fmt.Println(res)
	//		total++
	//
	//		if res < 100 {
	//			break
	//		}
	//	}
	//}

}
