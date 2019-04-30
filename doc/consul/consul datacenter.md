# Consul DataCenter

最新的secretID: 3b6d7fa6-0b58-b471-b3ec-68021012c64c

@[TOC]
- [Consul Reference Architecture](#Consul Reference Architecture)
- [Deployment Guide](#Deployment Guide)
    - [Download consul](#Download consul)
    - [Install consul](#Install consul)
    - [Configure system](#Configure system)
        - [Configure server](#Configure server)
    - [Start consul](#Start consul)
- [DataCenter Backups](#DataCenter Backups)
    - [Create Your First Backup](#Create Your First Backup)
    -[Restore from Backup](#Restore from Backup)
- [Bootstrapping the ACL System](#Bootstrapping the ACL System)
    - [Enable ACLs on all the Consul Servers](#Enable ACLs on all the Consul Servers)
    - [Create the Bootstrap Token](#Create the Bootstrap Token)
    - [Create an Agent Token Policy](#Create an Agent Token Policy)
    - [Create an Agent Token](#Create an Agent Token)
    - [Add the Agent Token to all the Servers](#Add the Agent Token to all the Servers)
    - []()
    - []()
    - []()
- []()
- [最终的配置文件样例](#最终的配置文件样例)

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
    sudo mkdir --parents /tmp/consul
    # 将权限转移给consul用户
    sudo chown --recursive consul:consul /tmp/consul
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
    ExecStart=/usr/bin/consul agent -bind=192.168.1.89 -config-dir=/etc/consul.d/ -data-dir=/tmp/consul
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
        data_dir = "/tmp/consul"
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

consul 提供了snapshot(快照)命令. 默认情况下, 所有快照使用`consistent mode`. 执行快照之前会经过leader校验
consul是否在线, 快照不会在数据中心降级及leader暂时不存在时执行

为了减少leader的压力, 可采取`stale consistency mode`在非leader机器上拍摄快照, 不过这样可能丢失少部分的最新数据.
通常丢失的是100ms以内的数据,但是在分布式服务上则无法保证该时间, 数据可能丢失更多.

### Create Your First Backup

通常备份应该写入脚本经常性的进行备份, 同时在datacenter升级之前也需要进行备份. 比如:
- 因为升级后导致的一些更改而使得不可能在降级回来
- datacenter失去仲裁所需最少节点而导致决策失效时, 你可以使用备份数据新增一个server节点从而使得服务恢复.

也无需在每个节点上都进行备份

备份的一些命令
- `consul snapshot save backup.snap` 会备份consul在执行该命令的目录下
    - `consul snapshot save -stale backup.snap`在非leader节点下specifying stale mode备份
    - ACL验证下备份
    ```bash
      export CONSUL_HTTP_TOKEN=<your ACL token>
      consul snapshot save -stale -ca-file=</path/to/file> backup.snap
    ``` 
- `consul snapshot inspect backup.snap`查看备份信息

### Restore from Backup 

为了确保从备份中恢复较平稳的运行, 请确保在leader上执行恢复备份的操作.
可使用`consul operator raft list-peers`查看节点状态.

恢复备份命令(此处是没有ACL验证的)
```bash
consul snapshot restore backup.snap
```

## Bootstrapping the ACL System

需要1.4+的集群, 使用前建议阅读[ACL 文档](https://www.consul.io/docs/acl/acl-system.html)

### Enable ACLs on all the Consul Servers
在`/etc/consul.d`目录下创建一个.json配置文件(agent使用, 这里我创建的是acl.json). 将以下配置写入文件
    ```bash
    {
      "acl": {
        "enabled": true,
        "default_policy": "deny",
        "down_policy": "extend-cache"
      }
    }
    ```
注意在consul集群中根据新的配置文件重启必须要一个接一个重启(确保之前重启的节点生效).创建成功则可在日志
中看到. ![](../../doc/picture/consul/acl.png)

### Create the Bootstrap Token 
bootstrap token is a management token with unrestricted privileges. 与集群中所有server共享

The bootstrap token can also be used in the server configuration file as the 
[master token](https://www.consul.io/docs/agent/options.html#acl_tokens_master)

bootstrap token只可创建一次, bootstrapping will disabled after the master token was created.
如果ACL system还处于bootstrapped状态, ACL tokens 可通过[ACL API](https://www.consul.io/api/acl/acl.html)更改.
```bash
# added to the state store
consul acl bootstrap
```
![](../../doc/picture/consul/acl%20token.png)

`注意`: 后续操作如果出现问题, 则可跳转至此按照步骤重新配置
```bash
# reset ACL system by reset index
consul acl bootstrap
        Failed ACL bootstrapping: Unexpected response code: 403 (Permission denied: ACL bootstrap no longer allowed (reset index: 13))
# 更改index, (此种方法暂时无效, 未理解）
# echo 13 >> <data-directory>/acl-bootstrap-reset
echo 13 >> /tmp/consul/acl-bootstrap-reset
```


一旦enable了acl, 那么我们所有的指令操作都要加上acl认证, 如下
- [command] -token SecretID方式, 不推荐
```bash
# before enable acl
consul members
# after enable acl, [command] -token SecretID
consul members -token 47a1551d-7ced-402b-78df-8ef98fd210ce
```
- 使用环境变量的方式自动认证acl, 然后即可使用`consul members`的方式照常访问了
```bash
export CONSUL_HTTP_TOKEN=47a1551d-7ced-402b-78df-8ef98fd210ce
```

### Create an Agent Token Policy
在创建token之前我们先需要创建一组策略来指定权限, [详情](https://www.consul.io/docs/acl/acl-rules.html)
- 编写策略文件`policy.hcl`(注意该文件最好别放在/etc/consul.d目录下, 会影响启动). 该策略允许任意设备可注册及可读(不建议在生产环境中使用)
```bash
node_prefix "" {
   policy = "write"
}
service_prefix "" {
   policy = "read"
}
```
- 只需生成一个policy策略就可使用在其他所有server节点上了, 注意: 如果未`export` CONSUL_HTTP_TOKEN, 
则需要在命令后面加上`- token SecretID`
```bash
consul acl policy create -name "agent-token" -description "Agent Token Policy" -rules @/tmp/consul/policy.hcl
```
![](../../doc/picture/consul/agent%20policy.png)
接下来可以使用这些策略去生成我们自己所需的agent token了

### Create an Agent Token
使用新创建的策略创建agent token. 可将token分享至所有server. 现在agent token中的SecretID可以用来认证API和CLI
命令了
```bash
consul acl token create -description "Agent Token" -policy-name "agent-token"
```
![](../../doc/picture/consul/agent%20token.png)

### Add the Agent Token to all the Servers
现在使用agent token配置consul server 并重启服务. 将以下配置写入`/etc/consul.d/server.json`中
```bash
{
  "primary_datacenter": "dc1",
  "acl": {
    "enabled": true,
    "default_policy": "deny",
    "down_policy": "extend-cache",
    "tokens": {
      "agent": "f5ccd945-cdea-1a92-465f-d6bda4fdfb0b"
    }
  }
}
```
检测ACL权限是否正常启用
```bash
curl http://127.0.0.1:8500/v1/catalog/nodes -H 'x-consul-token: f5ccd945-cdea-1a92-465f-d6bda4fdfb0b'
```

###

## 最终的配置文件样例
- `/etc/systemd/system/consul.service`,[文件](../../doc/consul/consul%20files/consul.service)
- `/etc/consul.d/consul.hcl`, [文件](../../doc/consul/consul%20files/consul.d/consul.hcl)
- `/etc/consul.d/server.hcl`, [文件](../../doc/consul/consul%20files/consul.d/server.hcl)
- `/etc/consul.d/acl.json`, [文件](../../doc/consul/consul%20files/consul.d/acl.json)
- `policy.hcl`, [文件](../../doc/consul/consul%20files/consul.d/policy.hcl)