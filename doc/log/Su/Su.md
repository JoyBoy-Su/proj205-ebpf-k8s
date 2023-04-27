# 工作记录 —— 苏佳迪

（本文档记录从2023-04-24开始的每日工作）

### 2023-04-24

阅读`linux6.0`源码与blog，整理eBPF的接口笔记，见[eBPF程序接口](../../note/Su/eBPF.md#2二、eBPF程序接口)；

使用`bpf()`系统调用书写简单的程序，见目录[bpf_syscall](../../note/Su/test/bpf_syscall)；

通过`libbpf`接口书写程序，但遇到了`__u32 type not defined`的bug；

### 2023-04-25

通过`libbpf`书写user和kernel的程序，监控`openat()`系统调用，并分别用gcc和clang进行编译运行成功；

通过libbpf库书写各种类型和功能的eBPF程序，见目录[libbpf](../../note/Su/test/libbpf)；

在通过`<bpf_tracing.h>`中的宏（如）运行时，编译会报错`Must specify a BPF target arch via __TARGET_ARCH_xxx`；

在`opensnoop`运行时遇到了`BTF is required, but is missing or corrupted.`的运行时异常；

### 2023-04-26

解决docker的网络问题：运行`~/fudan_net_auth.sh`脚本进行网络认证即可；

解决`__u32 type not defined`的bug，通过生成`vmlinux.h`头文件，解决过程见[解决__u32 type not defined](./解决__u32 type not defined.md)；