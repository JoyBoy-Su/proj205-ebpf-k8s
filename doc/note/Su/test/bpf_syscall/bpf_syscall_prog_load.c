/**
 * @file bpf_syscall_prog_load.c
 * @author JiadiSu (20302010043@fudan.edu.cn)
 * @brief 通过bpf系统调用接口加载bpf prog
 * @version 0.1
 * @date 2023-04-24
 * 
 * @copyright Copyright (c) 2023
 * 编译：gcc bpf_syscall_prog_load.c -o bpf_syscall_prog_load
 * （原生的gcc编译说明不需要任何其他依赖）
 * 运行环境：Ubuntu20.04（不要WSL！不要WSL！不要WSL！）
 */

#include <errno.h>
#include <linux/bpf.h>      /* bpf相关的定义，如map type，cmd type等，和下载下的源码linux6.0/include/uapi/linux/bpf.h一致 */
#include <stdio.h>
#include <stdlib.h>         /* exit()函数 */
#include <stdint.h>         /* uint64_t等type */
#include <sys/syscall.h>    /* SYS_xxx (SYS_bpf) */
#include <unistd.h>         /* syscall() */

#define ptr_to_u64(x) ((uint64_t)x)     /* ptr -> uint64_t */

/**
 * @brief 对bpf系统调用的一个包装，一切对bpf子系统的操作都通过这个函数交互
 * （这个封装是自定义的而不是系统的）
 * （似乎是因为linux并没有一个像read和write这种封装好的bpf系统调用）
 * @param cmd 
 * @param attr 
 * @param size 
 * @return int 
 */
int bpf(enum bpf_cmd cmd, union bpf_attr* attr, unsigned int size)
{
    return syscall(SYS_bpf, cmd, attr, size);
}

#define LOG_BUF_SIZE 0x1000
char bpf_log_buf[LOG_BUF_SIZE];

/**
 * @brief 对bpf(BPF_PROG_LOAD)的一个自定义封装，完成bpf_prog的加载
 * 
 * @param type 
 * @param insns 
 * @param insn_cnt 
 * @param license 
 * @return int 
 */
int bpf_prog_load(enum bpf_prog_type type, const struct bpf_insn* insns, int insn_cnt, char* license)
{
    union bpf_attr attr = 
    {
        .prog_type = type,
        .insns = ptr_to_u64(insns),
        .insn_cnt = insn_cnt,
        .license = ptr_to_u64(license),
        .log_buf = ptr_to_u64(bpf_log_buf),
        .log_level = 2,
        .log_size = LOG_BUF_SIZE
    };

    return bpf(BPF_PROG_LOAD, &attr, sizeof(attr));
}

int main(int argc, char const *argv[])
{
    /* bpf prog */
    struct bpf_insn bpf_prog[] = {
        { 0xb7, 0, 0, 0, 0x2 },     /* mov r0, 0x2; */
        { 0x95, 0, 0, 0, 0x0 },     /* exit; */
    };

    /* 加载prog */
    int prog_fd;
    if ((prog_fd = bpf_prog_load(BPF_PROG_TYPE_SOCKET_FILTER, bpf_prog, sizeof(bpf_prog) / sizeof(bpf_prog[0]), "GPL")) < 0)
    {
        perror("BPF prog load error");
        exit(-1);
    }
    printf("BPF prog fd: %d\n", prog_fd);

    /* 读取log */
    printf("%s\n", bpf_log_buf);

    exit(0);
}

