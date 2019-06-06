# Kubernetes文档

## 资料

- [官方文档](https://kubernetes.io/docs/home/)
    - [官方中文文档](http://docs.kubernetes.org.cn/)
- [Kubernetes权威指南](Kubernetes权威指南.pdf)

## 阅读摘要
- [kubernetes入门](kubernetes/readingNote/chapter1.md)
- [kubernetes实践指南](kubernetes/readingNote/chapter2.md)
- [kubernetes核心原理](kubernetes/readingNote/chapter3.md)

## 部署服务
[官方文档](https://kubernetes.io/docs/setup/)

### 自定义部署服务探索
有几种部署模式, [详情](https://kubernetes.io/docs/setup/pick-right-solution/#local-machine-solutions): 
- `local Docker-based solutions`: just for learning or checking, not for production. 
[本地学习部署](kubernetes/single%20deploy.md)(single master and etcd server), 
[部署高可用kubernetes集群](kubernetes/Highly%20available%20cluster.md)(multi master and etcd members)
- `hosted solution`: more machines and higher availability. [生产服务部署](kubernetes/production%20deploy.md)
- [通过docker快速部署](docker_practice.pdf)

### 生产环境部署服务
[阿里云部署](https://help.aliyun.com/product/85222.html?spm=a2c4g.11186623.6.540.6dec45e5aDPi0o)


## kubernetes使用教程
[RBAC 权限认证](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)


## kubernetes常见的一些错误
[troubleshooting](https://kubernetes.io/docs/setup/independent/troubleshooting-kubeadm/)


## kubernetes Yaml文件
[模板](../doc/kubernetes/k8s%20example.yaml)