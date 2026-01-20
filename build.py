# build.py 用于跨平台自动编译、打包 Txray 项目，支持多系统多架构
#!/usr/bin/env python
# -*- coding:utf-8 -*-

import sys
import os
import shutil
import zipfile

# INFO 字典定义了各主流操作系统支持的 CPU 架构
INFO = {
    'linux': ['386', 'amd64', 'arm','arm64'],
    'darwin': ['amd64', 'arm64'],
    'windows': ['386', 'amd64','arm','arm64'],
    'freebsd': ['386', 'amd64', 'arm', 'arm64'],
    'openbsd': ['386', 'amd64', 'arm', 'arm64'],
    'android':['arm64']
}

Name = 'Txray'  # 项目名称


# 打包目录为zip文件（未压缩）
def make_zip(source_dir, output_filename):
    """
    将 source_dir 目录下所有文件打包为 output_filename zip 文件
    :param source_dir: 需要打包的目录
    :param output_filename: 输出的 zip 文件名
    """
    zipf = zipfile.ZipFile(output_filename, 'w')
    pre_len = len(os.path.dirname(source_dir))
    for parent, dirnames, filenames in os.walk(source_dir):
        for filename in filenames:
            pathfile = os.path.join(parent, filename)
            arcname = pathfile[pre_len:].strip(os.path.sep)  # 相对路径
            zipf.write(pathfile, arcname)
    zipf.close()


def build(goos, goarch, path='Txray.go', cgo=0):
    """
    编译指定系统和架构的可执行文件，并打包为 zip
    :param goos: 目标操作系统
    :param goarch: 目标架构
    :param path: 主 go 文件路径
    :param cgo: 是否启用 CGO
    """
    e = '.exe' if goos == 'windows' else ''
    syst = 'macos' if goos == 'darwin' else goos
    arch = '64' if goarch == 'amd64' else goarch
    arch = '32' if goarch == '386' else arch
    os.environ["CGO_ENABLED"] = str(cgo)
    os.environ["GOOS"] = goos
    os.environ["GOARCH"] = goarch
    # 构造 go build 命令
    cmd = 'go build -o {} {} '.format("build/" +
                                      "-".join([Name, syst, arch]) + "/" + Name + e, path)
    os.system(cmd)
    # 拷贝 README.md 到目标目录
    shutil.copy("README.md", "build/"+"-".join([Name, syst, arch]))
    # 打包为 zip
    make_zip("build/"+"-".join([Name, syst, arch]),
             "build/"+"-".join([Name, syst, arch])+".zip")

# 常用系统架构编译
def default(info=INFO):
    """
    按预设的 INFO 字典批量编译所有主流系统和架构
    """
    for goos in info.keys():
        for goarch in INFO[goos]:
            print("正在编译", goos, '系统的', goarch, '架构版本中...')
            build(goos, goarch)

# 获取go支持的系统和架构
def get_all_info():
    """
    获取本地 go 环境支持的所有系统/架构组合
    :return: dict，key为系统，value为架构列表
    """
    result = os.popen('go tool dist list')
    res = result.read()
    d = dict()
    for line in res.splitlines():
        l = line.split("/")
        system = l[0]
        arch = l[1]
        if system in d.keys():
            d[system].append(arch)
        else:
            d.setdefault(system, [arch])
    return d

# 选择go支持的系统架构编译，不一定都能成功编译
def select_one():
    """
    交互式选择系统和架构进行编译
    """
    data = get_all_info()
    system_list = list(data.keys())
    i = 1
    for key in system_list:
        key = 'macos' if key == 'darwin' else key
        print("{}. {}".format(i, key))
        i += 1
    index = int(input("选择你要编译的操作系统序号(1-{}): ".format(len(system_list))))
    if 1 <= index <= len(system_list):
        system = system_list[index-1]
        arch_list = data.get(system)
        j = 1
        for arch in arch_list:
            print("{}. {}".format(j, arch))
            j += 1
        index = int(input("选择你要编译的架构序号(1-{}): ".format(len(arch_list))))
        if 1 <= index <= len(arch_list):
            arch = arch_list[index-1]
            print("正在编译", system, '系统的', arch, '架构版本中...')
            build(system, arch)


if __name__ == '__main__':
    # 命令行参数 -d 走批量编译，否则交互式选择
    if len(sys.argv) == 2 and sys.argv[1] == "-d":
        default()
    else:
        select_one()
