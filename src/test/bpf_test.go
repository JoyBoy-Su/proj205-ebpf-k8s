package test

import (
	"testing"

	"fudan.edu.cn/swz/bpf/bpf"
)

func TestAddSrc(t *testing.T) {
	var package_name string = "test_add_src_package"
	// 创建test_package
	bpf.PackageCreate(package_name)
	// 添加src到package
	bpf.PackageAddSrc(package_name, "/home/ubuntu/jiadisu/bpf/opensnoop.h")
}

func TestAddSrcList(t *testing.T) {
	var package_name string = "test_add_src_list_package"
	bpf.PackageCreate(package_name)
	// 添加src list到package
	var args []string
	args = append(args, "opensnoop.h")
	args = append(args, "opensnoop.bpf.c")
	bpf.PackageAddSrcList(package_name, "/home/ubuntu/jiadisu/bpf", args)
}

func TestMountPackage(t *testing.T) {
	bpf.MountPackageByConfigMap("7e8fa564-ccd1-4517-bbd4-e08aa82ba0e3")
}

func TestAddInst(t *testing.T) {
	bpf.InstAdd("test", "ae535746-43d1-4a5c-a3bb-f9e1516d6b25")
}
