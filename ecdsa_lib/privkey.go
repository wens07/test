/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: privkey.go, Date: 2018-09-21
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package ecdsa_lib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/wens07/util_lib/base58"
)

// PrivateKey wraps an ecdsa.PrivateKey as a convenience mainly for signing
// things with the the private key without having to directly import the ecdsa
// package.
type PrivateKey ecdsa.PrivateKey

// PrivKeyFromBytes returns a private and public key for `curve' based on the
// private key passed as an argument as a byte slice.
func PrivKeyFromBytes(curve elliptic.Curve, pk []byte) (*PrivateKey,
	*PublicKey) {
	x, y := curve.ScalarBaseMult(pk)

	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(pk),
	}

	return (*PrivateKey)(priv), (*PublicKey)(&priv.PublicKey)
}

// NewPrivateKey is a wrapper for ecdsa.GenerateKey that returns a PrivateKey
// instead of the normal ecdsa.PrivateKey.
func NewPrivateKey(curve elliptic.Curve) (*PrivateKey, error) {
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return (*PrivateKey)(key), nil
}

// PubKey returns the PublicKey corresponding to this private key.
func (p *PrivateKey) PubKey() *PublicKey {
	return (*PublicKey)(&p.PublicKey)
}

// ToECDSA returns the private key as a *ecdsa.PrivateKey.
func (p *PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	return (*ecdsa.PrivateKey)(p)
}

// Sign generates an ECDSA signature for the provided hash (which should be the result
// of hashing a larger message) using the private key. Produced signature
// is deterministic (same message and same key yield the same signature) and canonical
// in accordance with RFC6979 and BIP0062.
func (p *PrivateKey) Sign(hash []byte) (*Signature, error) {
	return signRFC6979(p, hash)
}

// PrivKeyBytesLen defines the length in bytes of a serialized private key.
const PrivKeyBytesLen = 32

// compressMagic is the magic byte used to identify a WIF encoding for
// an address created from a compressed serialized public key.
const compressMagic byte = 0x01

// Serialize returns the private key number d as a big-endian binary-encoded
// number, padded to a length of 32 bytes.
func (p *PrivateKey) Serialize() []byte {
	b := make([]byte, 0, PrivKeyBytesLen)
	return paddedAppend(PrivKeyBytesLen, b, p.ToECDSA().D.Bytes())
}

// Wif creates the Wallet Import Format string encoding of a WIF structure.
// See DecodeWIF for a detailed breakdown of the format and requirements of
// a valid WIF string.
func Wif(key *PrivateKey, version byte, compressed bool) string {
	// Precalculate size.  Maximum number of bytes before base58 encoding
	// is one byte for the network, 32 bytes of private key, possibly one
	// extra byte if the pubkey is to be compressed, and finally four
	// bytes of checksum.
	encodeLen := 1 + PrivKeyBytesLen + 4

	if compressed {
		encodeLen++
	}

	res := make([]byte, 0, encodeLen)

	res = append(res, version)
	res = paddedAppend(PrivKeyBytesLen, res, key.D.Bytes())
	if compressed {
		res = append(res, compressMagic)
	}

	//get checksum
	first := sha256.Sum256(res)
	second := sha256.Sum256(first[:])
	chk := second[:4]

	res = append(res, chk...)
	return base58.Encode(res)
}

//
// Wif string convert to ecdsa private key
//
func Wif2Privkey(wif string, curve elliptic.Curve) (*PrivateKey, *PublicKey) {

	wifOriginByte := base58.Decode(wif)
	privByte := wifOriginByte[1:33]

	return PrivKeyFromBytes(S256(), privByte)
}
