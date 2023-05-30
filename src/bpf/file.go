package bpf

import (
	"os"
)

func AddBPF(bpf_name string, pod_name string, src []string) {
	dirpath := BPF_HOME + bpf_name
	err := os.Mkdir(dirpath, 0777)
	if err != nil {
		panic(err)
	}
	// 创建并写入pod
	pod_path := dirpath + "/" + POD_FILE
	pod_file, err := os.OpenFile(pod_path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer pod_file.Close()
	pod_file.WriteString(pod_name)
	// 创建并写入src
	src_path := dirpath + "/" + SRC_FILE
	src_file, err := os.OpenFile(src_path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer src_file.Close()
	for _, value := range src {
		src_file.WriteString(value)
		src_file.WriteString(",")
	}
}

func ReadBPF(bpf_name string) (string, string) {
	var dir string = BPF_HOME + bpf_name
	// read pod
	pod_content, err := os.ReadFile(dir + "/" + POD_FILE)
	var pod_name string = string(pod_content)
	if err != nil {
		panic(err)
	}
	// read src
	src_content, err := os.ReadFile(dir + "/" + SRC_FILE)
	var src_name string = string(src_content)
	if err != nil {
		panic(err)
	}
	return pod_name, src_name
}

func ListBPF() []string {
	var bpfs []string
	files, err := os.ReadDir(BPF_HOME)
	if err != nil {
		panic(err)
	}
	for _, bpf_name := range files {
		bpfs = append(bpfs, bpf_name.Name())
	}
	return bpfs
}
