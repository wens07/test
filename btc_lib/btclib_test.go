/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: btclib_test.go, Date: 2018-11-08
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
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/txscript"
	"github.com/stretchr/testify/assert"
)

func TestBitcoin_sign(t *testing.T) {

	//check sign message
	//sig, err := Bitcoin_signmessage("cT95wcAqftAUE5GQ9iB5AvaNDNoNrJXt1J82qWQsV3aSsFAoYvus", "wens")
	//assert.Nil(t, err)
	//
	//sigStr := base64.StdEncoding.EncodeToString(sig)
	//assert.Equal(t, "HyjNZLF48xQ0KMYkcAsdiN4PW4VS7Yqv8kI1ps5IVse7Ul9W6HbGrkT90F5NHf2gbv35thSpYmwhepN7u+Fzl4k=", sigStr)

	//check sign
	hex_str := "02000000019c18d6cb7d23019bc5eaca919d84908172c62985566418aeddd183b81920c02d0000000000ffffffff01a0860100000000001976a914f77cdf657c80fafd14fccb30c8f0a47e354ed78b88ac00000000"
	hash_byte, err := hex.DecodeString(hex_str)
	assert.Nil(t, err)
	sig, err := Bitcoin_sign("cT95wcAqftAUE5GQ9iB5AvaNDNoNrJXt1J82qWQsV3aSsFAoYvus", hash_byte)
	assert.Nil(t, err)

	//signature, err := ecdsa_lib.ParseSignature(sig, ecdsa_lib.S256())
	//assert.Nil(t, err)
	//derSig := signature.Serialize()

	assert.Equal(t, "30450221008d9c74455580ac930e036358d36d9a2a366de31b3d7af49857a6a0de348ea1e502201582f5e5d08692e893fe2d1fbcec0901151cda3a3daca2f7788cd4f4ec698109", hex.EncodeToString(sig))
}

func TestBitcoin_TrxFromHexstring(t *testing.T) {

	hex_str := "02000000019c18d6cb7d23019bc5eaca919d84908172c62985566418aeddd183b81920c02d0000000000ffffffff01a0860100000000001976a914f77cdf657c80fafd14fccb30c8f0a47e354ed78b88ac00000000"
	//hex_str := "0200000000010137cb900af23010a57371e1433220ae27c50bac1c9ec8fe8b7e35c5a08c28cea6000000001716001486951a115a0350cefdc4857475ee33289f821c04feffffff02005a6202000000001976a9140add702a594c58a5ce8f0ce72232758f27b2c73c88ac2ac1a1b50000000017a914b1d5186a610b1af9841dcc8a9260bcb9ab254c30870247304402202d48dcc3edf4e568e8a0c366cc7da8a99febb01b70837c005a778f8bbcf0051d02203c81b5133061939ead48a382ec22290cb6edb1624aa4f8014e4a28f80e17a58f01210361a7d4459c653a5e8a06788bbb7f52ecb1e19c39193532465b1ec5b8c474602edb610800"
	msgTrx, err := Bitcoin_TrxFromHexstring(hex_str)
	assert.Nil(t, err)

	trxHash := msgTrx.TxHash()
	for i, j := 0, len(trxHash)-1; i < len(trxHash)/2; i, j = i+1, j-1 {
		trxHash[i], trxHash[j] = trxHash[j], trxHash[i]
	}
	assert.Equal(t, int32(2), msgTrx.Version)
	assert.Equal(t, "83a1581f853626b23db311e4ff1c6c2cd0ff2b065c96b107285e5a9a12a36183", hex.EncodeToString(trxHash[:]))

	var out bytes.Buffer
	err = msgTrx.Serialize(&out)
	assert.Nil(t, err)
	assert.Equal(t, "02000000019c18d6cb7d23019bc5eaca919d84908172c62985566418aeddd183b81920c02d0000000000ffffffff01a0860100000000001976a914f77cdf657c80fafd14fccb30c8f0a47e354ed78b88ac00000000",
		hex.EncodeToString(out.Bytes()))

}

func TestBitcoin_CalcSignatureHash(t *testing.T) {

	raw_trx := "02000000019c18d6cb7d23019bc5eaca919d84908172c62985566418aeddd183b81920c02d0000000000ffffffff01a0860100000000001976a914f77cdf657c80fafd14fccb30c8f0a47e354ed78b88ac00000000"
	msgTrx, err := Bitcoin_TrxFromHexstring(raw_trx)
	assert.Nil(t, err)

	prevout_script := "76a9142db6c06a1c4fd5548fdd9be22144ca346c92149588ac"
	script_byte, err := hex.DecodeString(prevout_script)
	assert.Nil(t, err)
	sigHash, err := Bitcoin_CalcSignatureHash(script_byte, txscript.SigHashAll, msgTrx, 0)
	assert.Nil(t, err)

	sig, err := Bitcoin_sign("cT95wcAqftAUE5GQ9iB5AvaNDNoNrJXt1J82qWQsV3aSsFAoYvus", sigHash)
	assert.Nil(t, err)

	assert.Equal(t, "30450221008d9c74455580ac930e036358d36d9a2a366de31b3d7af49857a6a0de348ea1e502201582f5e5d08692e893fe2d1fbcec0901151cda3a3daca2f7788cd4f4ec698109", hex.EncodeToString(sig))
}
