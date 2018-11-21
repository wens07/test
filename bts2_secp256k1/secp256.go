/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: secp256.go, Date: 2018-09-12
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package bts2_secp256k1

/*
#cgo CFLAGS: -I./libsecp256k1
#cgo CFLAGS: -I./libsecp256k1/src/
#define USE_NUM_NONE
#define USE_FIELD_10X26
#define USE_FIELD_INV_BUILTIN
#define USE_SCALAR_8X32
#define USE_SCALAR_INV_BUILTIN
#define NDEBUG
#include "./libsecp256k1/src/secp256k1.c"
#include "./libsecp256k1/include/secp256k1.h"

*/
import "C"

import (
	"errors"
	"unsafe"
)

var context *C.secp256k1_context_t

func init() {
	// around 20 ms on a modern CPU.
	//context = C.secp256k1_context_create_sign_verify()
	context = C.secp256k1_context_create_all()
}

var (
	ErrInvalidMsgLen       = errors.New("invalid message length, need 32 bytes")
	ErrInvalidSignatureLen = errors.New("invalid signature length")
	ErrInvalidRecoveryID   = errors.New("invalid signature recovery id")
	ErrInvalidKey          = errors.New("invalid private key")
	ErrSignFailed          = errors.New("signing failed")
	ErrRecoverFailed       = errors.New("recovery failed")
)

/** Create a compact ECDSA signature (64 byte + recovery id).
 *  Returns: 1: signature created
 *           0: the nonce generation function failed, or the secret key was invalid.
 *  In:      ctx:    pointer to a context object, initialized for signing (cannot be NULL)
 *           msg32:  the 32-byte message hash being signed (cannot be NULL)
 *           seckey: pointer to a 32-byte secret key (cannot be NULL)
 *           noncefp:pointer to a nonce generation function. If NULL, secp256k1_nonce_function_default is used
 *           ndata:  pointer to arbitrary data used by the nonce generation function (can be NULL)
 *  Out:     sig:    pointer to a 64-byte array where the signature will be placed (cannot be NULL)
 *                   In case 0 is returned, the returned signature length will be zero.
 *           recid:  pointer to an int, which will be updated to contain the recovery id (can be NULL)
 */
func SignCompact(msg []byte, seckey []byte, requireCanonical bool) ([]byte, error) {
	if len(msg) != 32 {
		return nil, ErrInvalidMsgLen
	}
	if len(seckey) != 32 {
		return nil, ErrInvalidKey
	}
	seckeydata := (*C.uchar)(unsafe.Pointer(&seckey[0]))
	if C.secp256k1_ec_seckey_verify(context, seckeydata) != 1 {
		return nil, ErrInvalidKey
	}

	var (
		sig       = make([]byte, 65)
		msgdata   = (*C.uchar)(unsafe.Pointer(&msg[0]))
		noncefunc = C.bts_extended_nonce_function
		sigstruct = (*C.uchar)(unsafe.Pointer(&sig[1]))
		recid     C.int
		count     C.uint
	)

	for {

		if C.secp256k1_ecdsa_sign_compact(context, msgdata, sigstruct, seckeydata, noncefunc, unsafe.Pointer(&count), &recid) == 0 {
			return nil, ErrSignFailed
		}

		//fmt.Println("in sign compact loop!")

		if IsCanonical(sig) {
			break
		}

	}

	sig[0] = 27 + 4 + byte(recid) // set head byte
	return sig, nil
}

// RecoverPubkey returns the public key of the signer.
// msg must be the 32-byte hash of the message to be signed.
// sig must be a 65-byte compact ECDSA signature containing the
// recovery id a
func RecoverPubkey(msg []byte, sig []byte, compressed bool) ([]byte, error) {
	if len(msg) != 32 {
		return nil, ErrInvalidMsgLen
	}
	if err := checkSignature(sig); err != nil {
		return nil, err
	}

	var pubkey []byte
	var compress_for_recover C.int

	if compressed {
		pubkey = make([]byte, 33)
		compress_for_recover = C.int(1)
	} else {
		pubkey = make([]byte, 65)
		compress_for_recover = C.int(0)
	}

	var (
		sigdata = (*C.uchar)(unsafe.Pointer(&sig[1]))
		msgdata = (*C.uchar)(unsafe.Pointer(&msg[0]))
		publen  C.int
		recid   C.int
	)

	recid = C.int((sig[0] - 27) & 0x03)

	if C.secp256k1_ecdsa_recover_compact(context, msgdata, sigdata, (*C.uchar)(unsafe.Pointer(&pubkey[0])), &publen, compress_for_recover, recid) == 0 {
		return nil, ErrRecoverFailed
	}
	return pubkey, nil
}

func checkSignature(sig []byte) error {
	if len(sig) != 65 {
		return ErrInvalidSignatureLen
	}
	if ((sig[0] - 27) & 0x03) >= 4 {
		return ErrInvalidRecoveryID
	}
	return nil
}

func IsCanonical(sig []byte) bool {

	tmp := (sig[0]&0x80 != 0) || (sig[0] == 0x0 && (sig[1]&0x80 != 0)) || (sig[32]&0x80 != 0) || (sig[32] == 0x0 && (sig[33]&0x80 != 0))

	return !tmp
}
