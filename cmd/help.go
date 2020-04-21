package cmd

var settingHelp = `
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
`

var nodeHelp = `
node {commands} [flags] ...

Commands:
    add     [flags]                     添加节点
    show    [索引范围]                  查看节点信息, 默认索引范围为all
    del     {索引范围}                  根据索引参数删除节点
    export  {索引范围}                  导出为vmess链接
    ping    [索引范围]                  ping指定索引节点, 默认索引范围为all

add Flags
    -v, --vmess {vmess链接}             导入vmess://数据
    -f, --file  {文件绝对路径}          从文件批量导入vmess://数据
`
var subHelp = `
sub {commands} [flags] ...

Commands:
    add   {订阅链接} [flags]             添加订阅
    alter {索引范围} [flags]             修改订阅
    show  [索引范围]                     查看订阅信息
    del   {索引范围}                     根据索引参数删除订阅
    update-node     [flags]              从订阅更新节点

add Flags
    -r, --remarks {别名}                 自定义别名

alter Flags
    -u, --url     {订阅链接}             订阅链接
    -r, --remarks {别名}                 自定义别名
    --using       {true|false}           是否使用此订阅

update-node Flags
    -p, --proxy {本地socks5端口}         从socks5代理更新节点
`

var dnsHelp = `
dns {commands} ...

Commands:
    add  {dns}                          添加dns
    show                                查看dns列表
    del  {索引范围}                     删除索引范围内的dns
`

var routeHelp = `
route {commands} ...

Commands:

    show-direct-ip                      查看直连ip规则
    show-direct-domain                  查看直连domain规则
    show-proxy-ip                       查看代理ip规则
    show-proxy-domain                   查看代理domain规则
    show-block-ip                       查看禁止ip规则
    show-block-domain                   查看禁止domain规则

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
`

var serviceHelp = `
service {commands} ...

Commands:
    start  [节点索引]                   启动v2ray-core服务
    stop                                停止v2ray-core服务
`

var help = `
Commands:
	setting      基础设置 		使用 setting 查看帮助信息
	node         节点管理 		使用 node 查看帮助信息
	sub          订阅管理 		使用 sub 查看帮助信息
	dns          DNS管理 		使用 dns 查看帮助信息
	route        路由管理		使用 route 查看帮助信息
	service      v2ray服务管理 	使用 service 查看帮助信息
	help         查看帮助信息
	clear        清屏
	exit         退出程序
`