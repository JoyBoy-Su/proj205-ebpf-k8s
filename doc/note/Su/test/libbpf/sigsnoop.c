/**
 * @file sigsnoop.c
 * @author JiadiSu (20302010043@fudan.edu.cn)
 * @brief 捕获进程发送信号的系统调用集合，使用 hash map 保存状态
 * @version 0.1
 * @date 2023-04-25
 * 
 * @copyright Copyright (c) 2023
 * 编译：clang -target bpf -Wall -O2 -c sigsnoop.c -o sigsnoop.o
 */

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

#define MAX_ENTRIES 10240
#define TASK_COMM_LEN 16

/* 定义一个event结构体，用来在map中存放sig事件数据 */
struct event {
    unsigned int pid;
    unsigned int tpid;
    int sig;
    int ret;
    char comm[TASK_COMM_LEN];
};

/* 定义一个bpf map，类型为hash，key为int(tid)，value为event */
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, MAX_ENTRIES);
    __type(key, __u32);
    __type(value, struct event);
} values SEC(".maps");

/* 进入tracepoint时在map中加入数据 */
static int probe_entry(pid_t tpid, int sig)
{
    struct event event = {};
    __u64 pid_tgid;
    __u32 tid;

    /* helpers获取id */
    pid_tgid = bpf_get_current_pid_tgid();
    tid = (__u32)pid_tgid;
    /* 设置要保存的数据 */
    event.pid = pid_tgid >> 32;
    event.tpid = tpid;
    event.sig = sig;
    bpf_get_current_comm(event.comm, sizeof(event.comm));
    /* 通过helpers向map中更新数据 */
    bpf_map_update_elem(&values, &tid, &event, BPF_ANY);
    return 0;
}

/* 退出tracepoint时从map中删除数据 */
static int probe_exit(void *ctx, int ret)
{
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    __u32 tid = (__u32)pid_tgid;
    struct event *eventp;

    /* 从map中查到tid对应的value */
    eventp = bpf_map_lookup_elem(&values, &tid);
    if (!eventp) return 0;

    eventp->ret = ret;
    bpf_printk("PID %d (%s) sent signal %d to PID %d, ret = %d",
        eventp->pid, eventp->comm, eventp->sig, eventp->tpid, ret);

    cleanup:
    bpf_map_delete_elem(&values, &tid); /* 从map中删除信息 */
    return 0;
}

SEC("tracepoint/syscalls/sys_enter_kill")
int kill_entry(struct trace_event_raw_sys_enter *ctx)
{
    pid_t tpid = (pid_t)ctx->args[0];
    int sig = (int)ctx->args[1];

    /* bpf-bpf call */
    return probe_entry(tpid, sig);
}

SEC("tracepoint/syscalls/sys_exit_kill")
int kill_exit(struct trace_event_raw_sys_exit *ctx)
{
    /* bpf-bpf call */
    return probe_exit(ctx, ctx->ret);
}

char _license[] SEC("license") = "Dual BSD/GPL";
