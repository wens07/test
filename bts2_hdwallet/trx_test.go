/**
  * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
  *  
  * Copyright © 2015--2018 . All rights reserved.
  *
  * File: trx_test.go 
  * Date: 2018-09-05
  *
  */

package bts

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

func TestTimeConvert(t *testing.T) {

	 tmp_time := time.Now().Unix()

	 fmt.Println(Time2Str(tmp_time))


}

func TestJson(t *testing.T) {

	transferOp := DefaultNoMemoTransferOperation()

	transferTrx := Transaction{
		1,
		2,
		"2018-09-04T08:16:25",
		[][]interface{}{{0, transferOp}},
		make([]interface{}, 0),
		[]string{"2018-09-04T08:16:25"},
		3,
		nil,

	}

	b, err := json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)


}


func TestSignature(t *testing.T) {

	chainidHex := "fe70279c1d9850d4ddb6ca1f00c577bc2e86bf33d54fafd4c606a6937b89ae32"
	//
	//seed := MnemonicToSeed("venture lazy digital aware plug hire acquire abuse chunk know gloom snow much employ glow rich exclude allow", "123")
	//
	//addrKey, _ := GetAddressKey(seed, 0, 0)

	wif := "5JWbFZMcweitDc11YyZenHdiCmLpRWv8mZ8kDeehsqvEprZY6jv"
	sig, _ := GetSignature(wif, []byte(chainidHex))

	fmt.Println(len(sig))
	fmt.Println(sig)




}

func TestPack(t *testing.T) {



	//transferTrx
	BuildTransferTransaction("5057,3759147062", "5K367bRWunQQibV53qiNnyyiDEiooro7TYhCMFZF9vUFHcNV4HU","1.2.18", "1.2.21", "", "",
		                   20000, 12000,"BTS", "4018d7844c78f6a6c41c6a552b898022310fc5dec06da467ee7905a8dad512c8")
	//ref_block_num, ref_block_prefix, _ := GetRefblockInfo("47926,2298351426")
	//
	//fmt.Printf("ref block num is %d\n", ref_block_num)
	//fmt.Printf("ref block prefix is %d\n", ref_block_prefix)

	//test binance deposit
	//BuildTransferTransaction("1.2.397681", "1.2.463816", "100827580", "BTS7yXA4684QDr9Rhz8JfTQEwcu8n8DJE4Z9WJ42GjpufCMVpfs6h",20000, 12000, "BTS")

	//createTrx
	//BuildLimitOrderCreateTransaction(10000000, 2000000, 578, "1.2.397681", "1.3.0", "1.3.113", "BTS")



}

func Test_Serialize(t *testing.T) {

/*	fmt.Printf("max int8 is %d\n", math.MaxInt8)
	fmt.Printf("max int16 is %d\n", math.MaxInt16)
	fmt.Printf("max int32 is %d\n", math.MaxInt32)

	fmt.Printf("max uint8 is %d\n", math.MaxUint8)
	fmt.Printf("max uint16 is %d\n", math.MaxUint16)
	fmt.Printf("max uint32 is %d\n", math.MaxUint32)

	c := math.Pow(0x80, 3)
	fmt.Printf("%.0f", c)*/

	pubkeystr := "BTS7wjh57oXRQTSPPegCfPBuUNnJA5hdShTGcPXJZiVVYxSCk5LRE"

	pubbyte := base58.Decode(pubkeystr[3:])

	fmt.Println(hex.EncodeToString(pubbyte[:len(pubbyte) - 4]))

}
