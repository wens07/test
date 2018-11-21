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

	runtime.GOMAXPROCS(4)

	var startblocknum uint64 = 4058000
	var endblocknum uint64 = 4080000
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
