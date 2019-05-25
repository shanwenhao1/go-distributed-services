# 生产部署方案

@[TOC]
- [组件部署](#组件部署)
    - [Dashboard](#Dashboard)

## 组件部署

### Dashboard
[官网](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/)
- 部署Dashboard UI
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/master/aio/deploy/recommended/kubernetes-dashboard.yaml
```
- Accessing the Dashboard UI: Dashboard默认开启了一个RBAC设置. Dashboard登录必须以`Bearer Token`凭证登录, 创建
[用户示例](https://github.com/kubernetes/dashboard/wiki/Creating-sample-user)
    - create a sample user
        - create service Account
        ```bash
          # admin.yaml
            apiVersion: v1
            kind: ServiceAccount
            metadata:
              name: admin-user
              namespace: kube-system
        ```
        - Create ClusterRoleBinding: 很多主流的工具如`kops`、`kubeadm`等在集群创建的时候已经创建了`ClusterRole`、
        `admin-Role`. 可用来创建`ClusterRoleBinding` for our `SeriveAccount`.
        ```bash
          # ClusterRoleBinding.yaml
            apiVersion: rbac.authorization.k8s.io/v1
            kind: ClusterRoleBinding
            metadata:
              name: admin-user
            roleRef:
              apiGroup: rbac.authorization.k8s.io
              kind: ClusterRole
              name: cluster-admin
            subjects:
            - kind: ServiceAccount
              name: admin-user
              namespace: kube-system
        ```
        - Bearer Token: 使用以下命令生成可用的token, 将该token复制至网页登录验证上即可访问了
        ```bash
        kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep admin-user | awk '{print $1}')
        # eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi11c2VyLXRva2VuLWZqNXBqIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImFkbWluLXVzZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI5ZWMwZmRhZC03Yjk5LTExZTktOTA2Yi0wMDBjMjliOGZhMjYiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06YWRtaW4tdXNlciJ9.ugeNAuzmii6nTCoMS7vfxgZzMB6IkJFsZI3JV7qo5BVKRG6-OrKBZ93LOnay1-SJ2qz2mFYUhyhKuzqnUnDbMWs_B8jumijL86e-YnsuUCaByDj3O6UkWVzvtC2N2hd_s6khPhW1GaOPhuo7_k8-J_79fjBWjwzxqgUoMhqFGGsQcceUgx9eYTY3MTQd8hbMKbBcKUiaz8KSGjEBu8oJBSO_nEtKmTHdRh4owWz4JW-P0LfAp3gL60JdV-uVMIbPRL7d8rKB0liqaSfg_v0-cc8Y7U7MfhntEr83w7t4tjNh2u7nKMKtzHmSy-sSLil-iV8RcYuUTa0sRhcx9XE1ag
        ```
        ![](../../doc/picture/kubernetes/Bearer%20Token.png)
    - Command line proxy
        - 可使用`kubectl proxy --address='192.168.80.133'  --accept-hosts='^*$'`启动服务
    
    

