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

	balance := eth_lib.ETH_getBalance_by_block("0x80e29acb842498fe6591f020bd82766dce619d43", 4797776)
	fmt.Println(balance)
	fmt.Println(balance[:len(balance)-18])

}
