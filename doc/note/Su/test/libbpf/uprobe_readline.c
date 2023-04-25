/**
 * @file uprobe_readline.c
 * @author JiadiSu (20302010043@fudan.edu.cn)
 * @brief 使用 uprobe 捕获 bash 的 readline 函数调用
 * @version 0.1
 * @date 2023-04-25
 * 
 * @copyright Copyright (c) 2023
 * 编译：clang -target bpf -Wall -O2 -c uprobe_readline.c -o uprobe_readline.o
 */

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

#define TASK_COMM_LEN 16
#define MAX_LINE_SIZE 80

/* Format of u[ret]probe section definition supporting auto-attach:
 * u[ret]probe/binary:function[+offset]
 *
 * binary can be an absolute/relative path or a filename; the latter is resolved to a
 * full binary path via bpf_program__attach_uprobe_opts.
 *
 * Specifying uprobe+ ensures we carry out strict matching; either "uprobe" must be
 * specified (and auto-attach is not possible) or the above format is specified for
 * auto-attach.
 */
SEC("uretprobe//bin/bash:readline")
int BPF_KRETPROBE(printret, const void* ret)
{
    char str[MAX_LINE_SIZE];
    char comm[TASK_COMM_LEN];

    if (!ret) return 0;

    bpf_get_current_comm(&comm, sizeof(comm));

    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    /* 复制sizeof(str)个字节的数据，从ret到str */
    bpf_probe_read_user_str(str, sizeof(str), ret);

    bpf_trace_printk("PID %d (%s) read: %s ", pid, comm, str);
    return 0;
}

char license[] SEC("license") = "GPL";
