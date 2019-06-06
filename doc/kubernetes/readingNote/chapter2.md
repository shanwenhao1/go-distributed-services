# Kubernetes实践指南

## k8s安装与配置
[详情](../../../doc/Kubernetes权威指南.pdf)第二章

主要是步骤为: 
- etcd服务安装
- k8s集群安装以及开启一些认证

## k8s版本升级

为了保证在当前集群中运行的容器不受影响, 对集群中的Node逐个进行隔离, 等待其上运行的容器
全部执行完成, 再更新该Node上的kubelet和kube-proxy服务, 将全部Node更新完成后, 最后更新Master服务

[使用自定义Docker镜像库](../../../doc/Kubernetes权威指南.pdf)第二章`2.1.5`

## 集群网络配置

## Pod

### 生命周期和重启策略
有以下重启策略:
- Always
- OnFailure
- Never

### pod的扩容和缩容
可利用RC的Scale机制完成pod的扩容和缩容
```bash
# 将redis-slave RC的pod副本数量由2更新为3
kubectl scale rc redis-slave --replicas=3
```
此外Horizontal Pod Autoscaler(HPA)的控制器可实现基于CPU使用率进行自动Pod扩容缩容功能.
kube-controller-manager服务启动参数--horizontal-pod-autoscaler-sync-period定义的时长(默认30秒)
周期性地监测pod的CPU使用率.
```bash
# 查看已创建的HPA
kubectl get hpa
```

### pod的滚动升级
为了保证服务可靠性, 需升级的服务采用滚动升级的方式. 使用命令`kubectl rolling-update`新建了一个RC
然后逐步控制旧的RC中的pod减少至0, 同时新的RC中的pod逐步增加至指定目标值.



## Service
Service 是k8s最核心的概念, 通过创建Service, 可以为一组具有相同功能的容器提供一个统一的入口地址, 并将请求进行负载
分发到后端的各个容器应用上.

引入的原因:
一般来说, 容器应用最简便的方式是通过TCP/IP机制及监听端口号来实现. 
- 例如: 定义一个提供web服务的RC, 由两个Tomcat容器副本组成, 这两个Pod自动生成的IP就可查到了. 可直接
通过Pod的IP地址和端口号访问容器应用, 但由于Pod的地址是不可靠的(比如pod宕掉重新调度时, IP就变了). 
- 加上容器应用如果是分布式的部署方式, 则需要负载均衡来进行分发请求.

```bash
# 最好使用配置文件定义Service
# 使用kubectl expose命令创建service, rcName是我们创建的RC的名称
kubectl expose rc rcName
```
接下来即可通过Service的IP和端口号访问该Service

### 负载均衡
两种负载分发策略: 
- RoundRobin(默认): 轮询模式, 即轮询将请求转发至后端的各个Pod上
- SessionAffinity: 基于客户端IP进行会话保持的模式, 即第一次将某个客户端发起的请求转发至后端的某个Pod上, 之后相同客户
端请求都发至相同的Pod上. 将service.spec.sessionAffinity设置为"ClientIP"可启用

### 集群外部访问Pod或Service
- 将Pod端口号映射到物理机
    - 通过设置容器级别的hostPort, 可将容器应用的端口号映射到物理机上.
    - 或者设置Pod级别的hostNetwork=true, 将Pod中所有容器的端口号直接映射到物理机上.(默认hostPort等于containerPort, 
    如果指定了hostPort, 则hostPort必须等于containerPort的值)
- 将Service的端口号映射到物理机
    - 通过设置nodePort映射到物理机, 同时设置Service的类型为NodePort
    - 通过设置LoadBalancer映射到云服务商提供的LoadBalancer地址. 这种用法仅用于在公有云服务提供商的云平台上
    设置Service的场景.

### DNS服务搭建
k8s提供的虚拟DNS服务名为skydns, 由一个RC和一个Service的定义组成. 详情P167

### Ingress 7层路由
