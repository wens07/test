/**
  * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
  *  
  * Copyright © 2015--2018 . All rights reserved.
  *
  * File: address_test 
  * Date: 2018-09-04
  *
  */

package bts

import (
	"testing"
	"fmt"
)

func TestKey(t *testing.T) {

	seed := MnemonicToSeed("venture lazy digital aware plug hire acquire abuse chunk know gloom snow much employ glow rich exclude allow", "123")
	wif, _ := ExportWif(seed,  0, 0)
	fmt.Println(wif)
	//privkey, _ := ImportWif("5JWbFZMcweitDc11YyZenHdiCmLpRWv8mZ8kDeehsqvEprZY6jv")

	privkey, _ := ImportWif(wif)
	pubKeystr, _ := GetPubKeyStr(privkey)

	fmt.Println(pubKeystr)

}
