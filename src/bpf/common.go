package bpf

import (
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

var BPF_HOME string = filepath.Join(homedir.HomeDir(), ".kube", "bpf") + "/"
var BPF_INST_HOME string = BPF_HOME + "instances/"
var BPF_PACKAGE_HOME string = BPF_HOME + "packages/"

var BPF_EMPTY_PACKAGE_NAME = ""
var BPF_EMPTY_INSTANCE_NAME = ""
var BPF_EMPTY_NODE_NAME = ""
var POD_FILE_NAME string = "pod"
var SRC_FILE_NAME string = "src"
var DATA_DIR_NAME string = "data"
var PACKAGE_FILE_NAME string = "package"
var INSTANCE_DIR_NAME string = "instance"

var INFO_SEPARATOR string = ":"

var BPF_NAMESPACE string = "bpf"
var CompileImage string = "jiadisu/ecc-min-ubuntu-x86:0.1"
var CompileMountPath string = "/code"
var RunImage string = "ngccc/ecli_x86_ubuntu"
var RunMountPath string = "/var/ebpfPackage/"
var RunCommand = []string{"/bin/sh", "-c", "./ecli run /var/ebpfPackage/package.json"}

type InstInfo struct {
	inst_name    string
	package_name string
	src_list     []string
	node         string
}

func InstInfoClear(inst_info *InstInfo) {
	inst_info.inst_name = ""
	inst_info.package_name = ""
	inst_info.node = ""
	inst_info.src_list = nil
}

func InstInfoGetInstName(inst_info *InstInfo) string {
	return inst_info.inst_name
}

func InstInfoGetPackageName(inst_info *InstInfo) string {
	return inst_info.package_name
}

func InstInfoGetNode(inst_info *InstInfo) string {
	return inst_info.node
}

func InstInfoGetSrcList(inst_info *InstInfo) []string {
	return inst_info.src_list
}

type PackageInfo struct {
	package_name string
	src_list     []string
	inst_list    []string
	size         int64
}

func PackageInfoClear(package_info *PackageInfo) {
	package_info.package_name = ""
	package_info.inst_list = nil
	package_info.src_list = nil
	package_info.size = 0
}

func PackageInfoGetPackageName(package_info *PackageInfo) string {
	return package_info.package_name
}

func PackageInfoGetSrcList(package_info *PackageInfo) []string {
	return package_info.src_list
}

func PackageInfoGetInstList(package_info *PackageInfo) []string {
	return package_info.inst_list
}

func PackageInfoGetSize(package_info *PackageInfo) int64 {
	return package_info.size
}
