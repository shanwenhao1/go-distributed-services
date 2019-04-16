# actor_rpc_cluster_examples

ProtoActor 节点通过[consul](https://github.com/hashicorp/consul)进行服务注册, 
cluster还是用了[grain框架](https://github.com/dianbaer/grain), protoactor内部重构了该grain框架使用.


Please check that the following package have been installed.
```bash
    go get -u github.com/gogo/protobuf/protoc-gen-gogoslick
    go get -u github.com/AsynkronIT/protoactor-go/protobuf/protoc-gen-gograin
```
