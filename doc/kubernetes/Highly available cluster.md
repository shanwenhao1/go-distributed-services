# 测试部署高可用集群

@[TOC]
- [kubeadm访问控制配置](#kubeadm访问控制配置)
- [高可用集群架构选择](#高可用集群架构选择)
    - [Stacked etcd topology](#Stacked-etcd-topology)
    - [External etcd topology](#External-etcd-topology)
- [创建一个高可用集群](#创建一个高可用集群)
    - [准备](#准备)
- []()
- []()

## kubeadm访问控制配置
[文档](https://kubernetes.io/docs/setup/independent/control-plane-flags/)

- APIServer参数设置(范例), [参数意义](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/)
```bash
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: v1.13.0
metadata:
  name: 1.13-sample
apiServer:
  extraArgs:
    advertise-address: 192.168.0.103
    anonymous-auth: false
    enable-admission-plugins: AlwaysPullImages,DefaultStorageClass
    audit-log-path: /home/johndoe/audit.log
```
- ControllerManager参数设置(范例), [参数意义](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/)
```bash
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: v1.13.0
metadata:
  name: 1.13-sample
controllerManager:
  extraArgs:
    cluster-signing-key-file: /home/johndoe/keys/ca.key
    bind-address: 0.0.0.0
    deployment-controller-sync-period: 50
```
-Scheduler参数设置(范例), [参数意义](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-scheduler/)
```bash
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: v1.13.0
metadata:
  name: 1.13-sample
scheduler:
  extraArgs:
    address: 0.0.0.0
    config: /home/johndoe/schedconfig.yaml
    kubeconfig: /home/johndoe/kubeconfig.yaml
```


## 高可用集群架构选择
[文档](https://kubernetes.io/docs/setup/independent/ha-topology/), 
我选择[Stacked etcd topology](#Stacked-etcd-topology)结构.

### Stacked etcd topology

A Stacked HA cluster由运行在`kubeadm`所创建的集群上etcd服务提供. 是kubeadm默认topology.
- 每个control plane node上都运行了`kube-apiserver`(通过负载均衡暴露给worker node)、`kube-scheduler`、
`kube-controller-manager`. 并且每个control plane node 上都创建一个本地etcd key-value存储服务, 只与本节点
的`kube-apiserver`通信.
- 该topology将control plan与etcd members部署在同一节点上, 方便扩展.但是这样做的坏处是当一个node挂掉时, etcd member
control plane都将丢失, 对冗余服务造成影响(可通过增加更多的control plane nodes规避该风险. 建议最少三个)
- local member当运行`kubeadm init`、`kubeadm join --experimental-control-plane`时会自动创建.
![](../../doc/picture/kubernetes/etcd%20with%20kubeadm%20cluster.png)

### External etcd topology

An HA cluster with external etcd是独立运行在kubeadm cluster之外的key-value存储.

较于`Stacked etcd topology`而言, 它能停供更稳定的服务, 对于集群冗余的影响较小. 但是至少需要两倍的机器进行部署
![](../../doc/picture/kubernetes/etcd%20cluster.png)


## 创建一个高可用集群
[文档](https://kubernetes.io/docs/setup/independent/high-availability/)

**注意: 通过kubeadm创建HA cluster一直在改进, 未来可能进一步简化**, 应时常关注更新文档.

### 准备
