# eBPF

## 一、eBPF是什么？

### BPF

BPF全称为Berkeley Packet Filter（伯克利包过滤器），起源于1992年的一篇论文，该论文提出了一种对网络包进行过滤的框架，如下：

<img src="img/bpf.png" style="zoom:100%;" />

早期从网卡中接收到很多的数据包，要想从中过滤出想要的数据包，就需要将网卡接收的数据包都要**从内核空间拷贝一份到用户空间**（因为过滤的业务逻辑是用户指定的，其运行在用户态），然后由**用户程序对这些进行过滤**。这其中存在的问题是：无论是否有效，网络的数据包必须全部拷贝，然后再过滤出所需的数据包，这其中就产生了对无效数据包的无效拷贝，浪费CPU资源。BPF技术产生在这个问题背景下，BPF会在内核中直接过滤，从而避免一些无用的、浪费的拷贝。其背后的思想其实就是：与其把数据包复制到用户空间执行用户态程序过滤，不如把**过滤程序注入到内核**去（即扩展了内核的功能，后面会将eBPF的方式与编写内核模块的方式进行对比）。

### eBPF

eBPF即extend BPF，是在BPF技术的基础上进行扩展，丰富了BPF的功能，使其**不仅能够进行网络数据包的过滤**，而是**基本上可以使用在Linux各个子系统中**（如监控系统调用，系统行为等等）。在eBPF出现后，过去的BPF被称为cBPF（classic BPF）。在eBPF技术的支持下，内核变得可编程。

详细介绍eBPF之前，先简单看一下Linux内核的架构：

<img src="img/kernel-arch.png" style="zoom: 67%;" />

Linux 内核的主要目的是**抽象硬件或虚拟硬件并提供一致的 API**（系统调用），允许应用程序运行和共享资源。为了实现这一点，Linux维护了广泛的子系统和层集来分配这些职责。一般来说，用户可以配置每个子系统的具体行为，如果无法配置所需的行为，则需要根据需求来扩展内核的功能。

传统的内核扩展方式有两种：

- **更改内核源代码**并说服 Linux 内核社区需要更改，并在若干时间后作为新的Linux版本发布；
- **开发内核模块**，利用Linux可以灵活地装载和卸载模块来实现功能扩展，但具有安全风险，且需要针对Linux的不同版本进行维护；

eBPF提供了这两种方式之外的一种方式，将要扩展的功能用**eBPF的指令**来实现（eBPF的指令是什么后面再说），并将该指令加载到内核中，从而实现内核功能的扩展。这里有几个问题：

**1、同样是把要写的功能挂载到内核中，这与内核模块有什么区别？**

到目前来看，eBPF和内核模块的行为方式很像，都是开发对应要扩展的功能，然后将其加载进内核使用。但这两者有一个根本的区别就是：内核模块是直接以硬件可以执行的字节码（如arm汇编）加载到内核，内核可以将其作为内核代码的一部分直接由处理器取指执行。因此，内核模块加载到内核之后执行时会很快，但它存在两个问题：第一，需要针对不同的Linux版本进行维护；第二、由于模块加载时没有经过安全检查，具有安全风险。而eBPF利用与内核模块不同的扩展机制解决了这两个问题，机制如下：

首先要明确，eBPF技术其实是一个沙箱（sandbox）技术，它与Java的运行机制很相似。Java是在本机上运行了jvm（Java virtual machine，Java虚拟机），将Java的字节码文件（也就是由.java文件经过javac指令后生成的.class文件）运行在jvm上，由jvm完成与物理机的交互，从而实现了跨平台的执行。eBPF也是类似，该技术是在**内核态**运行了一个**基于eBPF指令架构的虚拟机（暂时可以简单理解为一种指令集ISA，它与本机是arm还是x86_64还是risc-v都无关）**（这个虚拟机可以类比jvm），由eBPF指令编写好的字节码程序（也就是用户想要给内核扩展的功能）会由该虚拟机**解释执行**（划重点：这里是解释执行，后面会说），从而实现功能的扩展，因为该虚拟机的存在，eBPF不需要像内核模块一样需要根据不同的Linux内核版本维护不同的版本，因此eBPF解决了内核模块存在的第一个问题。接下来是第二个问题：安全性的保证。前面提到过，eBPF对内核的扩展需要将eBPF字节码加载到内核中，在加载时内核就会有一个Verification的过程，Verification会检查该eBPF是否会对系统作出损害且可以保证运行完成（不会以死循环的方式占据系统资源），Verification保证了加载到内核的eBPF的一定是安全的（这里的检查机制我还没详细地查过，简单知道有个安全检查的过程就好）。于是eBPF的安全性也得到了保证，内核模块的第二个问题解决。

好，现在回到刚才划重点的地方：eBPF虚拟机的**解释执行**，没错，这里的用词是解释执行，因此它在执行起来会很慢。为了加速这个过程，有一个大神（忘了是哪个大神了）提出了“即使编译”（JIT）技术，目前就个人理解来说，JIT其实就是eBPF字节码被加载并经过Verification后，会经过一个编译的过程，将eBPF字节码编译为物理机的机器码，从而加速了eBPF程序的执行。当然，Verification与JIT都是在内核里做的（后面会说这里的细节，其实是通过系统调用的方式，系统调用里会完成这两个过程），因此肯定需要内核提供的支持，Verification是肯定有的，而JIT是一个加速的过程，似乎有些Linux版本不支持，但因为JIT只是为了加速，不影响功能的实现，所以也问题不大。

**2、什么是eBPF指令？**

前面提到，eBPF其实是一个运行在内核的虚拟机（看到的文章都说它是个沙箱，我其实不清楚和虚拟机有什么区别，为了后面说的方便点就当虚拟机了），既然是虚拟机，肯定有**可以执行指令**以及**对应的寄存器**的规范。eBPF虚拟机的这些规范是自己设计的，与arm、risc-v这种没有关系，它是一个自己独立的指令集，规定了一些eBPF虚拟机才可以执行的指令。举个例子，这是一个cBPF的字节码，用来完成对tcp网络包的过滤（基于端口过滤）：

```assembly
ldh      [12]
jeq      #0x86dd          jt 2    jf 6
ldb      [20]
jeq      #0x6             jt 4    jf 15
ldh      [56]
jeq      #0x50            jt 14   jf 15
jeq      #0x800           jt 7    jf 15
ldb      [23]
jeq      #0x6             jt 9    jf 15
ldh      [20]
jset     #0x1fff          jt 15   jf 11
ldxb     4*([14]&0xf)
ldh      [x + 16]
jeq      #0x50            jt 14   jf 15
ret      #262144
ret      #0
```

不用懂这一段eBPF汇编具体每条是什么含义，只需要知道eBPF设计了一个独一无二的指令集，包括一些指令以及寄存器的数据（想了解详细的虚拟机架构可以去看[eBPF概述第二部分：机器和字节码](https://www.collabora.com/news-and-blog/blog/2019/04/15/an-ebpf-overview-part-2-machine-and-bytecode/)）。也就是说，eBPF程序的最原生的表现形式就是eBPF字节码。

另外提前插一句：在cBPF时期，由于程序的功能单一，只是做网络包的过滤，因此没有出现专门的编译器，cBPF程序的开发都是用字节码完成，因此编程难度很大。而到了eBPF时期，程序的功能越来越丰富，都用字节码来实现并不现实，于是出现了像`clang / llvm`这种把高级语言程序（如C）编译成eBPF字节码的编译器，从而可以利用高级语言开发eBPF程序。（这一部分在后面eBPF开发部分会说）

**3、内核模块可以通过Linux加载模块的指令加载进内核，那eBPF程序是怎么加载到内核的？**

eBPF程序加载到内核，一般是通过`bpf()`系统调用（现在的各种简化开发的框架，如bcc和libbpf，都在不同程度上封装了`bpf()`系统调用），该系统调用的参数以及具体功能后面会详细地说，现在只结合`bpf()`对eBPF的加载过程做一个简单的说明。

要加载一个eBPF程序，首先要有一个待加载的**eBPF字节码**。在最原生的方式下，eBPF的字节码不是一个单独的文件，而是一个结构体，该结构体被当作`bpf()`系统调用的参数传入后在内核中完成加载。为了便于理解，这里给出一个例子：

```c
// BPF程序就是一个bpf_insn数组, 一个struct bpf_insn代表一条bpf指令
struct bpf_insn bpf_prog[] = {
    { 0xb7, 0, 0, 0, 0x2 }, //初始化一个struct bpf_insn, 指令含义: mov r0, 0x2;
    { 0x95, 0, 0, 0, 0x0 }, //初始化一个struct bpf_insn, 指令含义: exit;
};
```

`struct bpf_insn`就是一条eBPF字节码指令，用上面的方式可以设置两条指令`mov r0 0x2; exit;`。这两条指令对应的数组就可以看作是一个eBPF字节码程序（这可能也是为什么这个数组命名为`bpf_prog`，含义应该就是bpf program）。最终这个eBPF字节码数组会被作为`bpf()`系统调用的参数，经系统调用把这段字节码加载到内核，完成内核功能的扩展。（另外要说明一下，这里的数组`bpf_prog`在被真正作为`bpf()`的参数时，其实会被另一个结构体在外面封装一层，封装这一层会包含另外一些需要用到的信息，如`bpf_type`，`bpf_license`这种，但最终这个字节码数组`bpf_prog`肯定会传给kernel的）。





## 参考

参考文章：

eBPF概念性介绍：[What is eBPF?](https://ebpf.io/what-is-ebpf/)、[eBPF介绍](http://kerneltravel.net/blog/2021/zxj-ebpf1/)

eBPF代码实现介绍：[eBPF概述第一部分：简介](https://www.collabora.com/news-and-blog/blog/2019/04/05/an-ebpf-overview-part-1-introduction/)、[eBPF概述第二部分：机器和字节码](https://www.collabora.com/news-and-blog/blog/2019/04/15/an-ebpf-overview-part-2-machine-and-bytecode/)、[XDP和eBPF简介](https://blogs.igalia.com/dpino/2019/01/07/introduction-to-xdp-and-ebpf/)、[BPF之路一bpf系统调用](https://www.anquanke.com/post/id/263803)

eBPF开发参考项目：[官方入门Lab](https://play.instruqt.com/embed/isovalent/tracks/ebpf-getting-started?token=em_9nxLzhlV41gb3rKM&show_challenges=true)、[awesome-ebpf](https://github.com/zoidbergwill/awesome-ebpf)、[bpf-developer-tutorial](https://github.com/eunomia-bpf/bpf-developer-tutorial)、[libbpf-bootstrap](https://github.com/libbpf/libbpf-bootstrap)、[bcc](https://github.com/iovisor/bcc)（该项目主要关注其中的文档：[tutorial_bcc_python_developer.md](https://github.com/iovisor/bcc/blob/master/docs/tutorial_bcc_python_developer.md)）

