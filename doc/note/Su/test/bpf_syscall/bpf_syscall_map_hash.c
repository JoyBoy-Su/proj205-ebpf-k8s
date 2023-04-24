/**
 * @file bpf_syscall_map_hash.c
 * @author JiadiSu (20302010043@fudan.edu.cn)
 * @brief 通过bpf系统调用接口操作map_hash
 * @version 0.1
 * @date 2023-04-24
 * 
 * @copyright Copyright (c) 2023
 * 编译：gcc bpf_syscall_map_hash.c -o bpf_syscall_map_hash
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

/**
 * @brief 对bpf(BPF_MAP_CREATE)的一个自定义封装，完成map的创建
 * 
 * @param map_type 
 * @param key_size 
 * @param value_size 
 * @param max_entries 
 * @return int 
 */
int bpf_create_map(enum bpf_map_type map_type, unsigned int key_size, unsigned int value_size, unsigned int max_entries)
{
    union bpf_attr attr = 
    {
        .map_type = map_type,
        .key_size = key_size,
        .value_size = value_size,
        .max_entries = max_entries
    };
    
    return bpf(BPF_MAP_CREATE, &attr, sizeof(attr));
}

/**
 * @brief 对bpf(BPF_MAP_UPDATE_ELEM)的一个自定义封装，完成map数据的更新与创建
 * 
 * @param fd        map fd
 * @param key       key     (const pointer)
 * @param value     value   (const pointer)
 * @param flags     flags   (ANY / EXIST / NOT EXIST)
 * @return int 
 */
int bpf_update_elem(int fd, const void* key, const void* value, uint64_t flags) 
{
    union bpf_attr attr = 
    {
        .map_fd = fd,
        .key = ptr_to_u64(key),
        .value = ptr_to_u64(value),
        .flags = flags
    };

    return bpf(BPF_MAP_UPDATE_ELEM, &attr, sizeof(attr));
}

/**
 * @brief 对bpf(BPF_MAP_LOOKUP_ELEM)的一个自定义封装，完成对数据的查找
 * 
 * @param fd        map fd
 * @param key       key     (const pointer)
 * @param value     value   (pointer)
 * @return int      
 */
int bpf_lookup_elem(int fd, const void* key, void* value)
{
    union bpf_attr attr = 
    {
        .map_fd = fd,
        .key = ptr_to_u64(key),
        .value = ptr_to_u64(value)
    };

    return bpf(BPF_MAP_LOOKUP_ELEM, &attr, sizeof(attr));
}

int main(int argc, char const *argv[])
{
    /* 创建一个 hash map，key为int，value为char* */
    int map_fd;
    if ((map_fd = bpf_create_map(BPF_MAP_TYPE_ARRAY, sizeof(int), sizeof(char*), 0x100)) < 0)
    {
        perror("BPF create map error");
        exit(-1);
    }
    printf("BPF map fd: %d\n", map_fd);

    /* 填充map的值 */
    char *strtab[] = {
        "This",
        "is",
        "eBPF",
        "hash",
        "map",
        "test",
    };
    for (int i = 0; i < 6; i++) 
    {
        char* value = strtab[i];
        /* 因为array的map是预定义max entries个项，因此都是exist的 */
        if (bpf_update_elem(map_fd, &i, &value, BPF_EXIST) < 0)
        {
            perror("BPF update create error");
            exit(-1);
        }
    }

    /* 查询map的值 */
    int key;
    char* value;
    printf("please input key to lookup:");
    scanf("%d", &key);

    if (bpf_lookup_elem(map_fd, &key, &value) < 0) 
    {
        perror("BPF lookup create error");
        exit(-1);
    }
    printf("BPF array map: key %d => value %s\n", key, value);

    exit(0);
}
