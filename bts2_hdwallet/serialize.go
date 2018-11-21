/**
  * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
  *  
  * Copyright © 2015--2018 . All rights reserved.
  *
  * File: serialize.go 
  * Date: 2018-09-07
  *
  */

package bts

import (
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

// inferface for serialize bts transaction
type BtsSearilze interface {

	Serialize() []byte

}



/**
 *  some basic type serialization function
 */
//func PackUint32(writer *bytes.Buffer, val uint32) ([]byte, error) {
//
//	uint64_val := uint64(val)
//
//	for {
//		uint8_val := uint8(uint64_val) & 0x7F
//
//		uint64_val >>= 7
//
//		if uint64_val > 0 {
//			uint8_val |= 0x1 << 7
//		} else {
//			uint8_val |= 0x0 << 7
//		}
//
//		err := writer.WriteByte(uint8_val)
//		if err != nil {
//			return nil, fmt.Errorf("in PackUint32 function, write byte failed: %v", err)
//		}
//
//		if uint64_val == 0 {
//			break
//		}
//
//	}
//
//	return writer.Bytes(), nil
//
//}
//
//
//func UnPackUint32(reader *bytes.Reader) (uint32, error) {
//
//	var uint32_val uint32 = 0
//	var by uint8 = 0
//	for {
//		uint8_val, err := reader.ReadByte()
//		if err != nil {
//			return 0, fmt.Errorf("in UnPackUint32 function, read byte failed: %v", err)
//		}
//
//		uint32_val |= uint32(uint8_val & 0x7F) << by
//
//		by += 7
//
//		if (uint8_val & 0x80) == 0 {
//			break
//		}
//
//	}
//
//	return uint32_val, nil
//}


func PackUint16(val uint16, isLittleEndian bool) []byte {

	res := make([]byte, 2)

	if isLittleEndian {
		binary.LittleEndian.PutUint16(res, val)
	} else {
		binary.BigEndian.PutUint16(res,val)
	}


	return res

}

func UnPackUint16(bytes []byte, isLittleEndian bool) uint16 {

	var res uint16

	if isLittleEndian {
		res = binary.LittleEndian.Uint16(bytes)
	} else {
		res = binary.BigEndian.Uint16(bytes)
	}

	return res
}

func PackUint32(val uint32, isLittleEndian bool) []byte {

	res := make([]byte, 4)

	if isLittleEndian {
		binary.LittleEndian.PutUint32(res, val)
	} else {
		binary.BigEndian.PutUint32(res,val)
	}


	return res

}

func UnPackUint32(bytes []byte, isLittleEndian bool) uint32 {

	var res uint32

	if isLittleEndian {
		res = binary.LittleEndian.Uint32(bytes)
	} else {
		res = binary.BigEndian.Uint32(bytes)
	}

	return res
}

func PackInt64(val int64, isLittleEndian bool) []byte {

	res := make([]byte, 8)

	if isLittleEndian {
		binary.LittleEndian.PutUint64(res, uint64(val))
	} else {
		binary.BigEndian.PutUint64(res, uint64(val))
	}

	return res
}

func UnPackInt64(bytes []byte, isLittleEndian bool) int64 {

	var res int64

	if isLittleEndian {
		res = int64(binary.LittleEndian.Uint64(bytes))
	} else {
		res = int64(binary.BigEndian.Uint64(bytes))
	}

	return res
}

func PackVarUint32(val uint32) []byte {

	res := make([]byte, 0)

	//one byte
	if val < 0x80 {

		res = append(res, byte(val))

		return res
	} else if val < 0x4000 { //two byte

		byte1 := val / 0x80
		byte2 := val % 0x80 + 0x80

		res = append(res, byte(byte2))
		res = append(res, byte(byte1))


	} else if val < 0x200000 { //three byte

		byte1 := val / 0x4000
		byte2 := val % 0x4000 / 0x80 + 0x80
		byte3 := val % 0x80 + 0x80

		res = append(res, byte(byte3))
		res = append(res, byte(byte2))
		res = append(res, byte(byte1))

	} else if val < 0x10000000 { //four byte

		byte1 := val / 0x200000
		byte2 := val % 0x200000 / 0x4000 + 0x80
		byte3 := val % 0x4000 / 0x80 + 0x80
		byte4 := val % 0x80 + 0x80

		res = append(res, byte(byte4))
		res = append(res, byte(byte3))
		res = append(res, byte(byte2))
		res = append(res, byte(byte1))
	} else {

		byte1 := val / 0x10000000
		byte2 := val % 0x10000000 / 0x200000 + 0x80
		byte3 := val % 0x200000 / 0x4000 + 0x80
		byte4 := val % 0x4000 / 0x80 + 0x80
		byte5 := val % 0x80 + 0x80

		res = append(res, byte(byte5))
		res = append(res, byte(byte4))
		res = append(res, byte(byte3))
		res = append(res, byte(byte2))
		res = append(res, byte(byte1))

	}


	return res
}


func (asset *Asset) Serialize() []byte {

	byte_int64 := PackInt64(asset.Bts_amount, true)

	//byte for asset_id_type, default to zero
	if asset.Bts_asset_id == "1.3.0" {
		byte_int64 = append(byte_int64, byte(0))
	} else if asset.Bts_asset_id == "1.3.1" {
		byte_int64 = append(byte_int64, byte(1))
	} else if asset.Bts_asset_id == "1.3.2" {
		byte_int64 = append(byte_int64, byte(2))
	}


	return byte_int64
}

func (memo *Memo) Serialize() []byte {

	if memo.IsEmpty {
		return []byte{0}
	} else {

		//byte for optional, have element default to one
		var res []byte
		res = append(res, byte(1))
		//for all null
		//byte_pub := make([]byte, 74)
		//res = append(res, byte_pub...)
		byte_pub := make([]byte, 33)
		res = append(res, byte_pub...)
		pubbyte := base58.Decode(memo.Bts_to[3:])
		res = append(res, pubbyte[:len(pubbyte) - 4]...)
		byte_pub = make([]byte, 8)
		res = append(res, byte_pub...)

		// memo message
		res = append(res, byte(len(memo.Message)+4))
		byte_pub = make([]byte, 4)
		res = append(res, byte_pub...)
		res = append(res, []byte(memo.Message)...)
		return res

	}

}

func (tranferOp  *NoMemoTransferOperation) Serialize() []byte {

	res := tranferOp.Bts_fee.Serialize()

	id_from, err := GetAccountId(tranferOp.Bts_from)
	if err != nil {
		fmt.Println(err)
		panic(id_from)
	}
	byte_uint32 := PackVarUint32(id_from)
	res = append(res, byte_uint32...)
	id_to, err := GetAccountId(tranferOp.Bts_to)
	if err != nil {
		fmt.Println(err)
		panic(id_to)
	}
	byte_uint32 = PackVarUint32(id_to)
	res = append(res, byte_uint32...)

	byteTmp := tranferOp.Bts_amount.Serialize()
	res = append(res, byteTmp...)

	res = append(res, byte(0))
	res = append(res, byte(0))

	return res

}

func (tranferOp  *MemoTransferOperation) Serialize() []byte {

	res := tranferOp.Bts_fee.Serialize()

	id_from, err := GetAccountId(tranferOp.Bts_from)
	if err != nil {
		fmt.Println(err)
		panic(id_from)
	}
	byte_uint32 := PackVarUint32(id_from)
	res = append(res, byte_uint32...)
	id_to, err := GetAccountId(tranferOp.Bts_to)
	if err != nil {
		fmt.Println(err)
		panic(id_to)
	}
	byte_uint32 = PackVarUint32(id_to)
	res = append(res, byte_uint32...)

	byteTmp := tranferOp.Bts_amount.Serialize()
	res = append(res, byteTmp...)

	byteTmp = tranferOp.Bts_memo.Serialize()
	res = append(res, byteTmp...)
	res = append(res, byte(0))

	return res

}

func (createOrderOp *LimitOrderCreateOperation) Serialize() []byte {

	res := createOrderOp.Bts_fee.Serialize()

	sell_id, err := GetAccountId(createOrderOp.Bts_seller)
	if err != nil {
		fmt.Println(err)
		panic(sell_id)
	}
	byte_uint32 := PackVarUint32(sell_id)
	res = append(res, byte_uint32...)

	byteTmp := createOrderOp.Bts_amount_to_sell.Serialize()
	res = append(res, byteTmp...)

	byteTmp = createOrderOp.Bts_min_to_receive.Serialize()
	res = append(res, byteTmp...)
	res = append(res, byte(0))
	res = append(res, byte(0))

	return res
}

func (cancelOrderOp *LimitOrderCancelOperation) Serialize() []byte {

	res := cancelOrderOp.Bts_fee.Serialize()

	order_id, err := GetAccountId(cancelOrderOp.Bts_order)
	if err != nil {
		fmt.Println(err)
		panic(order_id)
	}
	byte_uint32 := PackVarUint32(order_id)
	res = append(res, byte_uint32...)

	account_id, err := GetAccountId(cancelOrderOp.Bts_fee_paying_account)
	if err != nil {
		fmt.Println(err)
		panic(account_id)
	}
	byte_uint32 = PackVarUint32(account_id)
	res = append(res, byte_uint32...)
	res = append(res, byte(0))

	return res

}


func (trx  *Transaction) Serialize() []byte {

	var res []byte
	res = append(res, PackUint16(trx.Bts_ref_block_num, true)...)
	res = append(res, PackUint32(trx.Bts_ref_block_prefix, true)...)
	res = append(res, PackUint32(trx.Expiration, true)...)

	//operations
	res = append(res, byte(len(trx.Operations)))
	res = append(res, byte(0))
	for _, v := range trx.Operations {

		if noMemo, ok := v.(NoMemoTransferOperation); ok {
			res = append(res, noMemo.Serialize()...)
		} else if  memo, ok := v.(MemoTransferOperation); ok{
			res = append(res, memo.Serialize()...)
		} else if createOrder, ok := v.(LimitOrderCreateOperation); ok {
			res = append(res, createOrder.Serialize()...)
		}

	}

	//extension
	res = append(res, byte(0))

	//signature
	if len(trx.Bts_signatures) > 0 {
		res = append(res, byte(len(trx.Bts_signatures)))
	}


	return res
}

