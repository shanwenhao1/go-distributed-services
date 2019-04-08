# Go分布式服务

一个基础的、基于领域驱动设计(DDD)的Go分布式服务基础框架. 目前还在整理开发学习中.

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




## 资料来源
- [分布式服务框架](https://www.cnblogs.com/jiyukai/p/9459983.html)
- [微服务架构基础](https://blog.csdn.net/javaxuexi123/article/details/79500619#commentBox)
- [ddia数据密集型应用](https://github.com/Vonng/ddia/blob/master/preface.md)
- [skynete-archive](https://github.com/skynetservices/skynet-archive): go的分布式框架