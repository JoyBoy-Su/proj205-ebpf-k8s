package bpf

var BPF_HOME string = "/home/ubuntu/.kube/bpf/"
var BPF_INST_HOME string = BPF_HOME + "instances/"
var BPF_PACKAGE_HOME string = BPF_HOME + "packages/"

var BPF_EMPTY_PACKAGE_NAME = ""
var POD_FILE_NAME string = "pod"
var SRC_FILE_NAME string = "src"
var DATA_DIR_NAME string = "data"
var PACKAGE_FILE_NAME string = "package"

var INFO_SEPARATOR string = ":"

var BPF_NAMESPACE string = "default"
var CompileImage string = "jiadisu/ecc-min-ubuntu-x86:0.1"
var CompileMountPath string = "/code"
var RunImage string = "ngccc/ecli_x86_ubuntu"
var RunMountPath string = "/var/ebpfPackage/"
var RunCommand = []string{"/bin/sh", "-c", "./ecli run /var/ebpfPackage/package.json"}
