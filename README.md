# nging
基于 caddy 的网站服务程序，带图形化管理界面。

caddy 是由国外开发者开发的一套类似于nginx或apache的网站服务软件。
caddy的配置文件比nginx更简洁易用。但我相信事情还可以变得更简单，所以nging应运而生。

nging不仅仅包含了caddy的在线可视化配置，还包含了ftp服务的管理，下一步即将增加：

- [x] 文件在线管理
- [x] 计划任务
- [x] 数据库管理（adminer的Golang版）
- [x] 离线下载
- [x] 网站管理（支持caddy指令在线配置）
    - [x] fastcgi
    - [x] gzip
    - [x] header
    - [x] ipfilter
    - [x] log
    - [x] rewrite
    - [x] tls


# 下载地址

## 最新版本

### v1.3.3 下载

* MacOS64位版本：[nging_v1.3.3_darwin_amd64.tar.bz2](http://www.admpub.com:9000/api/file/getAttach?fileId=5c40259e04aa045d8a000041)

* Windows32位版本：[nging_v1.3.3_windows_386.zip](http://www.admpub.com:9000/api/file/getAttach?fileId=5c4025a704aa045d8a000044)

* Windows64位版本：[nging_v1.3.3_windows_amd64.zip](http://www.admpub.com:9000/api/file/getAttach?fileId=5c4025a104aa045d8a000042)

* Linux32位版本：[nging_v1.3.3_linux_386.tar.bz2](http://www.admpub.com:9000/api/file/getAttach?fileId=5c4025a504aa045d8a000043)

* Linux64位版本：[nging_v1.3.3_linux_amd64.tar.bz2](http://www.admpub.com:9000/api/file/getAttach?fileId=5c40259d04aa045d8a000040)


## 历史版本

### v1.3.2 下载

* MacOS64位版本：[nging_darwin_amd64.tar.bz2](http://www.admpub.com:9000/api/file/getAttach?fileId=5c40191404aa045d8a00003c)

* Windows32位版本：[nging_windows_386.zip](http://www.admpub.com:9000/api/file/getAttach?fileId=5c40191f04aa045d8a00003f)

* Windows64位版本：[nging_windows_amd64.zip](http://www.admpub.com:9000/api/file/getAttach?fileId=5c40191b04aa045d8a00003d)

* Linux32位版本：[nging_linux_386.tar.bz2](http://www.admpub.com:9000/api/file/getAttach?fileId=5c40191304aa045d8a00003b)

* Linux64位版本：[nging_linux_amd64.tar.bz2](http://www.admpub.com:9000/api/file/getAttach?fileId=5c40191e04aa045d8a00003e)

### v1.3.1 CSDN 下载

* MacOS64位版本：https://download.csdn.net/download/admpub/10868709

* Windows64位版本：https://download.csdn.net/download/admpub/10868722

* Linux32位版本：https://download.csdn.net/download/admpub/10868656

* Linux64位版本：https://download.csdn.net/download/admpub/10867479

# 运行方式
首先下载相应平台的程序，然后解压缩到当前目录，进入文件夹找到“nging-”开头的程序，在此程序所在的目录下执行此程序(在非windows系统里在执行之前请赋予该程序可执行权限)。
例如在Linux64位系统，分别执行以下命令：
```
chmod +x ./nging-linux-amd64
./nging-linux-amd64
```
打开浏览器，访问网址`http://localhost:9999/setup`，在页面中配置数据库和管理员账号信息进行安装。安装成功后会自动跳转到登录页面，使用安装时设置的管理员账号进行登录

# 先睹为快

### 安装：
[![](https://github.com/admpub/nging/blob/master/preview/preview_install.png?raw=true)](https://github.com/admpub/nging/blob/master/preview/preview_install.png)

### 登录：
[![](https://github.com/admpub/nging/blob/master/preview/preview_login.png?raw=true)](https://github.com/admpub/nging/blob/master/preview/preview_login.png)

### 系统信息：
[![](https://github.com/admpub/nging/blob/master/preview/preview_sysinfo.png?raw=true)](https://github.com/admpub/nging/blob/master/preview/preview_sysinfo.png)

### 在线编辑文件：
[![](https://github.com/admpub/nging/blob/master/preview/preview_editfile.png?raw=true)](https://github.com/admpub/nging/blob/master/preview/preview_editfile.png)

### 添加计划任务：
[![](https://github.com/admpub/nging/blob/master/preview/preview_task.png?raw=true)](https://github.com/admpub/nging/blob/master/preview/preview_task.png)

### MySQL数据库管理：
[![](https://github.com/admpub/nging/blob/master/preview/preview_listtable.png?raw=true)](https://github.com/admpub/nging/blob/master/preview/preview_listtable.png)

# 开发环境下的启动方式

- 第一步： 安装GO环境(建议1.7版以上)，配置GOPATH、GOROOT环境变量，并将`%GOROOT%/bin`和`%GOPATH%/bin`加入到PATH环境变量中
- 第二步： 执行命令`go get github.com/admpub/nging`
- 第三步： 进入`%GOPATH%/src/github.com/admpub/nging/`目录中启动`run_first_time.bat`(linux系统启动`run_first_time.sh`)
- 第四步： 打开浏览器，访问网址`http://localhost:8080/setup`，在页面中配置数据库账号和管理员账号信息进行安装
- 第五步： 安装成功后会自动跳转到登录页面，使用安装时设置的管理员账号进行登录
