/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: base58.go, Date: 2018-09-21
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

/*
Package base58 provides an API for working with modified base58 and Base58Check
encodings.

Modified Base58 Encoding

Standard base58 encoding is similar to standard base64 encoding except, as the
name implies, it uses a 58 character alphabet which results in an alphanumeric
string and allows some characters which are problematic for humans to be
excluded.  Due to this, there can be various base58 alphabets.

The modified base58 alphabet used by Bitcoin, and hence this package, omits the
0, O, I, and l characters that look the same in many fonts and are therefore
hard to humans to distinguish.

Base58Check Encoding Scheme

The Base58Check encoding scheme is primarily used for Bitcoin addresses at the
time of this writing, however it can be used to generically encode arbitrary
byte arrays into human-readable strings along with a version byte that can be
used to differentiate the same payload.  For Bitcoin addresses, the extra
version is used to differentiate the network of otherwise identical public keys
which helps prevent using an address intended for one network on another.
*/

package base58

import (
	"math/big"
)

var bigRadix = big.NewInt(58)
var bigZero = big.NewInt(0)

// Decode decodes a modified base58 string to a byte slice.
func Decode(b string) []byte {
	answer := big.NewInt(0)
	j := big.NewInt(1)

	scratch := new(big.Int)
	for i := len(b) - 1; i >= 0; i-- {
		tmp := b58[b[i]]
		if tmp == 255 {
			return []byte("")
		}
		scratch.SetInt64(int64(tmp))
		scratch.Mul(j, scratch)
		answer.Add(answer, scratch)
		j.Mul(j, bigRadix)
	}

	tmpval := answer.Bytes()

	var numZeros int
	for numZeros = 0; numZeros < len(b); numZeros++ {
		if b[numZeros] != alphabetIdx0 {
			break
		}
	}
	flen := numZeros + len(tmpval)
	val := make([]byte, flen)
	copy(val[numZeros:], tmpval)

	return val
}

// Encode encodes a byte slice to a modified base58 string.
func Encode(b []byte) string {
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0, len(b)*136/100)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		answer = append(answer, alphabet[mod.Int64()])
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, alphabetIdx0)
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
}
