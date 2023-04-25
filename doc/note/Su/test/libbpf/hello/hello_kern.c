/**
 * @file hello_kern.c
 * @author JiadiSu (20302010043@fudan.edu.cn)
 * @brief 通过libbpf接口实现一个hello world
 * @version 0.1
 * @date 2023-04-25
 * 
 * @copyright Copyright (c) 2023
 * 编译：
 * clang -target bpf -Wall -O2 -c hello_kern.c -o hello_kern.o
 * 产生汇编文件：
 * clang -target bpf -S -o hello_kern.S hello_kern.c
 */

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

SEC("tp/syscalls/sys_enter_openat")
int hello_world(void* ctx)
{
    char msg[] = "Hello world!\n";
    /* bpf-helpers 声明在<bpf/bpf_helper_defs.h> */
    bpf_trace_printk(msg, sizeof(msg));
    return 0;
}

char _license[] SEC("license") = "GPL";
