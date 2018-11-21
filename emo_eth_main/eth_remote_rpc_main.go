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

	//fmt.Println(eth_lib.WEB3_clientVersion())
	//res := eth_lib.ETH_check_transaction_successful("0xe7c41a57d6722bd01c26e69b70c72119f1aacbeec18cc070cc981e986ae34adc")
	res := eth_lib.ETH_getBalance("0x1c937d3b1d0a9ab6b8f10b2bf4733b8f5b82a4d3")
	fmt.Println(res)
}
