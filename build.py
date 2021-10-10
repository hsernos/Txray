#!/usr/bin/env python
# -*- coding:utf-8 -*-

import sys
import os
import shutil
import zipfile

INFO = {
    'linux': ['386', 'amd64', 'arm','arm64'],
    'darwin': ['amd64', 'arm64'],
    'windows': ['386', 'amd64'],
    'freebsd': ['386', 'amd64', 'arm', 'arm64'],
    'openbsd': ['386', 'amd64', 'arm', 'arm64'],
}

Name = 'Txray'


# 打包目录为zip文件（未压缩）
def make_zip(source_dir, output_filename):
    zipf = zipfile.ZipFile(output_filename, 'w')
    pre_len = len(os.path.dirname(source_dir))
    for parent, dirnames, filenames in os.walk(source_dir):
        for filename in filenames:
            pathfile = os.path.join(parent, filename)
            arcname = pathfile[pre_len:].strip(os.path.sep)  # 相对路径
            zipf.write(pathfile, arcname)
    zipf.close()


def build(goos, goarch, path='Txray.go', cgo=0):
    e = '.exe' if goos == 'windows' else ''
    syst = 'macos' if goos == 'darwin' else goos
    arch = '64' if goarch == 'amd64' else goarch
    arch = '32' if goarch == '386' else arch
    os.environ["CGO_ENABLED"] = str(cgo)
    os.environ["GOOS"] = goos
    os.environ["GOARCH"] = goarch
    cmd = 'go build -o {} {} '.format("build/" +
                                      "-".join([Name, syst, arch]) + "/" + Name + e, path)
    os.system(cmd)
    shutil.copy("README.md", "build/"+"-".join([Name, syst, arch]))
    make_zip("build/"+"-".join([Name, syst, arch]),
             "build/"+"-".join([Name, syst, arch])+".zip")

# 常用系统架构编译
def default(info=INFO):
    for goos in info.keys():
        for goarch in INFO[goos]:
            print("正在编译", goos, '系统的', goarch, '架构版本中...')
            build(goos, goarch)

# 获取go支持的系统和架构
def get_all_info():
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
    if len(sys.argv) == 2 and sys.argv[1] == "-d":
        default()
    else:
        select_one()
