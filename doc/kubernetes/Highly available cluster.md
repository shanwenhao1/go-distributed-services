# 测试部署高可用集群
做高可用kubernetes集群至少需要k8s master节点两个以上, 最好是三个(测试部署两个挂掉一个集群就Done掉了)

@[TOC]
- [kubeadm访问控制配置](#kubeadm访问控制配置)
- [高可用集群架构选择](#高可用集群架构选择)
    - [Stacked etcd topology](#Stacked-etcd-topology)
    - [External etcd topology](#External-etcd-topology)
- [创建一个高可用集群](#创建一个高可用集群)
    - [准备](#准备)
    - [步骤](#步骤)
- [kubernetes使用](#kubernetes使用)

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

[如何创建独立的etcd cluster集群](https://kubernetes.io/docs/setup/independent/setup-ha-etcd-with-kubeadm/)


## 创建一个高可用集群
[文档](https://kubernetes.io/docs/setup/independent/high-availability/)

**注意: 通过kubeadm创建HA cluster一直在改进, 未来可能进一步简化**, 应时常关注更新文档.

### 准备
- 三台符合kubeadm要求的机器并且拥有Full network
- 拥有root权限
- kubeadm和kubelet已经安装在机器上

### 步骤
- 大致步骤:
    - 创建kube-apiserver 负载均衡(使用域名)
        - 默认健康检查端口为6443
        - 不建议在云环境中直接使用IP地址(最好使用域名)
        - 该load balance must be able to communicate with all control plan nodes on the apiServer
        - 确保负载均衡的地址与kubeadm配置中`controlPlaneEndpoint`匹配.
    - 将第一个节点加入the load balancer并测试连通性
    ```bash
    nc -v LOAD_BALANCER_IP PORT
    # 启动apiServer后才可运行该选项, 否则会失败
    nc -v 192.168.80.129 6443
    ```
    - 将剩下的control plane node 加入the load balancer
    
- Stacked control plane and etcd nodes(control plane与etcd 服务在同一节点上)
    - 在第一个control plane节点上, 创建配置文件`kubeadm-config.yaml`
    ```bash
    apiVersion: kubeadm.k8s.io/v1beta1
    kind: ClusterConfiguration
    kubernetesVersion: stable                                     # 设置使用的kubernetes版本, 这里使用v1.14.0
    controlPlaneEndpoint: "LOAD_BALANCER_DNS:LOAD_BALANCER_PORT"  # 负载均衡的地址(域名)和端口, 这里使用"192.168.80.129:6443
    ```
    - 初始化control plane
        - kubeadm init
        ```bash
        # experimental-upload-cert表示将所有control plane共享的证书上传至集群
        sudo kubeadm init --config=kubeadm-config.yaml --experimental-upload-certs
        ```
        获得如下结果
        ![](../../doc/picture/kubernetes/kubeadm%20init%20cluster%20master.png)
        ```bash
        # control plane join
        kubeadm join 192.168.80.129:6443 --token py4fq8.ocm3cety6u2uhr7v \
          --discovery-token-ca-cert-hash sha256:88becb99285b01e5ee0fbc39197f6946512c010afb010df89c22204a77ef24ca \
          --experimental-control-plane --certificate-key d5049926888228af9a2c747c071f16216f7b77e7daf0bb0059d242b2dcb67efd
        # worker node join
        kubeadm join 192.168.80.129:6443 --token py4fq8.ocm3cety6u2uhr7v \
            --discovery-token-ca-cert-hash sha256:88becb99285b01e5ee0fbc39197f6946512c010afb010df89c22204a77ef24ca
        # 如果需要更新证书, 只需要在集群中任意一个control plane上运行以下命令. 默认两小时过期
        kubeadm init phase upload-certs --experimental-upload-certs
        ```
        - 剩余步骤
        ```bash
        # 添加policy
        kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
        # 验证control plane是否已启动
        kubectl get pods --all-namespaces
        kubectl get nodes
        ```
    - 将剩下的control plane加入至集群, ~~注意更改主机名称~~([例子](single%20deploy.md#其它注意的事项))
        ```bash
        # 执行第一个control plane上生成的control plane join命令
          #   --experimental-control-plane要求kubeadm join 创建一个新的control plane
          #   --certificate-key从集群中下载证书并利用指定的key进行解密
         kubeadm join 192.168.80.129:6443 --token py4fq8.ocm3cety6u2uhr7v \
             --discovery-token-ca-cert-hash sha256:88becb99285b01e5ee0fbc39197f6946512c010afb010df89c22204a77ef24ca \
             --experimental-control-plane --certificate-key 6b6b3402020a861e9a10e74c4afc1ff8dba5b73a711d864cabe6537c096c3d92
        ```
        - 添加kubectl的环境变量`KUBECONFIG=/etc/kubernetes/admin.conf`, [示例](single%20deploy.md#step-2)
- External etcd nodes(etcd 服务独立出来的部署方式): 我们不采用这种方式, 因此就不写步骤了. 
[官方文档](https://kubernetes.io/docs/setup/independent/high-availability/#external-etcd-nodes)
- 接下来就是work node的加入了, [参考](single%20deploy.md#init-worker-node)
- 如果需要手动管理证书, [参考](https://kubernetes.io/docs/setup/independent/high-availability/#manual-certs)
- [使用kubeadm设置kubelet](https://kubernetes.io/docs/setup/independent/kubelet-integration/): kubernetes载入配置
都设置在`/etc/systemd/system/kubelet.service.d/10-kubeadm.conf`中

## kubernetes使用

[kubernetes使用](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/)