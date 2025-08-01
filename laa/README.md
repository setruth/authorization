# 本地鉴权代理模块

> (Local Authentication Agent - LAA)

此服务部署在离线业务程序所在的机器上，负责生成设备唯一标识、存储授权文件，并执行所有本地的授权验证逻辑，提供授权状态APi供三方程序查询。

## 开发环境
* Golang: 1.24.5
* Gin 框架

## 支持功能
- [x] 为每台设备生成专属标识
- [x] 提供WebGUI进行可视化操作
- [x] 通过SSE即时推送最新授权变化
- [x] RSA签名与AES/MGC确保授权安全
- [x] 设备码识别，阻止在非授权设备上使用
- [x]  支持灵活控制授权有效期。
## 平台适配
> 目前没做其他平台的设备码读取的适配，其余功能不影响
- [x] Windows
- [ ] MacOS
- [ ] Linux 

## 使用模块
使用`GoLand`IDE或手动Clone后请操作`laa`文件夹
### 1.配置密钥
模块内置默认密钥，如果不需要自定义可以跳过此步骤，自定义的密钥请遵守下列的规范。相关密钥可自行生成或利用`ams`模块的WebGUI进行生成
#### 配置AES密钥

此密钥用于对授权密文解密。需自定义(无需可跳过)密钥时请在`model/constant.go`中设置`aesKeyBASE64`变量。

**请遵循AES密钥规范，确保设置的是AES密钥的`BASE64`编码内容且和`ams`模块中的AES密钥相同**

> 默认aesKeyBASE64: `lBbYXsW3hbBc6IyUOXAPelaWB7t+lsqLbzyaO1oM+uU=`

#### 配置RSA公钥
此密钥用于对授权码进行签名校验确保密文一致性。需自定义(无需可跳过)密钥时请在`model/constant.go`中设置`rsaPublicKeyBASE64`变量。

**请遵循RSA密钥规范，确保设置的是RSA公钥的`BASE64`编码内容且和`ams`模块中的RSA私钥为一对**

> 默认rsaPublicKeyBASE64: 请查看`model/constant.go`文件中`rsaPublicKeyBASE64`
### 2.构建执行程序
使用`GoLand`IDE或在项目根目录执行指令构建程序
```shell
//构建指令
go build -o 自定义生成的exe文件名.exe  ./core 
```
### 3.使用WebGUI
默认的地址端口是`localhost:1022`，直接访问`localhost:1022/`可以使用内置的WebGUI来进行相关的配置

### 4.接入授权状态变化
#### Server-Send Events
使用`SSE`请求ip:port/status/subscribe，订阅后在默认事件中会得到最新的授权状态，往后授权状态变化时会进行主动推送

> 推送的消息内容
>```json
>{"tag":0,"endTimestamp":-1}
>```
- `endTimestamp`：授权结束的毫秒级时间戳。
- `tag`：当前授权状态码，具体对照如下：

| Tag | 状态     |
|:----| :------- |
| 0   | 未授权   |
| 1   | 已授权   |
| 2   | 授权过期 |

#### HTTP(GET)
使用`GET`请求ip:port/status，可拿到当前的授权状态，获取成功时`HttpStatus`为`200`，失败时会返回不同的`HttpStatus`，请记得处理。

>响应Body
> ```json
> {
>    "msg": "获取成功",
>    "data": {
>        "tag": 2,
>        "endTimestamp": 1754006400000
>    }
> }
> ```

- `data`:授权信息的结构与上面SSE推送内容的相同，字段说明参考上述
- `msg`:服务器的处理结果消息

## LAA模块结构图
![结构图](https://github.com/setruth/authorization/blob/master/laa/authorization-laa.png)

