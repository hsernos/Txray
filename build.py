#!/usr/bin/env python3
# -*- coding:utf-8 -*-
import sys
import os

INFO = {
    'linux': ['386', 'amd64', 'arm'],
    'darwin': ['386', 'amd64'],
    'freebsd': ['386', 'amd64'],
    'windows': ['386', 'amd64']
}

Version = 'v1.0.0'
Name = 'v3ray'


def get_os():
    platform = sys.platform
    if platform == 'linux2':
        return 'linux'
    elif platform == 'win32':
        return 'win'
    elif platform == 'darwin':
        return 'mac'
    else:
        return 'other'


def build(goos, goarch, path, cgo=0):
    if get_os() == 'win':
        cmd1 = 'SET CGO_ENABLED={}'.format(cgo)
        cmd2 = 'SET GOOS={}'.format(goos)
        cmd3 = 'SET GOARCH={}'.format(goarch)
        e = '.exe' if goos == 'windows' else ''
        cmd4 = 'go build -o {} {}'.format("build/"+"-".join([Name, Version, goos, goarch]) + e, path)
        cmd = '\n'.join([cmd1, cmd2, cmd3, cmd4])
        os.system(cmd)
    else:
        e = '.exe' if goos == 'windows' else ''
        cmd = 'CGO_ENABLED={} GOOS={} GOARCH={} go build -o {} {} '.format(cgo, goos, goarch,
                                                                           "build/"+"-".join([Name, Version, goos, goarch]) + e,
                                                                           path)
        os.system(cmd)


if __name__ == '__main__':
    for goos in INFO.keys():
        for goarch in INFO[goos]:
            print("正在编译", goos, '系统的', goarch, '架构版本中...')
            build(goos, goarch, 'v3ray.go')
