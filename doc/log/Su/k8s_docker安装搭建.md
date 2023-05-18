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

## K8S集群搭建（kubeadm方式）

准备三台Ubuntu虚拟机（一台master两台node），均为2核，内核版本5.15.0

```bash
$ uname -r
5.15.0-1031-aws
```

搭建过程参考[Ubuntu 22.04 搭建K8s集群](https://www.cnblogs.com/way2backend/p/16970506.html)

### 一、网络初始化配置

（以下配置需要在每台机器上进行）

### 1、设置主机名

```bash
$ hostnamectl set-hostname k8s-master		# node为k8s-node01/02
# 查看hostname
$ hostname
k8s-master
```

#### 2、配置hosts

编辑`/etc/hosts`配置ip与hostname直接的匹配关系：

```bash
$ vi /etc/hosts
```

添加三行内容：

```tex
ip hostname

```

之后可以通过ping指令判断是否成功。

#### 3、安装ssh

```bash
$ apt install openssh-server
```

这一步一般都可以忽略。

### 二、系统初始化配置

（以下配置需要在每台机器上执行）

#### 1、关闭防火墙与selinux

Ubuntu22.04默认关闭，不需要设置

```bash
$ ufw disable
```

#### 2、禁用swap分区

```bash
$ sed -i '/ swap / s/^(.*)$/#1/g' /etc/fstab
```

#### 3、系统时间设置

```bash
$ timedatectl set-timezone Asia/Shanghai
#同时使系统日志时间戳也立即生效
$ systemctl restart rsyslog
```

#### 4、修改内核参数

载入如下内核模块

```bash
$ tee /etc/modules-load.d/containerd.conf <<EOF
> overlay
> br_netfilter
> EOF

$ modprobe overlay
$ modprobe br_netfilter
```

配置网络参数

```bash
$ tee /etc/sysctl.d/kubernetes.conf <<EOF
> net.bridge.bridge-nf-call-ip6tables = 1
> net.bridge.bridge-nf-call-iptables = 1
> net.ipv4.ip_forward = 1
> EOF

# 执行如下指令确保修改生效
$ sysctl --system
```

### 三、安装containerd

（以下配置需要在每台机器上执行）

Docker与Kubernetes在控制组件时都会使用containerd

#### 1、安装依赖

```bash
$ apt install -y curl gnupg2 software-properties-common apt-transport-https ca-certificates
```

#### 2、添加docker repo到apt

```bash
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmour -o /etc/apt/trusted.gpg.d/docker.gpg
$ add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
$ apt update
```

#### 3、安装containerd

```bash
$ apt install -y containerd.io
```

#### 4、配置containerd使用systemd作为cgroup

```bash
$ containerd config default | sudo tee /etc/containerd/config.toml >/dev/null 2>&1
$ sed -i 's/SystemdCgroup \= false/SystemdCgroup \= true/g' /etc/containerd/config.toml
```

#### 5、重启并设置开机自启

```bash
$ systemctl restart containerd
$ systemctl enable containerd
```

### 四、安装kube组件

（以下配置需要在每台机器上执行）

安装kubernetes需要的组件：kubelet（控制pod生命周期）、kubeadm（管理kubernetes集群）与kubectl（交互命令行）

#### 1、添加repo到apt

```bash
$ curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
$ apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
$ apt update
```

#### 2、安装kube组件

```bash
$ apt install -y kubelet kubeadm kubectl
# 防止自动安装
$ apt-mark hold kubelet kubeadm kubectl
```

### 五、初始化master节点

（以下配置需要在master节点上执行）

#### 1、kubeadm初始化

在master节点上执行如下命令即可完成初始化：

```bash
$ kubeadm init --control-plane-endpoint=your-ip
# 很多调试信息

Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:	# start

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

Alternatively, if you are the root user, you can run:							# 非root

  export KUBECONFIG=/etc/kubernetes/admin.conf

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

You can now join any number of control-plane nodes by copying certificate authorities
and service account keys on each node and then running the following as root:

  kubeadm join 172.31.22.63:6443 --token g4oadc.afzv2tkc5d8c5ewm \
        --discovery-token-ca-cert-hash sha256:13174d3db9b1756971d94fe0679307e687743aea65554f0a570ee4de02df55d4 \
        --control-plane 

Then you can join any number of worker nodes by running the following on each as root:	# join node

kubeadm join 172.31.22.63:6443 --token g4oadc.afzv2tkc5d8c5ewm \
        --discovery-token-ca-cert-hash sha256:13174d3db9b1756971d94fe0679307e687743aea65554f0a570ee4de02df55d4 
```

#### 2、配置config文件

按照提示信息，进行后续config初始化：

```bash
$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

这时已经可以通过kubectl查看集群状态：

```bash
$ kubectl get nodes
NAME     STATUS     ROLES           AGE    VERSION
master   NotReady   control-plane   110s   v1.27.2
```

### 六、加入node节点

TODO：...