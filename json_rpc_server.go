//package main
//
//import (
//	"log"
//	"net/http"
//
//	"github.com/gorilla/mux"
//	"github.com/gorilla/rpc"
//	"github.com/gorilla/rpc/json"
//)
//
//type Args struct {
//	A, B int
//}
//
//type Arith int
//
//type Result int
//
//func (t *Arith) Multiply(r *http.Request, args *Args, result *Result) error {
//	log.Printf("Multiplying %d with %d\n", args.A, args.B)
//	*result = Result(args.A * args.B)
//	return nil
//}
//
//func main() {
//	s := rpc.NewServer()
//	s.RegisterCodec(json.NewCodec(), "application/json")
//	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
//	arith := new(Arith)
//	s.RegisterService(arith, "")
//	r := mux.NewRouter()
//	r.Handle("/rpc", s)
//	http.ListenAndServe(":1234", r)
//}

package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Args struct {
	X, Y int
}

type Calculator struct{}

func (t *Calculator) Add(args *Args, reply *int) error {
	*reply = args.X + args.Y
	return nil
}

func main() {
	cal := new(Calculator)
	server := rpc.NewServer()
	server.Register(cal)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Printf("new connection established\n")
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}
