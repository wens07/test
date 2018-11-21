/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2017 . All rights reserved.
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package eth_lib

import (
	"fmt"
	"log"

	"github.com/bitly/go-simplejson"
	"github.com/parnurzeal/gorequest"
)

const api_url_for_test string = "http://192.168.1.178:8546"

//const test_remote_api_url string = "http://192.168.1.124:8588"

func Test_rpc_call(args ...string) *simplejson.Json {

	id := "1"
	method := args[0]
	params := "[]"
	if len(args) > 1 {
		params = args[1]
	}

	postBody := `{"jsonrpc":"2.0","method":"` + method + `","params":` + params + `,"id":` + id + `}`
	fmt.Println("postBody: " + postBody)

	_, body, errs := gorequest.New().Post(api_url_for_test).
		Send(postBody).
		End()

	if errs != nil {
		fmt.Println(errs)
		panic(errs)
	}

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		log.Fatalln(err)
	}

	return js
}
