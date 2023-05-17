# ecc Dockerfile

## 包含所有依赖的完整版镜像

先使用[compiler-dockerfile](https://github.com/eunomia-bpf/eunomia-bpf/blob/master/compiler/dockerfile)构建镜像，成功构建出如下镜像：

```bash
$ docker build -t ecc:1.0.1 .					# 在eunomia-bpf目录下构建
[+] Building 1052.9s (14/14) FINISHED
$ docker images
REPOSITORY   TAG       IMAGE ID       CREATED             SIZE
ecc          1.0.1     f53efee0dd9f   6 minutes ago       3.11GB
```

若运行让其编译当前目录下的`opensnoop.bpf.c`：

```bash
$ docker run -it -v `pwd`/:/src/ --name ecc-01 ecc:1.0.1
ls: cannot access '/src/*.h': No such file or directory
export PATH=ATH:~/.eunomia/bin
make: workspace/bin/ecc: No such file or directory			# 目录问题
make: *** [Makefile:101: build] Error 127
```

出现了如上报错，原因是容器最终启动的指令是：

```dockerfile
WORKDIR /usr/local/src/compiler

ENTRYPOINT ["make"]
CMD ["build"]
```

也就是会执行make build操作（eunomia-bpf的compiler下定义）：

```makefile
build:
	export PATH=$PATH:~/.eunomia/bin
	$(Q)workspace/bin/ecc $(shell ls $(SOURCE_DIR)*.bpf.c) $(shell ls -h1 $(SOURCE_DIR)*.h | grep -v .*\.bpf\.h)
```

可以发现build中并不是通过`ecc-rs`指令直接编译指定的文件，一方面需要配置目录，另一方面是即使已经安装完成了却还依赖与eunomia-bpf的makefile，于是考虑修改原有的dockerfile，前面的大篇幅为配置依赖并编译产生`ecc-rs`，这部分不需要修改。修改make后的容器启动指令，将其修改为以entry point执行`ecc-rs`指令，参数由docker run传入，见[Dockerfile](../../containers/Dockerfile.ecc)：

```dockerfile
FROM ubuntu:22.04						# base

WORKDIR /root/eunomia					# 设定后续的cwd为/root/eunomia
COPY . /root/eunomia					# 将当前目录（即eunomia-bpf仓库根目录）的内容复制，以供后续进行ecc编译

# =============== 以下为依赖配置过程 ===============

RUN apt-get update -y && \
    apt-get install -y --no-install-recommends \
        libelf1 libelf-dev zlib1g-dev libclang-13-dev \
        make wget curl python2 clang llvm pkg-config build-essential git && \
    apt-get install -y --no-install-recommends ca-certificates	&& \
	update-ca-certificates	&& \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*			# apt添加依赖

RUN wget --progress=dot:giga --no-check-certificate \
        https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-17/wasi-sdk-17.0-linux.tar.gz && \
	tar -zxf wasi-sdk-17.0-linux.tar.gz && \
    rm wasi-sdk-17.0-linux.tar.gz   && \
	mkdir -p /opt/wasi-sdk/ && \
    mv wasi-sdk-17.0/* /opt/wasi-sdk/	# 下载make需要的资源

RUN cp /usr/bin/python2 /usr/bin/python

RUN wget -nv -O - https://sh.rustup.rs | sh -s -- -y	# make需要的脚本

ENV PATH="/root/.cargo/bin:${PATH}"
ARG CARGO_REGISTRIES_CRATES_IO_PROTOCOL=sparse

# =============== 以下为编译源码产生ecc的过程 ===============

RUN make ecc    && \
    rm -rf /root/.eunomia && cp -r compiler/workspace /root/.eunomia    && \
    cd compiler/cmd && cargo clean

# =============== 以下为配置路径与容器入口的过程 ===============

ENV PATH="/root/.eunomia/bin:${PATH}"				# 添加ecc到path

WORKDIR /code										# 设置工作目录为code，运行时将bpf映射到这个目录即可

ENTRYPOINT ["ecc-rs"]								# 执行ecc-rs指令，参数由docker run传递

```

按上述过程完成镜像构建后，启动一个容器，挂载对应bpf程序到docker内，并执行测试：

```bash
$ docker run -it -v `pwd`/:/code --name ecc-rs ecc-rs:1.0 opensnoop.bpf.c opensnoop.h
Compiling bpf object...
warning: text is not json: Process ID to trace use it as a string
warning: text is not json: Thread ID to trace use it as a string
warning: text is not json: User ID to trace use it as a string
warning: text is not json: trace only failed events use it as a string
warning: text is not json: Trace open family syscalls. use it as a string
Generating export types...
Packing ebpf object and config into /code/package.json...
$ ls -l
total 36
-rw-rw-r-- 1 ubuntu ubuntu  3702 May 17 05:09 opensnoop.bpf.c
-rw-r--r-- 1 root   root   12912 May 17 17:49 opensnoop.bpf.o
-rw-rw-r-- 1 ubuntu ubuntu   427 May 17 05:10 opensnoop.h
-rw-r--r-- 1 root   root    1525 May 17 17:49 opensnoop.skel.json
-rw-r--r-- 1 root   root    5934 May 17 17:49 package.json
```

成功产生package.json，经过ecli测试可以正确执行，至此编译使用的镜像构建完成。

### image构建

拉取eunomia-bpf仓库（包含子模块）：

```bash
git clone --recursive https://github.com/eunomia-bpf/eunomia-bpf.git
```

进入**仓库根目录**（一定要在这个目录下构建，否则镜像中文件会有问题），并创建Dockerfile：

```bash
$ cd eunomia-bpf
$ vi Dockerfile
```

将Dockerfile的内容填写完成，然后在根目录下执行`docker build`指令，按dockerfile的逻辑构建image：

```bash
$ docker build -t ecc-rs:latest .
```

### image执行

检查是否已经存在ecc-rs镜像：

```bash
$ docker images | grep ecc-rs
REPOSITORY   TAG       IMAGE ID       CREATED          SIZE
ecc-rs       latest    6312288dcc53   26 minutes ago   3.11GB
```

以ecc-rs为image启动一个容器，启动时需要：

1. 挂载目录：将bpf程序所在目录挂载到容器的`/code`目录
2. 指定bpf文件（包括头文件）：即要编译的bpf文件名

指令格式如下：

```bash
$ docker run -it -v /your_bpf_path/:/code/ --name your_container_name ecc-rs:version bpf.c bpf.h
```

例如，宿主机当前目录下有两个文件`opensnoop.bpf.c`与`opensnoop.h`，想要编译它就执行如下指令：

```bash
$ docker run -it -v `pwd`/:/code/ --name ecc-rs ecc-rs:latest opensnoop.bpf.c opensnoop.h
```

容器便会编译产生package.json存放在当前目录下，完成编译的过程，为下一步分发做准备。

注：[ecc-rs镜像](https://hub.docker.com/r/jiadisu/ecc-rs)已经推送到docker hub上了，可以通过如下指令拉取：

```bash
$ docker pull jiadisu/ecc-rs
```

## 仅包含ecc与ecc运行最小依赖的迷你版镜像

上述过程由于是用源码编译产生ecc，镜像会包含很多编译使用但运行不需要的软件包，现在考虑去除这些依赖包。

逻辑就是，我们用的镜像不再从源代码编译得到ecc，而是从对应机器上直接复制ecc到容器中，并拉取运行需要的最小依赖包，修改后的dockerfile框架如下：

```dockerfile
FROM ubuntu:22.04

COPY . /root/.eunomia/bin

ENV PATH="/root/.eunomia/bin:${PATH}"

WORKDIR /code

ENTRYPOINT ["ecc-rs"]

```

下面需要为其添加ecc-rs运行时需要用到的所有库。
