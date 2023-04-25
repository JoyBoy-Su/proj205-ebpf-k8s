/* Copyright (c) 2015 PLUMgrid, http://plumgrid.com
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of version 2 of the GNU General Public
 * License as published by the Free Software Foundation.
 */
#include <linux/ptrace.h>
#include <linux/version.h>
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

struct pair {
	u64 val;
	u64 ip;
};

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__type(key, long);
	__type(value, struct pair);
	__uint(max_entries, 1000000);
} my_map SEC(".maps");				/* 通过SEC(".map")创建一个map用来存取数据 */

/* kprobe is NOT a stable ABI. If kernel internals change this bpf+kprobe
 * example will no longer be meaningful
 */
SEC("kprobe/kmem_cache_free")		/* 进入kmem_cache_free时的探针 */
int bpf_prog1(struct pt_regs *ctx)
{
    /* 声明在<bpf/bpf_tracing.h> #define PT_REGS_PARM2(x) ((x)->si) */
	long ptr = PT_REGS_PARM2(ctx);
	
    /* bpf-helpers，删除map中ptr这个键 */
	bpf_map_delete_elem(&my_map, &ptr);
	return 0;
}

SEC("kretprobe/kmem_cache_alloc_node")
int bpf_prog2(struct pt_regs *ctx)
{
    /* 声明在<bpf/bpf_tracing.h> #define PT_REGS_RC(x) ((x)->ax) */
	long ptr = PT_REGS_RC(ctx);
	long ip = 0;

	/**
	 * 声明在<bpf/bpf_tracing.h>
	 * #define BPF_KRETPROBE_READ_RET_IP ({ (ip) = (ctx)->link; })
	 * 获得kmem_cache_alloc_node调用者的ip 
	 */
	BPF_KRETPROBE_READ_RET_IP(ip, ctx);

	struct pair v = {
		.val = bpf_ktime_get_ns(),	/* bpf-helpers，获取系统启动以来的时间 */
		.ip = ip,
	};
	
    /* bpf-helpers，更新map中ptr这个键的值为v */
	bpf_map_update_elem(&my_map, &ptr, &v, BPF_ANY);
	return 0;
}

/* 设置eBPF的license和version */
char _license[] SEC("license") = "GPL";
u32 _version SEC("version") = LINUX_VERSION_CODE;