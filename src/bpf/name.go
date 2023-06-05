package bpf

import (
	"strings"

	"github.com/google/uuid"
)

// 产生一个可用的package的name
func packageName(package_name string) string {
	// 若未指定package的name
	if strings.Compare(package_name, BPF_EMPTY_PACKAGE_NAME) == 0 {
		package_name = uuid.NewString()
	}
	return package_name
}

func instanceName(inst_name string, serial bool) string {
	// 判断是否需要添加serial number
	exist, err := InstExist(inst_name)
	if err != nil {
		panic(err)
	}
	// 已存在或本身设置需要 两个条件都不满足则直接返回inst name
	if !(exist || serial) {
		return inst_name
	}
	// 添加serial number
	return inst_name + "-" + uuid.NewString()
}
