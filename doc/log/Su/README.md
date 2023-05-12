# 工作记录 —— 苏佳迪

（本文档记录从2023-04-24开始的每日工作）

### 2023-04-24

阅读`linux6.0`源码与blog，整理eBPF的接口[笔记](../../note/Su/eBPF.md#bpf-syscall)；

使用`bpf()`系统调用书写简单的程序，见[目录](../../note/Su/test/bpf_syscall)；

通过`libbpf`接口书写程序，但遇到了`__u32 type not defined`的bug；

### 2023-04-25

通过`libbpf`书写user和kernel的程序，监控`openat()`系统调用，并分别用gcc和clang进行编译运行成功；

通过libbpf库书写各种类型和功能的eBPF程序，见[目录](../../note/Su/test/libbpf)；

在通过`<bpf_tracing.h>`中的宏（如）运行时，编译会报错`Must specify a BPF target arch via __TARGET_ARCH_xxx`；

在`opensnoop`运行时遇到了`BTF is required, but is missing or corrupted.`的运行时异常；

### 2023-04-26

解决docker的网络问题：运行`~/fudan_net_auth.sh`脚本进行网络认证即可；

解决`__u32 type not defined`的bug，通过生成`vmlinux.h`头文件，[解决过程](./问题与解决.md#error:unknown_type_name_'__u64')；

### 2023-04-27

在os课程的docker中安装docker并运行，过程见[详情](./k8s_docker安装搭建.md)，为下一步搭建k8s集群做准备；

在docker中安装libbpf和bpftool工具，过程见[详情](./bpf环境搭建.md)；

通过建立软链接解决fatal error: 'asm/types.h' file not found的问题，过程见[详情](./问题与解决.md#fatal_error:_'asm/types.h'_file_not_found)

### 2023-04-28

阅读[使用libbpf-bootstrap构建BPF程序](https://forsworns.github.io/zh/blogs/20210627/)，了解了在原生的libbpf库下，怎样通过clang的方式编译构建bpf程序，见[详情](./libbpf编译构建BPF过程.md)；

### 2023-05-07

查看[kubectl-trace](https://github.com/iovisor/kubectl-trace)项目，该项目是作为一个k8s的插件存在，可以在k8s集群上方便地编排部署bpftrace程序；

观察其架构图，初步构想我们项目的结构：master进行编译并创建pod，挂载config-map，将bpf分发到不同的node上执行；

### 2023-05-08

解决了`Must specify a BPF target arch via __TARGET_ARCH_xxx`的问题，通过修改clang的编译指令，过程见[详情](./问题与解决.md#fatal_error:_'asm/types.h'_file_not_found)；

由于os课程的docker的内核问题（并没有开启BTF支持），尝试重新编译内核，过程见[详情](./重新编译linux内核.md)，最终在重新启动docker的时候失败；

### 2023-05-09

讨论整理过去两周的工作，组内讨论得出项目的大体思路和存在的疑惑，内容见[详情](../../会议记录.md#2023-05-09 20 : 00)；

### 2023-05-12

在自己的云服务器（huawei cloud）上配置eunomia-bpf的框架环境；

配置好环境后，用eunomia-bpf框架运行`opensnoop.bpf.c`程序，也通过观察该过程的结果解决了之前的一些疑问；

搭建与使用的过程见[详情](./搭建与使用eunomia-bpf框架环境.md)；