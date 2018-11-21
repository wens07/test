/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: my_go.go
 * Date: 2018-04-11
 *
 */

package main

import (
	"fmt"

	"github.com/wens07/util_lib"
)

func main() {

	src := []byte{1, 2, 3, 4}
	out := util_lib.ReverseSlice(src)

	fmt.Printf("%v", out)

}
