/**
  * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
  *  
  * Copyright © 2015--2018 . All rights reserved.
  *
  * File: operation.go, Date: 2018-10-24
  *
  *
  * This library is free software under the terms of the GNU General Public License 
  * as published by the Free Software Foundation; either version 3 of the License, 
  * or (at your option) any later version.
  *
  */

package bts


type Asset struct {
	Bts_amount      int64 `json:"amount"`
	Bts_asset_id    string `json:"asset_id"`
}

//
// bts  "1.3.0"
func DefaultAsset() Asset {
	return Asset{
		0,
		"1.3.0",
	}
}

type Extension struct {
	extension  []string
}

type Memo struct {
	Bts_from     string `json:"from"`   //public_key_type  33
	Bts_to       string `json:"to"`     //public_key_type  33
	Bts_nonce    uint64 `json:"nonce"`
	Bts_message  string `json:"message"`

	IsEmpty     bool    `json:"-"`
	Message     string  `json:"-"`

}

func DefaultMemo() Memo {

	return Memo{
		"BTS1111111111111111111111111111111114T1Anm",
		"BTS1111111111111111111111111111111114T1Anm",
		0,
		"",
		true,
		"",
	}

}


// transfer operation tag is  0
type NoMemoTransferOperation struct {

	Bts_fee        Asset  `json:"fee"`
	Bts_from       string `json:"from"`
	Bts_to         string `json:"to"`

	Bts_amount     Asset   `json:"amount"`

	Bts_extensions []interface{} `json:"extensions"`
}

type MemoTransferOperation struct {

	Bts_fee        Asset  `json:"fee"`
	Bts_from       string `json:"from"`
	Bts_to         string `json:"to"`

	Bts_amount     Asset   `json:"amount"`
	Bts_memo       Memo    `json:"memo"`

	Bts_extensions []interface{} `json:"extensions"`
}

func DefaultNoMemoTransferOperation() *NoMemoTransferOperation {

	return &NoMemoTransferOperation{
		DefaultAsset(),
		"1.2.0",
		"1.2.0",
		DefaultAsset(),
		make([]interface{}, 0),
	}
}

func DefaultMemoTransferTransaction() *MemoTransferOperation {

	return &MemoTransferOperation{
		DefaultAsset(),
		"1.2.0",
		"1.2.0",

		DefaultAsset(),
		DefaultMemo(),
		make([]interface{}, 0),
	}
}


type LimitOrderCreateOperation struct {

	Bts_fee                Asset  `json:"fee"`
	Bts_seller             string `json:"seller"`
	Bts_amount_to_sell     Asset  `json:"amount_to_sell"`
	Bts_min_to_receive     Asset  `json:"min_to_receive"`
	Bts_expiration         string `json:"expiration"`
    Bts_fill_or_kill       bool `json:"fill_or_kill"`
	Bts_extensions         []interface{} `json:"extensions"`

	Expiration             uint32   `json:"-"`

}

func DefaultLimitOrderCreateOperation() *LimitOrderCreateOperation {

	return &LimitOrderCreateOperation{
		DefaultAsset(),
		"1.2.0",
		DefaultAsset(),
		DefaultAsset(),
		"",
		false,
		make([]interface{}, 0),
		0,
	}
}

type LimitOrderCancelOperation struct {
	Bts_fee        				Asset   `json:"fee"`
	Bts_order      				string  `json:"order"`
	Bts_fee_paying_account      string  `json:"fee_paying_account"`
	Bts_extensions              []interface{} `json:"extensions"`

}

func DefaultLimitOrderCancelOperation()  *LimitOrderCancelOperation {

	return &LimitOrderCancelOperation{
		DefaultAsset(),
		"1.7.0",
		"1.2.0",
		make([]interface{}, 0),
	}
}


