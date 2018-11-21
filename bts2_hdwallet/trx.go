/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: trx.go
 * Date: 2018-09-04
 *
 */

package bts

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	bts "anybit.io/spv/bts2_secp256k1"
)

/*const (
	//chain_id string = "1b318870be374fa1234d40e39f87769f7d1fdfd6dffb02fbfc6a09cb14749b88"
	chain_id string = "4018d7844c78f6a6c41c6a552b898022310fc5dec06da467ee7905a8dad512c8"
)*/

// define transfer trx structure
type Transaction struct {
	Bts_ref_block_num    uint16 `json:"ref_block_num"`
	Bts_ref_block_prefix uint32 `json:"ref_block_prefix"`
	Bts_expiration       string `json:"expiration"`

	Bts_operations [][]interface{} `json:"operations"`
	Bts_extensions []interface{}   `json:"extensions"`
	Bts_signatures []string        `json:"signatures"`

	Expiration uint32        `json:"-"`
	Operations []interface{} `json:"-"`
}

func DefaultTransferTransaction() *Transaction {

	return &Transaction{
		0,
		0,
		"",
		nil,
		nil,
		nil,
		0,
		nil,
	}
}

func GetAccountId(id string) (uint32, error) {

	idSlice := strings.Split(id, ".")

	if len(idSlice) != 3 {
		return 0, fmt.Errorf("in GetAccountId function, get account id failed")
	}

	res, err := strconv.ParseUint(idSlice[2], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("in GetAccountId function, Parse id error %v", err)
	}

	return uint32(res), nil

}

func Str2Time(str string) int64 {

	str += "Z"
	t, err := time.Parse(time.RFC3339, str)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	return t.Unix()

}

func Time2Str(t int64) string {

	l_time := time.Unix(t, 0).UTC()
	timestr := l_time.Format(time.RFC3339)

	timestr = timestr[:len(timestr)-1]

	return timestr
}

// in multiple precision mode
func CalculateFee(basic_op_fee int64, len_memo int64) int64 {

	if len_memo == 0 {
		return 2000000
	}

	var basic_memo_fee int64 = 100000
	return basic_op_fee + len_memo*basic_memo_fee
}

func (asset *Asset) SetAssetBySymbol(symbol string) {

	if symbol == "BTS" {
		asset.Bts_asset_id = "1.3.0"
	}

}

func GetRefblockInfo(info string) (uint16, uint32, error) {

	refinfo := strings.Split(info, ",")

	if len(refinfo) != 2 {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, get refblockinfo failed")
	}
	ref_block_num_str, ref_block_prefix_str := refinfo[0], refinfo[1]
	ref_block_num, err := strconv.ParseUint(ref_block_num_str, 10, 16)
	if err != nil {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, convert ref_block_num failed: %v", err)
	}

	ref_block_prefix, err := strconv.ParseUint(ref_block_prefix_str, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, convert ref_block_prefix failed: %v", err)
	}

	return uint16(ref_block_num), uint32(ref_block_prefix), nil
}

func GetSignature(wif string, hash []byte) ([]byte, error) {

	ecPrivkey, err := ImportWif(wif)
	if err != nil {
		return nil, fmt.Errorf("in GetSignature function, get ecprivkey failed: %v", err)
	}

	ecPrivkeyByte := ecPrivkey.Serialize()
	fmt.Println("the uncompressed pubkey is: ", hex.EncodeToString(ecPrivkey.PubKey().SerializeUncompressed()))
	fmt.Println("the compressed pubkey is: ", hex.EncodeToString(ecPrivkey.PubKey().SerializeCompressed()))

	for {
		sig, err := bts.SignCompact(hash, ecPrivkeyByte, true)
		if err != nil {
			return nil, fmt.Errorf("in GetSignature function, sign compact failed: %v", err)
		}

		pubkey_byte, err := bts.RecoverPubkey(hash, sig, true)
		if err != nil {
			return nil, fmt.Errorf("in GetSignature function, sign compact failed: %v", err)
		}
		fmt.Println("recoverd pubkey is: ", hex.EncodeToString(pubkey_byte))

		if bytes.Compare(ecPrivkey.PubKey().SerializeCompressed(), pubkey_byte) == 0 {
			return sig, nil
		}

	}
}

func BuildTransferTransaction(refinfo, wif string, from, to, memo, to_pubkey string, amount, fee int64, symbol string, chain_id string) {

	asset_amount := DefaultAsset()
	asset_amount.Bts_amount = amount
	asset_amount.SetAssetBySymbol(symbol)

	asset_fee := DefaultAsset()
	//asset_fee.Bts_amount = CalculateFee(2000000, int64(len(memo)))
	asset_fee.Bts_amount = fee
	asset_fee.SetAssetBySymbol(symbol)

	expir_sec := time.Now().Unix() + 600
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-10-11T01:54:25"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		panic("get refinfo failed!")
	}
	memo_trx := DefaultMemo()
	var transferTrx Transaction
	if memo != "" {
		memo_trx.Message = memo
		memo_trx.IsEmpty = false
		memo_trx.Bts_message = hex.EncodeToString(append(make([]byte, 4), []byte(memo_trx.Message)...))

		memotransferOp := DefaultMemoTransferTransaction()
		memotransferOp.Bts_fee = asset_fee
		memotransferOp.Bts_from = from
		memotransferOp.Bts_to = to

		memotransferOp.Bts_amount = asset_amount
		memotransferOp.Bts_memo = memo_trx
		memotransferOp.Bts_memo.Bts_to = to_pubkey

		transferTrx = Transaction{
			ref_block_num,
			ref_block_prefix,
			expir_str,
			[][]interface{}{{0, memotransferOp}},
			make([]interface{}, 0),
			make([]string, 0),
			uint32(expir_sec),
			[]interface{}{*memotransferOp},
		}

	} else {

		noMemotransferOp := DefaultNoMemoTransferOperation()
		noMemotransferOp.Bts_fee = asset_fee
		noMemotransferOp.Bts_from = from
		noMemotransferOp.Bts_to = to
		noMemotransferOp.Bts_amount = asset_amount

		transferTrx = Transaction{
			ref_block_num,
			ref_block_prefix,
			expir_str,
			[][]interface{}{{0, noMemotransferOp}},
			make([]interface{}, 0),
			make([]string, 0),
			uint32(expir_sec),
			[]interface{}{*noMemotransferOp},
		}

	}

	res := transferTrx.Serialize()
	fmt.Println("the serialized trx is: ", hex.EncodeToString(res))

	//seed := MnemonicToSeed("venture lazy digital aware plug hire acquire abuse chunk know gloom snow much employ glow rich exclude allow", "123")
	//addrkey, _:= GetAddressKey(seed, 0, 0)
	//addr, _ := GetAddress(seed,0,0)
	//fmt.Println("addr is: ", addr)
	//wif, _ := ExportWif(seed, 0, 0)
	//fmt.Println("wif is: ", wif)

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))
	fmt.Println("serialized for sign is ", hex.EncodeToString(toSign[:]))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
	}

	transferTrx.Bts_signatures = append(transferTrx.Bts_signatures, hex.EncodeToString(sig))
	fmt.Println("found canonical signature")
	fmt.Println(hex.EncodeToString(sig))

	b, err := json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	fmt.Println()

	/*	transferTrx.Bts_signatures = append(transferTrx.Bts_signatures, "2035d9698e310b2ea83506e11c3374981261fd854cbd0a03a6209f8588fc2358ce34e4ea4a05441a57b8f1fc9793a1a5ff5d88e8e306773e739b42a829452f4266")

		b, err := json.Marshal(transferTrx)
		if err != nil {
			fmt.Println("error:", err)
		}
		os.Stdout.Write(b)

		fmt.Println()
		trxSerialize := transferTrx.Serialize()

		fmt.Println("the serialize trx is: ", hex.EncodeToString(trxSerialize))*/

}

func BuildLimitOrderCreateTransaction(refinfo, wif string, sell_amount, buy_amount, fee int64, sell_id, sell_asset_id, buy_asset_id, fee_symbol string, chain_id string) {

	asset_fee := DefaultAsset()
	//asset_fee.Bts_amount = CalculateFee(2000000, int64(len(memo)))
	asset_fee.Bts_amount = fee
	asset_fee.SetAssetBySymbol(fee_symbol)

	asset_sell := DefaultAsset()
	asset_sell.Bts_amount = sell_amount
	asset_sell.Bts_asset_id = sell_asset_id

	asset_buy := DefaultAsset()
	asset_buy.Bts_amount = buy_amount
	asset_buy.Bts_asset_id = buy_asset_id

	expir_sec := time.Now().Unix() + 600
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-10-25T02:11:51"
	//expir_sec := Str2Time(expir_str)

	limitOrderCreateOp := DefaultLimitOrderCreateOperation()

	limitOrderCreateOp.Bts_fee = asset_fee
	limitOrderCreateOp.Bts_seller = sell_id
	limitOrderCreateOp.Bts_amount_to_sell = asset_sell
	limitOrderCreateOp.Bts_min_to_receive = asset_buy

	op_expir_str := Time2Str(expir_sec + 31536000)
	op_expir_sec := uint32(expir_sec + 31536000)
	//op_expir_str := "2018-10-25T03:11:21"
	//op_expir_sec := Str2Time(expir_str)
	//limitOrderCreateOp.Expiration = uint32(expir_sec + 31536000)
	//limitOrderCreateOp.Bts_expiration = Time2Str(expir_sec + 31536000)
	limitOrderCreateOp.Expiration = uint32(op_expir_sec)
	limitOrderCreateOp.Bts_expiration = op_expir_str

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		panic("get refinfo failed!")
	}

	createTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{1, limitOrderCreateOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*limitOrderCreateOp},
	}

	res := createTrx.Serialize()
	fmt.Println("the serialized trx is: ", hex.EncodeToString(res))

	//seed := MnemonicToSeed("venture lazy digital aware plug hire acquire abuse chunk know gloom snow much employ glow rich exclude allow", "123")
	//addrkey, _:= GetAddressKey(seed, 0, 0)
	//addr, _ := GetAddress(seed,0,0)
	//fmt.Println("addr is: ", addr)
	//wif, _ := ExportWif(seed, 0, 0)
	//fmt.Println("wif is: ", wif)

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))
	fmt.Println("serialized for sign is ", hex.EncodeToString(toSign[:]))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
	}

	createTrx.Bts_signatures = append(createTrx.Bts_signatures, hex.EncodeToString(sig))
	fmt.Println("found canonical signature")
	fmt.Println(hex.EncodeToString(sig))

	b, err := json.Marshal(createTrx)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	fmt.Println()

}

func BuildLimitOrderCancelTransaction(refinfo, wif string, account_id, cancel_id string, fee int64, fee_symbol string, chain_id string) {

	asset_fee := DefaultAsset()
	//asset_fee.Bts_amount = CalculateFee(2000000, int64(len(memo)))
	asset_fee.Bts_amount = fee
	asset_fee.SetAssetBySymbol(fee_symbol)

	expir_sec := time.Now().Unix() + 600
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-10-25T02:11:51"
	//expir_sec := Str2Time(expir_str)

	limitOrderCancelOp := DefaultLimitOrderCancelOperation()
	limitOrderCancelOp.Bts_fee = asset_fee
	limitOrderCancelOp.Bts_order = cancel_id
	limitOrderCancelOp.Bts_fee_paying_account = account_id

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		panic("get refinfo failed!")
	}
	cancelTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{1, limitOrderCancelOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*limitOrderCancelOp},
	}

	res := cancelTrx.Serialize()
	fmt.Println("the serialized trx is: ", hex.EncodeToString(res))

	//seed := MnemonicToSeed("venture lazy digital aware plug hire acquire abuse chunk know gloom snow much employ glow rich exclude allow", "123")
	//addrkey, _:= GetAddressKey(seed, 0, 0)
	//addr, _ := GetAddress(seed,0,0)
	//fmt.Println("addr is: ", addr)
	//wif, _ := ExportWif(seed, 0, 0)
	//fmt.Println("wif is: ", wif)

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))
	fmt.Println("serialized for sign is ", hex.EncodeToString(toSign[:]))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
	}

	cancelTrx.Bts_signatures = append(cancelTrx.Bts_signatures, hex.EncodeToString(sig))
	fmt.Println("found canonical signature")
	fmt.Println(hex.EncodeToString(sig))

	b, err := json.Marshal(cancelTrx)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	fmt.Println()

}
