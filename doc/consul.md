# consul工具

@[TOC]
- [what is consul](#what-is-consul)
- [consul部署](#consul部署)
    - [部署注意事项](#部署注意事项)
    - [基础使用](#基础使用)

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

### 基础使用
本人使用:
- [下载](https://www.consul.io/downloads.html)linux压缩包, 在指定目录下解压会生成可执行文件, 将可执行文件拷贝至bin目录下
- 启动单点server
    ```bash
         consul agent -dev -config-dir=/home/swh/consul/consul.d
    ```
    - ctrl c 或者consul leave 命令停止agent
- 注册服务
    - 创建配置文件目录, 并配置agent启动配置(配置文件以.json命名), 通过更改配置文件并发送信号至agent可实现不停机修改.
        ```
          echo '{"service": {"name": "web", "tags": ["rails"], "port": 80}}' > /home/swh/consul/consul.d/web.json
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
          # 模拟service 
          socat -v tcp-l:8181,fork exec:"/bin/cat"
        ```
    - 连接service
        ```bash
          # 启动代理进程(直联方式)
          consul connect proxy -sidecar-for socat
          # 启用代理进程(使用本地命令行方式配置本地代理进程, 此处将端口从8181替换为9191)
          # 这种方式更好, 因为它可以帮助你伪装成任何服务
          consul connect proxy -service web -upstream socat:9191
        ```
        - 使用以下命令测试, 分别使用8181和9191端口配合直联、command方式体会下两种方式的不同
            ```bash
              nc 127.0.0.1 9191
            ```
    
    
    ```bash
      consul agent -server -bootstrap-expect 1 -data-dir /home/swh/consul/data -node=s1 -bind=192.168.1.89 -ui-dir /home/swh/consul/ui -rejoin -config-dir=/home/swh/consul/config -client 0.0.0.0
    ```