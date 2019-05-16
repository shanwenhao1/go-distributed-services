# consul deploy with kubernetes

@[TOC]
- [脚本方式](#Running-Consul-on-kubernetes)
- [官方方式](#官方方式)
- [参考](#参考)


## Running Consul on kubernetes

### 事先准备
- kubernetes 1.14.x
- 安装[consul1.4.4](consul%20learning.md#官网教程使用-basic-use)
- [cfssl](https://pkg.cfssl.org/)和[cfssljson](https://pkg.cfssl.org/)

```bash
wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
# 赋予可执行权限并移入/usr/local/bin下
chmod +x cfssljson_linux-amd64 cfssl_linux-amd64
mv cfssljson_linux-amd64 /usr/local/bin/cfssljson
mv cfssl_linux-amd64 /usr/local/bin/cfssl
```

- 下载repo
```bash
git clone https://github.com/shanwenhao1/consul-on-kubernetes.git
cd consul-on-kubernetes
```
- 生成TLS 认证(consul member之间使用RPC进行通信), 可根据自己需要对以下文件进行修改
    - [ca-csr.json](../../doc/consul/consul%20kubernetes/ca/ca-csr.json)、
    - [ca-config.json](../../doc/consul/consul%20kubernetes/ca/ca-config.json)、
    - [consul-csr.json](../../doc/consul/consul%20kubernetes/ca/consul-csr.json)
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
GOSSIP_ENCRYPTION_KEY=$(consul keygen)
```
- 生成consul密钥和configmap(包含一些CLI flags, TLS认证和配置文件)
    - 将gossip encryption key and TLS certificates存储在一个密钥内
    ```bash
    kubectl create secret generic consul \
      --from-literal="gossip-encryption-key=${GOSSIP_ENCRYPTION_KEY}" \
      --from-file=ca.pem \
      --from-file=consul.pem \
      --from-file=consul-key.pem
    ```
    - 将consul server配置存入`configs/server.json`中
    ```bash
    kubectl create configmap consul --from-file=configs/server.json
    ```
- 创建consul service
```bash
# 可更改consul.yaml修改启动配置
kubectl create -f services/consul.yaml
```
- 创建consul service账号
```bash
kubectl apply -f serviceaccounts/consul.yaml
kubectl apply -f clusterroles/consul.yaml
```
- 创建consul StatefulSet, 在三个节点上使用以下命令依次创建consul member
```bash
kubectl create -f statefulsets/consul.yaml
# 使用命令验证consul member是否安装完毕
kubectl get nodes
```

**未完待续**: 正常来讲到此consul服务就应该搭建起来了, 但在我本地的环境下, consul集群一直处于
pending状态, 还未找到原因




## 官方方式

当启用kubernetes后, 我们安装consul并部署他们, 然后依靠TLS via使得多节点的consul互相发现对方并组成集群.


### 前言

官方部署是通过kubernetes `Minikube`部署, [文档](https://learn.hashicorp.com/consul/getting-started-k8s/minikube).
我们这里使用本地搭建的[kubernetes集群](../../doc/kubernets.md)进行搭建consul集群学习.

### 部署

- install consul
    - step 1: init helm
    ```bash
    # clone  hashicorp/demo-consul-101
    git clone https://github.com/hashicorp/demo-consul-101.git
    cd demo-consul-101/k8s
    # start install by using `helm` tool, 注意默认是unauthenticated的, 增加参数--tiller-tls-verify使用认证方式
    # 详情: https://docs.helm.sh/using_helm/#securing-your-helm-installation
    # helm 需要外网环境, 可使用helm version 验证pod是否启动, 如未启动则可运行helm init --upgrade再次尝试
    helm init
    ```
    ![](../../doc/picture/consul/helm%20init.png)
    - step 2: install consul with helm
        ```bash
        # 在demo-consul-101/k8s/helm-consul-values.yaml下是官方重写的配置
        # (详细的参数代表信息可在consul-helm/values.yaml里面看到)
        # 适合本地环境.    注意: 目录不变
        git clone https://github.com/hashicorp/consul-helm.git
         
        # 使用helm 一起安装我们重载的`helm-consul-values.yaml`和`consul-helm`文件
        # --name如果不提供, the chart会随机一个名称for the installation
      
        helm install -f helm-consul-values.yaml --name hedgehog ./consul-helm
        ```
        - 出现以下错误时
            ```bash
            Error: release hedgehog failed: namespaces "default" is forbidden: User "system:serviceaccount:kube-system:default" 
            cannot get resource "namespaces" in API group "" in the namespace "default"
            ```
            - 错误解决:(这里采用的是第二种方法Helm权限配置)
                - 官方方法(请仔细阅读文档), [stack over flow](https://stackoverflow.com/questions/47973570/kubernetes-log-user-systemserviceaccountdefaultdefault-cannot-get-services)、
                [官方RBAC Authorization文档](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
                ```bash
                # 创建一个role.yaml文件, 将以下内容写入文件
                    kind: ClusterRole
                    apiVersion: rbac.authorization.k8s.io/v1
                    metadata:
                      namespace: default
                      name: service-reader
                    rules:
                    - apiGroups: [""] # "" indicates the core API group
                      resources: ["services"]
                      verbs: ["get", "watch", "list"]
                # 启用该策略
                  kubectl apply -f role.yaml
                #  使用该cluster role创建一个clusterRoleBinding
                     kubectl create clusterrolebinding service-reader-pod \
                        --clusterrole=service-reader  \
                        --serviceaccount=default:default
                # 创建的clusterRoleBinding类似于    
                    apiVersion: rbac.authorization.k8s.io/v1
                    # This cluster role binding allows anyone in the "manager" group to read secrets in any namespace.
                    kind: ClusterRoleBinding
                    metadata:
                      name: read-service-global
                    subjects:
                    - kind: Group
                      name: manager # Name is case sensitive
                      apiGroup: rbac.authorization.k8s.io
                    roleRef:
                      kind: ClusterRole
                      name: service-reader
                      apiGroup: rbac.authorization.k8s.io
                ```
                - Helm权限配置, [参考](http://www.libaibai.net/node/357)
                ```bash
                 kubectl create serviceaccount --namespace kube-system tiller
                 kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
                 kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'
                ```
- 部署Consul-aware application






















## 参考
- 脚本方式
    - [consul-on-kubernetes](https://github.com/shanwenhao1/consul-on-kubernetes.git): fork from 
    [kelseyhightower](https://github.com/kelseyhightower/consul-on-kubernetes.git)
- 官方方式
- [consul-k8s](https://github.com/hashicorp/consul-k8s)










### 前言

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