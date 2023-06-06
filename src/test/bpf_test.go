package test

import (
	"testing"

	"fudan.edu.cn/swz/bpf/bpf"
)

// func TestAddSrc(t *testing.T) {
// 	var package_name string = "test_add_src_package"
// 	// 创建test_package
// 	bpf.PackageCreate(package_name)
// 	// 添加src到package
// 	bpf.PackageAddSrc(package_name, "/home/ubuntu/jiadisu/bpf/opensnoop.h")
// }

// func TestAddSrcList(t *testing.T) {
// 	var package_name string = "test_add_src_list_package"
// 	bpf.PackageCreate(package_name)
// 	// 添加src list到package
// 	var args []string
// 	args = append(args, "opensnoop.h")
// 	args = append(args, "opensnoop.bpf.c")
// 	bpf.PackageAddSrcList(package_name, "/home/ubuntu/jiadisu/bpf", args)
// }

// func TestMountPackage(t *testing.T) {
// 	bpf.MountPackageByConfigMap("7e8fa564-ccd1-4517-bbd4-e08aa82ba0e3")
// }

// func TestAddInst(t *testing.T) {
// 	bpf.InstAdd("test", "ae535746-43d1-4a5c-a3bb-f9e1516d6b25")
// }

// func TestDeletePackage(t *testing.T) {
// 	// 创建package
// 	bpf.PackageCreate("test_package")
// 	// 创建instance
// 	bpf.InstAdd("test_inst", "test_package")
// 	bpf.PackageDelete("test_package", true)
// }

// func TestListPackage(t *testing.T) {
// 	packages := bpf.PackageList()
// 	for _, package_name := range packages {
// 		t.Logf("package name: %s\n", package_name)
// 	}
// }

// func TestReadPackage(t *testing.T) {
// 	package_name := "2a9d4745-3555-4912-9863-bf373ee69b28"
// 	var package_info bpf.PackageInfo
// 	bpf.PackageRead(package_name, &package_info)
// 	t.Logf("package_name: %s\n", package_name)
// 	t.Logf("package_src_list: %q\n", bpf.PackageInfoGetSrcList(&package_info))
// 	t.Logf("package_inst_list: %q\n", bpf.PackageInfoGetInstList(&package_info))
// }

func TestReadSrc(t *testing.T) {
	package_name := "2a9d4745-3555-4912-9863-bf373ee69b28"
	src_name := "opensnoop.h"
	content := bpf.SrcRead(package_name, src_name)
	t.Logf("%s", content)
}
