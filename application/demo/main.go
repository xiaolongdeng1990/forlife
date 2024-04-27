package main

import (
	"context"
	"flag"
	"fmt"

	fllog "github.com/xiaolongdeng1990/forlife/MSF/log"
	flsvr "github.com/xiaolongdeng1990/forlife/MSF/server"
	math "github.com/xiaolongdeng1990/forlife/protocol/json/math"
)

var (
	cfg string
)

func init() {
	flag.StringVar(&cfg, "c", "../conf/rpcx_demo.toml", "config file path, default ../conf/rpcx_demo.toml")
}

func Mul(ctx context.Context, args *math.Args, reply *math.Reply) error {
	reply.C = args.A * args.B
	fllog.Log().Debug("req=", args, "reply=", reply)
	return nil
}

func Add(ctx context.Context, args *math.Args, reply *math.Reply) error {
	reply.C = args.A + args.B
	fllog.Log().Debug("req=", args, "reply=", reply)

	// client rpc
	return nil
}

func main() {
	flag.Parse()
	// log init
	if err := fllog.Init(cfg); err != nil {
		fmt.Printf("fllog init failed. err:%+v", err)
		return
	}
	// config init.

	fllog.Log().Debug("test fllog debug cfg:", cfg)
	// server init
	svr := flsvr.NewFLServer(cfg)
	svr.RegisterFunc(Mul)
	svr.RegisterFunc(Add)
	svr.StartServer()
}
