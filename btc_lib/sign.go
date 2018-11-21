/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: sign.go, Date: 2018-11-07
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package btc_lib

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/wens07/ecdsa_lib"
	"github.com/wens07/util_lib/serialize"
)

const (
	MessageMagic string = "Bitcoin Signed Message:\n"
)

func Bitcoin_doublehash256(origin []byte) []byte {
	fisrt := sha256.Sum256(origin)
	second := sha256.Sum256(fisrt[:])
	return second[:]
}

func Bitcoin_sign(wif string, hash []byte) ([]byte, error) {

	privKey, _ := ecdsa_lib.Wif2Privkey(wif, ecdsa_lib.S256())

	sig, err := privKey.Sign(hash)
	if err != nil {
		return nil, fmt.Errorf("private key sign error")
	}

	return sig.Serialize(), nil

}

func Bitcoin_signmessage(wif string, message string) ([]byte, error) {

	privKey, _ := ecdsa_lib.Wif2Privkey(wif, ecdsa_lib.S256())

	var buf bytes.Buffer
	serialize.WriteVarString(&buf, 0, MessageMagic)
	serialize.WriteVarString(&buf, 0, message)

	hashByte := Bitcoin_doublehash256(buf.Bytes())

	return ecdsa_lib.SignCompact(ecdsa_lib.S256(), privKey, hashByte, true)

}
