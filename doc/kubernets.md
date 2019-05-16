# Kubernetes文档

## 资料

- [官方文档](https://kubernetes.io/docs/home/)
- [Kubernetes权威指南](Kubernetes权威指南.pdf)

## 阅读摘要

- [kubernetes入门](kubernetes/readingNote/chapter1.md)

## 部署服务
[官方文档](https://kubernetes.io/docs/setup/)

有几种部署模式, [详情](https://kubernetes.io/docs/setup/pick-right-solution/#local-machine-solutions): 
- `local Docker-based solutions`: just for learning or checking, not for production. 
[本地学习部署](kubernetes/single%20deploy.md)(single master and etcd server), 
[部署高可用kubernetes集群](kubernetes/Highly%20available%20cluster.md)(multi master and etcd members)
- `hosted solution`: more machines and higher availability. [生产服务部署](kubernetes/production%20deploy.md)
- [通过docker快速部署](docker_practice.pdf)

## kubernetes使用教程
[RBAC 权限认证](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)



## kubernetes常见的一些错误
[troubleshooting](https://kubernetes.io/docs/setup/independent/troubleshooting-kubeadm/)
