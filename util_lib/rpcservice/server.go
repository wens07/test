/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: server.go, Date: 2018-10-26
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package rpcservice

import (
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

type MyArgs struct {
	Who string
}

type MyReply struct {
	HelloMsg string
}

type MyService struct {
}

var MyServer *rpc.Server

func init() {
	MyServer = rpc.NewServer()
	MyServer.RegisterCodec(json.NewCodec(), "application/json")
	MyServer.RegisterService(new(MyService), "")
	http.Handle("/rpc", MyServer)
}

func (h *MyService) Info(r *http.Request, args *MyArgs, reply *MyReply) error {
	reply.HelloMsg = "Hello " + args.Who + ", Welcome to my rpc service!"
	return nil
}
