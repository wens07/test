/**
* Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
*
* Copyright © 2015--2018 . All rights reserved.
*
* File: subset.go, Date: 2018-07-12
*
*
* This library is free software under the terms of the GNU General Public License
* as published by the Free Software Foundation; either version 3 of the License,
* or (at your option) any later version.
*
 */

package algorithm

func Subset(src []int) [][]int {
	var res [][]int

	arr_len := len(src)

	for i := 0; i < (1 << uint8(arr_len)); i++ {

		var tmp []int

		for j := 0; j < arr_len; j++ {

			if i&(1<<uint8(j)) > 0 {
				tmp = append(tmp, src[j])
			}

		}

		res = append(res, tmp)

	}

	return res

}
