# Readme

`doc/note/xxx`目录中存放队员xxx的笔记，

该Readme文档的存在是为了保证该目录并非空目录，从而可以由git提交到托管网站上，创建好自己记录笔记的文档（保证该目录不是空目录）后可删除该readme。

[基于 ebpf 和 vxlan 实现一个 k8s 网络插件（一）](https://zhuanlan.zhihu.com/p/565254116)

[k8s管理员要懂eBPF](https://www.jianshu.com/p/3e21bb174445)

[calico](https://www.tigera.io/blog/introducing-the-calico-ebpf-dataplane/)

[基于 eBPF 的 Kubernetes 问题排查全景图发布](https://developer.aliyun.com/article/879258)

![](https://pic2.zhimg.com/80/v2-c1ea9720cb0892a5b02ce3e91556d94d_1440w.webp)

1. [在linux上安装kubectl](http://kubernetes.p2hp.com/docs/tasks/tools/install-kubectl-linux.html)

   1. ```bash
      ~#：uname -i    
      aarch64
      //平台不一样，需要将amd64换成arm64
      ```

   2. curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubectl"

   3. curl -LO "https://dl.k8s.io/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubectl.sha256"

   4. 安装补全工具失败 type _init_completion:放弃

   > 下载 Google Cloud 公开签名秘钥：访问外网失败

2. 安装 kubeadm、kubelet 

   ```bash
   curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubeadm"
   curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubeadm.sha256"
   echo "$(cat kubeadm.sha256)  kubeadm" | sha256sum --check
    install -o root -g root -m 0755 kubeadm /usr/local/bin/kubeadm
   ```

   ```bash
   curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubelet"
   curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubelet.sha256"
   echo "$(cat kubelet.sha256)  kubelet" | sha256sum --check
    install -o root -g root -m 0755 kubelet /usr/local/bin/kubelet
   ```

   

3. apt安装iptables  apparmor docker-ce-rootless-extras  pigz

   1. [安装yum](https://mirrors.tuna.tsinghua.edu.cn/help/ubuntu/)放弃。
   2.  注意架构，我们是[arm64](https://mirrors.tuna.tsinghua.edu.cn/help/ubuntu-ports/)，[不能用amd64/32](https://mirrors.tuna.tsinghua.edu.cn/help/ubuntu/)

   ```
   deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal main restricted universe multiverse
   # deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal main restricted universe multiverse
   deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-updates main restricted universe multiverse
   # deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-updates main restricted universe multiverse
   deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-backports main restricted universe multiverse
   # deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-backports main restricted universe multiverse
   
   # deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-security main restricted universe multiverse
   # # deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu-ports/ focal-security main restricted universe multiverse
   
   deb http://ports.ubuntu.com/ubuntu-ports/ focal-security main restricted universe multiverse
   ```

4. [minikube](https://minikube.sigs.k8s.io/docs/start/) 安装

   1. 安装

      ```
      curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-arm64
      sudo install minikube-linux-arm64 /usr/local/bin/minikube
      ```

      

   2. 安装docker，使用docker为dirver会出现，结果失败，更换driver

      原因：

      ```
      sudo systemctl restart docker会出错，可以让minikude重启docker使用sudo service docker start，但是不会修改
      ```

      > https://blog.csdn.net/weixin_39266845/article/details/127327423 修改daemon的类型失败
      >
      > https://blog.csdn.net/u012833399/article/details/128533933

   

   1. 根据文章创建用户

      [Exiting due to DRV_AS_ROOT: The "docker" driver should not be used with root privileges.](https://github.com/kubernetes/minikube/issues/7903)

   2. 使用qumu为driver

      ```
      sudo chown root:kvm /dev/kvm
      sudo chmod 666 /dev/kvm
      //可能第一个不成功，我是第二个成功了
      ```

      

   3. minikube start --driver=qemu

      ![](https://p.ipic.vip/uuc4y6.png)

      

