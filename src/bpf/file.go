package bpf

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// 根据package添加一个与之对应的inst到文件管理中
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
	// 在package的instance中添加一个软链接
	src = BPF_INST_HOME + inst_name
	dst = BPF_PACKAGE_HOME + package_name + "/" + INSTANCE_DIR_NAME + "/" + inst_name
	err = os.Symlink(src, dst)
	if err != nil {
		panic(err)
	}
}

// 根据inst查找与之对应的信息：这里只包括可以通过文件读到的pacakge name和src list
func InstRead(inst_name string, inst_info *InstInfo) {
	// 待返回的对象
	inst_info.inst_name = inst_name
	var inst_dir string = BPF_INST_HOME + inst_name
	var package_dir string = inst_dir + "/" + PACKAGE_FILE_NAME
	// 解析软链接，读取package_name
	link_info, err := os.Lstat(package_dir)
	if err != nil {
		panic(err)
	}
	if link_info.Mode()&os.ModeSymlink != 0 {
		package_path, err := os.Readlink(package_dir)
		if err != nil {
			panic(err)
		}
		inst_info.package_name = path.Base(package_path)
	}
	// 通过package_name读取src
	src_content, err := os.ReadFile(package_dir + "/" + SRC_FILE_NAME)
	if err != nil {
		panic(err)
	}
	var src_str string = string(src_content)
	inst_info.src_list = strings.Split(src_str, INFO_SEPARATOR)[1:]
}

// 读取dir获取所有inst的name
func InstList() []string {
	var insts []string
	files, err := os.ReadDir(BPF_INST_HOME)
	if err != nil {
		panic(err)
	}
	for _, inst_file := range files {
		insts = append(insts, inst_file.Name())
	}
	return insts
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

// 根据inst的name删除inst
func InstDelete(inst_name string) {
	// 获取package name
	var inst_dir string = BPF_INST_HOME + inst_name
	var package_dir string = inst_dir + "/" + PACKAGE_FILE_NAME
	// 解析软链接，读取package_name
	link_info, err := os.Lstat(package_dir)
	if err != nil {
		panic(err)
	}
	if link_info.Mode()&os.ModeSymlink != 0 {
		package_path, err := os.Readlink(package_dir)
		if err != nil {
			panic(err)
		}
		// 删除pacakge中对应该inst的软链接
		os.Remove(package_path + "/" + INSTANCE_DIR_NAME + "/" + inst_name)
	}
	// 删除inst的文件信息
	os.RemoveAll(inst_dir)
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
	os.Mkdir(dirpath+"/"+INSTANCE_DIR_NAME, 0777)
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

func PackageExist(package_name string) (bool, error) {
	pathname := BPF_PACKAGE_HOME + package_name
	_, err := os.Stat(pathname)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 删除package
func PackageDelete(package_name string, force bool) {
	// 先判断是否存在
	exist, err := PackageExist(package_name)
	if err != nil {
		panic(err)
	}
	if !exist {
		err := fmt.Errorf("package %s not exist", package_name)
		if err != nil {
			panic(err)
		}
		return
	}
	// 存在则先检查与package相关的instance
	instance_path := BPF_PACKAGE_HOME + package_name + "/" + INSTANCE_DIR_NAME
	inst_dir, err := os.ReadDir(instance_path)
	if err != nil {
		panic(err)
	}
	// 若存在instance
	if len(inst_dir) != 0 {
		// 非强制删除
		if !force {
			err := fmt.Errorf("there are instances that depend on package: %s", package_name)
			if err != nil {
				panic(err)
			}
			return
		}
		// 强制删除instance
		for _, inst_file := range inst_dir {
			InstDelete(inst_file.Name())
		}
	}
	// 删除package的文件夹
	os.RemoveAll(BPF_PACKAGE_HOME + package_name)
}

// 列出当前所有的package
func PackageList() []string {
	var packages []string
	files, err := os.ReadDir(BPF_PACKAGE_HOME)
	if err != nil {
		panic(err)
	}
	for _, package_file := range files {
		packages = append(packages, package_file.Name())
	}
	return packages
}

// 读取package的信息
func PackageRead(package_name string, package_info *PackageInfo) {
	package_info.package_name = package_name
	// 读取src list(file)
	src_path := BPF_PACKAGE_HOME + package_name + "/" + SRC_FILE_NAME
	src_content, err := os.ReadFile(src_path)
	if err != nil {
		panic(err)
	}
	var src_str string = string(src_content)
	package_info.src_list = strings.Split(src_str, INFO_SEPARATOR)[1:]
	// 读取inst list
	var instances []string
	inst_path := BPF_PACKAGE_HOME + package_name + "/" + INSTANCE_DIR_NAME
	insts, err := os.ReadDir(inst_path)
	if err != nil {
		panic(err)
	}
	for _, inst := range insts {
		instances = append(instances, inst.Name())
	}
	package_info.inst_list = instances
	// 读取data size
	var size int64 = 0
	data_dir := BPF_PACKAGE_HOME + package_name + "/" + DATA_DIR_NAME
	data, err := os.ReadDir(data_dir)
	if err != nil {
		panic(err)
	}
	for _, item := range data {
		state, err := os.Stat(data_dir + "/" + item.Name())
		if err != nil {
			panic(err)
		}
		size += state.Size()
	}
	package_info.size = size
}

func SrcRead(package_name string, src_name string) string {
	src_path := BPF_PACKAGE_HOME + package_name + "/" + DATA_DIR_NAME + "/" + src_name
	src_content, err := os.ReadFile(src_path)
	if err != nil {
		panic(err)
	}
	return string(src_content)
}
