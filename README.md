# GolangStudyNotes
**目的就是记录一下学习go语言过程中遇到的一些问题和思考，能自己的学习历程留下点什么东西**

[TOC]

## 工具和环境
### [关于VS code 安装go 插件失败解决方案](https://github.com/zhangCan112/GolangStudyNotes/blob/master/NO1.md)

### [GitHub上优秀的Go开源项目](https://studygolang.com/articles/10217)

### [Go 资源大全中文版](https://github.com/jobbole/awesome-go-cn)



## 关于抄写学习go语言实践的一点体会建议

 1. 学习和熟悉如何使用golang标准库net包中自带的api创建一个简单的web服务
 2. 开始抄[urfave/negroni](https://github.com/urfave/negroni),这个库的主要功能就是扩展了http.Handler，为其添加了中间件能力。我们学习web服务很容易在掌握基本的web服务创建后就想到这个问题。这个库很顺滑的将你带入到抄写步骤中来
 3. 待续

## MAC上设置永久添加环境变量
Mac系统的环境变量，加载顺序为： 
a. /etc/profile 
b. /etc/paths 
c. ~/.bash_profile 
d. ~/.bash_login 
e. ~/.profile 
f. ~/.bashrc 
其中a和b是系统级别的，系统启动就会加载，其余是用户接别的。c,d,e按照从前往后的顺序读取，如果c文件存在，则后面的几个文件就会被忽略不读了，以此类推。~/.bashrc没有上述规则，它是bash shell打开的时候载入的。这里建议在c中添加环境变量，以下也是以在c中添加环境变量来演示的。
Go 一般需要设置的2个环境变量：
//#用来设置工作区
export GOPATH=/Users/xxxxusername/xxxpath/goxxx 
//#用来使用go module
export GO111MODULE=on

## GoModule代理的设置，可以解决下载不了golang/x/xxx的问题
go mod的代理比较出名的有微软的athens，可以基于它搭建一个私有的代理，管理内部的私有代码，而且微软提供了一个公共的代理，我们可以直接使用

Linux export GOPROXY="https://athens.azurefd.net"

Windows 设置GOPROXY环境变量

这样google下的包可以顺利下载了,速度还挺快的。