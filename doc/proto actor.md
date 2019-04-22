# ProtoActor

[Proto Actor模型框架](https://github.com/AsynkronIT/protoactor-go): 基于Actor模型的消息驱动的高并发框架, 
使用rpc通信, 目前还处于开发当中. [使用文档说明](https://github.com/AsynkronIT/protoactor-go)
- actor介绍:
    - [基于go的actor框架-protoactor使用小结](https://studygolang.com/articles/12302)
    - [下一代的 Actor 模型框架 Proto Actor](https://www.oschina.net/p/protoactor)
    - [高并发解决方案之Actor](https://www.cnblogs.com/gengzhe/p/6561655.html)
    - [为什么Actor模型是高并发事务的终极解决方案?](https://www.jdon.com/45728)
- Actor模型说明, [Actor文档](http://proto.actor/docs/actors):
    - proto actor运用的协议是google的rpc协议, [gRPC文档](http://doc.oschina.net/grpc?t=60133)
    - 适用于高并发分布式服务, 每个线程都是一个actor, actor之间并不共享任何内存.因此不必使用分布式锁.
    - Actor模型=数据+行为+消息: Actor模型内部的状态由自己的行为维护,外部线程不能直接调用对象的行为,
    必须通过消息才能激发行为,这样就保证Actor内部数据只有被自己修改
    - CQRS架构中的写操作就尽量使用actor模型实现.
    - Actor 关注实例状态,重,默认在实例外加了个壳子(pid+context,AKKA的是actorRef+Context),
    封装两个队列(MSPC(multi produce single consumer)、RingBuffer)三个集合(children、watchers、watchees),
    更能处理复杂业务逻辑.
- 安装及依赖, [可参考](https://travis-ci.org/AsynkronIT/protoactor-go/jobs/516220191)
    - rpc和protobuf安装
        - [go Protocol google安装](https://blog.csdn.net/u010230794/article/details/78606021), 
        [官方安装文档](https://github.com/golang/protobuf)
            - .exe放在GOPATH\bin下面
            - include下的google文件夹放在GOPATH\src下面
        - actor 生成rpc文件必须使用--gogoslick_out, 而不是--go_out. [安装命令](doc/updateproto.bat)
        - 需要该项目的工具gograin生成必要文件, 在protobuf目录下
            ```bash
              go get -u github.com/AsynkronIT/protoactor-go/protobuf/protoc-gen-gograin
            ```