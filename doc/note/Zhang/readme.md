# Readme

`doc/note/xxx`目录中存放队员xxx的笔记，

该Readme文档的存在是为了保证该目录并非空目录，从而可以由git提交到托管网站上，创建好自己记录笔记的文档（保证该目录不是空目录）后可删除该readme。

[基于 ebpf 和 vxlan 实现一个 k8s 网络插件（一）](https://zhuanlan.zhihu.com/p/565254116)

[k8s管理员要懂eBPF](https://www.jianshu.com/p/3e21bb174445)

[calico](https://www.tigera.io/blog/introducing-the-calico-ebpf-dataplane/)

[基于 eBPF 的 Kubernetes 问题排查全景图发布](https://developer.aliyun.com/article/879258)

![](https://pic2.zhimg.com/80/v2-c1ea9720cb0892a5b02ce3e91556d94d_1440w.webp)

1. [在linux上安装kubectl](http://kubernetes.p2hp.com/docs/tasks/tools/install-kubectl-linux.html)

   1. ```bash
      ~#：uname -i    
      aarch64
      //平台不一样，需要将amd64换成arm64
      ```

   2. curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubectl"

   3. curl -LO "https://dl.k8s.io/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubectl.sha256"

   4. 安装补全工具失败 type _init_completion:放弃

   > 下载 Google Cloud 公开签名秘钥：访问外网失败

2. 安装 kubeadm、kubelet 

   ```bash
   curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubeadm"
   curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubeadm.sha256"
   echo "$(cat kubeadm.sha256)  kubeadm" | sha256sum --check
    install -o root -g root -m 0755 kubeadm /usr/local/bin/kubeadm
   ```

   ```bash
   curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubelet"
   curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubelet.sha256"
   echo "$(cat kubelet.sha256)  kubelet" | sha256sum --check
    install -o root -g root -m 0755 kubelet /usr/local/bin/kubelet
   ```

   

3. apt安装iptables  apparmor docker-ce-rootless-extras  pigz

   1. [安装yum](https://mirrors.tuna.tsinghua.edu.cn/help/ubuntu/)放弃。
   2.  注意架构，我们是[arm64](https://mirrors.tuna.tsinghua.edu.cn/help/ubuntu-ports/)，[不能用amd64/32](https://mirrors.tuna.tsinghua.edu.cn/help/ubuntu/)

   ```
   deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal main restricted universe multiverse
   # deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal main restricted universe multiverse
   deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-updates main restricted universe multiverse
   # deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-updates main restricted universe multiverse
   deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-backports main restricted universe multiverse
   # deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-backports main restricted universe multiverse
   
   # deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-security main restricted universe multiverse
   # # deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-security main restricted universe multiverse
   
   deb http://ports.ubuntu.com/ubuntu-ports/ focal-security main restricted universe multiverse
   ```

4. [minikube](https://minikube.sigs.k8s.io/docs/start/) 安装

   1. 安装

      ```
      curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-arm64
      sudo install minikube-linux-arm64 /usr/local/bin/minikube
      ```

      

   2. 安装docker，使用docker为dirver会出现，结果失败，更换driver

      原因：

      ```
      sudo systemctl restart docker会出错，可以让minikude重启docker使用sudo service docker start，但是不会修改
      ```

      > https://blog.csdn.net/weixin_39266845/article/details/127327423 修改daemon的类型失败
      >
      > https://blog.csdn.net/u012833399/article/details/128533933

   

   1. 根据文章创建用户

      [Exiting due to DRV_AS_ROOT: The "docker" driver should not be used with root privileges.](https://github.com/kubernetes/minikube/issues/7903)

   2. 使用qumu为driver

      ```
      sudo chown root:kvm /dev/kvm
      sudo chmod 666 /dev/kvm
      //可能第一个不成功，我是第二个成功了
      ```

      

   3. minikube start --driver=qemu

      ![](https://p.ipic.vip/uuc4y6.png)

      

## bpf_helpers.h

```c
/*
 * Note that bpf programs need to include either
 * vmlinux.h (auto-generated from BTF) or linux/types.h
 * in advance since bpf_helper_defs.h uses such types
 * as __u64.
 * bpf_helper_defs.h使用了__u64类型
 */
/* vmlinux.h
*typedef long long unsigned int __u64;
*/
```

1. BTF与[vmlinux.h](https://blog.csdn.net/m0_37824308/article/details/122337077)

   Vmlinux.h是使用工具生成的代码文件。包含了内核运行时使用的所有数据类型的定义。 vmlinux 文件通常会打包在Linux发行版中。

   bpftool可以根据vmlinux文件生成vmlinux.h文件。

   ```bash
    bpftool btf dump file /sys/kernel/btf/vmlinux format c > vmlinux.h
   ```

   ```
   WARNING: bpftool not found for kernel 5.19.0-1024
   
     You may need to install the following packages for this specific kernel:
       linux-tools-5.19.0-1024-aws
       linux-cloud-tools-5.19.0-1024-aws
   
     You may also want to install one of the following packages to keep up to date:
       linux-tools-aws
       linux-cloud-tools-aws
   ```

   ![](https://p.ipic.vip/nxp1xx.png)

   ```
   //*vmlinx是elf文件？
   file /sys/kernel/btf/vmlinux
   /sys/kernel/btf/vmlinux: data
   ```

   使得ebpg读取内核内存数据时根据数据结构读取和理解数据。

   不同机器Linux数据结构可能不同。但是使用libbpf库可以实现“CO-RE”。

   

   > /sys/kernel/btf不一定有
   >
   > ![](https://p.ipic.vip/xyak7a.png)
   >
   > ![](https://p.ipic.vip/elta1t.png)

   

2. /usr/include/bpf/bpf_helper_defs.h

   /* Forward declarations of BPF structs */里面生命的结构体在vmlinux.h里面可以找到

   可能是 auto-generated file. 

   也可能是linux-headers-$(uname -r)/tool/bpf里面

   ```
   ls /usr/include/bpf
   ```

   

3. linux/types.h 有一个树形结构的依赖，很难找到所有文件。

4. bpf_helpers.h

   ```c
   /* SPDX-License-Identifier: (LGPL-2.1 OR BSD-2-Clause) */
   #ifndef __BPF_HELPERS__
   #define __BPF_HELPERS__
   
   /*
    * Note that bpf programs need to include either
    * vmlinux.h (auto-generated from BTF) or linux/types.h
    * in advance since bpf_helper_defs.h uses such types
    * as __u64.
    */
   #include "bpf_helper_defs.h"
   
   #define __uint(name, val) int (*name)[val]
   #define __type(name, val) typeof(val) *name
   #define __array(name, val) typeof(val) *name[]
   
   /* Helper macro to print out debug messages */
   #define bpf_printk(fmt, ...)				\
   ({							\
   	char ____fmt[] = fmt;				\
   	bpf_trace_printk(____fmt, sizeof(____fmt),	\
   			 ##__VA_ARGS__);		\
   })
   
   /*
    * Helper macro to place programs, maps, license in
    * different sections in elf_bpf file. Section names
    * are interpreted by libbpf depending on the context (BPF programs, BPF maps,
    * extern variables, etc).
    * To allow use of SEC() with externs (e.g., for extern .maps declarations),
    * make sure __attribute__((unused)) doesn't trigger compilation warning.
    Helper宏将程序、地图、许可证放在elf_bpf文件的不同部分。部分名称由libbpf根据上下文（BPF程序、BPF映射、外部变量等）进行解释。要允许将SEC[（）]与extern一起使用（例如，对于extern.maps声明），请确保__attribute__（未使用））不会触发编译警告。
    */
   #define SEC(name) \
   	_Pragma("GCC diagnostic push")					    \
   	_Pragma("GCC diagnostic ignored \"-Wignored-attributes\"")	    \
   	__attribute__((section(name), used))				    \
   	_Pragma("GCC diagnostic pop")					    \
   
   /* Avoid 'linux/stddef.h' definition of '__always_inline'. */
   #undef __always_inline
   #define __always_inline inline __attribute__((always_inline))
   
   #ifndef __noinline
   #define __noinline __attribute__((noinline))
   #endif
   #ifndef __weak
   #define __weak __attribute__((weak))
   #endif
   
   /*
    * Use __hidden attribute to mark a non-static BPF subprogram effectively
    * static for BPF verifier's verification algorithm purposes, allowing more
    * extensive and permissive BPF verification process, taking into account
    * subprogram's caller context.
    */
   #define __hidden __attribute__((visibility("hidden")))
   
   /* When utilizing vmlinux.h with BPF CO-RE, user BPF programs can't include
    * any system-level headers (such as stddef.h, linux/version.h, etc), and
    * commonly-used macros like NULL and KERNEL_VERSION aren't available through
    * vmlinux.h. This just adds unnecessary hurdles and forces users to re-define
    * them on their own. So as a convenience, provide such definitions here.
    */
   #ifndef NULL
   #define NULL ((void *)0)
   #endif
   
   #ifndef KERNEL_VERSION
   #define KERNEL_VERSION(a, b, c) (((a) << 16) + ((b) << 8) + ((c) > 255 ? 255 : (c)))
   #endif
   
   /*
    * Helper macros to manipulate data structures
    */
   #ifndef offsetof
   #define offsetof(TYPE, MEMBER)	((unsigned long)&((TYPE *)0)->MEMBER)
   #endif
   #ifndef container_of
   #define container_of(ptr, type, member)				\
   	({							\
   		void *__mptr = (void *)(ptr);			\
   		((type *)(__mptr - offsetof(type, member)));	\
   	})
   #endif
   
   /*
    * Helper macro to throw a compilation error if __bpf_unreachable() gets
    * built into the resulting code. This works given BPF back end does not
    * implement __builtin_trap(). This is useful to assert that certain paths
    * of the program code are never used and hence eliminated by the compiler.
    *
    * For example, consider a switch statement that covers known cases used by
    * the program. __bpf_unreachable() can then reside in the default case. If
    * the program gets extended such that a case is not covered in the switch
    * statement, then it will throw a build error due to the default case not
    * being compiled out.
    */
   #ifndef __bpf_unreachable
   # define __bpf_unreachable()	__builtin_trap()
   #endif
   
   /*
    * Helper function to perform a tail call with a constant/immediate map slot.
    */
   #if __clang_major__ >= 8 && defined(__bpf__)
   static __always_inline void
   bpf_tail_call_static(void *ctx, const void *map, const __u32 slot)
   {
   	if (!__builtin_constant_p(slot))
   		__bpf_unreachable();
   
   	/*
   	 * Provide a hard guarantee that LLVM won't optimize setting r2 (map
   	 * pointer) and r3 (constant map index) from _different paths_ ending
   	 * up at the _same_ call insn as otherwise we won't be able to use the
   	 * jmpq/nopl retpoline-free patching by the x86-64 JIT in the kernel
   	 * given they mismatch. See also d2e4c1e6c294 ("bpf: Constant map key
   	 * tracking for prog array pokes") for details on verifier tracking.
   	 *
   	 * Note on clobber list: we need to stay in-line with BPF calling
   	 * convention, so even if we don't end up using r0, r4, r5, we need
   	 * to mark them as clobber so that LLVM doesn't end up using them
   	 * before / after the call.
   	 */
   	asm volatile("r1 = %[ctx]\n\t"
   		     "r2 = %[map]\n\t"
   		     "r3 = %[slot]\n\t"
   		     "call 12"
   		     :: [ctx]"r"(ctx), [map]"r"(map), [slot]"i"(slot)
   		     : "r0", "r1", "r2", "r3", "r4", "r5");
   }
   #endif
   
   /*
    * Helper structure used by eBPF C program
    * to describe BPF map attributes to libbpf loader
    */
   struct bpf_map_def {
   	unsigned int type;
   	unsigned int key_size;
   	unsigned int value_size;
   	unsigned int max_entries;
   	unsigned int map_flags;
   };
   
   enum libbpf_pin_type {
   	LIBBPF_PIN_NONE,
   	/* PIN_BY_NAME: pin maps by name (in /sys/fs/bpf by default) */
   	LIBBPF_PIN_BY_NAME,
   };
   
   enum libbpf_tristate {
   	TRI_NO = 0,
   	TRI_YES = 1,
   	TRI_MODULE = 2,
   };
   
   #define __kconfig __attribute__((section(".kconfig")))
   #define __ksym __attribute__((section(".ksyms")))
   
   #ifndef ___bpf_concat
   #define ___bpf_concat(a, b) a ## b
   #endif
   #ifndef ___bpf_apply
   #define ___bpf_apply(fn, n) ___bpf_concat(fn, n)
   #endif
   #ifndef ___bpf_nth
   #define ___bpf_nth(_, _1, _2, _3, _4, _5, _6, _7, _8, _9, _a, _b, _c, N, ...) N
   #endif
   #ifndef ___bpf_narg
   #define ___bpf_narg(...) \
   	___bpf_nth(_, ##__VA_ARGS__, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0)
   #endif
   
   #define ___bpf_fill0(arr, p, x) do {} while (0)
   #define ___bpf_fill1(arr, p, x) arr[p] = x
   #define ___bpf_fill2(arr, p, x, args...) arr[p] = x; ___bpf_fill1(arr, p + 1, args)
   #define ___bpf_fill3(arr, p, x, args...) arr[p] = x; ___bpf_fill2(arr, p + 1, args)
   #define ___bpf_fill4(arr, p, x, args...) arr[p] = x; ___bpf_fill3(arr, p + 1, args)
   #define ___bpf_fill5(arr, p, x, args...) arr[p] = x; ___bpf_fill4(arr, p + 1, args)
   #define ___bpf_fill6(arr, p, x, args...) arr[p] = x; ___bpf_fill5(arr, p + 1, args)
   #define ___bpf_fill7(arr, p, x, args...) arr[p] = x; ___bpf_fill6(arr, p + 1, args)
   #define ___bpf_fill8(arr, p, x, args...) arr[p] = x; ___bpf_fill7(arr, p + 1, args)
   #define ___bpf_fill9(arr, p, x, args...) arr[p] = x; ___bpf_fill8(arr, p + 1, args)
   #define ___bpf_fill10(arr, p, x, args...) arr[p] = x; ___bpf_fill9(arr, p + 1, args)
   #define ___bpf_fill11(arr, p, x, args...) arr[p] = x; ___bpf_fill10(arr, p + 1, args)
   #define ___bpf_fill12(arr, p, x, args...) arr[p] = x; ___bpf_fill11(arr, p + 1, args)
   #define ___bpf_fill(arr, args...) \
   	___bpf_apply(___bpf_fill, ___bpf_narg(args))(arr, 0, args)
   
   /*
    * BPF_SEQ_PRINTF to wrap bpf_seq_printf to-be-printed values
    * in a structure.
    */
   #define BPF_SEQ_PRINTF(seq, fmt, args...)			\
   ({								\
   	static const char ___fmt[] = fmt;			\
   	unsigned long long ___param[___bpf_narg(args)];		\
   								\
   	_Pragma("GCC diagnostic push")				\
   	_Pragma("GCC diagnostic ignored \"-Wint-conversion\"")	\
   	___bpf_fill(___param, args);				\
   	_Pragma("GCC diagnostic pop")				\
   								\
   	bpf_seq_printf(seq, ___fmt, sizeof(___fmt),		\
   		       ___param, sizeof(___param));		\
   })
   
   /*
    * BPF_SNPRINTF wraps the bpf_snprintf helper with variadic arguments instead of
    * an array of u64.
    */
   #define BPF_SNPRINTF(out, out_size, fmt, args...)		\
   ({								\
   	static const char ___fmt[] = fmt;			\
   	unsigned long long ___param[___bpf_narg(args)];		\
   								\
   	_Pragma("GCC diagnostic push")				\
   	_Pragma("GCC diagnostic ignored \"-Wint-conversion\"")	\
   	___bpf_fill(___param, args);				\
   	_Pragma("GCC diagnostic pop")				\
   								\
   	bpf_snprintf(out, out_size, ___fmt,			\
   		     ___param, sizeof(___param));		\
   })
   
   #endif
   
   ```

   > \__attribute__(())
   >
   > 通知编译器对变量、函数做一些特性的检查。失败不通过编译
   >
   > \__attribute__((section("section_name"))) 
   > 其作用是将作用的函数或数据放入指定名为"section_name"输入段。 
   >
   > \__attribute__((used))用于告诉编译器在目标文件中保留一个静态函数或者静态变量，即使它没有被引用。
   >
   > 可以用于屏蔽局部代码的警告
   >
   > ```
   > #pragma GCC diagnostic push 
   > #pragma GCC diagnostic ignored "-Wformat" 
   > //code
   > #pragma GCC diagnostic pop
   > 
   > ```
   >
   > 

5. linux/bpf.h

   定义了可以使用的eBPF Syscall Commands

   ```c
   enum bpf_cmd {
   	BPF_MAP_CREATE,
   	BPF_MAP_LOOKUP_ELEM,
   	BPF_MAP_UPDATE_ELEM,
   	BPF_MAP_DELETE_ELEM,
   	BPF_MAP_GET_NEXT_KEY,
   	BPF_PROG_LOAD,
   	BPF_OBJ_PIN,
   	BPF_OBJ_GET,
   	BPF_PROG_ATTACH,
   	BPF_PROG_DETACH,
   	BPF_PROG_TEST_RUN,
   	BPF_PROG_RUN = BPF_PROG_TEST_RUN,
   	BPF_PROG_GET_NEXT_ID,
   	BPF_MAP_GET_NEXT_ID,
   	BPF_PROG_GET_FD_BY_ID,
   	BPF_MAP_GET_FD_BY_ID,
   	BPF_OBJ_GET_INFO_BY_FD,
   	BPF_PROG_QUERY,
   	BPF_RAW_TRACEPOINT_OPEN,
   	BPF_BTF_LOAD,
   	BPF_BTF_GET_FD_BY_ID,
   	BPF_TASK_FD_QUERY,
   	BPF_MAP_LOOKUP_AND_DELETE_ELEM,
   	BPF_MAP_FREEZE,
   	BPF_BTF_GET_NEXT_ID,
   	BPF_MAP_LOOKUP_BATCH,
   	BPF_MAP_LOOKUP_AND_DELETE_BATCH,
   	BPF_MAP_UPDATE_BATCH,
   	BPF_MAP_DELETE_BATCH,
   	BPF_LINK_CREATE,
   	BPF_LINK_UPDATE,
   	BPF_LINK_GET_FD_BY_ID,
   	BPF_LINK_GET_NEXT_ID,
   	BPF_ENABLE_STATS,
   	BPF_ITER_CREATE,
   	BPF_LINK_DETACH,
   	BPF_PROG_BIND_MAP,
   };
   
   enum bpf_map_type {
   	BPF_MAP_TYPE_UNSPEC,
   	BPF_MAP_TYPE_HASH,
   	BPF_MAP_TYPE_ARRAY,
   	BPF_MAP_TYPE_PROG_ARRAY,
   	BPF_MAP_TYPE_PERF_EVENT_ARRAY,
   	BPF_MAP_TYPE_PERCPU_HASH,
   	BPF_MAP_TYPE_PERCPU_ARRAY,
   	BPF_MAP_TYPE_STACK_TRACE,
   	BPF_MAP_TYPE_CGROUP_ARRAY,
   	BPF_MAP_TYPE_LRU_HASH,
   	BPF_MAP_TYPE_LRU_PERCPU_HASH,
   	BPF_MAP_TYPE_LPM_TRIE,
   	BPF_MAP_TYPE_ARRAY_OF_MAPS,
   	BPF_MAP_TYPE_HASH_OF_MAPS,
   	BPF_MAP_TYPE_DEVMAP,
   	BPF_MAP_TYPE_SOCKMAP,
   	BPF_MAP_TYPE_CPUMAP,
   	BPF_MAP_TYPE_XSKMAP,
   	BPF_MAP_TYPE_SOCKHASH,
   	BPF_MAP_TYPE_CGROUP_STORAGE,
   	BPF_MAP_TYPE_REUSEPORT_SOCKARRAY,
   	BPF_MAP_TYPE_PERCPU_CGROUP_STORAGE,
   	BPF_MAP_TYPE_QUEUE,
   	BPF_MAP_TYPE_STACK,
   	BPF_MAP_TYPE_SK_STORAGE,
   	BPF_MAP_TYPE_DEVMAP_HASH,
   	BPF_MAP_TYPE_STRUCT_OPS,
   	BPF_MAP_TYPE_RINGBUF,
   	BPF_MAP_TYPE_INODE_STORAGE,
   	BPF_MAP_TYPE_TASK_STORAGE,
   };
   
   /* Note that tracing related programs such as
    * BPF_PROG_TYPE_{KPROBE,TRACEPOINT,PERF_EVENT,RAW_TRACEPOINT}
    * are not subject to a stable API since kernel internal data
    * structures can change from release to release and may
    * therefore break existing tracing BPF programs. Tracing BPF
    * programs correspond to /a/ specific kernel which is to be
    * analyzed, and not /a/ specific kernel /and/ all future ones.
    */
   enum bpf_prog_type {
   	BPF_PROG_TYPE_UNSPEC,
   	BPF_PROG_TYPE_SOCKET_FILTER,
   	BPF_PROG_TYPE_KPROBE,
   	BPF_PROG_TYPE_SCHED_CLS,
   	BPF_PROG_TYPE_SCHED_ACT,
   	BPF_PROG_TYPE_TRACEPOINT,
   	BPF_PROG_TYPE_XDP,
   	BPF_PROG_TYPE_PERF_EVENT,
   	BPF_PROG_TYPE_CGROUP_SKB,
   	BPF_PROG_TYPE_CGROUP_SOCK,
   	BPF_PROG_TYPE_LWT_IN,
   	BPF_PROG_TYPE_LWT_OUT,
   	BPF_PROG_TYPE_LWT_XMIT,
   	BPF_PROG_TYPE_SOCK_OPS,
   	BPF_PROG_TYPE_SK_SKB,
   	BPF_PROG_TYPE_CGROUP_DEVICE,
   	BPF_PROG_TYPE_SK_MSG,
   	BPF_PROG_TYPE_RAW_TRACEPOINT,
   	BPF_PROG_TYPE_CGROUP_SOCK_ADDR,
   	BPF_PROG_TYPE_LWT_SEG6LOCAL,
   	BPF_PROG_TYPE_LIRC_MODE2,
   	BPF_PROG_TYPE_SK_REUSEPORT,
   	BPF_PROG_TYPE_FLOW_DISSECTOR,
   	BPF_PROG_TYPE_CGROUP_SYSCTL,
   	BPF_PROG_TYPE_RAW_TRACEPOINT_WRITABLE,
   	BPF_PROG_TYPE_CGROUP_SOCKOPT,
   	BPF_PROG_TYPE_TRACING,
   	BPF_PROG_TYPE_STRUCT_OPS,
   	BPF_PROG_TYPE_EXT,
   	BPF_PROG_TYPE_LSM,
   	BPF_PROG_TYPE_SK_LOOKUP,
   	BPF_PROG_TYPE_SYSCALL, /* a program that can execute syscalls */
   };
   
   enum bpf_attach_type {
   	BPF_CGROUP_INET_INGRESS,
   	BPF_CGROUP_INET_EGRESS,
   	BPF_CGROUP_INET_SOCK_CREATE,
   	BPF_CGROUP_SOCK_OPS,
   	BPF_SK_SKB_STREAM_PARSER,
   	BPF_SK_SKB_STREAM_VERDICT,
   	BPF_CGROUP_DEVICE,
   	BPF_SK_MSG_VERDICT,
   	BPF_CGROUP_INET4_BIND,
   	BPF_CGROUP_INET6_BIND,
   	BPF_CGROUP_INET4_CONNECT,
   	BPF_CGROUP_INET6_CONNECT,
   	BPF_CGROUP_INET4_POST_BIND,
   	BPF_CGROUP_INET6_POST_BIND,
   	BPF_CGROUP_UDP4_SENDMSG,
   	BPF_CGROUP_UDP6_SENDMSG,
   	BPF_LIRC_MODE2,
   	BPF_FLOW_DISSECTOR,
   	BPF_CGROUP_SYSCTL,
   	BPF_CGROUP_UDP4_RECVMSG,
   	BPF_CGROUP_UDP6_RECVMSG,
   	BPF_CGROUP_GETSOCKOPT,
   	BPF_CGROUP_SETSOCKOPT,
   	BPF_TRACE_RAW_TP,
   	BPF_TRACE_FENTRY,
   	BPF_TRACE_FEXIT,
   	BPF_MODIFY_RETURN,
   	BPF_LSM_MAC,
   	BPF_TRACE_ITER,
   	BPF_CGROUP_INET4_GETPEERNAME,
   	BPF_CGROUP_INET6_GETPEERNAME,
   	BPF_CGROUP_INET4_GETSOCKNAME,
   	BPF_CGROUP_INET6_GETSOCKNAME,
   	BPF_XDP_DEVMAP,
   	BPF_CGROUP_INET_SOCK_RELEASE,
   	BPF_XDP_CPUMAP,
   	BPF_SK_LOOKUP,
   	BPF_XDP,
   	BPF_SK_SKB_VERDICT,
   	BPF_SK_REUSEPORT_SELECT,
   	BPF_SK_REUSEPORT_SELECT_OR_MIGRATE,
   	BPF_PERF_EVENT,
   	__MAX_BPF_ATTACH_TYPE
   };
   
   #define MAX_BPF_ATTACH_TYPE __MAX_BPF_ATTACH_TYPE
   
   enum bpf_link_type {
   	BPF_LINK_TYPE_UNSPEC = 0,
   	BPF_LINK_TYPE_RAW_TRACEPOINT = 1,
   	BPF_LINK_TYPE_TRACING = 2,
   	BPF_LINK_TYPE_CGROUP = 3,
   	BPF_LINK_TYPE_ITER = 4,
   	BPF_LINK_TYPE_NETNS = 5,
   	BPF_LINK_TYPE_XDP = 6,
   	BPF_LINK_TYPE_PERF_EVENT = 7,
   
   	MAX_BPF_LINK_TYPE,
   };
   
   ```

   ```bash
   ubuntu@ip-172-31-28-100:~$ /usr/src/linux-headers-5.15.0-1031-aws/scripts/bpf_doc.py --filename /usr/src/linux-aws-headers-5.15.0-1031/include/uapi/linux/bpf.h > /tmp/bpf-helpers.rst
   ubuntu@ip-172-31-28-100:~$ rst2man /tmp/bpf-helpers.rst > /tmp/bpf-helpers.7
   ubuntu@ip-172-31-28-100:~$ man /tmp/bpf-helpers.7
   ```

   

6. Bpf/libbpf.h

   ```c
   /*
    * Libbpf allows callers to adjust BPF programs before being loaded
    * into kernel. One program in an object file can be transformed into
    * multiple variants to be attached to different hooks.
    /
   ```

   

7. 