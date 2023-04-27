## [使用libbpf写一个bpf程序](https://blog.qsliu.dev/posts/writing_an_ebpf_application/)

### 安装Linux headers 失败



```bash
root@07b6e94d459c:~# apt install linux-headers-$(uname -r)
Reading package lists... Done
Building dependency tree       
Reading state information... Done
E: Unable to locate package linux-headers-4.19.90-2204.4.0.0146.oe1.aarch64
E: Couldn't find any package by glob 'linux-headers-4.19.90-2204.4.0.0146.oe1.aarch64'
```

```bash
root@07b6e94d459c:~# uname -r
4.19.90-2204.4.0.0146.oe1.aarch64
```

解决方法：[Docker从不使用其他内核：内核始终是您的宿主内核。](https://blog.csdn.net/weixin_36401868/article/details/116661211)换主机

### asm/types.h找不到

1. clang

   ```
   ubuntu@ip-172-31-28-100:~$ clang \
     -target bpf \
     -g -O2 \
     -o hello_world.bpf.o \
     -c hello_world.bpf.c
   In file included from hello_world.bpf.c:1:
   In file included from /usr/include/linux/bpf.h:11:
   /usr/include/linux/types.h:5:10: fatal error: 'asm/types.h' file not found
   #include <asm/types.h>
            ^~~~~~~~~~~~~
   1 error generated.
   ```

   ![](https://p.ipic.vip/l5acaa.png)

   ```
   sudo ln -s usr/include/aarch64-linux-gnu/asm asm
   ```

2. bpf/libbpf.h第三方库找不到
       1 | #include <bpf/libbpf.h>

   1. 使用sudo apt-get install libbpf-dev
   2. 使用 https://libbpf.readthedocs.io/en/latest/libbpf_build.html

3.  

   ```bash
   gcc helloworld.c  -l bpf -o helloworld && ./helloworld.c
   ```