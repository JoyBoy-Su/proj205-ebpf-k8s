/**
 * @file hello_user.c
 * @author JiadiSu (20302010043@fudan.edu.cn)
 * @brief 通过libbpf接口实现一个hello world
 * @version 0.1
 * @date 2023-04-25
 * 
 * @copyright Copyright (c) 2023
 * 编译：gcc hello_user.c -o hello -lbpf
 */

#include <bpf/bpf.h>            /* 一些对bpf syscall的封装 */
#include <bpf/libbpf.h>         /* 一些bpf_progam和bpf_object相关的函数 */
#include <fcntl.h>
#include <unistd.h>

int main(int argc, char const *argv[])
{
    struct bpf_link* link;
	struct bpf_program *prog;
	struct bpf_object *obj;
	char filename[256];

	snprintf(filename, sizeof(filename), "%s_kern.o", argv[0]);
    
    /* 声明在<bpf/libbpf.h>，打开filename对应的eBPF */
	obj = bpf_object__open_file(filename, NULL);
	if (libbpf_get_error(obj)) {
		fprintf(stderr, "ERROR: opening BPF object file failed\n");
		return 0;
	}

	/* 声明在<bpf/libbpf.h>，加载到内核 */
	if (bpf_object__load(obj)) {
		fprintf(stderr, "ERROR: loading BPF object file failed\n");
		goto cleanup;
	}
	
    /* 循环每一个bpf_program，其实就是循环SEC("tp/syscalls/sys_enter_open") */
	bpf_object__for_each_program(prog, obj) {
		link = bpf_program__attach(prog);	/* attach eBPF的程序 */
		if (libbpf_get_error(link)) {
			fprintf(stderr, "ERROR: bpf_program__attach failed\n");
			link = NULL;
			goto cleanup;
		}
	}

    /* 死循环不退出 */
	while (1) {}

cleanup:
    bpf_link__destroy(link);
	bpf_object__close(obj);
	return 0;
}

