/**
  * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
  *  
  * Copyright © 2015--2018 . All rights reserved.
  *
  * File: key.go
  * Date: 2018-08-31
  *
  */

package bts

import (
	"crypto/sha512"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/golangcrypto/ripemd160"
	"github.com/tyler-smith/go-bip39"
)


const (
	CoinBTS string = "BTS"
)

/**
 * BTS address struct
 */
type btsPubkeyStr struct {
	pubkeyStr  string   //bts main chain address string
}


// using BIP44 to manage bts address
// m / purpose' / coin_type' / account' / change / address_index
// https://github.com/satoshilabs/slips/blob/master/slip-0044.md

func MnemonicToSeed(mnemonic, password string) []byte {
	return bip39.NewSeed(mnemonic, password)
}


func getMasterkey(seed []byte, mainnet bool) (*hdkeychain.ExtendedKey, error) {

	// main net
	if mainnet {

		return hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)

	} else { //test net

		return hdkeychain.NewMaster(seed, &chaincfg.TestNet3Params)

	}


}

func getAccountExtentkey(masterKey *hdkeychain.ExtendedKey, account uint32, addrIndex uint32) (*hdkeychain.ExtendedKey, error) {

	// purpose & coin_tyep & change
	purpose := uint32(0x8000002C)
	coinType := uint32(0x80000000)
	change := uint32(0)


	purposeKey, err := masterKey.Child(purpose)
	if err != nil {
		return nil, fmt.Errorf("create purpose key failed: %v", err)

	}

	coinTypeKey, err := purposeKey.Child(coinType)
	if err != nil {
		return nil, fmt.Errorf("create coin type key failed: %v", err)

	}

	accountKey, err := coinTypeKey.Child(account)
	if err != nil {
		return nil, fmt.Errorf("create account key failed: %v", err)

	}


	changeKey, err := accountKey.Child(change)
	if err != nil {
		return nil, fmt.Errorf("create change key failed: %v", err)

	}

	addressKey, err := changeKey.Child(addrIndex)

	if err != nil {
		return nil, fmt.Errorf("create address key failed: %v", err)

	}

	return addressKey, err

}

func getAccountAddr(addressKey *hdkeychain.ExtendedKey) (string, error) {

	ecPubkey, err := addressKey.ECPubKey()
	if err != nil {
		return "", fmt.Errorf("get ecPubkey failed: %v", err)
	}

	// sha512 & ripemd160
	pubkeyByte := ecPubkey.SerializeCompressed()
	myRipemd := ripemd160.New()
	sha512Byte := sha512.Sum512(pubkeyByte)
	myRipemd.Write(sha512Byte[:])
	addrByte := myRipemd.Sum(nil)
	//fmt.Println(len(addrByte))

	myRipemd.Reset()
	myRipemd.Write(addrByte)
	addrByteChecksum := myRipemd.Sum(nil)

	addrByte = append(addrByte, addrByteChecksum[0:4]...)
	//fmt.Println(len(addrByte))

	return CoinBTS + string(base58.Encode(addrByte)), nil

}

func getAddressBytes(addressKey *hdkeychain.ExtendedKey) ([]byte, error) {

	ecPubkey, err := addressKey.ECPubKey()
	if err != nil {
		return nil, fmt.Errorf("get ecPubkey failed: %v", err)
	}

	// sha512 & ripemd160
	pubkeyByte := ecPubkey.SerializeCompressed()
	myRipemd := ripemd160.New()
	sha512Byte := sha512.Sum512(pubkeyByte)
	myRipemd.Write(sha512Byte[:])
	addrByte := myRipemd.Sum(nil)

	return addrByte, nil
}


func GetAddress(seed []byte, account uint32, addrIndex uint32) (string, error) {

	mastkey, err := getMasterkey(seed, true) //bts using btc mainchain cfg
	if err != nil {
		return "", fmt.Errorf("in GetAddress function, get mastkey failed: %v", err)
	}

	accountExtendKey, err := getAccountExtentkey(mastkey, account, addrIndex)
	if err != nil {
		return "", fmt.Errorf("in GetAddress function, get accountExtensionKey failed: %v", err)
	}

	accountAddr, err := getAccountAddr(accountExtendKey)
	if err != nil {
		return "", fmt.Errorf("in GetAddress function, get accountAddr failed: %v", err)
	}

	return  accountAddr, nil

}

func GetAddressBytes(addr string) ([]byte, error) {

	if len(addr) <= 2 {
		return nil, fmt.Errorf("in GetAddressBytes function, wrong addr format")
	}

	base58_addr := addr[2:]

	addrBytes := base58.Decode(base58_addr)


	return addrBytes[:len(addrBytes)-4], nil
}

func getWifkey(addressKey *hdkeychain.ExtendedKey) (string, error) {

	ecPrivkey, err := addressKey.ECPrivKey()
	if err != nil {
		return "", fmt.Errorf("get ecPrivkey failed: %v", err)
	}

	wif, err := btcutil.NewWIF(ecPrivkey, &chaincfg.MainNetParams, false)
	if err != nil {
		return "", fmt.Errorf("get wif failed: %v", err)
	}

	return wif.String(), nil
}

func getPrivKey(wif string) (*btcec.PrivateKey, error) {

	wifstruct, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, fmt.Errorf("decode wif string failed: %v", err)
	}

	return wifstruct.PrivKey, err

}


func GetAddressKey(seed []byte, account uint32, addrIndex uint32) (*hdkeychain.ExtendedKey, error) {

	mastkey, err := getMasterkey(seed, true) //bts using btc mainchain cfg
	if err != nil {
		return nil, fmt.Errorf("in GetAddress function, get mastkey failed: %v", err)
	}

	accountExtendKey, err := getAccountExtentkey(mastkey, account, addrIndex)
	if err != nil {
		return nil, fmt.Errorf("in GetAddress function, get accountExtensionKey failed: %v", err)
	}

	return accountExtendKey, nil
}

func GetPubKeyStr(key *btcec.PrivateKey) (string, error) {

	ecPubkey := key.PubKey()

	// ripemd160
	pubkeyByte := ecPubkey.SerializeCompressed()
	myRipemd := ripemd160.New()
	myRipemd.Write(pubkeyByte[:])
	outByte := myRipemd.Sum(nil)

	pubByte := outByte[:4]
	pubByte = append(pubkeyByte, pubByte...)


	return CoinBTS + string(base58.Encode(pubByte)), nil


}

func ExportWif(seed []byte, account uint32, addrIndex uint32) (string, error) {

	mastkey, err := getMasterkey(seed, true) //bts using btc mainchain cfg
	if err != nil {
		return "", fmt.Errorf("in ExportWif function, get mastkey failed: %v", err)
	}

	accountExtendKey, err := getAccountExtentkey(mastkey, account, addrIndex)
	if err != nil {
		return "", fmt.Errorf("in ExportWif function, get accountExtensionKey failed: %v", err)
	}

	wifKey, err := getWifkey(accountExtendKey)
	if err != nil {
		return "", fmt.Errorf("in ExportWif function, get wif failed: %v", err)
	}

	return wifKey, nil

}

func ImportWif(wifstr string) (*btcec.PrivateKey, error) {

	return getPrivKey(wifstr)
}

func IsCanonical(sig []byte) bool {

	tmp := (sig[1] & 0x80 != 0) || ( sig[1] == 0x0 && (sig[2] & 0x80 != 0) ) || (sig[33] & 0x80 != 0) || (sig[33] == 0x0 && (sig[34] & 0x80 != 0))

	return !tmp
}










