# Consul DataCenter

@[TOC]
- [Consul Reference Architecture](#Consul Reference Architecture)
- [Deployment Guide](#Deployment Guide)
    - [Download consul](#Download consul)
    - [Install consul](#Install consul)
    - [Configure system](#Configure system)
        - [Configure server](#Configure server)
    - [Start consul](#Start consul)
- [DataCenter Backups](#DataCenter Backups)
- []()
- []()
- []()

本文档为建立consul DataCenter的指引文档

[官方文档](https://learn.hashicorp.com/consul/advanced/day-1-operations/path-intro)

## Consul Reference Architecture

consul data center 架构, ![](../../doc/picture/consul/infrastructure%20diagram.png)

### DataCenter design

一个简单的consul DataCenter通常由3-5个节点构成. 对于读写要求很高的集群, 可通过将cluster server部署在同一物理地址上提高性能

#### Single DataCenter

#### Multiple DataCenters

You can join Consul clusters running the same service in different datacenters by WAN links.

在Multiple DataCenter中, consul数据并不会跨数据中心同步, 
如有需要请使用[consul-replicate](https://github.com/hashicorp/consul-replicate)工具定期同步.

## Deployment Guide

Single DataCenter, ![](../../doc/picture/consul/single%20datacenter.png)

### Download consul
consul下载地址[https://releases.hashicorp.com/consul/](https://releases.hashicorp.com/consul/)


### Install consul
解压缩文件`unzip consul_${CONSUL_VERSION}_linux_amd64.zip`, 并将可执行文件`consul`移至`usr/local/bin`中
- Enable autocompletion, autocompletion is `consul` command features
    ```bash
    consul -autocomplete-install
    complete -C /usr/bin/consul consul
    ```
- 在系统中创建一个唯一的、非特权的名为consul的用户, 并为其创建所属数据文件夹. 
可用`sudo userdel consul`删除新建用户
    ```bash
    # 创建consul用户
    sudo useradd --system --home /etc/consul.d --shell /bin/false consul
    # 创建consul用户所需datacenter文件夹
    sudo mkdir --parents /home/swh/consul/datacenter
    # 将权限转移给consul用户
    sudo chown --recursive consul:consul /home/swh/consul/datacenter
    ```

### Configure system
Systemd 一般使用[默认](https://www.freedesktop.org/software/systemd/man/systemd.directives.html), 因此一些非
默认参数必须写入配置文件中`/etc/systemd/system/consul.service`
- 将以下配置加入consul.service中
    ```bash
    [Unit]
    # Free-form string describing the consul service 
    Description="HashiCorp Consul - A service mesh solution"
    # Link to the consul documentation 
    Documentation=https://www.consul.io/
    # Configure a requirement dependency on the network service
    Requires=network-online.target
    # Configure an ordering dependency on the network service being started before the consul service 
    After=network-online.target
    # Check for a non-zero sized configuration file before consul is started
    ConditionFileNotEmpty=/etc/consul.d/consul.hcl
    
    [Service]
    # Run consul as the consul user 
    User=consul
    Group=consul
    # Start consul with the `agent` argument and path to the configuration file 
    # bind is default to 0.0.0.0, 本不用配置, 但由于可能存在多个IP, 因此此处加上bind参数
    ExecStart=/usr/bin/consul agent -bind=192.168.1.89 -config-dir=/etc/consul.d/
    # 开放客户端端口
    # ExecStart=/usr/bin/consul agent -bind=192.168.1.89 -config-dir=/etc/consul.d/ -client=0.0.0.0
    # Send consul a reload signal to trigger a configuration reload in consul 
    ExecReload=/usr/bin/consul reload
    # Treat consul as a single process
    KillMode=process
    # Restart consul unless it returned a clean exit code 
    Restart=on-failure
    # Set an increased Limit for File Descriptors
    LimitNOFILE=65536
    
    [Install]
    # Creates a weak dependency on consul being started by the multi-user run level 
    WantedBy=multi-user.target
    ```
    
#### Configure server
Configuration可按照lexical order从多个文件总加载. 关于configuration loading和合并的
[详情](https://www.consul.io/docs/agent/options.html)

以下为配置所有consul agent使用的配置`consul.hcl` 及server所需的配置`server.hcl`
- 配置`/etc/consul.d/consul.hcl`
    - 创建consul.hcl
        ```bash
        sudo mkdir --parents /etc/consul.d
        sudo touch /etc/consul.d/consul.hcl
        sudo chown --recursive consul:consul /etc/consul.d
        sudo chmod 640 /etc/consul.d/consul.hcl
        ```
    - 将以下配置写入`consul.hcl`
        ```bash
        # The datacenter in which the agent is running
        datacenter = "dc1"
        # The data directory for the agent to store state
        data_dir = "/home/swh/consul/datacenter"
        # Specifies the secret key to use for encryption of Consul network traffic
        encrypt = "Luj2FZWwlt8475wD1WtwUQ=="
        ```
        - cluster auto-join: `retry_join` 参数允许使用ip, DNS等自动加入集群
            ```bash
            # Address of another agent to join upon starting up
            retry_join = ["192.168.1.89"]
            ```
        - Performance stanza: 允许调节不同子系统的性能. 
        将`raft_multiplier`设置为1时表示启用最高性能(也是生产中推荐配置)
            ```bash
            # raft_multiplier: An integer multiplier used by Consul servers to scale key Raft timing parameters
            performance{
              raft_multiplier = 1
            }
        ```
- 配置`/etc/consul.d/server.hcl`
    - 创建server.hcl
        ```bash
        sudo touch /etc/consul.d/server.hcl
        sudo chmod 640 /etc/consul.d/server.hcl
        ```
    - 将以下配置加入`server.hcl`
        ```bash
        # This flag is used to control if an agent is in server or client mode
        server = true
        # This flag provides the number of expected servers in the datacenter. 
        # Either this value should not be provided or the value must agree with other servers in the cluster
        bootstrap_expect = 3
        ```
        - enable ui
            ```bash
            ui = true
            ```          
### Start consul

启动consul
```bash
sudo systemctl enable consul
sudo systemctl start consul
sudo systemctl status consul
```  

## DataCenter Backups