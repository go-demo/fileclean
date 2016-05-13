# Golang file clean

> 基于Golang的文件清理小程序

## 获取

``` bash
$ go get -v github.com/LyricTian/fileclean
```

## 使用

``` bash
$ fileclean --help
```

```
NAME:
   fileclean - 文件清理程序

USAGE:
   fileclean [global options] command [command options] [arguments...]

VERSION:
   0.1.1

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dir, -d "."                        文件目录
   --file, -f                           从指定的文件中读取要移除的文件列表
   --recur, -r                          递归处理
   --name [--name option --name option] 文件名
   --reg [--reg option --reg option]    正则过滤的文件名
   --exclude, -e                        排除当前指定的文件名
   --all, -a                            清理所有文件
   --help, -h                           show help
   --version, -v                        print the version
```