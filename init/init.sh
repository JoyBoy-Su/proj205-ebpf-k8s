#!/bin/bash
# bpf file
cd ~/.kube
# create bpf home
if [ ! -d  bpf  ];then
  mkdir bpf
fi
cd bpf
# create instances and packages
if [ ! -d  packages  ];then
  mkdir packages
fi
if [ ! -d  instances  ];then
  mkdir instances
fi

# bpf namespace
kubectl create namespace bpf

