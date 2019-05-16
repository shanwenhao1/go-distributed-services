# kubernetes入门
                                                 
## kubernetes架构导读

- services是分布式集群架构的核心
- k8s设计了`Pod`对象, 将服务进程包装至相应的Pod中, 使其成为Pod中运行的一个container. Pod也是k8s管理运行的最小运行单元
    - 每个pod都有自己的`Label`
    - 每个Pod中都运行了名为Pause的容器, 其余的容器则为业务容器
- k8s集群中的机器划分为一个Master节点和一群工作节点(Node)
    - Master通常占用一个独立的机器(太重要, 不允许宕机). Master节点上运行以下关键进程
        - kube-apiserver: k8s内所有资源控制的唯一入口
        - kube-controller-manager: k8s所有资源对象的自动化控制中心
        - kube-scheduler: 负责资源调度(Pod调度)
        - etcd: 高可用的key-value store 负责存储k8s集群数据
     
     
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

# 删除所有已退出的docker容器    谨慎使用
docker rm `docker ps -a|grep Exited|awk '{print $1}'`

```