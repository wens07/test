/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: ecdsa_test.go, Date: 2018-11-08
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package ecdsa_lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWif(t *testing.T) {

	priv, _ := Wif2Privkey("cT95wcAqftAUE5GQ9iB5AvaNDNoNrJXt1J82qWQsV3aSsFAoYvus", S256())

	wif := Wif(priv, 239, true)

	assert.Equal(t, "cT95wcAqftAUE5GQ9iB5AvaNDNoNrJXt1J82qWQsV3aSsFAoYvua", wif)
}
