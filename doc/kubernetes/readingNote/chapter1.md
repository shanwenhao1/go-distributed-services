# kubernetes入门
                                                 
## kubernetes架构导读

- services是分布式集群架构的核心
- Kubernetes 基本概念和术语: k8s集群中的机器划分为一个Master节点和一群工作节点(Node)
    - `Master`通常占用一个独立的机器(太重要, 不允许宕机). Master节点上运行以下关键进程
        - kube-apiserver: k8s内所有资源控制的唯一入口
        - kube-controller-manager: k8s所有资源对象的自动化控制中心
        - kube-scheduler: 负责资源调度(Pod调度)
        - etcd: 高可用的key-value store 负责存储k8s集群数据
    - `Node`节点是k8s集群中的工作负载节点, 每个Node都被Master分配一些工作负载(Docker容器等). Node运行着以下关键进程
        - kubelet: 负责pod对应的容器的创建、启停等, 同时与Master节点密切协作, 实现集群管理的基本功能.
        - kube-proxy: 实现K8s service的通信与负载均衡机智的重要组件.
        - Docker Engine: docker引擎, 负责本机的容器创建和管理工作.
    - `Pod`对象, 将服务进程包装至相应的Pod中, 使其成为Pod中运行的一个container. Pod也是k8s管理运行的最小运行单元
        - 每个pod都有自己的`Label`
        - 每个Pod中都运行了名为Pause的"根容器", 其余的容器则为业务容器,. 
            - pause容器可用来判断Pod容器状态(因其与业务无关且不易死亡)
            - Pod内多个业务容器共享pause容器的IP和挂接的volume, 简化了业务容器之间的通信和文件共享问题
    - `Label`标签, 一个Label是一个key=value键值对, 一个对象可以定义任意数量的Label. Label Selector作用于pod时主要是用于
      `kube-controller`、`kube-proxy`或者`kube-scheduler`"定向调度"等场景
    - `Replication Controller`(RC): RC是kubernetes系统中的核心概念之一, 按照预期设定的值自动化的维护Pod期待的副本数. 维持
    服务高可用
    - `Deployment`: 内部使用了Replica Set来实现目的. 相对于RC, 我们可随时知道当前"Pod"部署的进度.
    - `Horizontal Pod`Autoscaler(HPA): Pod横向扩容, 用来智能化动态管理分配CPU资源等
    - `Service`: K8s中的service是微服务架构中的一个"微服务", 它定义了一个服务的访问入口地址. 前端应用(Pod)通过该入口
    地址访问其背后一组由Pod副本组成的集群实例. Service与后端Pod副本集群通过Label Selector实现"无缝对接"
        - service一旦创建, k8s会自动为其创建一个可用的Cluster IP(是虚拟IP, 只能结合Service Port组成一个具体的通信端口),
        属于K8s集群内部的地址.
        - 外部系统访问service
            - 有三种IP:
                - Node IP: Node节点的IP地址
                - Pod IP: Pod的IP地址
                - Cluster IP: Service的IP地址
            - 使用NodePort为需要外部访问的Service开启对应的TCP监听端口, 外部可使用Node的 IP地址 + 具体的NodePort
            即可访问此服务.
    - `Namespace`:  通过将集群内部资源的资源对象"分配"到不同的Namespace形成逻辑上的分组. 如果不特别致命Namespace, 
    K8s默认将用户创建的Pod、service等放入`default`的Namespace中.
        
     
     
## kubernetes配置
kubernetes里所有的资源对象都可以采用yaml或者JSON格式的文件来定义或描述
- pod资源定义文件:
```bash

``` 
     
## 常用命令


```bash
# 查看集群node数量
kubectl get nodes
# 获取node详细信息
kubectl describe node <node_name>
# 查看所有pod运行状态
kubectl get pods --all-namespaces
kubectl get pod --all-namespaces -o wide
# 获取pod详细信息
kubectl describe pod your_pod_name -n your_namespace
kubectl describe pod coredns-fb8b8dccf-ngsh5 -n kube-system
# 删除pod
kubectl delete pod pod_name
# 创建token(默认过期时间是24h)
kubeadm token generate
kubeadm token create

# 删除某个节点, --ignore-daemonsets为master节点使用
kubectl drain ubuntu-node-1 --delete-local-data --ignore-daemonsets
kubectl delete node ubuntu-node-1
kubect get nodes

# 查看service信息
kubectl get svc kubernetes -o yaml

# 删除所有已退出的docker容器    谨慎使用
docker rm `docker ps -a|grep Exited|awk '{print $1}'`

```