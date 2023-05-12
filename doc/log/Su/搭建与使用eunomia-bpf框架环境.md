# 搭建与使用eunomia-bpf框架

服务器环境：

- Ubuntu 20.04.4 LTS
- Linux hecs-67681 5.4.0-100-generic #113-Ubuntu SMP Thu Feb 3 18:43:29 UTC 2022 x86_64 x86_64 x86_64 GNU/Linux

### 搭建过程



### 使用过程

使用eunomia-bpf编译构建的bpf程序是：[opensnoop.bpf.c](../../note/Su/test/eunomia/opensnoop.bpf.c)

#### 1、使用ecc编译

在`ecc`所在的路径下（如果已经将ecc加到了系统路径则忽略这个限制），执行如下指令：

```bash
$ ./ecc ../opensnoop.bpf.c ../opensnoop.h
# 结果如下，生成了package.json
Compiling bpf object...
warning: text is not json: Process ID to trace use it as a string
warning: text is not json: Thread ID to trace use it as a string
warning: text is not json: User ID to trace use it as a string
warning: text is not json: trace only failed events use it as a string
warning: text is not json: Trace open family syscalls. use it as a string
Generating export types...
Packing ebpf object and config into ../package.json...
```

编译后多出了如下三个文件：

```bash
-rw-r--r-- 1 root root   12936 May 13 00:33 opensnoop.bpf.o
-rw-r--r-- 1 root root    1525 May 13 00:33 opensnoop.skel.json
-rw-r--r-- 1 root root    5954 May 13 00:33 package.json
```

- opensnoop.bpf.o：opensnoop.bpf.c的编译结果，见[eunomia](../../note/Su/test/eunomia/)目录下，是一个ELF格式字节码文件；
- [opensnoop.skel.json](../../note/Su/test/eunomia/opensnoop.skel.json)：描述bpf程序结构的json文件，包括了ELF的描述性信息；
- [package.json](../../note/Su/test/eunomia/package.json)：最终生成的中间产物，该文件可以被ecli执行，完成bpf的安装；

#### 2、使用ecli执行

在`ecli`所在的路径下（如果已经将ecc加到了系统路径则忽略这个限制），执行如下指令：

```bash
$ ./ecli run ../package.json
# 结果如下，会不断输出被检测到的系统调用
INFO [faerie::elf] strtab: 0xe4c symtab 0xe88 relocs 0xed0 sh_offset 0xed0
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable pid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable tgid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable uid_target, use the default one in ELF
INFO [bpf_loader_lib::skeleton::preload::section_loader] load runtime arg (user specified the value through cli, or predefined in the skeleton) for targ_failed: Bool(false), real_type=<INT> '_Bool' bits:8 off:0 enc:bool, btf_type=BtfVar { name: "targ_failed", type_id: 46, kind: GlobalAlloc }
INFO [bpf_loader_lib::skeleton::preload::section_loader] User didn't specify custom value for variable __eunomia_dummy_event_ptr, use the default one in ELF
INFO [bpf_loader_lib::skeleton::poller] Running ebpf program...
<1683909690> <PLAIN> TIME     TS     PID    UID    RET    FLAGS  COMM   FNAME
```

此时也可以通过`bpftool prog list`观察系统挂载的bpf程序，存在刚才通过eunomia挂载的4个程序：

```bash
$ bpftool prog list
...
94: tracepoint  name tracepoint__sys  tag 07014be5359438f8  gpl
        loaded_at 2023-05-13T00:43:03+0800  uid 0
        xlated 288B  jited 167B  memlock 4096B  map_ids 51,48
        btf_id 36
96: tracepoint  name tracepoint__sys  tag 8ee3432dcd98ffc3  gpl
        loaded_at 2023-05-13T00:43:03+0800  uid 0
        xlated 288B  jited 167B  memlock 4096B  map_ids 51,48
        btf_id 36
97: tracepoint  name tracepoint__sys  tag 541339de114a40e6  gpl
        loaded_at 2023-05-13T00:43:03+0800  uid 0
        xlated 696B  jited 477B  memlock 4096B  map_ids 48,51,49
        btf_id 36
98: tracepoint  name tracepoint__sys  tag 541339de114a40e6  gpl
        loaded_at 2023-05-13T00:43:03+0800  uid 0
        xlated 696B  jited 477B  memlock 4096B  map_ids 48,51,49
        btf_id 36
```

#### 3、解决的一些疑惑

##### 用户态怎么没的？

个人理解是，`ecli`程序代替了user态bpf程序的功能，且它更灵活更便于控制。

##### ecli执行依赖什么？

ecli执行bpf只依赖package.json，这点可以通过删除opensnoop.bpf.o和opensnoop.skel.json文件后仍然正常执行得出结论。

##### package.json的内容是什么？

观察package.json的内容，发现它包含了opensnoop.skel.json的内容，在meta字段的bpf_skel结构中，整个包含了opensnoop.skel.json，这也是为什么ecli执行不依赖opensnoop.skel.json；

同时它还有一个字段bpf_object，猜测是bpf程序，它以一种方式保存了opensnoop.bpf.o这个ELF格式文件（推断的依赖是，后面有一个bpf_object_size为12396，正好对应了opensnoop.bpf.o的大小）。