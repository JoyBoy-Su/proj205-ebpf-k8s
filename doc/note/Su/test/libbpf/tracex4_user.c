// SPDX-License-Identifier: GPL-2.0-only
/* Copyright (c) 2015 PLUMgrid, http://plumgrid.com
 */
#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <unistd.h>
#include <stdbool.h>
#include <string.h>
#include <time.h>

#include <bpf/bpf.h>			/* 一些对bpf syscall的封装 */
#include <bpf/libbpf.h>			/* 一些bpf_progam和bpf_object相关的函数 */

struct pair {
	long long val;
	__u64 ip;
};

/* 获得当前时钟的毫秒值 */
static __u64 time_get_ns(void)
{
	struct timespec ts;

	clock_gettime(CLOCK_MONOTONIC, &ts);
	return ts.tv_sec * 1000000000ull + ts.tv_nsec;
}

/* 输出fd对应的map中存储的数据 */
static void print_old_objects(int fd)
{
	long long val = time_get_ns();
	__u64 key, next_key;
	struct pair v;

	key = write(1, "\e[1;1H\e[2J", 11); /* clear screen */

	key = -1;
    /* bpf_map_get_next_key()声明在<bpf/bpf.h>，封装了bpf(BPF_MAP_GET_NEXT_KEY) */
	while (bpf_map_get_next_key(fd, &key, &next_key) == 0) {
        /* 同样声明在<bpf/bpf.h> ，获取查询值 */
		bpf_map_lookup_elem(fd, &next_key, &v);
		key = next_key;
		if (val - v.val < 1000000000ll)
			/* object was allocated more then 1 sec ago */
			continue;
		printf("obj 0x%llx is %2lldsec old was allocated at ip %llx\n",
		       next_key, (val - v.val) / 1000000000ll, v.ip);
	}
}

int main(int ac, char **argv)
{
	struct bpf_link *links[2];
	struct bpf_program *prog;
	struct bpf_object *obj;
	char filename[256];
	int map_fd, i, j = 0;

	snprintf(filename, sizeof(filename), "%s_kern.o", argv[0]);
    
    /* bpf_object__open_file()声明在<bpf/libbpf.h> */
	obj = bpf_object__open_file(filename, NULL);
	if (libbpf_get_error(obj)) {
		fprintf(stderr, "ERROR: opening BPF object file failed\n");
		return 0;
	}

	/* load BPF program */
	if (bpf_object__load(obj)) {
		fprintf(stderr, "ERROR: loading BPF object file failed\n");
		goto cleanup;
	}
	
    /* 根据map的名字查询fd，这里的map名字是由tracex4_kern.c中my_map SEC(".map")处设置的 */
	map_fd = bpf_object__find_map_fd_by_name(obj, "my_map");
	if (map_fd < 0) {
		fprintf(stderr, "ERROR: finding a map in obj file failed\n");
		goto cleanup;
	}
	
    /* 循环每一个bpf_program，其实就是循环SEC("kprobe")和SEC("kretprobe") */
	bpf_object__for_each_program(prog, obj) {
		links[j] = bpf_program__attach(prog);	/* attach eBPF的程序 */
		if (libbpf_get_error(links[j])) {
			fprintf(stderr, "ERROR: bpf_program__attach failed\n");
			links[j] = NULL;
			goto cleanup;
		}
		j++;
	}
	
    /* 每秒输出map的内容 */
	for (i = 0; ; i++) {
		print_old_objects(map_fd);
		sleep(1);
	}

cleanup:
	for (j--; j >= 0; j--)
		bpf_link__destroy(links[j]);

	bpf_object__close(obj);
	return 0;
}