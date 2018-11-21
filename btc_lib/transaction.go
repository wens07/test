/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: transaction.go, Date: 2018-11-09
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
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func Bitcoin_TrxFromHexstring(trx string) (*wire.MsgTx, error) {

	var msgTrx wire.MsgTx
	var trxByte bytes.Buffer
	tmpByte, err := hex.DecodeString(trx)
	if err != nil {
		return nil, fmt.Errorf("decode hex trx string error")
	}

	trxByte.Write(tmpByte)

	//msgTrx.DeserializeNoWitness(&trxByte)
	msgTrx.Deserialize(&trxByte)

	return &msgTrx, nil
}

func Bitcoin_CalcSignatureHash(script []byte, hashType txscript.SigHashType, tx *wire.MsgTx, idx int) ([]byte, error) {

	return txscript.CalcSignatureHash(script, hashType, tx, idx)

}
