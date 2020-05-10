#!/usr/bin/env python3
# -*- coding:utf-8 -*-
import sys
import os
import shutil
import zipfile

INFO = {
    'linux': ['386', 'amd64', 'arm'],
    'darwin': ['amd64'],
    'freebsd': ['386', 'amd64'],
    'windows': ['386', 'amd64']
}

Name = 'v3ray'




#打包目录为zip文件（未压缩）
def make_zip(source_dir, output_filename):
    zipf = zipfile.ZipFile(output_filename, 'w')
    pre_len = len(os.path.dirname(source_dir))
    for parent, dirnames, filenames in os.walk(source_dir):
        for filename in filenames:
            pathfile = os.path.join(parent, filename)
            arcname = pathfile[pre_len:].strip(os.path.sep)     #相对路径
            zipf.write(pathfile, arcname)
    zipf.close()

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
    e = '.exe' if goos == 'windows' else ''
    syst = 'macos' if goos == 'darwin' else goos
    arch = '64' if goarch == 'amd64' else goarch
    arch = '32' if goarch == '386' else arch
    if get_os() == 'win':
        cmd1 = 'SET CGO_ENABLED={}'.format(cgo)
        cmd2 = 'SET GOOS={}'.format(goos)
        cmd3 = 'SET GOARCH={}'.format(goarch)
        cmd4 = 'go build -o {} {}'.format("build/"+"-".join([Name, syst, arch]) + e, path)
        cmd = '  '.join([cmd1, cmd2, cmd3, cmd4])
        os.system(cmd)
    else:
        cmd = 'CGO_ENABLED={} GOOS={} GOARCH={} go build -o {} {} '.format(cgo, goos, goarch,
                                                                           "build/"+"-".join([Name, syst, arch]) + "/v3ray" + e,
                                                                           path)
        os.system(cmd)
    shutil.copy("README.md", "build/"+"-".join([Name, syst, arch]))
    make_zip("build/"+"-".join([Name, syst, arch]),"build/"+"-".join([Name, syst, arch])+".zip")

if __name__ == '__main__':

    for goos in INFO.keys():
        for goarch in INFO[goos]:
            print("正在编译", goos, '系统的', goarch, '架构版本中...')
            build(goos, goarch, 'v3ray.go')
