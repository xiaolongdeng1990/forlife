package flsvr

import (
	"errors"
	"fmt"
	"strings"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	cserver "github.com/rpcxio/rpcx-consul/serverplugin"
	"github.com/smallnest/rpcx/server"

	"github.com/xiaolongdeng1990/forlife/MSF/config"
	fllog "github.com/xiaolongdeng1990/forlife/MSF/log"
)

type SvrCfg struct {
	Server struct {
		Name       string `default:""`
		Address    string `default:""`
		ConsulAddr string `default:""`
	}
}

func Server(cfg string, svrHandle interface{}) error {
	svrAddr, consulAddr, basePath, svrName, err := loadSvrCfgInfo(cfg)
	if err != nil {
		return err
	}
	s := server.NewServer()
	registerConuslPlugin(s, svrAddr, consulAddr, basePath)
	s.RegisterName(svrName, svrHandle, "")

	if err := s.Serve("tcp", svrAddr); err != nil {
		fllog.Error("serve failed. err:%+v", err)
		return err
	}
	fmt.Println("start server succ")
	fllog.Error("start server succ")
	return nil
}

func loadSvrCfgInfo(cfg string) (string, string, string, string, error) {
	svrCfg := SvrCfg{}
	if err := config.ParseConfigWithPath(&svrCfg, cfg); err != nil {
		fmt.Printf("load svr logcfg failed. err:%+v cfg:%s", err, cfg)
		fllog.Error("load svr logcfg failed. err:%+v cfg:%s", err, cfg)
		return "", "", "", "", err
	}
	fmt.Printf("svrCfg:%+v", svrCfg)
	fllog.Debug("svrCfg:%+v", svrCfg)
	if len(svrCfg.Server.Name) == 0 || len(svrCfg.Server.Address) == 0 {
		return "", "", "", "", errors.New("svrcfg invalid")
	}

	basePath, svrName := parseSvrName(svrCfg.Server.Name)
	if len(basePath) == 0 || len(svrName) == 0 {
		fllog.Error("basePath:%s svrName:%s", basePath, svrName)
		return "", "", "", "", errors.New("parse server name failed")
	}
	// to-do
	// if len(svrCfg.Server.ConsulAddr) == 0 {
	// }
	return svrCfg.Server.Address, svrCfg.Server.ConsulAddr, basePath, svrName, nil
}

func registerConuslPlugin(s *server.Server, svrAddr, conuslAddr, basePath string) {
	r := &cserver.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + svrAddr,
		ConsulServers:  []string{conuslAddr},
		BasePath:       basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		fmt.Println(err)
		fllog.Error("register consul failed. err:%+v", err)
	}

	s.Plugins.Add(r)
	fmt.Println("add register succ")
	fllog.Debug("register consul succ")
	return
}

func parseSvrName(name string) (string, string) {
	vecSplit := strings.Split(name, ".")
	if len(vecSplit) != 2 {
		return "", ""
	}
	return vecSplit[0], vecSplit[1]
}
