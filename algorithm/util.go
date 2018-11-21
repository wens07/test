/**
* Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
*
* Copyright © 2015--2018 . All rights reserved.
*
* File: util.go, Date: 2018-04-12
*
*
* This library is free software under the terms of the GNU General Public License
* as published by the Free Software Foundation; either version 3 of the License,
* or (at your option) any later version.
*
 */

package algorithm

import (
	"regexp"
	"strconv"
)

func Int2string(i int) string {

	return strconv.Itoa(i)
}

func Str_reverse(src string) string {

	temp := []rune(src)

	for i, j := 0, len(src)-1; i < j; i, j = i+1, j-1 {
		temp[i], temp[j] = temp[j], temp[i]
	}

	return string(temp)
}

func Is_alphanumeric(src string) bool {

	reg := regexp.MustCompile("[[:alnum:]]")

	return reg.MatchString(src)

}
