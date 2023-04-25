/**
 * @file kprobe_unlink.c
 * @author JiadiSu (20302010043@fudan.edu.cn)
 * @brief 通过kprobe捕获unlink系统调用
 * @version 0.1
 * @date 2023-04-25
 * 
 * @copyright Copyright (c) 2023
 * 编译：clang -target bpf -Wall -O2 -c kprobe_unlink.c -o kprobe_unlink.o
 */

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>

/* BPF_KRPOBE宏定义在<bpf/bpf_tracing.h> */
SEC("kprobe/do_unlinkat")
int do_unlinkat(int dfd, struct filename* name)
{
    /* helper: bpf_get_current_pid_tgid() */
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    /* BPF_CORE_READ宏定义在<bpf/bpf_core_name.h> */
    // const char* filename = BPF_CORE_READ(name, name);
    /* helper: bpf_trace_printk() */
    bpf_trace_printk("KPROBE ENTRY pid = %d\n", pid);
    return 0;
}

/* BPF_KRETPROBE宏定义在<bpf/bpf_tracing.h> */
SEC("kretprobe/do_unlinkat")
int do_unlinkat_exit(long ret)
{
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    bpf_trace_printk("KPROBE EXIT: pid = %d, ret = %ld\n", pid, ret);
    return 0;
}

char _license[] SEC("license") = "Dual BSD/GPL";
