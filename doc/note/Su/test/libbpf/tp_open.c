/**
 * @file tp_open.c
 * @author JiadiSu (20302010043@fudan.edu.cn)
 * @brief 捕获进程打开文件的系统调用
 * @version 0.1
 * @date 2023-04-25
 * 
 * @copyright Copyright (c) 2023
 * 编译：clang -target bpf -Wall -O2 -c tp_open.c -o tp_open.o
 */

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>

const volatile pid_target = 0;

SEC("tp/syscalls/sys_enter_openat")
int openat_handler(struct trace_event_raw_sys_enter* ctx)
{
    pid_t pid = bpf_get_current_pid_tgid();

    bpf_trace_printk("Process ID: %d enter sys openat\n", pid);
    return 0;
}

char _license[] SEC("license") = "GPL";
