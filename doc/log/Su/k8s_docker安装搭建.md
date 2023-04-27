# k8s与docker安装搭建

## docker的安装

### 安装过程

1、删除旧版本

```bash
$ sudo apt-get remove docker docker-engine docker.io containerd runc
```

2、设置仓库

更新apt

```bash
$ sudo apt-get update
```

安装依赖包

```bash
$ sudo apt-get install apt-transport-https ca-certificates curl gnupg-agent software-properties-common
```

添加Docker官方密钥

```bash
$ curl -fsSL https://mirrors.ustc.edu.cn/docker-ce/linux/ubuntu/gpg | sudo apt-key add -
OK
```

验证密钥

```bash
$ apt-key fingerprint 0EBFCD88
pub   rsa4096 2017-02-22 [SCEA]
      9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88
uid           [ unknown] Docker Release (CE deb) <docker@docker.com>
sub   rsa4096 2017-02-22 [S]

```

设置稳定版仓库

```bash
$ add-apt-repository \
   "deb [arch=amd64] https://mirrors.ustc.edu.cn/docker-ce/linux/ubuntu/ \
  $(lsb_release -cs) \
  stable"
```

3、安装docker（最新版本）

```bash
$ sudo apt-get install docker-ce docker-ce-cli containerd.io
```

### 遇到的问题

"Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"

原因是docker未启动，使用如下指令启动

```bash
$ systemctl start docker
```

报错："System has not been booted with systemd as init system (PID 1). Can't operate."

原因是不是用`systemctl`管理服务，修改指令：

```bash
$ service docker start
 * Starting Docker: docker           [ OK ]
```

启动后尝试运行hello-world容器：

```bash
$ sudo docker run hello-world

Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
1b930d010525: Pull complete                                                                                                                                  Digest: sha256:c3b4ada4687bbaa170745b3e4dd8ac3f194ca95b2d0518b417fb47e5879d9b5f
Status: Downloaded newer image for hello-world:latest


Hello from Docker!
This message shows that your installation appears to be working correctly.


To generate this message, Docker took the following steps:
 1. The Docker client contacted the Docker daemon.
 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (amd64)
 3. The Docker daemon created a new container from that image which runs the
    executable that produces the output you are currently reading.
 4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.


To try something more ambitious, you can run an Ubuntu container with:
 $ docker run -it ubuntu bash


Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/


For more examples and ideas, visit:
 https://docs.docker.com/get-started/
```

出现如上结果，安装成功。