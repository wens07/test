/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2017 . All rights reserved.
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package eth_lib

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/bitly/go-simplejson"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/parnurzeal/gorequest"
)

const test_remote_api_url string = "http://192.168.1.178:8546"

//const test_remote_api_url string = "http://192.168.1.124:8588"

func Remote_rpc_call(args ...string) *simplejson.Json {

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

func Remote_eth_getBalance(addr string, blocknum uint64) string {

	hexblocknum := Int2hexstr(blocknum)
	fmt.Println(hexblocknum)

	callArgs := `["` + addr + `",` + `"` + hexblocknum + `"` + "]"

	res := Remote_rpc_call("eth_getBalance", callArgs).Get("result").MustString()

	//fmt.Println(res)

	balance, err := hexutil.DecodeBig(res)

	CheckErr(err)

	res = fmt.Sprintf("%s", balance)

	//fmt.Println(res)

	return res

}

//func Remote_eth_sendTransaction(sender_addr string, contract_addr string, data string) string {
//
//	gas := Int2hexstr(4000000)
//	gasPrice := Int2hexstr(4000000000)
//
//	callArgs := `[{` + `"from": ` + `"` + sender_addr + `", ` + `"to": ` + `"` + contract_addr + `", ` + `"gas": ` + `"` + gas + `", `
//	callArgs += `"gasPrice": ` + `"` + gasPrice + `", ` + `"data": ` + `"` + data + `"` + "}]"
//
//	res := Remote_rpc_call("eth_sendTransaction", callArgs)
//
//	fmt.Println(res)
//
//	trxid := res.Get("result").MustString()
//
//	fmt.Println(trxid)
//	return trxid
//
//}

func Multigoroutine_update_balance(start uint64, offset uint64, wg *sync.WaitGroup) {

	defer wg.Done()

	Multigoroutine_update_balance_helper(start, offset)
}

func Multigoroutine_update_balance_helper(start uint64, offset uint64) {

	mysql_conn_str_tmp := "wq:123456@tcp(192.168.1.123:3306)/ethAddrBalance?charset=utf8"
	db := Connect_db(mysql_conn_str_tmp)
	defer db.Close()

	fmt.Println(start, offset)
	select_str := "select addr from addr_balance limit " + fmt.Sprintf("%d", start) + `, ` + fmt.Sprintf("%d", offset)

	fmt.Println(select_str)

	rows, err := db.Query(select_str)
	CheckErr(err)

	var addr string

	for rows.Next() {

		if err := rows.Scan(&addr); err != nil {
			CheckErr(err)
		}

		fmt.Println(addr)

		Update_balance_and_insert_addr(addr, db)

		check := "1"
		update_str := "update addr_balance set `check` = " + `"` + check + `"` + " where addr = " + `"` + addr + `"`
		fmt.Println(update_str)

		_, err := db.Exec(update_str)
		CheckErr(err)

	}

}

func Update_balance_and_insert_addr(addr string, db *sql.DB) {

	//to do should get by blocknum
	var block_num uint64 = 4730666
	balance_int_string := Remote_eth_getBalance(addr, block_num)
	precision := 19

	//fmt.Println(balance_int_string)

	//update_str := "update test_addr_balance set balance = " + `"` + balance + `"` + " where addr = " + `"` + addr + `"`
	// the balance should larger than "1000000000000000000"
	if len(balance_int_string) >= precision {

		newbalance := balance_int_string[:len(balance_int_string)-(precision-1)]

		insert_str := "insert ignore into air_drop (addr, balance) VALUES (" + `"` + addr + `"` + `, "` + newbalance + `"` + ")"

		fmt.Println(insert_str)

		_, err := db.Exec(insert_str)
		CheckErr(err)

	}

}

func Multigoroutine_balance_stat(start uint64, offset uint64, db *sql.DB) uint64 {

	var sum uint64 = 0
	select_str := "select addr from addr_balance limit " + fmt.Sprintf("%d", start) + `, ` + fmt.Sprintf("%d", offset)

	fmt.Println(select_str)

	rows, err := db.Query(select_str)
	CheckErr(err)

	var addr string

	for rows.Next() {

		if err := rows.Scan(&addr); err != nil {
			CheckErr(err)
		}

		//fmt.Println(addr)

		balance := ETH_getBalance(addr)
		precision := 19

		//fmt.Println(balance)

		//update_str := "update test_addr_balance set balance = " + `"` + balance + `"` + " where addr = " + `"` + addr + `"`
		// the balance should larger than "1000000000000000000"
		if len(balance) >= precision {

			newbalance := balance[:len(balance)-(precision-1)]

			intbalance, err := strconv.Atoi(newbalance)
			CheckErr(err)

			fmt.Println("string :", newbalance)
			fmt.Println("int : ", intbalance)

			sum += uint64(intbalance)

		}

	}

	return sum

}

// RightPad2Len
func RightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func LeftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func Encode_param(hex_str string, left bool) string {

	enc := hex_str[2:]
	result := ""

	if left {
		result = LeftPad2Len(enc, "0", 64)
	} else {
		result = LeftPad2Len(enc, "0", 64)
	}

	//fmt.Println(result)

	return result

}

func Get_parameter_for_call_contract(dests []string, vals []string) string {

	var result string = ""

	func_byte := "aad41a41"

	// first param
	first_param := ""

	first_param_loc := Encode_param("0x40", true)
	//fmt.Println(first_param_loc)
	first_param_len := Encode_param(Int2hexstr(uint64(len(dests))), true)
	//fmt.Println(first_param_len)

	first_param += first_param_len
	for _, v := range dests {
		first_param += Encode_param(v, true)
	}
	//fmt.Println(first_param)

	// second param
	second_param := ""

	second_param_loc := Encode_param(Int2hexstr(uint64(96+len(dests)*32)), true)
	//fmt.Println(second_param_loc)
	second_param_len := Encode_param(Int2hexstr(uint64(len(vals))), true)
	//fmt.Println(second_param_len)

	second_param += second_param_len
	for _, v := range vals {

		num_big := new(big.Int)
		num_big, ok := num_big.SetString(v, 10)
		if !ok {
			panic("SetString: error")
		}

		balance_int_string := hexutil.EncodeBig(num_big)

		second_param += Encode_param(balance_int_string, true)
	}
	//fmt.Println(second_param)

	result += "0x" + func_byte + first_param_loc + second_param_loc + first_param + second_param

	fmt.Println(result)

	return result
}
