# actor_rpc_cluster_examples(集群分片样例)

ProtoActor 节点通过[consul](https://github.com/hashicorp/consul)进行服务注册, 
cluster还是用了[grain框架](https://github.com/dianbaer/grain), protoactor内部重构了该grain框架使用.


Please check that the following package have been installed.
```bash
    go get -u github.com/gogo/protobuf/protoc-gen-gogoslick
    go get -u github.com/AsynkronIT/protoactor-go/protobuf/protoc-gen-gograin
```

cluster使用Consul 管理集群信息:
- Each member gets a hash-code, this hash-code is based on host + port + unique id of the member

## [官方文档](proto.actor/docs/grains#proto-cluster)

[cluster 参考](https://blog.csdn.net/TIGER_XC/article/details/86129995)

[cluster example其他例子](https://github.com/gviedma/cluster-examples)

## 注意

运行该cluster例子之前必须先搭建consul服务. 

使用命令启动测试所需单点consul
```bash
consul agent -server -bootstrap-expect 1 -data-dir /home/swh/consul/data -node=s1 -bind=192.168.1.89 -client 0.0.0.0
```
[详情](../../../doc/consul/consul.md)