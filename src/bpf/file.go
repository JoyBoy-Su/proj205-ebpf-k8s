package bpf

import (
	"io"
	"os"
	"path"
)

func InstAdd(inst_name string, package_name string) {
	dirpath := BPF_INST_HOME + inst_name
	err := os.Mkdir(dirpath, 0777)
	if err != nil {
		panic(err)
	}
	// 创建与package之间的软连接
	src := BPF_PACKAGE_HOME + package_name
	dst := BPF_INST_HOME + inst_name + "/" + PACKAGE_FILE_NAME
	err = os.Symlink(src, dst)
	if err != nil {
		panic(err)
	}
}

func InstRead(bpf_name string) (string, string) {
	var dir string = BPF_INST_HOME + bpf_name
	// read pod
	pod_content, err := os.ReadFile(dir + "/" + POD_FILE_NAME)
	var pod_name string = string(pod_content)
	if err != nil {
		panic(err)
	}
	// read src
	src_content, err := os.ReadFile(dir + "/" + SRC_FILE_NAME)
	var src_name string = string(src_content)
	if err != nil {
		panic(err)
	}
	return pod_name, src_name
}

func InstList() []string {
	var bpfs []string
	files, err := os.ReadDir(BPF_INST_HOME)
	if err != nil {
		panic(err)
	}
	for _, bpf_name := range files {
		bpfs = append(bpfs, bpf_name.Name())
	}
	return bpfs
}

// 校验bpf name是否可用
func InstExist(name string) (bool, error) {
	pathname := BPF_INST_HOME + name
	_, err := os.Stat(pathname)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 创建package目录，管理package对应的信息
func PackageCreate(package_name string) {
	dirpath := BPF_PACKAGE_HOME + package_name
	err := os.Mkdir(dirpath, 0777)
	if err != nil {
		panic(err)
	}
	_, err = os.Create(dirpath + "/" + SRC_FILE_NAME)
	if err != nil {
		panic(err)
	}
	err = os.Mkdir(dirpath+"/"+DATA_DIR_NAME, 0777)
	if err != nil {
		panic(err)
	}
}

// 根据src的路径，添加src到package中
func PackageAddSrc(package_name string, src_path string) {
	// 打开源文件
	src_file, err := os.Open(src_path)
	if err != nil {
		panic(err)
	}
	defer src_file.Close()

	// 确定dest的path
	dst_name := path.Base(src_path)
	dst_path := BPF_PACKAGE_HOME + package_name + "/" + DATA_DIR_NAME + "/" + dst_name
	// 打开目标文件
	dst_file, err := os.Create(dst_path)
	if err != nil {
		panic(err)
	}
	defer dst_file.Close()

	// 复制文件
	_, err = io.Copy(dst_file, src_file)
	if err != nil {
		panic(err)
	}
	// 添加name到src
	file, err := os.OpenFile(BPF_PACKAGE_HOME+package_name+"/"+SRC_FILE_NAME, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// 写入file
	_, err = file.WriteString(INFO_SEPARATOR + dst_name)
	if err != nil {
		panic(err)
	}
}

func PackageAddSrcList(package_name string, base_path string, files_list []string) {
	// 打开src
	file, err := os.OpenFile(BPF_PACKAGE_HOME+package_name+"/"+SRC_FILE_NAME, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 循环处理所有file
	for _, file_name := range files_list {
		// 打开源文件
		src_path := base_path + "/" + file_name
		src_file, err := os.Open(src_path)
		if err != nil {
			panic(err)
		}

		// 确定dest的path
		dst_path := BPF_PACKAGE_HOME + package_name + "/" + DATA_DIR_NAME + "/" + file_name
		// 打开目标文件
		dst_file, err := os.Create(dst_path)
		if err != nil {
			panic(err)
		}

		// 复制文件
		_, err = io.Copy(dst_file, src_file)
		if err != nil {
			panic(err)
		}

		// 关闭文件
		src_file.Close()
		dst_file.Close()

		// 写入src
		_, err = file.WriteString(INFO_SEPARATOR + file_name)
		if err != nil {
			panic(err)
		}
	}
}

func DeletePackage() {

}

func ListPackage() {

}

func ReadPackage() {

}
