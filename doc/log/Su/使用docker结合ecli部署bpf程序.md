# 使用docker结合ecli部署bpf程序

服务器环境：

- Ubuntu 20.04.4 LTS
- Linux hecs-67681 5.4.0-100-generic #113-Ubuntu SMP Thu Feb 3 18:43:29 UTC 2022 x86_64 x86_64 x86_64 GNU/Linux

## 试错...

### 1、查看已有镜像

```bash
$ docker images
REPOSITORY                 TAG            IMAGE ID       CREATED         SIZE
maven_tomcat               latest         6e062b2caab3   6 weeks ago     540MB
tomcat                     latest         608294908754   6 weeks ago     475MB
node                       latest         0e0ab07dbedd   7 weeks ago     999MB
rofrano/vagrant-provider   ubuntu-22.04   0ebddde03671   16 months ago   242MB
hello-world                latest         feb5d9fea6a5   19 months ago   13.3kB
```

### 2、拉取ubuntu镜像（已存在不需要再拉取）

```bash
$ docker pull ubuntu:20.04
20.04: Pulling from library/ubuntu
ca1778b69356: Pull complete 
Digest: sha256:db8bf6f4fb351aa7a26e27ba2686cf35a6a409f65603e59d4c203e58387dc6b3
Status: Downloaded newer image for ubuntu:20.04
docker.io/library/ubuntu:20.04
$ docker images
REPOSITORY                 TAG            IMAGE ID       CREATED         SIZE
ubuntu                     20.04          88bd68917189   4 weeks ago     72.8MB
maven_tomcat               latest         6e062b2caab3   6 weeks ago     540MB
tomcat                     latest         608294908754   6 weeks ago     475MB
node                       latest         0e0ab07dbedd   7 weeks ago     999MB
rofrano/vagrant-provider   ubuntu-22.04   0ebddde03671   16 months ago   242MB
hello-world                latest         feb5d9fea6a5   19 months ago   13.3kB
```

### 3、运行ubuntu容器

```bash
# 启动ubuntu docker
$ docker run --name ubuntu-test -it -v /usr:/usr -v ${PWD}:/eunomia -d -p 20143:22 ubuntu:20.04
c0daabd10f68eea7ea5dc49d9bcd6044a83542e0efc814e2d2cb0a3b9e23a7d2
# 查看是否启动成功
$ docker ps
CONTAINER ID   IMAGE          COMMAND       CREATED          STATUS          PORTS                                     NAMES
c0daabd10f68   ubuntu:20.04   "/bin/bash"   23 seconds ago   Up 22 seconds   0.0.0.0:20143->22/tcp, :::20143->22/tcp   ubuntu-test
```

上述指令有三个映射关系：

1. 物理机`/usr`目录映射到docker的`/usr`目录（环境）
2. 物理机`${PWD}`目录映射到docker的`/eunomia`目录下（ecc和ecli）
3. 物理机的20143端口映射到docker的22端口（ssh）

### 4、进入ubuntu容器

```bash
$ docker exec -it c0daabd10f68 /bin/bash
root@c0daabd10f68:/# 
# 检查上述文件的映射关系是否成功
# clang
$ ls -l /usr/lib | grep clang
drwxr-xr-x   4 root root   4096 May 12 14:01 clang
lrwxrwxrwx   1 root root     11 May 12 12:11 libclang-10.so.1 -> libclang.so
lrwxrwxrwx   1 root root     32 May 12 12:12 libclang-13.so.13 -> /usr/lib/llvm-10/lib/libclang.so
lrwxrwxrwx   1 root root     32 May 12 12:10 libclang.so -> /usr/lib/llvm-10/lib/libclang.so
$ clang -v
Ubuntu clang version 14.0.0-1ubuntu1
Target: x86_64-pc-linux-gnu
Thread model: posix
InstalledDir: /usr/bin
Found candidate GCC installation: /usr/bin/../lib/gcc/x86_64-linux-gnu/11
Found candidate GCC installation: /usr/bin/../lib/gcc/x86_64-linux-gnu/9
Selected GCC installation: /usr/bin/../lib/gcc/x86_64-linux-gnu/11
Candidate multilib: .;@m64
Selected multilib: .;@m64
# llvm
$ ls -l /usr/lib | grep llvm
lrwxrwxrwx   1 root root     32 May 12 12:12 libclang-13.so.13 -> /usr/lib/llvm-10/lib/libclang.so
lrwxrwxrwx   1 root root     32 May 12 12:10 libclang.so -> /usr/lib/llvm-10/lib/libclang.so
drwxr-xr-x   7 root root   4096 Apr 17 07:05 llvm-10
drwxr-xr-x   7 root root   4096 May 12 13:22 llvm-13
drwxr-xr-x   7 root root   4096 May 12 14:01 llvm-14
$ llvm-as --version
Ubuntu LLVM version 14.0.0
  
  Optimized build.
  Default target: x86_64-pc-linux-gnu
  Host CPU: cascadelake
# eunomia
$ ls -l /eunomia                
total 41236
-rwxr-xr-x 1 root root 22213696 May  5 11:16 ecc
-rwxr-xr-x 1 root root 20005104 May  5 11:21 ecli
$ cd /eunomia
$ ./ecli -h							# 测试./ecli -h指令
ecli subcommands, including run, push, pull, login, logout

Usage: ecli <COMMAND>

Commands:
  run     run ebpf program
  client  Client operations
  push    
  pull    pull oci image from registry
  login   login to oci registry
  logout  logout from registry
  help    Print this message or the help of the given subcommand(s)

Options:
  -h, --help  Print help
```

可以发现通过文件目录的挂载，内部的ubuntu已经具备了clang / llvm环境，且具有完备的ecc和ecli指令（在`/eunomia`目录下）。

### 5、共享package.json

在物理机上，将package.json复制到docker挂载的目录中，以供docker执行：

```bash
# node
$ mkdir /usr/local/temp						# 创建temp目录
$ cp package.json /usr/local/temp/			# 移动package.json到temp目录
$ ls -l /usr/local/temp/
total 8
-rw-r--r-- 1 root root 5954 May 13 15:21 package.json
$ docker exec -it c0daabd10f68 /bin/bash	# 进入docker
root@c0daabd10f68:/#

# docker
$ ls -l /usr/local/temp/					# 检验package.json是否已挂载
total 8
-rw-r--r-- 1 root root 5954 May 13 07:21 package.json
```

### 6、通过ecli执行bpf程序

```shell
# docker
$ cd /eunomia
$ ls -l      
total 41236
-rwxr-xr-x 1 root root 22213696 May  5 11:16 ecc
-rwxr-xr-x 1 root root 20005104 May  5 11:21 ecli
$ ./ecli run /usr/local/temp/package.json
INFO [faerie::elf] strtab: 0xe4c symtab 0xe88 relocs 0xed0 sh_offset 0xed0
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable pid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable tgid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable uid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] load runtime arg (user specified the value through cli, or predefined in the skeleton) for targ_failed: Bool(false), real_type=<INT> '_Bool' bits:8 off:0 enc:bool, btf_type=BtfVar { name: "targ_failed", type_id: 46, kind: GlobalAlloc }
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable __eunomia_dummy_event_ptr, use the default one in ELF
libbpf: Failed to bump RLIMIT_MEMLOCK (err = -1), you might need to do it explicitly!
libbpf: Error in bpf_object__probe_loading():Operation not permitted(1). Couldn't load trivial BPF program. Make sure your kernel supports BPF (CONFIG_BPF_SYSCALL=y) and/or that RLIMIT_MEMLOCK is set to big enough value.
libbpf: failed to load object 'opensnoop_bpf'
Error: Bpf("Failed to start polling: Bpf(\"Failed to load and attach: Failed to load bpf object\"), receiving on a closed channel")
```

出现上述报错

### 7、解决RLIMIT_MEMLOCK限制问题

报错内容包括`Failed to bump RLIMIT_MEMLOCK (err = -1)`，包括后面也都可以发现与`RLIMIT_MEMLOCK`相关，猜测是`RLIMIT_MEMLOCK`导致的资源限制。这是由于shell默认的资源不足造成的，尝试使用`ulimite`命令解决：

在docker中执行如下指令：

```bash
# 查看当前的默认分配的资源
$ ulimit -a
core file size          (blocks, -c) unlimited
data seg size           (kbytes, -d) unlimited
scheduling priority             (-e) 0
file size               (blocks, -f) unlimited
pending signals                 (-i) 7580
max locked memory       (kbytes, -l) 64					# << 报错的位置，locked memory只有64k
max memory size         (kbytes, -m) unlimited
open files                      (-n) 1048576
pipe size            (512 bytes, -p) 8
POSIX message queues     (bytes, -q) 819200
real-time priority              (-r) 0
stack size              (kbytes, -s) 8192
cpu time               (seconds, -t) unlimited
max user processes              (-u) unlimited
virtual memory          (kbytes, -v) unlimited
file locks                      (-x) unlimited
# 设置max locked memory为不受限
$ ulimit -l unlimited
bash: ulimit: max locked memory: cannot modify limit: Operation not permitted	# 出现了权限不足的问题
```

以上结果判断为：docker中无法修改ulimit。

参考[如何调整 docker 下 linux 的 ulimit 大小设置？](https://gorden5566.com/post/1089.html)按上面的步骤，重新启动一个容器，在`docker run`时添加参数`--ulimit memlock=xxx`设置`memlock`参数：

```bash
# node
$ docker run --name ubuntu-test -it -v /usr:/usr -v ${PWD}:/eunomia -d -p 20143:22 --ulimit memlock=8192 ubuntu:20.04		# memlock = 8129是说设置为了8192Bytes即8KB
6fe157343c7927a111806ed7adab4d52616a69b15377b0233d9ac82496841adf
$ docker ps
CONTAINER ID   IMAGE          COMMAND       CREATED          STATUS          PORTS                                     NAMES
6fe157343c79   ubuntu:20.04   "/bin/bash"   24 seconds ago   Up 23 seconds   0.0.0.0:20143->22/tcp, :::20143->22/tcp   ubuntu-test
$ docker exec -it 6fe157343c79 /bin/bash
root@6fe157343c79:/#

# docker
$ ulimit -a
core file size          (blocks, -c) unlimited
data seg size           (kbytes, -d) unlimited
scheduling priority             (-e) 0
file size               (blocks, -f) unlimited
pending signals                 (-i) 7580
max locked memory       (kbytes, -l) 8				# << 这里由刚才的ulimit指令设置为8192Bytes即8K
max memory size         (kbytes, -m) unlimited
open files                      (-n) 1048576
pipe size            (512 bytes, -p) 8
POSIX message queues     (bytes, -q) 819200
real-time priority              (-r) 0
stack size              (kbytes, -s) 8192
cpu time               (seconds, -t) unlimited
max user processes              (-u) unlimited
virtual memory          (kbytes, -v) unlimited
file locks                      (-x) unlimited
```

经上述测试，`--ulimit memlock=xxx`设置生效，则按照上述步骤，设置`memlock=-1`（即unlimited）：

```bash
# node
$ docker run --name ubuntu-test -it -v /usr:/usr -v ${PWD}:/eunomia -d -p 20143:22 --ulimit memlock=-1 ubuntu:20.04
d6ddad1b8651ceb04baf19219fca548fedd1fd989122f9455ef0c5bb3e6e9bca
$ docker ps
CONTAINER ID   IMAGE          COMMAND       CREATED          STATUS          PORTS                                     NAMES
d6ddad1b8651   ubuntu:20.04   "/bin/bash"   24 seconds ago   Up 23 seconds   0.0.0.0:20143->22/tcp, :::20143->22/tcp   ubuntu-test
$ docker exec -it d6ddad1b8651 /bin/bash
root@d6ddad1b8651:/# 

# docker
$ ulimit -a
core file size          (blocks, -c) unlimited
data seg size           (kbytes, -d) unlimited
scheduling priority             (-e) 0
file size               (blocks, -f) unlimited
pending signals                 (-i) 7580
max locked memory       (kbytes, -l) unlimited				# << 设置为了1024K
max memory size         (kbytes, -m) unlimited
open files                      (-n) 1048576
pipe size            (512 bytes, -p) 8
POSIX message queues     (bytes, -q) 819200
real-time priority              (-r) 0
stack size              (kbytes, -s) 8192
cpu time               (seconds, -t) unlimited
max user processes              (-u) unlimited
virtual memory          (kbytes, -v) unlimited
file locks                      (-x) unlimited
# 检查资源
$ clang -v
Ubuntu clang version 14.0.0-1ubuntu1
Target: x86_64-pc-linux-gnu
Thread model: posix
InstalledDir: /usr/bin
Found candidate GCC installation: /usr/bin/../lib/gcc/x86_64-linux-gnu/11
Found candidate GCC installation: /usr/bin/../lib/gcc/x86_64-linux-gnu/9
Selected GCC installation: /usr/bin/../lib/gcc/x86_64-linux-gnu/11
Candidate multilib: .;@m64
Selected multilib: .;@m64
$ llvm-as --version
Ubuntu LLVM version 14.0.0
  
  Optimized build.
  Default target: x86_64-pc-linux-gnu
  Host CPU: cascadelake
$ ls -l /usr/local/temp 
total 8
-rw-r--r-- 1 root root 5954 May 13 07:21 package.json
$ ls -l /eunomia 
total 41236
-rwxr-xr-x 1 root root 22213696 May  5 11:16 ecc
-rwxr-xr-x 1 root root 20005104 May  5 11:21 ecli
```

尝试使用`ecli`执行`package.json`：

```bash
$ cd /eumomia
$ ./ecli run /usr/local/temp/package.json
INFO [faerie::elf] strtab: 0xe4c symtab 0xe88 relocs 0xed0 sh_offset 0xed0
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable pid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable tgid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable uid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] load runtime arg (user specified the value through cli, or predefined in the skeleton) for targ_failed: Bool(false), real_type=<INT> '_Bool' bits:8 off:0 enc:bool, btf_type=BtfVar { name: "targ_failed", type_id: 46, kind: GlobalAlloc }
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable __eunomia_dummy_event_ptr, use the default one in ELF
libbpf: Error in bpf_object__probe_loading():Operation not permitted(1). Couldn't load trivial BPF program. Make sure your kernel supports BPF (CONFIG_BPF_SYSCALL=y) and/or that RLIMIT_MEMLOCK is set to big enough value.
libbpf: failed to load object 'opensnoop_bpf'
Error: Bpf("Failed to start polling: Bpf(\"Failed to load and attach: Failed to load bpf object\"), receiving on a closed channel")
```

阿里源镜像源：

```tex
deb http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse
```

### 8、最终docker run

```bash
$ docker run --name ubuntu-test -it -v ${PWD}:/eunomia -v /usr:/usr -v /boot:/boot -v /sys:/sys --privileged -d -p 20143:22 --ulimit memlock=-1 ubuntu:20.04
```

启动后进入docker，并启动`ecli run`执行pacakge.json，在node上通过`bpftool prog list`发现成功挂载bpf程序。

## 正确方式

需要的内容：

- docker image：在ubuntu的基础上添加需要的依赖；
- ecli：用来执行package.json；
- sources.list：镜像源；
- package.json：待执行的bpf；

### 1、构建docker image

通过docker file构建出包含依赖的docker image：

```dockerfile
FROM ubuntu:20.04						# 在ubuntu:20.04的基础上构建

ENV UBUNTU_SOURCE /etc/apt				# 设置环境变量UBUNTU_SOURCE

COPY ./ /root							# 将宿主机（执行docker命令的机器）当前工作目录，复制到docker /root下

WORKDIR /root							# 设置docker的工作目录为/root

ADD sources.list $UBUNTU_SOURCE/		# 为docker的apt添加镜像源

RUN apt-get update && \					# 在docker build时执行apt指令，安装依赖
    apt-get -y install gcc libelf-dev

#CMD ./ecli run /root/my/package.json
CMD ["/bin/bash"]						# 在docker run时执行"/bin/bash"指令，切换到docker的终端

```

sources.list配置镜像源：

```tex
deb http://mirrors.aliyun.com/ubuntu/ jammy main restricted universe multiverse
#deb-src http://mirrors.aliyun.com/ubuntu/ jammy main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ jammy-security main restricted universe multiverse
#deb-src http://mirrors.aliyun.com/ubuntu/ jammy-security main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ jammy-updates main restricted universe multiverse
#deb-src http://mirrors.aliyun.com/ubuntu/ jammy-updates main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ jammy-proposed main restricted universe multiverse
#deb-src http://mirrors.aliyun.com/ubuntu/ jammy-proposed main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ jammy-backports main restricted universe multiverse
#deb-src http://mirrors.aliyun.com/ubuntu/ jammy-backports main restricted universe multiverse
```

将如上将如上dockerfile与ecli和sources.list放在同一目录，并在该目录下执行docker build构建docker镜像：

```bash
$ docker build -t ecli:1.0.1 .	# 指定image名为ecli 版本为1.0.1 构建目录.（即当前目录）
```

### 2、启动docker

构建完成image后，查看已有镜像：

```bash
$ docker images
REPOSITORY                 TAG            IMAGE ID       CREATED          SIZE
ecli                       1.0.1          1e89a8629648   20 minutes ago   327MB
```

根据ecli image启动container：

```bash
$ docker run -it -v /usr/src:/usr/src:ro \				# 映射node:/usr/src到container:/usr/src
       -v /lib/modules/:/lib/modules:ro \				# 映射node:/lib/modules到container:/lib/modules
       -v /sys/kernel/debug/:/sys/kernel/debug:rw \		# 映射node:/sys/kernel/debug到container:/sys/kernel/debug
       -v /root/mydev:/root/mydev \						# 映射node:/root/mydev到/root/mydev
       --net=host --pid=host --privileged \				# 设置网络与权限
       --name ecli ecli:1.0.1							# name image
root@hecs-67681:~# #此时已经进入了docker中

# 另起一个shell查看当前的container
$ docker ps
CONTAINER ID   IMAGE        COMMAND       CREATED          STATUS          PORTS     NAMES
14630ab79551   ecli:1.0.1   "/bin/bash"   49 seconds ago   Up 48 seconds             ecli	# ecli
```

### 3、执行ecli指令

启动并进入docker后，执行`/root`下的ecli完成bpf的执行：

```bash
$ ls -l
total 19552
-rw-r--r-- 1 root root      230 May 13 12:54 Dockerfile
-rwxr-xr-x 1 root root 20005104 May 13 12:54 ecli
drwxr-xr-x 6 root root     4096 Apr 17 08:55 mydev
-rw-r--r-- 1 root root      897 May 13 12:55 sources.list
# 执行ecli
$ ./ecli run mydev/ebpf/eunomia/package.json 
INFO [faerie::elf] strtab: 0xe4c symtab 0xe88 relocs 0xed0 sh_offset 0xed0
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable pid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable tgid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable uid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] load runtime arg (user specified the value through cli, or predefined in the skeleton) for targ_failed: Bool(false), real_type=<INT> '_Bool' bits:8 off:0 enc:bool, btf_type=BtfVar { name: "targ_failed", type_id: 46, kind: GlobalAlloc }
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable __eunomia_dummy_event_ptr, use the default one in ELF
INFO [bpf_loader_lib::skeleton::poller] Running ebpf program...
<1683984446> <PLAIN> TIME     TS     PID    UID    RET    FLAGS  COMM   FNAME 

# 此时在物理机上查看挂载的程序，如下：
$ bpftool prog list
134: tracepoint  name tracepoint__sys  tag 07014be5359438f8  gpl
        loaded_at 2023-05-13T21:27:26+0800  uid 0
        xlated 288B  jited 167B  memlock 4096B  map_ids 75,72
        btf_id 52
136: tracepoint  name tracepoint__sys  tag 8ee3432dcd98ffc3  gpl
        loaded_at 2023-05-13T21:27:26+0800  uid 0
        xlated 288B  jited 167B  memlock 4096B  map_ids 75,72
        btf_id 52
137: tracepoint  name tracepoint__sys  tag 541339de114a40e6  gpl
        loaded_at 2023-05-13T21:27:26+0800  uid 0
        xlated 696B  jited 477B  memlock 4096B  map_ids 72,75,73
        btf_id 52
138: tracepoint  name tracepoint__sys  tag 541339de114a40e6  gpl
        loaded_at 2023-05-13T21:27:26+0800  uid 0
        xlated 696B  jited 477B  memlock 4096B  map_ids 72,75,73
        btf_id 52
# 同时docker中会有显示对应监测到的系统调用
<1683984446> <PLAIN> TIME     TS     PID    UID    RET    FLAGS  COMM   FNAME  
<1683984463> <PLAIN> 13:27:43  0     39004  0      3      524288 "bpftool" "/etc/ld.so.cache"
<1683984463> <PLAIN> 13:27:43  0     39004  0      3      524288 "bpftool" "/lib/x86_64-linux-gnu/libtinfo.so.6"
<1683984463> <PLAIN> 13:27:43  0     39004  0      3      524288 "bpftool" "/lib/x86_64-linux-gnu/libdl.so.2"
<1683984463> <PLAIN> 13:27:43  0     39004  0      3      524288 "bpftool" "/lib/x86_64-linux-gnu/libc.so.6"
<1683984463> <PLAIN> 13:27:43  0     39004  0      3      2050   "bpftool" "/dev/tty"
<1683984463> <PLAIN> 13:27:43  0     39004  0      3      524288 "bpftool" "/usr/lib/locale/locale-archive"
<1683984463> <PLAIN> 13:27:43  0     38999  0      -2     524288 "ecli" "/etc/localtime"
<1683984463> <PLAIN> 13:27:43  0     38999  0      -2     524288 "ecli" "/etc/timezone"
```

至此通过docker封装环境依赖，并使用ecli安装ebpf程序成功。