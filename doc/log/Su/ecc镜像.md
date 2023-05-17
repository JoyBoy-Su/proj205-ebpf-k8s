# ecc Dockerfile

## �������������������澵��

��ʹ��[compiler-dockerfile](https://github.com/eunomia-bpf/eunomia-bpf/blob/master/compiler/dockerfile)�������񣬳ɹ����������¾���

```bash
$ docker build -t ecc:1.0.1 .					# ��eunomia-bpfĿ¼�¹���
[+] Building 1052.9s (14/14) FINISHED
$ docker images
REPOSITORY   TAG       IMAGE ID       CREATED             SIZE
ecc          1.0.1     f53efee0dd9f   6 minutes ago       3.11GB
```

������������뵱ǰĿ¼�µ�`opensnoop.bpf.c`��

```bash
$ docker run -it -v `pwd`/:/src/ --name ecc-01 ecc:1.0.1
ls: cannot access '/src/*.h': No such file or directory
export PATH=ATH:~/.eunomia/bin
make: workspace/bin/ecc: No such file or directory			# Ŀ¼����
make: *** [Makefile:101: build] Error 127
```

���������ϱ���ԭ������������������ָ���ǣ�

```dockerfile
WORKDIR /usr/local/src/compiler

ENTRYPOINT ["make"]
CMD ["build"]
```

Ҳ���ǻ�ִ��make build������eunomia-bpf��compiler�¶��壩��

```makefile
build:
	export PATH=$PATH:~/.eunomia/bin
	$(Q)workspace/bin/ecc $(shell ls $(SOURCE_DIR)*.bpf.c) $(shell ls -h1 $(SOURCE_DIR)*.h | grep -v .*\.bpf\.h)
```

���Է���build�в�����ͨ��`ecc-rs`ָ��ֱ�ӱ���ָ�����ļ���һ������Ҫ����Ŀ¼����һ�����Ǽ�ʹ�Ѿ���װ�����ȴ��������eunomia-bpf��makefile�����ǿ����޸�ԭ�е�dockerfile��ǰ��Ĵ�ƪ��Ϊ�����������������`ecc-rs`���ⲿ�ֲ���Ҫ�޸ġ��޸�make�����������ָ������޸�Ϊ��entry pointִ��`ecc-rs`ָ�������docker run���룬��[Dockerfile](../../containers/Dockerfile.ecc)��

```dockerfile
FROM ubuntu:22.04						# base

WORKDIR /root/eunomia					# �趨������cwdΪ/root/eunomia
COPY . /root/eunomia					# ����ǰĿ¼����eunomia-bpf�ֿ��Ŀ¼�������ݸ��ƣ��Թ���������ecc����

# =============== ����Ϊ�������ù��� ===============

RUN apt-get update -y && \
    apt-get install -y --no-install-recommends \
        libelf1 libelf-dev zlib1g-dev libclang-13-dev \
        make wget curl python2 clang llvm pkg-config build-essential git && \
    apt-get install -y --no-install-recommends ca-certificates	&& \
	update-ca-certificates	&& \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*			# apt�������

RUN wget --progress=dot:giga --no-check-certificate \
        https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-17/wasi-sdk-17.0-linux.tar.gz && \
	tar -zxf wasi-sdk-17.0-linux.tar.gz && \
    rm wasi-sdk-17.0-linux.tar.gz   && \
	mkdir -p /opt/wasi-sdk/ && \
    mv wasi-sdk-17.0/* /opt/wasi-sdk/	# ����make��Ҫ����Դ

RUN cp /usr/bin/python2 /usr/bin/python

RUN wget -nv -O - https://sh.rustup.rs | sh -s -- -y	# make��Ҫ�Ľű�

ENV PATH="/root/.cargo/bin:${PATH}"
ARG CARGO_REGISTRIES_CRATES_IO_PROTOCOL=sparse

# =============== ����Ϊ����Դ�����ecc�Ĺ��� ===============

RUN make ecc    && \
    rm -rf /root/.eunomia && cp -r compiler/workspace /root/.eunomia    && \
    cd compiler/cmd && cargo clean

# =============== ����Ϊ����·����������ڵĹ��� ===============

ENV PATH="/root/.eunomia/bin:${PATH}"				# ���ecc��path

WORKDIR /code										# ���ù���Ŀ¼Ϊcode������ʱ��bpfӳ�䵽���Ŀ¼����

ENTRYPOINT ["ecc-rs"]								# ִ��ecc-rsָ�������docker run����

```

������������ɾ��񹹽�������һ�����������ض�Ӧbpf����docker�ڣ���ִ�в��ԣ�

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

�ɹ�����package.json������ecli���Կ�����ȷִ�У����˱���ʹ�õľ��񹹽���ɡ�

### image����

��ȡeunomia-bpf�ֿ⣨������ģ�飩��

```bash
git clone --recursive https://github.com/eunomia-bpf/eunomia-bpf.git
```

����**�ֿ��Ŀ¼**��һ��Ҫ�����Ŀ¼�¹��������������ļ��������⣩��������Dockerfile��

```bash
$ cd eunomia-bpf
$ vi Dockerfile
```

��Dockerfile��������д��ɣ�Ȼ���ڸ�Ŀ¼��ִ��`docker build`ָ���dockerfile���߼�����image��

```bash
$ docker build -t ecc-rs:latest .
```

### imageִ��

����Ƿ��Ѿ�����ecc-rs����

```bash
$ docker images | grep ecc-rs
REPOSITORY   TAG       IMAGE ID       CREATED          SIZE
ecc-rs       latest    6312288dcc53   26 minutes ago   3.11GB
```

��ecc-rsΪimage����һ������������ʱ��Ҫ��

1. ����Ŀ¼����bpf��������Ŀ¼���ص�������`/code`Ŀ¼
2. ָ��bpf�ļ�������ͷ�ļ�������Ҫ�����bpf�ļ���

ָ���ʽ���£�

```bash
$ docker run -it -v /your_bpf_path/:/code/ --name your_container_name ecc-rs:version bpf.c bpf.h
```

���磬��������ǰĿ¼���������ļ�`opensnoop.bpf.c`��`opensnoop.h`����Ҫ��������ִ������ָ�

```bash
$ docker run -it -v `pwd`/:/code/ --name ecc-rs ecc-rs:latest opensnoop.bpf.c opensnoop.h
```

�������������package.json����ڵ�ǰĿ¼�£���ɱ���Ĺ��̣�Ϊ��һ���ַ���׼����

ע��[ecc-rs����](https://hub.docker.com/r/jiadisu/ecc-rs)�Ѿ����͵�docker hub���ˣ�����ͨ������ָ����ȡ��

```bash
$ docker pull jiadisu/ecc-rs
```

## ������ecc��ecc������С����������澵��

����������������Դ��������ecc�����������ܶ����ʹ�õ����в���Ҫ������������ڿ���ȥ����Щ��������

�߼����ǣ������õľ����ٴ�Դ�������õ�ecc�����ǴӶ�Ӧ������ֱ�Ӹ���ecc�������У�����ȡ������Ҫ����С���������޸ĺ��dockerfile������£�

```dockerfile
FROM ubuntu:22.04

COPY . /root/.eunomia/bin

ENV PATH="/root/.eunomia/bin:${PATH}"

WORKDIR /code

ENTRYPOINT ["ecc-rs"]

```

������ҪΪ�����ecc-rs����ʱ��Ҫ�õ������п⡣
