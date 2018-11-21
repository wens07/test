/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: eth_rpc_test.go
 * Date: 2018-06-28
 *
 */

package main

import (
	"fmt"
	"reflect"

	"github.com/wens07/eth_lib"
)

func main() {

	result := eth_lib.ETH_getRawTransactionByHash("0x4074eea5b9a83f5c08b8cb2eb19ee5703ffe9878b32335e5828c510202deb6de")
	fmt.Println(result)
	fmt.Println(reflect.TypeOf(result))
}
