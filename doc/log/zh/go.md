# go环境配置

## 源代码安装

根据x86-64或者ARM64选择对应版本

https://golang.google.cn/dl/

## ****解压安装包****

```bash
cd /usr/local/src

tar -xzf go1.19.linux-amd64.tar.gz
```

## 环境变量

```bash
export PATH=$PATH:/usr/local/src/go/bin
source /etc/profile
```

## 设置GOPROXY

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

1. 系统变量 `GOROOT` :GO安装的根目录
2. 用户变量 `GOPATH`:用来设置工作目录，即编写代码的地方。包也都是从 `GOPATH`设置的路径中寻找。
3. 系统变量 `PATH` :各个操作系统都存在的环境变量，用于指定系统可执行命令的默认查找路径
