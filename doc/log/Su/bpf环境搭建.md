# bpf环境搭建

## 1、bpftool

直接apt的方式

```bash
$ apt-get install linux-tools-common
```

通过kernel source生成

```bash
# 安装依赖
$ apt-get install libelf-dev libbfd-dev libcap-dev clang make gcc
# 创建文件夹存放git仓库
$ mkdir git_repo
$ cd git_repo
# 拉取源码
$ git clone -b v5.11 --depth 1 https://github.com/torvalds/linux
$ cd linux/tools/bpf/bpftool/
# 安装
$ make && make install
```

## 2、libbpf库

安装libbpf依赖

```bash
# for ubuntu
$ apt install clang llvm libelf-dev iproute2
# test clang
$ clang -v
# test llvm
$ llc --version
# test iproute2
$ ip link
```

直接apt的方式（有点不靠谱）

```bash
$ apt-get install libbpf-dev
```

通过libbpf项目安装

```bash
$ cd git_repo
$ git clone https://github.com/libbpf/libbpf.git
$ cd libbpf/src
$ make && make install
```

可能会报错需要安装`pkg_config`：

```bash
$ apt-get install pkg_config
```

