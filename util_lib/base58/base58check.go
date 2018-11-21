/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: base58check.go, Date: 2018-09-21
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package base58

import (
	"crypto/sha256"
	"errors"
)

// ErrChecksum indicates that the checksum of a check-encoded string does not verify against
// the checksum.
var ErrChecksum = errors.New("checksum error")

// ErrInvalidFormat indicates that the check-encoded string has an invalid format.
var ErrInvalidFormat = errors.New("invalid format: version and/or checksum bytes missing")

// checksum: first four bytes of sha256^2
func checksum(input []byte) (cksum [4]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	copy(cksum[:], h2[:4])
	return
}

// CheckEncode prepends a version byte and appends a four byte checksum.
func CheckEncode(input []byte, version byte) string {
	b := make([]byte, 0, 1+len(input)+4)
	b = append(b, version)
	b = append(b, input[:]...)
	cksum := checksum(b)
	b = append(b, cksum[:]...)
	return Encode(b)
}

// CheckDecode decodes a string that was encoded with CheckEncode and verifies the checksum.
func CheckDecode(input string) (result []byte, version byte, err error) {
	decoded := Decode(input)
	if len(decoded) < 5 {
		return nil, 0, ErrInvalidFormat
	}
	version = decoded[0]
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if checksum(decoded[:len(decoded)-4]) != cksum {
		return nil, 0, ErrChecksum
	}
	payload := decoded[1 : len(decoded)-4]
	result = append(result, payload...)
	return
}
