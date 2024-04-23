package consul

import (
	"fmt"
	"sync"
)

type ConsulConfig struct {
	ConsulAddr string
}

type ConsulUtils struct {
	ConsulCfg ConsulConfig
}

var onceConusl sync.Once
var ConsulInstance *ConsulUtils

func NewConsulUtils() *ConsulUtils {
	onceConusl.Do(func() {
		ConsulInstance = &ConsulUtils{}
	})
	return ConsulInstance
}

func SetConsulAddr(addr string) {
	utils := NewConsulUtils()
	if utils != nil {
		utils.ConsulCfg.ConsulAddr = addr
		fmt.Printf("set  consulAddr=%s suuc \n", addr)
	}
}

func GetConsulAddr() string {
	utils := NewConsulUtils()
	if utils != nil {
		fmt.Printf("get  consulAddr=%s suuc \n", utils.ConsulCfg.ConsulAddr)
		return utils.ConsulCfg.ConsulAddr
	}
	fmt.Printf("utils:%+v", utils)
	return ""
}
