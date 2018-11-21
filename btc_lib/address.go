/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: address.go, Date: 2018-09-19
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package btc_lib

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/wens07/ecdsa_lib"
	"github.com/wens07/util_lib/base58"
	"github.com/wens07/util_lib/ripemd160"
)

const (
	MAIN_VERSION_P2PKH = 0
	MAIN_VERSION_P2SH  = 5
	MAIN_VERSION_WIF   = 128

	TEST_VERSION_P2PKH = 111
	TEST_VERSION_P2SH  = 196
	TEST_VERSION_WIF   = 239
)

type AddressDetail struct {
	Ver    string `json:"version"`
	Rip160 string `json:"ripemd160"`

	Ck string `json:"checksum"`
}

func GetSecp256k1Params() *ecdsa_lib.KoblitzCurve {

	return ecdsa_lib.S256()

}

func DecodeAddress(addr string) ([]byte, error) {

	outstr := hex.EncodeToString(base58.Decode(addr))
	//ver, err := strconv.ParseInt(outstr[0:2], 16, 64)
	//if err != nil {
	//	return "", fmt.Errorf("in DecodeAddress function, get version error: %v", err)
	//}

	res := AddressDetail{
		outstr[0:2],
		outstr[2:42],
		outstr[42:],
	}

	jsonres, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("in DecodeAddress function, get json error: %v", err)
	}

	return jsonres, nil

}

func Address(version byte, pubkeybytes []byte) string {

	h := sha256.Sum256(pubkeybytes)
	rip160 := ripemd160.New()
	rip160.Write(h[:])
	res := rip160.Sum(nil)
	return base58.CheckEncode(res, version)
}
