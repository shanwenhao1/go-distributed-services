# consul deploy with kubernetes

@[TOC]
- [Running Consul on kubernetes](#Running-Consul-on-kubernetes)
- []()
- [参考](#参考)


## Running Consul on kubernetes

### 事先准备
- 安装[consul](consul%20learning.md#官网教程使用-basic-use)
- [cfssl](https://pkg.cfssl.org/)和[cfssljson](https://pkg.cfssl.org/)

```bash
wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
wget https://pkg.cfssl.org/R1.1/cfssljson_linux-amd64
# 赋予可执行权限并移入/usr/local/bin下
chmod +x cfssljson_linux-amd64 cfssl_linux-amd64
mv cfssljson_linux-amd64 /usr/local/bin/cfssljson
mv cfssl_linux-amd64 /usr/local/bin/cfssl
```

- 下载repo
```bash
git clone https://github.com/kelseyhightower/consul-on-kubernetes.git
cd consul-on-kubernetes
```
- 生成TLS 认证(consul member之间使用RPC进行通信), [ca-csr.json](../../doc/consul/consul%20kubernetes/ca/ca-csr.json)、
[ca-config.json](../../doc/consul/consul%20kubernetes/ca/ca-config.json)、
[consul-csr.json](../../doc/consul/consul%20kubernetes/ca/consul-csr.json)
需要修改下
    ```bash
    # 生成CA证书
    cfssl gencert -initca ca/ca-csr.json | cfssljson -bare ca
    # 生成TLS认证和密钥
    cfssl gencert \
      -ca=ca.pem \
      -ca-key=ca-key.pem \
      -config=ca/ca-config.json \
      -profile=default \
      ca/consul-csr.json | cfssljson -bare consul
    ```
    会生成以下四个文件
    ![](../../doc/picture/consul/consul%20k8s%20ca.png)
- 生成consul Gossip Encryption key用于consul member之间的通信, 
[Gossip in consul](https://www.consul.io/docs/internals/gossip.html)
```bash
# 要安装consul ===================================
GOSSIP_ENCRYPTION_KEY=$(consul keygen)
```


## 参考
- [consul-on-kubernetes](https://github.com/shanwenhao1/consul-on-kubernetes.git): fork from 
[kelseyhightower](https://github.com/kelseyhightower/consul-on-kubernetes.git)




















@[TOC]
- [前言](#前言)
- [consul cluster](#consul-cluster)

## 前言

官方部署是通过kubernetes `Minikube`部署, [文档](https://learn.hashicorp.com/consul/getting-started-k8s/minikube).
我们这里使用本地搭建的[kubernetes集群](../../doc/kubernets.md)进行搭建consul集群学习.

[Helm Chart官方文档](https://www.consul.io/docs/platform/k8s/helm.html)

使用前需要先安装kubernetes应用管理工具Helm, [安装](https://www.jianshu.com/p/ab26b5762cf5)
- step-1:
    ```bash
    # 使用官方脚本安装最新版本
    curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh
    chmod 700 get_helm.sh
    ./get_helm.sh
    helm init
    # 更新至最新版本(第一次helm init不用)
    helm init --upgrade
    # 更新charts列表
    helm repo update
    ```
- step-2: 为Tilier添加权限. [参考](https://helm.sh/docs/using_helm/#role-based-access-control)
    - 创建`rbac-config.yaml`[文件](../../doc/kubernetes/install/helm.sh) 
    ```bash
    kubectl create -f rbac-config.yaml
    ```
 

## Use Cases
`Running a consul server cluster`: Consul cluster可直接运行在kubernetes集群上并且可随时扩展

`Running Consul clients`: 作为pods运行在kubernetes上, 内部包含很多consul 工具

`Service sync to enable Kubernetes and non-Kubernetes services to communicate`: 这允许k8s services利用native service
去连接外部service, 外部services 可以使用consul 服务发现去连接k8s services

`Automatic encryption and authorization for pod network connections`(pod网络自动加密和授权): 使用TLS  

## Running Consul
[Helm chart](https://www.consul.io/docs/platform/k8s/helm.html)通过命令行的方式运行consul和kubernetes.
当Helm chart自动导出了所需配置后, 它并不会自动操作consul, 还是需要你手动操作备份升级等操作.

**注意: Helm chart默认配置的consul cluster是insecure的**

```bash
git clone https://github.com/hashicorp/consul-helm.git
cd consul-helm
# 切换至最新release版本
git checkout v0.8.1
# run helm, 几分钟后一个consul cluster就建立起来了, 可用helm status consul查看
helm install --name consul ./
# 通过查看UI检测consul cluster是否建立成功(因为Helm chart默认会建立Consul UI)
# 但由于默认UI并不暴露端口, 因此
kubectl port-forward consul-server-0 8500:8500
```


## consul cluster


```bash
# 下载consul Helm Chart
git clone https://github.com/hashicorp/demo-consul-101.git
cd demo-consul-101/k8s
# 下载consul-helm
git clone https://github.com/hashicorp/consul-helm.git

```


## 参考
- [consul-k8s](https://github.com/hashicorp/consul-k8s)