# 授权管理服务模块

> (Authentication Management Service - AMS)

此服务一般用于部署到云服务中，集中进行管理所有的授权机器记录，提供生成授权码和更新到期时间功能，提供生成新的系统锁需的相关密钥。

## 开发环境

* Golang: 1.24.5
* Gin 框架
* GORM 框架
* SQLite 数据库

## 支持功能

- [x] 对授权记录进行增删查询
- [x] 更新授权记录的到期时间
- [x] 生成对应授权记录的授权码
- [x] 支持SQLite的持久性存储
- [x] 提供WebGUI进行可视化操作
- [x] 提供生成系统相关密钥能力

## 平台适配

> 可自行编译所需的平台，`ams`未做特定的平台功能，所以无需担心平台适配问题

- [x] Windows
- [x] MacOS
- [x] Linux

## 使用模块

使用`GoLand`IDE或手动Clone后请操作`ams`文件夹

### 1.配置密钥
模块内置默认密钥，如果不需要自定义只需设置默认的私钥即可，自定义的密钥请遵守下列的规范。相关密钥可自行生成或利用`ams`模块的WebGUI进行生成
#### 配置AES密钥

此密钥用于对授权信息加密。需自定义(无需可跳过)密钥时请在`model/constant.go`中设置`aesKeyBASE64`变量。

**请遵循AES密钥规范，确保设置的是AES密钥的`BASE64`编码内容且和`laa`模块中的AES密钥相同**

> 默认aesKeyBASE64: `lBbYXsW3hbBc6IyUOXAPelaWB7t+lsqLbzyaO1oM+uU=`

#### 配置RSA私钥

此密钥用于对密文进行签名确保密文一致性。请在部署此模块的系统的环境变量，以`model/constant.RsaPrivateEnvKey`
的值为key，Value为Rsa私钥的`BASE64`编码内容。

**请遵循RSA密钥规范，确保设置的是RSA私钥的`BASE64`编码内容且和`laa`模块中的RSA公钥为一对**

> 默认RSA私钥在`ams/RsaPrivateKeyBase64.txt`文件中，它与`laa`模块中的RSA公钥是一对，可直接填入环境变量的值中

### 2.构建执行程序

使用`GoLand`IDE或在`ams`文件夹根目录执行构建指令编译为不同平台执行程序

### 3.使用WebGUI

默认的地址端口是`localhost:1023`，直接访问`localhost:1023/web`可以使用内置的WebGUI来进行相关的配置

## AMS模块结构图

![结构图](https://github.com/setruth/authorization/blob/master/ams/authorization-ams.png)