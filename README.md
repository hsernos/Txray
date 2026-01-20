# Txray

Txray 是一款 xray 终端版客户端，使用 go 编写。

项目地址：https://github.com/hsernos/Txray

Project X core： https://github.com/XTLS/Xray-core

## 注意

此文档只针对于最新的 commit, 可能不适用于已发布的最新版本.

## 特色

1.多平台支持, 支持 Windows, macOS, linux.

2.Tab 键命令补齐

3.支持 VMess、Shadowsocks、Trojan、VLESS、VMessAEAD、Socks 协议

## 安装和使用

1.下载对应平台架构的[Txray](https://github.com/hsernos/Txray/releases)和[xray](https://github.com/XTLS/Xray-core/releases),按如下目录结构解压放置

```
Txray（目录命名随意）
│   Txray
│   README.md
│
└───xray-core （目录命名随意）
    │   xray
    │   geoip.dat
    |   geosite.dat
    │   ...

```

2.将 Txray 程序所在目录添加到 PATH 环境变量（添加环境变量请自行谷歌或百度）

3.非 Windows 平台用户添加`Txray`和`xray`可执行权限

```
# 进入 Txray 所在目录执行 'chmod u+x Txray'
# 为 Txray 添加可执行权限
[xxx@xxx Txray-linux-64]$ chmod u+x Txray
```

```
# 进入xray所在目录执行 'chmod u+x xray'
# 为 xray 添加可执行权限
[xxx@xxx Xray-linux-64]$ chmod u+x xray
```

4.打开终端输入`Txray`回车进入 Shell 交互 或 继续在末尾添加命令直接运行

```
# 1.shell交互运行，命令可Tab补齐
[xxx@xxx  xxxx]$ Txray

# 2.直接运行，如更新节点 (不会进入shell交互)
[xxx@xxx  xxxx]$ Txray sub update-node
```

## 稍高级使用

1.单电脑多系统共用同一份配置文件（配置环境变量`TXRAY_HOME`）

- Txray 检测 xray-core 所在的优先级: 环境变量 `CORE_HOME` > Txray 所在目录 > 环境变量 `PATH` 中的目录
- Txray 检测 geosite.dat和geoip.dat文件目录 环境变量 `XRAY_LOCATION_ASSET` > xray所在目录
- 配置文件目录优先级: 环境变量 `TXRAY_HOME` > Txray 所在目录

  2.开机自启，请自行谷歌或百度查找对应系统的开机自启脚本的写法

PS：开机自启推荐搭配[命令别名](#查看命令别名帮助文档)使用

<!-- toc -->

## 目录

- [编译/交叉编译 说明](#编译交叉编译-说明)
- [命令列表及说明 (以下全为Shell交互下的演示说明)](#命令列表及说明)
  - [命令总览](#命令总览)
  - [查看基本设置帮助文档](#查看基本设置帮助文档)
    - [查看基本设置](#查看基本设置)
    - [修改基本设置](#修改基本设置)
  - [查看订阅帮助文档](#查看订阅帮助文档)
    - [添加订阅](#添加订阅)
    - [查看订阅](#查看订阅)
    - [修改订阅](#修改订阅)
    - [删除订阅](#删除订阅)
    - [从订阅更新节点](#从订阅更新节点)
  - [查看节点帮助文档](#查看节点帮助文档)
    - [添加节点](#添加节点)
    - [查看节点](#查看节点)
    - [删除节点](#删除节点)
    - [tcping 测试](#tcping测试)
    - [节点查找](#节点查找)
    - [导出节点](#导出节点)
    - [节点排序](#节点排序)
  - [查看节点过滤器帮助文档](#查看节点过滤器帮助文档)
    - [添加过滤器规则](#添加过滤器规则)
    - [过滤节点](#过滤节点)
  - [查看节点回收站帮助文档](#查看节点回收站帮助文档)
  - [查看命令别名帮助文档](#查看命令别名帮助文档)
    - [添加和修改别名](#添加和修改别名)
  - [查看路由帮助文档](#查看路由帮助文档)
    - [添加路由](#添加路由)
    - [domain 路由规则](#domain路由规则)
    - [ip 路由规则](#ip路由规则)
  - [启动或重启 xray-core 服务](#启动或重启xray-core服务)
  - [停止 xray-core 服务](#停止xray-core服务)
  - [显示运行时 xray-core 的日志](#显示运行时xray-core的日志)
- [已知问题](#已知问题)
- [交流反馈](#交流反馈)

<!-- tocstop -->

# 编译/交叉编译 说明

1. 在终端下进入项目目录

2. 设置`GOPROXY`,提高编译所需依赖的下载速度
   ## Linux/Mac
   启用 Go Modules 功能
   `go env -w GO111MODULE=on`
   
   设置 Go 模块代理，选一即可
   1. 七牛 CDN 
   `go env -w  GOPROXY=https://goproxy.cn,direct`
   2. 阿里云
   `go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct`
   3. 腾讯云
   `go env -w GOPROXY=https://mirrors.cloud.tencent.com/go/,direct`
   4. 官方
   `go env -w  GOPROXY=https://goproxy.io,direct`

   ## windows
   启用 Go Modules 功能
   `env:GO111MODULE="on"`

   设置 Go 模块代理，选一即可
   1. 七牛 CDN 
   `set GOPROXY=https://goproxy.cn,direct`
   2. 阿里云
   `env:GOPROXY="https://mirrors.aliyun.com/goproxy/,direct`
   3. 腾讯云
   `env:GOPROXY=https://mirrors.cloud.tencent.com/go/,direct`
   4. 官方
    `env:GOPROXY="https://goproxy.io,direct`

   ## 检查
   `go env | grep GOPROXY`


3. 编译常用平台
   运行 `go build Txray.go`, 可编译当前平台的版本
   运行 `python3 build.py`, 可编译常用平台的版本

4. 编译其他平台
   运行 `go tool dist list` 查看所有支持的 GOOS/GOARCH

   Linux/Darwin 例子: 编译 Windows 下的 64 位程序

   `GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build Txray.go`

   Windows 例子: 编译 Linux 下的 32 位程序

   `set GOOS=linux`
   `set GOARCH=386`
   `set CGO_ENABLED=0`
   `go build Txray.go`

# 命令列表及说明

> 在终端中运行 Txray 进入 shell 交互

## 命令总览

```
Commands:
    setting                  基础设置             使用 'setting help' 查看详细用法
    node                     节点管理             使用 'node help' 查看详细用法
    sub                      订阅管理             使用 'sub help' 查看详细用法
    routing                  路由管理             使用 'routing help' 查看详细用法
    filter                   节点过滤             使用 'filter help' 查看详细用法
    recycle                  回收站               使用 'recycle help' 查看详细用法
    alias                    命令别名             使用 'alias help' 查看详细用法
    help, -h                 查看帮助信息
    version, -v              查看版本
    clear                    清屏
    exit                     退出程序
    run                      启动或重启节点
    stop                     关闭节点
    log                      查看运行时xray日志

Usage: run [索引式]
    run [索引式]      默认为上一次运行节点，如果选中多个节点，则选择访问 'setting' 中测试国外URL延迟最小的


说明：
一、索引式：更简单的批量选择
1.选择前6个：'1,2,3,4,5,6' 或 '1-3,4-6' 或 '1-6' 或 '-6'
2.选择第6个及后面的所有：'6-'
3.选择第6个：'6'
4.选择所有：'all' 或 '-'
注意：超出部分会被忽略，'all' 只能单独使用

二、[] 和 {}：帮助说明中的中括号和大括号
1. []: 表示该选项可忽略
2. {}: 表示该选项为必须，不可忽略
```

## 查看基本设置帮助文档

```
>>> setting help
setting {commands}

Commands:
                                  查看所有设置
    help                          查看帮助
    
    mixed [port]                  设置混合端口
    socks [port]                  设置socks端口
    http [port]                   设置http端口, 0为关闭http监听
    udp [y|n]                     是否启用udp转发
    sniffing [y|n]                是否启用流量地址监听
    from_lan_conn [y|n]           是否启用来自局域网连接
    mux [y|n]                     是否启用多路复用（下载和看视频时建议关闭）

    dns.port [port]               设置DNS端口
    dns.foreign [dns]             设置国外DNS
    dns.domestic [dns]            设置国内DNS
    dns.backup [dns]              设置国内备用DNS

    routing.strategy {1|2|3}      设置路由策略为{AsIs|IPIfNonMatch|IPOnDemand}
    routing.bypass {y|n}          是否绕过局域网及大陆

    version.min [version]         设置版本最小值，为空表示不限制
    version.max [version]         设置版本最大值，为空表示不限制

    test.url [url]                设置外网测试URL
    test.timeout [time]           设置外网测试超时时间 (秒)
    test.mintime [time]           设置批量测试终止时间 (毫秒)

    run_before [命令组] [flags]    程序启动时执行命令或命令组，可与命令别名搭配


run_before Flags
    -c, --close                   启动时不执行任何命令

说明：
1.命令，如 'node' 'node tcping' 'sub update-node' 这样的单条命令。
2.命令组，形如 'sub update-node | node tcping | run' 这样的多条命令，以 '|' 分隔，顺序执行。
PS：命令组包含命令，即命令组也可以设置单条命令
```

### 查看基本设置

```
>>> setting
+-----------+----------+---------+--------------+--------------------+----------+
| SOCKS端口 | HTTP端口  | UDP转发 | 流量地址监听  | 允许来自局域网连接   | 多路复用  |
+-----------+----------+---------+--------------+--------------------+----------+
|   1080    |    0     |  true   |     true     |        true        |   true   |
+-----------+----------+---------+--------------+--------------------+----------+
+---------+---------+--------------+-----------------+------------+------------------+
| DNS端口  | 国外DNS |   国内DNS    |   备用国内DNS    |  路由策略   | 绕过局域网和大陆  |
+---------+---------+--------------+-----------------+------------+------------------+
|  1351   | 1.1.1.1 | 119.29.29.29 | 114.114.114.114 | IPOnDemand |       true       |
+---------+---------+--------------+-----------------+------------+------------------+
+-------------------------+-------------------+-------------------------+------------+
|       测试国外URL        | 测试超时时间 (秒)  | 批量测试终止时间 (毫秒)   | 启动时执行  |
+-------------------------+-------------------+-------------------------+------------+
| https://www.youtube.com |         6         |          1000           |            |
+-------------------------+-------------------+-------------------------+------------+
```

### 修改基本设置

```
# 修改socks监听端口为3333
>>> setting socks 3333

# 修改http监听端口为3334
>>> setting http 3334

# 修改不绕过局域网和大陆
>>> setting routing.bypass n

# 修改路由策略为IPIfNonMatch, {1|2|3}=>{AsIs|IPIfNonMatch|IPOnDemand}
>>> setting routing.strategy 2

# 启动时运行仅针对进入shell交互才会触发
# 设置启动时从订阅更新节点
>>> setting run_before "sub update-node"

# 设置启动时对节点进行tcp测试，然后运行延迟最小的那个
>>> setting run_before "node tcping | run"

# 设置启动时不执行任何命令
>>> setting run_before -c

# 设置批量测试终止时间为1000,即节点测试延迟在1~1000ms内就会停止，不会继续测试后续节点
>>> setting test.mintime 1000
```

## 查看订阅帮助文档

```
>>> sub help
sub {commands} [flags] ...

Commands:
                                 查看订阅信息
    help                         查看帮助
    rm {索引式}                   删除订阅
    add {订阅url} [flags]         添加订阅
    mv {索引式} {flags}           修改订阅
    update-node [索引式] [flags]  从订阅更新节点, 索引式会忽略是否启用

add Flags
    -r, --remarks {别名}          定义别名

mv Flags
    -u, --url {订阅url}           修改订阅链接
    -r, --remarks {别名}          定义别名
    --using {y|n}                 是否启用此订阅

update-node Flags
    -s, --socks5 [port]           通过本地的socks5代理更新, 默认为设置中的socks5端口
    -h, --http [port]             通过本地的http代理更新, 默认为设置中的http端口
    -a, --addr {address}          对上面两个参数的补充, 修改代理地址
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
>>> sub
+------+-------+---------------------+----------+
| 索引  | 别名   |       URL          |  是否启用  |
+------+-------+---------------------+----------+
|  1   | test1 | https://sublink.com |   true   |
|  2   | test2 | https://sublink.com |   true   |
|  3   | test3 | https://sublink.com |   true   |
|  4   | test4 | https://sublink.com |   true   |
|  5   | test5 | https://sublink.com |   true   |
|  6   | test6 | https://sublink.com |   true   |
+------+-------+---------------------+----------+

# 查看索引为2,3,4的订阅
>>> sub 2-4
+------+-------+---------------------+----------+
| 索引  | 别名  |         URL         |  是否启用  |
+------+-------+---------------------+----------+
|  2   | test2 | https://sublink.com |   true   |
|  3   | test3 | https://sublink.com |   true   |
|  4   | test4 | https://sublink.com |   true   |
+------+-------+---------------------+----------+
```

### 修改订阅

```
# 修改索引为1的订阅链接为https://test.com，别名为test8
>>> sub mv 1 -u https://test.com -r test8
>>> sub 1
+------+-------+------------------+----------+
| 索引  | 别名  |        URL        |  是否启用 |
+------+-------+------------------+----------+
|  1   | test8 | https://test.com |   true   |
+------+-------+------------------+----------+

# 禁用索引为3和5的订阅链接
>>> sub mv 3,5 --using n
>>> sub
+------+-------+---------------------+----------+
| 索引 |  别名  |         URL         |  是否启用  |
+------+-------+---------------------+----------+
|  1   | test8 | https://sublink.com |   true   |
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
>>> sub rm 3,5

# 删除所有订阅
>>> sub rm all
```

### 从订阅更新节点

```
# 从启用的订阅且不使用代理更新节点
>>> sub update-node

# 从索引范围更新节点，无论是否启用
>>> sub update-node 1,3,6

# 使用端口为2333的本地socks5代理更新节点
>>> sub update-node -s 2333

# 使用设置中的socks端口通过本地socks5代理更新节点
>>> sub update-node -s

# 使用端口为2334的本地http代理更新节点
>>> sub update-node -h 2334

# 使用端口为2333，地址为1.2.3.4的socks代理更新节点
>>> sub update-node -s 2333 -a 1.2.3.4
```

## 查看节点帮助文档

```
>>> node help
node {commands} [flags] ...

Commands:
    [索引式] [flags]              查看节点信息, 默认 'all'
    help                          查看帮助
    tcping                        测试节点tcp延迟
    sort {0|1|2|3|4|5}            排序方式，分别按{逆转|协议|别名|地址|端口|测试结果}排序
    info {索引}                   查看单个节点详细信息
    rm {索引式}                   删除节点
    find {关键词}                 查找节点（按别名）
    add [flags]                   添加节点
    export [索引式] [flags]       导出节点链接, 默认'all'

Flags
    -d, --desc                    降序查看

add Flags
    -l, --link {link}             从链接导入一条节点
    -f, --file {path}             从节点链接文件或订阅文件导入节点
    -c, --clipboard               从剪贴板读取的节点链接或订阅文本导入节点

export Flags
    -c, --clipboard               导出节点链接到剪贴板
```

### 添加节点

```
# 添加一个vmess节点
>>> node add  -l vmess://xxxxxxXXXXxxxxxXX

# 添加一个trojan节点
>>> node add  -l trojan://xxxxxxXXXXxxxxxXX

# 由链接文件批量添加节点
>>> node add -f /home/links.txt

# 解析订阅文件添加节点，可以将订阅文件下载下来然后从本地导入
>>> node add -f /home/subtext.txt

# 从剪贴板读取的节点链接或订阅文本导入节点, 功效和上面从文件导入一样
>>> node add -c

# 手动添加一个节点
>>> node add
```

### 查看节点

```
# 查看前20个节点
>>> node 1-20

# 降序查看所有节点
>>> node -d

# 查看某个节点的全部信息
>>> node info 1

```

### 删除节点

```
# 删除前20个节点
>>> node rm 1-20
```

### tcping 测试

```
# tcping测试所有节点
>>> node tcping
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
>>> node export -20

# 导出前20个节点到剪贴板
>>> node export -20 -c
```

### 节点排序

```
# 逆转节点顺序
>>> node sort 0

# 按别名排序
>>> node sort 2
```

## 查看节点过滤器帮助文档

> 过滤器会在添加节点的时候自动运行，也可以使用 'filter run' 手动运行

```
>>> filter help
filter {commands} ...

Commands:
                                 查看过滤规则
    help                         查看帮助
    rm {索引式}                   删除过滤规则
    open {索引式}                 开启过滤规则
    close {索引式}                关闭过滤规则
    add {过滤规则}                添加过滤规则
    run [过滤规则]                手动过滤节点，默认使用内置规则

PS: 过滤规则==> '过滤范围:正则表达式'
过滤范围可选值 proto:|name:|addr:|port: 分别代表 协议|别名|地址|端口
默认为 'name:'
```

### 添加过滤器规则

```
# 添加地址为baidu.com的过滤规则
>>> filter add addr:baidu.com

# 添加协议为VMess的过滤规则
>>> filter add proto:VMess

# 添加别名含有'美国'的过滤规则
>>> filter add "美国"
```

### 过滤节点

```
# 删除地址为baidu.com的节点
>>> filter run addr:baidu.com

# 运行已有规则
>>> filter run
```

## 查看节点回收站帮助文档

> 数据仅保存当次交互中

```
>>> recycle help
recycle {commands} ...

Commands:
    {索引式}                      查看节点回收站
    help                         查看帮助
    restore {索引式}              恢复节点
    clear                        清空节点回收站

PS: 回收站的数据仅运行存在 (仅存储在内存中)
```

## 查看命令别名帮助文档

> 不能覆盖自带命令，小心使用，不要弄成死循环了

```
>>> alias help
alias {commands} ...

Commands:
                                 查看命令别名
    help                         查看帮助
    set {别名} {命令组}           开启过滤规则
    rm {索引式}                   删除命令别名

说明：
1.命令，如 'node' 'node tcping' 'sub update-node' 这样的单条命令。
2.命令组，形如 'sub update-node | node tcping | run' 这样的多条命令，以 '|' 分隔，顺序执行。
PS：命令组包含命令，即命令组也可以设置单条命令
```

### 添加和修改别名

```
# 设置别名 'one' 为更新订阅，然后tcping，最后运行延迟最小的那个
>>> alias set one "sub update-node | node tcping | run"

# 运行
>>> one

# 等价于 'sub update-node | node tcping | run 0-10'
>>> one 0-10
```

## 查看路由帮助文档

```
>>> route help
routing {commands} [flags] ...

Commands:
    help                          查看帮助
    block [索引式] | [flags]      查看或管理禁止路由规则
    direct [索引式] | [flags]     查看或管理直连路由规则
    proxy [索引式] | [flags]      查看或管理代理路由规则

block, direct, proxy Flags
    -a, --add {规则}              添加路由规则
    -r, --rm {索引式}             删除路由规则
    -f, --file {path}             从文件导入规则
    -c, --clipboard               从剪贴板导入规则

PS: 规则详情请访问 https://xtls.github.io/config/routing.html#ruleobject
```

### 添加路由

```
# 添加www.baidu.com到黑名单
>>> routing block -a www.baidu.com

# 添加www.google.com到代理名单
>>> routing proxy -a www.google.com

# 从文件批量导入到黑名单
>>> routing block -f /home/xxx/block.txt

# 从剪贴板导入到黑名单
>>> routing block -c

```

### domain路由规则

- 纯字符串: 当此字符串匹配目标域名中任意部分，该规则生效。比如`sina.com`可以匹配`sina.com`、sina.com.cn 和www.sina.com，但不匹配`sina.cn`。
- 正则表达式: 由`regexp:`开始，余下部分是一个正则表达式。当此正则表达式匹配目标域名时，该规则生效。例如`regexp:\\.goo.*\\.com$`匹配`www.google.com`、`fonts.googleapis.com`，但不匹配`google.com`。
- 子域名 (推荐): 由`domain:`开始，余下部分是一个域名。当此域名是目标域名或其子域名时，该规则生效。例如`domain:xray.com`匹配`www.xray.com`、`xray.com`，但不匹配`xxray.com`。
- 完整匹配: 由`full:`开始，余下部分是一个域名。当此域名完整匹配目标域名时，该规则生效。例如`full:xray.com`匹配`xray.com`但不匹配`www.xray.com`。
- 预定义域名列表：由`"geosite:"`开头，余下部分是一个名称，如`geosite:google`或者`geosite:cn`。名称及域名列表参考[预定义域名列表](https://www.v2ray.com/chapter_02/03_routing.html#dlc)。
- 从文件中加载域名: 形如`ext:file:tag`，必须以`ext:`（小写）开头，后面跟文件名和标签，文件存放在[资源目录](https://www.v2ray.com/chapter_02/env.html#asset-location)中，文件格式与`geosite.dat`相同，标签必须在文件中存在。

### ip路由规则

- IP: 形如`127.0.0.1`。
- [CIDR](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing): 形如`10.0.0.0/8`.
- GeoIP: 形如`geoip:cn`，必须以`geoip:`（小写）开头，后面跟双字符国家代码，支持几乎所有可以上网的国家。
- 特殊值：`geoip:private` (xray 3.5+)，包含所有私有地址，如`127.0.0.1`。
- 从文件中加载 IP: 形如`ext:file:tag`，必须以`ext:`（小写）开头，后面跟文件名和标签，文件存放在[资源目录](https://www.v2ray.com/chapter_02/env.html#asset-location)中，文件格式与`geoip.dat`相同标签必须在文件中存在。

## 启动或重启xray-core服务

```
# 启动或重启索引为3的节点
>>> run 3

# 自动选择所有节点中访问YouTube延迟最小的那个节点
>>> run all

# 自动选择1-10中访问YouTube延迟最小的那个节点
>>> run 1-10

# 自动选择tcp延迟最小的10个中访问YouTube延迟最小的那个节点
>>> run -t -10
```

## 停止xray-core服务

```
# 停止上次启动的xray-core进程
>>> stop

```

## 显示运行时xray-core的日志

```
# 显示运行时 xray-core 的日志
>>> log

```

# 已知问题

- 有时直接从订阅更新节点失败，可以用浏览器下载订阅文本然后使用 'node add -f {绝对路径}' 导入，或者使用代理导入（sub update-node -s [端口]）

# 交流反馈

提交 Issue: [Issues](https://github.com/hsernos/Txray/issues)
