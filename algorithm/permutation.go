/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: permutation.go, Date: 2018-04-14
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package algorithm

import "fmt"

func print_result(src []int) {

	for i := 0; i < len(src); i++ {
		fmt.Print(src[i])
	}

	fmt.Println()
}

/**
 *  Heap's algorithm
 *  https://en.wikipedia.org/wiki/Heap%27s_algorithm
 */
func Permutation(src []int, size int) {

	if size == 1 {
		print_result(src)
		return
	}

	for i := 0; i < size; i++ {
		Permutation(src, size-1)

		if size%2 == 0 {
			src[i], src[size-1] = src[size-1], src[i]
		} else {
			src[0], src[size-1] = src[size-1], src[0]
		}

	}

}
