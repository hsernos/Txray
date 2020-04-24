# v3ray
v3ray是一款v2ray 终端版客户端，使用go编写.
项目地址：https://github.com/hsernos/v3ray

## 注意
此文档只针对于最新的commit, 可能不适用于已发布的最新版本.

<!-- toc -->

## 目录

- [特色](#特色)
- [编译/交叉编译 说明](#编译交叉编译-说明)
- [下载/运行 说明](#下载运行-说明)
  * [注意!!!](#注意---) 
- [命令列表及说明](#命令列表及说明)
  * [查看基本设置帮助文档](#查看基本设置帮助文档)
    + [查看基本设置](#查看基本设置)
    + [修改基本设置](#修改基本设置)
  * [查看订阅帮助文档](#查看订阅帮助文档)
    + [添加订阅](#添加订阅)
    + [查看订阅](#查看订阅)
    + [修改订阅](#修改订阅)
    + [删除订阅](#删除订阅)
    + [从订阅更新节点](#从订阅更新节点)
  * [查看节点帮助文档](#查看节点帮助文档)
    + [添加节点](#添加节点)
    + [查看节点](#查看节点)
    + [删除节点](#删除节点)
    + [tcping测试](#tcping测试)
    + [节点查找](#节点查找)
    + [导出节点](#导出节点)
  * [查看DNS帮助文档](#查看DNS帮助文档)
    + [添加DNS](#添加DNS)
    + [查看DNS](#查看DNS)
    + [删除DNS](#删除DNS)
  * [查看路由帮助文档](#查看路由帮助文档)
    + [domain路由规则](#domain路由规则)
    + [ip路由规则](#ip路由规则)
  * [查看服务帮助文档](#查看服务帮助文档)
    + [启动或重启v2ray-core服务](#启动或重启v2ray-core服务)
    + [停止v2ray-core服务](#停止v2ray-core服务)
- [已知问题](#已知问题)
- [交流反馈](#交流反馈)

<!-- tocstop -->

# 特色

1. 多平台支持, 支持 Windows, macOS, linux.
2. Tab键命令补齐

#  编译/交叉编译 说明

1. 在终端下进入项目目录

2. 设置`GOPROXY`,提高编译所需依赖的下载速度
   Linux/Mac下，运行 `GOPROXY=https://goproxy.cn,direct`
   Windows下,运行 `set GOPROXY=https://goproxy.cn,direct`

3. 编译常用平台
   运行 `go build v3ray.go`, 可编译当前平台的版本
   运行 `python3 build.py`, 可编译常用平台的版本

4. 编译其他平台
   运行 `go tool dist list` 查看所有支持的 GOOS/GOARCH

   Linux/Darwin 例子: 编译 Windows 下的 64 位程序

   `GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build v3ray.go`
   
   Windows 例子: 编译 Linux 下的 32 位程序
   
   `set GOOS=linux`
   `set GOARCH=386`
   `set CGO_ENABLED=0`
   `go build v3ray.go`


# 下载/运行 说明

Go语言程序, 可直接在[发布页](https://github.com/hsernos/v3ray/releases)下载使用

## 注意！！！

运行前需要下载[v2ray-core](https://github.com/v2ray/v2ray-core/releases), 下载解压后将`v2ray`程序所在的目录添加到`PATH`环境变量中

`v3ray`生成的配置文件以及日志优先放在`V3RAY`环境变量指向的文件夹中

其次放在当前用户家目录下的`v3ray文件夹`下



# 命令列表及说明

> 在终端中运行v3ray进入shell交互 (可以将v3ray放在v2ray所在目录下，就可以在终端直接运行进入)



## 查看基本设置帮助文档

```
>>> setting 

setting {commands} [flags] ...
    
Commands:
    show                               查看基本设置
    alter [flags]                      修改基础设置
    
alter Flags
    -p, --port       {port}            设置监听端口
    -u, --udp        {true|false}      是否启用udp
    -s, --sniffing   {true|false}      是否启用流量监听
    -l, --lanconn    {true|false}      是否启用局域网连接
    -m, --mux        {true|false}      是否启用局域网连接
    -b, --bypass     {true|false}      是否启用绕过局域网及大陆
    -r, --route      {1|2|3}           设置路由策略为{AsIs|IPIfNonMatch|IPOnDemand}
```

### 查看基本设置

```
>>> setting show
+--------+---------+------------+--------+--------------+---------------+---------------+
| 监听端口 | UDP转发 | 启用流量监听 | 多路复用 | 允许局域网连接 | 绕过局域网和大陆 |    路由策略    |
+---------+--------+------------+---------+-------------+----------------+--------------+
|  2333   |  true  |    true    |   true  |    false    |      true      | IPIfNonMatch |
+---------+--------+------------+---------+-------------+----------------+--------------+
```

### 修改基本设置

```
# 修改监听端口为3333
>>> setting alter -p 3333

# 修改不绕过局域网和大陆
>>> setting alter -b false

# 修改路由策略为IPIfNonMatch, {1|2|3}=>{AsIs|IPIfNonMatch|IPOnDemand}
>>> setting alter -r 2
```



## 查看订阅帮助文档

```
>>> sub

sub {commands} [flags] ...

Commands:
    add   {订阅链接} [flags]             添加订阅
    alter {索引范围} [flags]             修改订阅
    show  [索引范围]                     查看订阅信息
    del   {索引范围}                     根据索引参数删除订阅
    update-node     [flags]             从订阅更新节点

add Flags
    -r, --remarks {别名}                 自定义别名

alter Flags
    -u, --url     {订阅链接}             订阅链接
    -r, --remarks {别名}                 自定义别名
    --using       {true|false}          是否使用此订阅

update-node Flags
    -p, --proxy {本地socks5端口}         从socks5代理更新节点
```

### 添加订阅

```
# 添加订阅链接为https://sublink.com
>>> sub add https://sublink.com

# 添加订阅链接为https://sublink.com，并命名为test
>>> sub add https://sublink.com -r test
```

### 查看订阅

```
# 查看全部订阅
>>> sub show
+------+-------+---------------------+----------+
| 索引  | 别名   |       URL          |  是否启用  |
+------+-------+---------------------+----------+
|  0   | test  | https://sublink.com |   true   |
|  1   | test1 | https://sublink.com |   true   |
|  2   | test2 | https://sublink.com |   true   |
|  3   | test3 | https://sublink.com |   true   |
|  4   | test4 | https://sublink.com |   true   |
|  5   | test5 | https://sublink.com |   true   |
|  6   | test6 | https://sublink.com |   true   |
+------+-------+---------------------+----------+

# 查看索引为2,3,4的订阅
>>> sub show 2-4
+------+-------+---------------------+----------+
| 索引  | 别名  |         URL         |  是否启用  |
+------+-------+---------------------+----------+
|  2   | test2 | https://sublink.com |   true   |
|  3   | test3 | https://sublink.com |   true   |
|  4   | test4 | https://sublink.com |   true   |
+------+-------+---------------------+----------+

# 查看索引为0,1,2,3,6的订阅
>>> sub show 0-3,6
+------+-------+---------------------+----------+
| 索引  | 别名  |         URL         |  是否启用 |
+------+-------+---------------------+----------+
|  0   | test  | https://sublink.com |   true   |
|  1   | test1 | https://sublink.com |   true   |
|  2   | test2 | https://sublink.com |   true   |
|  3   | test3 | https://sublink.com |   true   |
|  6   | test6 | https://sublink.com |   true   |
+------+-------+---------------------+----------+

# 查看索引为3以及后面的的订阅
>>> sub show 3-1000
+------+-------+---------------------+----------+
| 索引  | 别名  |         URL         |  是否启用  |
+------+-------+---------------------+----------+
|  3   | test3 | https://sublink.com |   true   |
|  4   | test4 | https://sublink.com |   true   |
|  5   | test5 | https://sublink.com |   true   |
|  6   | test6 | https://sublink.com |   true   |
+------+-------+---------------------+----------+
```

### 修改订阅

```
# 修改索引为0的订阅链接为https://test.com，别名为test0
>>> sub atler 0 -u https://test.com -r test0
>>> sub show 0
+------+-------+------------------+----------+
| 索引  | 别名  |        URL        |  是否启用 |
+------+-------+------------------+----------+
|  0   | test0 | https://test.com |   true   |
+------+-------+------------------+----------+

# 禁用索引为3和5的订阅链接
>>> sub atler 3,5 --using false
>>> sub show 
+------+-------+---------------------+----------+
| 索引 |  别名  |         URL         |  是否启用  |
+------+-------+---------------------+----------+
|  0   | test0 |  https://test.com   |   true   |
|  1   | test1 | https://sublink.com |   true   |
|  2   | test2 | https://sublink.com |   true   |
|  3   | test3 | https://sublink.com |  false   |
|  4   | test4 | https://sublink.com |   true   |
|  5   | test5 | https://sublink.com |  false   |
|  6   | test6 | https://sublink.com |   true   |
+------+-------+---------------------+----------+
```

### 删除订阅

```
# 删除索引为3和5的订阅
>>> sub del 3,5

# 删除所有订阅
>>> sub del all
```

### 从订阅更新节点

```
# 不使用代理更新节点
>>> sub update-node

# 使用端口为2333的本地sock5代理更新节点
>>> sub update-node -p 2333
```



## 查看节点帮助文档

```
>>> node

node {commands} [flags] ...

Commands:
    add     [flags]                    添加节点
    show    [索引范围 | tcping]         查看节点信息, 默认索引范围为all, tcping可以按延迟排序查看
    del     {索引范围}                  根据索引参数删除节点
    export  {索引范围}                  导出为vmess链接
    tcping  [索引范围]                  tcping指定索引节点, 默认索引范围为all
    find    {关键词}                    查找节点，有中文关键词需要用单引号或双引号括起来

add Flags
    -v, --vmess {vmess链接}             导入vmess://数据
    -f, --file  {文件绝对路径}           从文件批量导入vmess://数据
```

### 添加节点

```
# 由vmess链接添加节点
>>> node add  -v vmess://xxxxxxXXXXxxxxxXX

# 由vmess链接文件批量添加节点
>>> node add -f /home/vmess.txt

# 手动添加一个节点
>>> node add
```

### 查看节点

```
# 查看前20个节点
>>> node show 0-19

# 查看按tcp延迟排序的节点
>>> node show tcping
```

### 删除节点

```
# 删除前20个节点
>>> node del 0-19
```

### tcping测试

```
# tcping前20个节点
>>> node tcping 0-19
```

### 节点查找

```
# 查找关键词为'vip'的节点
>>> node find vip

# 查找关键词为'香港'的节点
>>> node find "香港"
```

### 导出节点

```
# 导出前20个节点到终端
>>> node export 0-19
```

## 查看DNS帮助文档

```
dns {commands} ...

Commands:
    add  {dns}                         添加dns
    show                               查看dns列表
    del  {索引范围}                     删除索引范围内的dns
```

### 添加DNS

```
# 添加DNS
>>> node add 114.114.114.114
```

### 查看DNS

```
# 查看所有的DNS
>>> node show
```

### 删除DNS

```
# 删除前2个dns
>>> node del 0,1
```

### 

## 查看路由帮助文档

```
>>> route

route {commands} ...

Commands:

    show-direct-ip                     查看直连ip规则
    show-direct-domain                 查看直连domain规则
    show-proxy-ip                      查看代理ip规则
    show-proxy-domain                  查看代理domain规则
    show-block-ip                      查看禁止ip规则
    show-block-domain                  查看禁止domain规则
    
    add-direct-ip       {路由规则}      添加一条直连ip规则
    add-direct-domain   {路由规则}      添加一条直连domain规则
    add-proxy-ip        {路由规则}      添加一条代理ip规则
    add-proxy-domain    {路由规则}      添加一条代理domain规则
    add-block-ip        {路由规则}      添加一条禁止ip规则
    add-block-domain    {路由规则}      添加一条禁止domain规则
    
    del-direct-ip       {索引范围}      删除索引范围的直连ip路由规则
    del-direct-domain   {索引范围}      删除索引范围的直连domain路由规则
    del-proxy-ip        {索引范围}      删除索引范围的代理ip路由规则
    del-proxy-domain    {索引范围}      删除索引范围的代理domain路由规则
    del-block-ip        {索引范围}      删除索引范围的禁止ip路由规则
    del-block-domain    {索引范围}      删除索引范围的禁止domain路由规则
```

### domain路由规则

- 纯字符串: 当此字符串匹配目标域名中任意部分，该规则生效。比如`sina.com`可以匹配`sina.com`、sina.com.cn和www.sina.com，但不匹配`sina.cn`。
- 正则表达式: 由`regexp:`开始，余下部分是一个正则表达式。当此正则表达式匹配目标域名时，该规则生效。例如`regexp:\\.goo.*\\.com$`匹配`www.google.com`、`fonts.googleapis.com`，但不匹配`google.com`。
- 子域名 (推荐): 由`domain:`开始，余下部分是一个域名。当此域名是目标域名或其子域名时，该规则生效。例如`domain:v2ray.com`匹配`www.v2ray.com`、`v2ray.com`，但不匹配`xv2ray.com`。
- 完整匹配: 由`full:`开始，余下部分是一个域名。当此域名完整匹配目标域名时，该规则生效。例如`full:v2ray.com`匹配`v2ray.com`但不匹配`www.v2ray.com`。
- 预定义域名列表：由`"geosite:"`开头，余下部分是一个名称，如`geosite:google`或者`geosite:cn`。名称及域名列表参考[预定义域名列表](https://www.v2ray.com/chapter_02/03_routing.html#dlc)。
- 从文件中加载域名: 形如`ext:file:tag`，必须以`ext:`（小写）开头，后面跟文件名和标签，文件存放在[资源目录](https://www.v2ray.com/chapter_02/env.html#asset-location)中，文件格式与`geosite.dat`相同，标签必须在文件中存在。

### ip路由规则

- IP: 形如`127.0.0.1`。
- [CIDR](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing): 形如`10.0.0.0/8`.
- GeoIP: 形如`geoip:cn`，必须以`geoip:`（小写）开头，后面跟双字符国家代码，支持几乎所有可以上网的国家。
- 特殊值：`geoip:private` (V2Ray 3.5+)，包含所有私有地址，如`127.0.0.1`。
- 从文件中加载 IP: 形如`ext:file:tag`，必须以`ext:`（小写）开头，后面跟文件名和标签，文件存放在[资源目录](https://www.v2ray.com/chapter_02/env.html#asset-location)中，文件格式与`geoip.dat`相同标签必须在文件中存在。



## 查看服务帮助文档

  ```
>>> service

service {commands} ...

Commands:
    start  [节点索引]                   启动或重启v2ray-core服务
    stop                               停止v2ray-core服务
  ```

### 启动或重启v2ray-core服务

```
# 启动或重启索引为3的节点
>>> service start 3
```

### 停止v2ray-core服务

```
>>> service stop
```



# 已知问题

- Mac/Linux系统下节点不能使用ping

# 交流反馈

提交Issue: [Issues](https://github.com/hsernos/v3ray/issues)

邮箱: hsernos@163.com
