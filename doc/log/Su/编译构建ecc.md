# ���빹��ecc

����Դ�빹��ecc�Ĺ��̣�Ubuntu

### ��װ����

```bash
# ��װ����
$ sudo apt-get update
# �ȴ���װ������ɣ�û��git�Ļ���Ҫ��װgit
$ sudo apt install gcc clang libelf1 libelf-dev zlib1g-dev llvm cargo
```

### ���빹��

##### clone�ֿ�

```bash
$ git clone --recursive https://github.com/eunomia-bpf/eunomia-bpf.git	# ����������ģ��
Cloning into 'eunomia-bpf'...
remote: Enumerating objects: 7407, done.
remote: Counting objects: 100% (957/957), done.
remote: Compressing objects: 100% (346/346), done.
remote: Total 7407 (delta 660), reused 786 (delta 591), pack-reused 6450
Receiving objects: 100% (7407/7407), 15.86 MiB | 26.03 MiB/s, done.
Resolving deltas: 100% (4367/4367), done.
Submodule 'third_party/bpftool' (https://github.com/eunomia-bpf/bpftool) registered for path 'third_party/bpftool'
Submodule 'third_party/vmlinux' (git@github.com:eunomia-bpf/vmlinux.git) registered for path 'third_party/vmlinux'
Cloning into '/home/ubuntu/test/eunomia-bpf/third_party/bpftool'...
remote: Enumerating objects: 1610, done.        
remote: Total 1610 (delta 0), reused 0 (delta 0), pack-reused 1610        
Receiving objects: 100% (1610/1610), 732.92 KiB | 13.33 MiB/s, done.
Resolving deltas: 100% (952/952), done.
Cloning into '/home/ubuntu/test/eunomia-bpf/third_party/vmlinux'...
remote: Enumerating objects: 18, done.        
remote: Counting objects: 100% (18/18), done.        
remote: Compressing objects: 100% (8/8), done.        
remote: Total 18 (delta 4), reused 18 (delta 4), pack-reused 0        
Receiving objects: 100% (18/18), 1.46 MiB | 22.71 MiB/s, done.
Resolving deltas: 100% (4/4), done.
Submodule path 'third_party/bpftool': checked out '05940344f5db18d0cb1bc1c42e628f132bc93123'
Submodule 'libbpf' (https://github.com/libbpf/libbpf.git) registered for path 'third_party/bpftool/libbpf'
Cloning into '/home/ubuntu/test/eunomia-bpf/third_party/bpftool/libbpf'...
remote: Enumerating objects: 10593, done.        
remote: Counting objects: 100% (1239/1239), done.        
remote: Compressing objects: 100% (394/394), done.        
remote: Total 10593 (delta 855), reused 906 (delta 825), pack-reused 9354        
Receiving objects: 100% (10593/10593), 9.09 MiB | 22.22 MiB/s, done.
Resolving deltas: 100% (7100/7100), done.
Submodule path 'third_party/bpftool/libbpf': checked out 'e3a40329bb05a333fc588e3bf50365a554fda0a6'
Submodule path 'third_party/vmlinux': checked out '933f83becb45f5586ed5fd089e60d382aeefb409'
```

##### make����

make���ʱ�ϳ����������Ļ���һ�����ļ��Ҳ�����������ʾ��װ���ɣ�����submodule gcc cargo��

```bash
# ���밲װ
$ cd compiler
$ make
```

makeʱ�����������⣺

```bash
$ make # ����
error: failed to run custom build command for `clang-sys v1.4.0`

Caused by:
  process didn't exit successfully: `/home/ubuntu/test/eunomia-bpf/compiler/cmd/target/release/build/clang-sys-0d97fb534e5efb11/build-script-build` (exit status: 101)
  --- stderr
  thread 'main' panicked at 'called `Result::unwrap()` on an `Err` value: "couldn't find any valid shared libraries matching: ['libclang.so', 'libclang-*.so'], set the `LIBCLANG_PATH` environment variable to a path where one of these files can be found (invalid: [])"', /home/ubuntu/.cargo/registry/src/github.com-1ecc6299db9ec823/clang-sys-1.4.0/build/dynamic.rs:211:45
  note: run with `RUST_BACKTRACE=1` environment variable to display a backtrace
make: *** [Makefile:66: cmd/target/release/ecc-rs] Error 101
```

��ȡ�ؼ���Ϣ��couldn't find any valid shared libraries matching: ['libclang.so', 'libclang-*.so']��Ҳ������`/usr/lib`��ȱ��.so�ļ���

```bash
# ȥ��һ��lib ��Ȼû��
$ ls -l /usr/lib/ | grep clang
drwxr-xr-x  4 root root   4096 May 17 03:30 clang
```

���Ѿ���װ��clang��Ӧ��Ҳ����libclang-\*.so�ģ��������ң���`x86_64-linux-gnu`�¿����ҵ���Ӧ�汾��libclang-\*.so�ļ���

```bash
# �����ң���װclangʱ��.so��x86_64-linux-gnu��
$ ls -l /usr/lib/x86_64-linux-gnu/ | grep clang
lrwxrwxrwx  1 root root        17 Mar 24  2022 libclang-14.so.1 -> libclang-14.so.13
lrwxrwxrwx  1 root root        21 Mar 24  2022 libclang-14.so.13 -> libclang-14.so.14.0.0
-rw-r--r--  1 root root  30580864 Mar 24  2022 libclang-14.so.14.0.0
lrwxrwxrwx  1 root root        33 Mar 24  2022 libclang-cpp.so.14 -> ../llvm-14/lib/libclang-cpp.so.14
```

������`/usr/lib`�´��������Ӽ��ɣ�

```bash
# ����������
$ sudo ln -s /usr/lib/x86_64-linux-gnu/libclang-14.so.14.0.0 /usr/lib/libclang-14.so
# ���lib
$ ls -l /usr/lib | grep clang
drwxr-xr-x  4 root root   4096 May 17 03:30 clang
lrwxrwxrwx  1 root root     47 May 17 04:28 libclang-14.so -> /usr/lib/x86_64-linux-gnu/libclang-14.so.14.0.0
```

��ʱ����libclang�⣬���±��룬makeͨ����

```bash
# ����libclang-14.so������make
$ make
...
Finished dev [unoptimized + debuginfo] target(s) in 2m 57s		# ����ɹ�
```

##### make��װ

���ڲ������ļ�·�����ӣ�������ʹ�ã�ͨ��make install�޸���һ�㣨make install��ʵ����ִ������������cpָ���

```bash
$ make install
rm -rf ~/.eunomia && cp -r workspace ~/.eunomia					# ����~/.eunomia
# ���Լ򵥿�һ��Ŀ¼
$ ls -l ~/.eunomia/
total 8
drwxrwxr-x 2 ubuntu ubuntu 4096 May 17 05:03 bin
drwxrwxr-x 4 ubuntu ubuntu 4096 May 17 05:03 include
$ ls -l ~/.eunomia/bin/
total 31520
-rwxrwxr-x 1 ubuntu ubuntu  1799208 May 17 05:03 bpftool
-rwxrwxr-x 1 ubuntu ubuntu 30472288 May 17 05:03 ecc-rs			# compiler
```

##### ���û�������

Ϊ��֧��������Ŀ¼�¿���ִ��ecc-rs��������ӵ���������PATH�У�

```bash
# ��ӻ�������
$ export PATH=$PATH:~/.eunomia/bin								# ��ʱ��ӣ����������/etc/profile��
```

##### ����

```bash
# ����һ��
$ ecc-rs -h
eunomia-bpf compiler

Usage: ecc-rs [OPTIONS] <SOURCE_PATH> [EXPORT_EVENT_HEADER]

Arguments:
  <SOURCE_PATH>          path of the bpf.c file to compile
  [EXPORT_EVENT_HEADER]  path of the bpf.h header for defining event struct [default: ]

Options:
  -o, --output-path <OUTPUT_PATH>
          path of output bpf object [default: ]
  -w, --workspace-path <WORKSPACE_PATH>
          custom workspace path
  -a, --additional-cflags <ADDITIONAL_CFLAGS>
          additional c flags for clang [default: ]
  -c, --clang-bin <CLANG_BIN>
          path of clang binary [default: clang]
  -l, --llvm-strip-bin <LLVM_STRIP_BIN>
          path of llvm strip binary [default: llvm-strip]
  -s, --subskeleton
          do not pack bpf object in config file
  -v, --verbose
          print the command execution
  -y, --yaml
          output config skel file in yaml
      --header-only
          generate a bpf object for struct definition in header file only
      --wasm-header
          generate wasm include header
  -b, --btfgen
          fetch custom btfhub archive file
      --btfhub-archive <BTFHUB_ARCHIVE>
          directory to save btfhub archive file [default: /home/ubuntu/.eunomia/btfhub-archive]
  -h, --help
          Print help (see more with '--help')
  -V, --version
          Print version
# ����һ��opensnoop.bpf.c��һ��
$ ecc-rs opensnoop.bpf.c opensnoop.h 
Compiling bpf object...
warning: text is not json: Process ID to trace use it as a string
warning: text is not json: Thread ID to trace use it as a string
warning: text is not json: User ID to trace use it as a string
warning: text is not json: trace only failed events use it as a string
warning: text is not json: Trace open family syscalls. use it as a string
Generating export types...
Packing ebpf object and config into package.json...
$ ls -l
total 36
-rw-rw-r-- 1 ubuntu ubuntu  3702 May 17 05:09 opensnoop.bpf.c
-rw-rw-r-- 1 ubuntu ubuntu 12928 May 17 05:10 opensnoop.bpf.o
-rw-rw-r-- 1 ubuntu ubuntu   427 May 17 05:10 opensnoop.h
-rw-rw-r-- 1 ubuntu ubuntu  1525 May 17 05:10 opensnoop.skel.json
-rw-rw-r-- 1 ubuntu ubuntu  5954 May 17 05:10 package.json			# �ɹ�����package.json
```

������json��ecli���к󣬳ɹ����bpf��������compiler�����óɹ���