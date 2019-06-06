# Kubernetes核心原理

## Kubernetes APIServer

提供了k8s各类资源对象(如Pod、RC、Service等)的增、删、改、查及watch等HTTP Rest接口. 有以下功能特性:
- 是集群管理的API入口
- 是资源配额控制的入口
- 提供了完备的集群安全机制

kube-apiserver进程在本机的8080端口

独特的Kubernetes Proxy API接口, 用作代理REST请求. 
- 在kubernetes集群之外访问某个Pod容器的服务, 可以用Proxy API实现

## Controller Manager
作为集群内部的管理控制中心, 负责集群内的Node、Pod副本、服务端点(Endpoint)、命名空间(Namespace)、服务账号(ServiceAccount)
、资源定额(ResourceQuota)等的管理

- Replication Controller(副本控制器): 核心作用是确保任何时候集群中一个RC所关联的副本数量保持预设值.
- Node Controller
- ResourceQuota Controller(资源配额管理)
- Namespace Controller: 用户通过API Server可以创建新的Namespace并保存在etcd中, 也可设置DeletionTimestamp属性
以优雅删除Namespace(通过设置删除期限)
- Service Controller:
- Endpoint Controller

## Scheduler

Kubernetes Scheduler是k8s中负责Pod调度的模块, 在整个系统中承担了"承上启下"的重要功能.
"承上"指它负责接收Controller Manager创建的新Pod, 为其安排Node节点部署;
"启下"指安置工作完成后, 目标Node上的kubelet服务进程接管后继工作, 负责Pod生命周期的"下半生".

## Kubelet运行机制
k8s中, 每个Node节点上都会启动一个kubelet服务进程, 该进程用于处理Master节点下发到本节点的任务, 管理Pod及
Pod中的容器. 每个kubelet进程会在API Server上注册节点自身信息, 定期向Master节点汇报节点资源的使用情况, 
并通过cAdvisor监控容器和节点资源.


## Kube-proxy运行机制

kube-proxy进程可以看作Service的透明代理兼负载均衡器, 核心功能是将某个Service的访问请求转发到后端的多个Pod
实例上. 

## 深入分析集群安全机制
集群的安全性须考虑以下几个目标:
- 保证容器与其所在的宿主机的隔离
- 限制容器给基础设施及其他容器带来的消极影响的能力
- 最小权限原则----合理限制所有组件的权限, 确保组件只执行它被授权的行为, 通过限制单个组件的能力来限制它所能到达的权限范围
- 明确组件间边界的划分
- 划分普通用户和管理员的角色
- 在必要的时候允许将管理员权限赋给普通用户
- 允许拥有"Secret"数据(Keys、Certs、Passwords)的应用在集群中运行.

k8s集群中所有资源的访问和变更都是通过Kubernetes API Server的REST API来实现的. 因此集群安全的关键点就在于如何识别并认证
客户端身份(Authentication), 以及随后的访问权限的授权(Authorization)这两个关键问题. k8s目前提供了三种级别的客户端
身份认证方式:
- 最严格的HTTPS证书认证: 基于CA根证书签名的双向数字证书认证方式
- HTTP Token认证: 通过一个Token来识别合法用户
- HTTP Base认证: 通过用户名 + 密码的方式认证

## 网络原理

kubernetes网络模型被称作IP-per-Pod模型, IP以Pod为单位进行分配, 一个Pod内部的所有容器共享一个网络堆栈.

Docker本身技术依赖于Linux内核虚拟化技术, 因此Dockere对Linux内核的特性有很强的依赖



