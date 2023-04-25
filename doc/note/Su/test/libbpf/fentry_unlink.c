/**
 * @file fentry_unlink.c
 * @author JiadiSu (20302010043@fudan.edu.cn)
 * @brief 通过fentry(function entry)捕获unlink系统调用
 * @version 0.1
 * @date 2023-04-25
 * 
 * @copyright Copyright (c) 2023
 * 编译：clang -target bpf -Wall -O2 -c fentry_unlink.c -o fentry_unlink.o
 */

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

SEC("fentry/do_unlinkat")
int BPF_PROG(do_unlinkat, int dfd, struct filename* name)
{
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    bpf_trace_printk("fentry: pid = %d, filename = %s\n", pid, name->name);
    return 0;
}

SEC("fexit/do_unlinkat")
int BPF_PROG(do_unlinkat_exit, int dfd, struct filename* name, long ret)
{
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    bpf_trace_printk("fexit: pid = %d, filename = %s, ret = %ld\n", pid, name->name, ret);
    return 0;
}


char _license[] SEC("license") = "GPL";
