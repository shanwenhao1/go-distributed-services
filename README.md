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
- [大型网站系统架构演化](http://www.cnblogs.com/leefreeman/p/3993449.html)

### 目前考量的
- [Proto Actor模型框架](https://github.com/AsynkronIT/protoactor-go): 基于Actor模型的消息驱动的高并发框架, 
使用rpc通信, 目前还处于开发当中. [使用文档说明](https://github.com/AsynkronIT/protoactor-go)
    - Actor模型说明:
        - proto actor运用的协议是google的rpc协议, [gRPC文档](http://doc.oschina.net/grpc?t=60133), 
            - [go Protocol google安装](https://blog.csdn.net/u010230794/article/details/78606021)
                - .exe放在GOPATH\bin下面
                - include下的google文件夹放在GOPATH\src下面
            - [官方安装文档](https://github.com/golang/protobuf)
        - 适用于高并发分布式服务, 每个线程都是一个actor, actor之间并不共享任何内存.因此不必使用分布式锁.
        - Actor模型=数据+行为+消息: Actor模型内部的状态由自己的行为维护,外部线程不能直接调用对象的行为,
        必须通过消息才能激发行为,这样就保证Actor内部数据只有被自己修改
        - CQRS架构中的写操作就尽量使用actor模型实现.
    - actor介绍:
        - [基于go的actor框架-protoactor使用小结](https://studygolang.com/articles/12302)
        - [下一代的 Actor 模型框架 Proto Actor](https://www.oschina.net/p/protoactor)
        - [高并发解决方案之Actor](https://www.cnblogs.com/gengzhe/p/6561655.html)
        - [为什么Actor模型是高并发事务的终极解决方案?](https://www.jdon.com/45728)
- [skynete-archive](https://github.com/skynetservices/skynet-archive): go的分布式框架