# consul工具

@[TOC]
- [what is consul](#what-is-consul)
- [consul部署](#consul部署)
    - [部署注意事项](#部署注意事项)
    - [官网教程使用 basic use](#官网教程使用-basic-use)
        - [基础教程-single agent](#基础教程-single-agent)
            - [初步了解](#初步了解)
            - [晋升使用](#晋升使用)
        - [基础教程-consul cluster](#基础教程-consul-cluster)
        - [Health Check](#Health Check)
        - [KV Data](#KV Data)
        - [Web UI](#Web UI)
    - [Kubernetes部署](#Kubernetes部署)

[consul](https://github.com/hashicorp/consul)主要能实现三种服务:
- service discovery
- service configuration
- service segmentation.

## what is consul

Consul是一个用来实现分布式系统的服务发现与配置的开源工具. [consul简介](https://www.consul.io/intro/index.html), 
[中文版本](https://blog.51cto.com/firephoenix/2131616)


[官方文档](https://www.consul.io/docs/install/index.html)

## consul部署

[deploy consul service](https://learn.hashicorp.com/consul/#advanced), 
[中文版本](http://www.liangxiansen.cn/2017/04/06/consul/)

### 部署注意事项
- 部署至少3-5台server

### 官网教程使用 basic use

[下载](https://www.consul.io/downloads.html)linux压缩包, 在指定目录下解压会生成可执行文件, 将可执行文件拷贝至bin目录下

##### 基础教程-single agent

###### 初步了解
- 启动单点server
    ```bash
         consul agent -dev -config-dir=/home/swh/consul/consul.d
    ```
    - ctrl c 或者consul leave 命令停止agent
- 注册服务
    - 创建配置文件目录, 并配置agent启动配置(配置文件以.json命名), 通过更改配置文件并发送信号至agent可实现不停机修改.
        ```
            cat <<EOF | sudo tee /home/swh/consul/consul.d/web.json
            {
              "service": {
                "name": "web",
                "tags": ["rails"],
                "port": 80
              }
            }
            EOF
        ```
        - name: 服务名称
        - port: 监听端口
        - tags: 标签, 可用于查找该服务
- 注册服务并连接
    - 在consul.d目录下添加配置文件, 并使用consul reload重载consul
        ```bash
        cat <<EOF | sudo tee /home/swh/consul/consul.d/socat.json
        {
          "service": {
            "name": "socat",
            "port": 8181,
            "connect": { "sidecar_service": {} }
          }
        }
        EOF
        ```
        - 注意与之前注册的服务的区别是"connect"行, 该行只是告诉Consul这里有个注册的服务需要代理, 
        consul不会帮助你启动该service. 因此首先你需要启动该注册的服务.本样例中为
        ```bash
          # 模拟service(我们需要代理的服务)
          socat -v tcp-l:8181,fork exec:"/bin/cat"
        ```
    - 连接service
        ```bash
          # 启动代理进程(直联方式)
          consul connect proxy -sidecar-for socat
          # 启用代理进程(使用本地命令行方式配置本地代理进程, 此处将端口从8181替换为9191)
          # 这种封装方式更好, 因为它可以帮助你伪装成任何服务
          # 需注意此处需要socat代理进程已启动
          consul connect proxy -service web -upstream socat:9191
        ```
        - 使用以下命令测试, 分别使用8181和9191端口配合直联、command方式体会下两种方式的不同
            ```bash
              nc 127.0.0.1 9191
            ```
###### 晋升使用        
与[初步了解](#初步了解)功能相同
  
事实上, service与connect是上下游的依赖关系, service依赖于connect. 以下配置文件, 另建9191端口监听请求, 实际服务
部署在8181端口的socat服务.

- 创建配置文件, 调用console reload更新服务(或`consul agent -dev -config-dir=/home/swh/consul/consul.d`启动服务)        
    ```bash
    cat <<EOF | sudo tee /home/swh/consul/consul.d/web.json
    {
        "service": {
         "name": "web",
         "port": 8181,
         "connect": {
           "sidecar_service": {
             "proxy": {
               "upstreams": [{
                  "destination_name": "socat",
                  "local_bind_port": 9191
               }]
             }
           }
         }
        }
    }
    EOF
    ```
- 使用`socat -v tcp-l:8181,fork exec:"/bin/cat"`启动实际服务
- socat.json依旧需要保留, 使用命令`consul connect proxy -sidecar-for socat`启用socat代理
- 启动consul service, `consul connect proxy -sidecar-for web`, 使用`nc 127.0.0.1 9191`测试服务是否可用
    
###### Controlling access with Intentions  

consul中intentions是用来控制访问, [详情](https://www.consul.io/docs/connect/intentions.html). 例如下例
- 创建一条规则拒绝web请求socat的服务, 创建该规则后, `nc 127.0.0.1 9191`测试将会失败
```bash
consul intention create -deny web socat
```
- 删除该规则, 恢复连接
```bash
consul intention delete web socat
```

##### 基础教程-consul cluster

consul agent启动时就是一个单点集群. 如果要组成集群的话, 则需要加入已存在的集群, 只需要join集群中的一个节点就可以了
后续通过gossip协议可以迅速发现集群中的其他节点. agent加入集群不受限于agent类型

###### 模拟分布式服务(未进行该example实际操作)
- 通过Vagrant 模拟生成两个节点的集群
```bash
# 下载脚本
curl -O https://raw.githubusercontent.com/hashicorp/consul/master/demo/vagrant-cluster/Vagrantfile
# 启动
vagrant up
```
- `vagrant ssh n1`连接n1节点, 并consul agent server样例
```bash
consul agent -server -bootstrap-expect=1 -data-dir=/tmp/consul -node=agent-one -bind=172.20.20.10 -enable-script-checks=true -config-dir=/etc/consul.d
```
- 同理, 对n2节点也做以上操作, consul agent client样例
```bash
consul agent -data-dir=/tmp/consul -node=agent-two -bind=172.20.20.11 -enable-script-checks=true -config-dir=/etc/consul.d
```
- 至此, 两个agent目前还不知道对方, 可使用`consul members`查看. 在n1中按以下操作即可将n1中的agent加入至n2中集群
```bash
consul join 172.20.20.11
```

##### Health Check

 健康检查对于服务发现至关重要, 而consul中添加健康检查非常简单, 
 [详情](https://learn.hashicorp.com/consul/getting-started/checks)
 
##### KV Data
consul 提供了自带的KV store用来支持service discovery和health checking等服务. 
- put data, 如果是已存在的key则为更新操作
```bash
consul kv put redis/config/minconns 1
consul kv get -detailed redis/config/minconns
```
- get data
```bash
consul kv get redis/config/minconns
consul kv get -detailed redis/config/minconns
```
- delete key
```bash
consul kv delete redis/config/minconns
# 根据前缀删除所有该类型的key
consul kv delete -recurse redis
```
- `consul kv get -recurse`显示所有keys

##### Web UI

启动agent并允许UI访问, 访问`http://localhost:8500/ui`即可后台管理services、nodes、key/value等
```bash
consul agent -dev -ui
```

注意如果要查看的话, 还需要在后面添加`-client`参数
```bash
consul agent -dev -ui -client 0.0.0.0
```

### Kubernetes部署