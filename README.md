# Go分布式服务

一个基础的、基于领域驱动设计(DDD)和消息驱动的Go分布式服务基础框架. 目前还在整理开发学习中.

@[TOC]
- [分布式架构基础知识](#分布式架构基础知识)
    - [SOA架构](#SOA架构)
    - [微服务架构](#微服务架构)
- []()
- [待引入](#待引入)
- [资料来源](#资料来源)
    - [目前考量的](#目前考量的)

## 分布式架构基础知识

### SOA架构

![](doc/picture/distributed%20frame/soa.png)
最基础的几个模块
- Provider: 服务提供者, 比如业务服务
- Consumer: 发起调用的客户端
- Registry: 服务注册中心, 是分布式服务系统中的一个重要组成模块, 管理Provider的Manager, 在实际的运行环境中,
服务注册中心的Registry被动通知或Consumer主动询问. 在Provider有节点宕机或新增时, 客户端可以实时感知到, 从而
避免某个Provider被无限调用或是无限闲置.
- Gateway: 网关主要进行接受外部HTTP请求, 校验权限等, 路由转发至对应的Provider, 再将Provider传回的结果传递给
Consumer.
- 负载均衡: 服务分流, 比如Nginx
- 监控: 接受来自Consumer和Provider异步上报的性能监控数据, 对有风险的节点发出告警.


### 微服务架构

属于分布式架构的一种, 但更加强调单一职责、轻量级通信(HTTP或gRPC)、独立性且进程隔离

## 使用工具

依赖工具
- [ProtoActor](doc/proto%20actor.md)
- [consul](doc/consul/consul.md): 分布式集群服务发现工具
- [Kubernetes](doc/kubernets.md): docker集群部署管理工具

框架使用
- [gOrm](http://gorm.book.jasperxu.com/)

## 待引入

- ProtoActor: 新型的actor框架, 事件驱动
- consul: 分布式集群服务发现依赖
    - kubernetes部署服务: docker集群管理
    - 负载均衡: ![](doc/picture/distributed%20frame/Load%20Balance.jpg)
- kafka(考虑)




## 资料来源
- [3分钟读懂何为分布式、微服务和集群](http://server.51cto.com/News-557053.htm)
- [分布式服务框架](https://www.cnblogs.com/jiyukai/p/9459983.html)
- [微服务架构基础](https://blog.csdn.net/javaxuexi123/article/details/79500619#commentBox)
- [ddia数据密集型应用](https://github.com/Vonng/ddia/blob/master/preface.md)
- [大型网站系统架构演化](http://www.cnblogs.com/leefreeman/p/3993449.html)

### 目前考量的
- [Proto Actor模型框架](doc/proto%20actor.md): 注意ProtoActor还未发布release版本
    - protoactor分布式集群依赖[consul](doc/consul/consul.md)服务
- [skynete-archive](https://github.com/skynetservices/skynet-archive): go的分布式框架