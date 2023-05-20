# ecli Dockerfile

目前采用[ecli-dockerfile](https://github.com/eunomia-bpf/eunomia-bpf/blob/master/documents/ecli-dockerfile-usage.md)的内容，并在其基础上进行了修改：

```dockerfile
FROM ubuntu:latest

ENV UBUNTU_SOURCE /etc/apt

COPY ./ /root

WORKDIR /root

ADD sources.list $UBUNTU_SOURCE/

RUN apt-get update && \
    apt-get -y install gcc libelf-dev

#CMD ./ecli run /root/my/package.json
CMD ["/bin/bash"]
```

迷你版框架

```dockerfile
FROM ubuntu:22.04

COPY . /root/.eunomia/bin

ENV PATH="/root/.eunomia/bin:${PATH}"

WORKDIR /code

ENTRYPOINT ["ecli", "run"]
```
