# NetCore

>  自己用的脚手架(已经作废啦, 重新规划ing)

## Menu 目录

[TOC]
## How to use?  如何使用?

### STEP1 第一步: 配置文件

```toml

[Sys]
  num_cpu = 2 

[WebSocket]
  ip = "0.0.0.0" 
  port = 5678
  name = "ws"
  [WebSocket.log]
    WriteToFile = true
    LogFilePath = "./logs/websocket/"
    ZipTime = 2
    ChannelSize = 25

[Http]
  ip = "0.0.0.0"
  port = 5679
  name = "http"
  [Http.log]
    WriteToFile = true
    LogFilePath = "./logs/http/"
    ZipTime = 2
    ChannelSize = 25
```

##### Sys

 	1. num_cpu : 设定使用多少个核心来运行本程序,设定为0表示自动分配,此时程序将根据运行设备的总CPU数采取尽可能占用的策略分配CPU.

##### Websocket

1. ip : 服务器运行的IP,绝大多数请填写0.0.0.0
2. port : websocket 服务监听的端口,0-65535 请不要重复已占用的端口
3. name : 给 websocket 起( mei )个( qiu )名( luan )字( yong )呗 

##### Http

1. ip : 服务器运行的IP,绝大多数请填写0.0.0.0
2. port : http 服务监听的端口,0-65535 请不要重复已占用的端口
3. name : 给 http 起( mei )个( qiu )名( luan )字( yong )呗 

#####  <span id="logconfig">Log</span>

1. WriteToFile : 是否写入文件,如果这里是 ```false``` 的话,将只会把日志输出到**控制台**
2. LogFilePath: 如果 ```WriteToFile``` 为 ```false``` 此项无效, 为```true``` 时为日志输出**目录**
3. ZipTime: 如果 ```WriteToFile``` 为 ```false``` 此项无效, 为```true```时为执行压缩日志文件的时间**小时数**
4. ChannelSize: 如果 ```WriteToFile``` 为 ```false``` 此项无效, 为```true```时为接受日志的通道大小.

### STEP2 第二步: 配置通讯协议
> 此步骤可以在创建完服务器之后再做,此处只是为了书写方便

本项目采用Google的[ProtoBuf](https://github.com/golang/protobuf)协议

> protobuf 1.4.1
>
> protoc 1.11.4

*范例*

```protobuf
syntax = "proto3";

package msg;

message PingRequest {
    int64 time = 1;
}

message PingResponse {
    int64 time = 1;
}
```

> 建议性设置所有请求都增加```Request``` 作为标志, 所有返回都增加 ```Response``` 作为标志.



### STEP3 第三步: 创建服务器

> 在```main.go```的```main```程序入口当中写下面的代码可以创建一个服务器

```go
  if e, err := entry.Create(); err != nil {
		fmt.Printf("启动异常: %s", err.Error())
	} else {
		e.Start()
		e.ExitSignalMonitor()
	}
```

### STEP4 第四部: 编写任务处理函数

#### 函数结构体

```go
// websocket 任务处理函数 
type RecvHandler func(message *MessageData) *MessageData
// http(gin) 任务处理函数
type HandlerFunc func(c *gin.Context)
```
#### *```MessageData```* 参数结构

```go
type MessageData struct {
	MessageType uint16
	Message     []byte
}
```

1. ```MessageType``` 的含义表示消息的指令ID
2. ```Message``` 消息体

#### 如何将处理函数添加至服务

```go
if e, err := entry.Create(); err != nil {
		fmt.Printf("启动异常: %s", err.Error())
	} else {
    // 添加HTTP处理函数
		hs := e.GetHttp()
		hs.AddHandler("GET", "/ping", func(c *gin.Context) {
			message, exists := c.Get("data")
			if exists {
				fmt.Printf("handler --> %s \n", message)
			}
			data := &protoexample.Test{
				Label: &label,
				Reps:  reps,
			}
			c.ProtoBuf(http.StatusOK, data)
		})
    // 添加WebSocket处理函数
		wss := e.GetWebSocket()
		wss.AddHandler(uint16(1), func(message *network.MessageData) *network.MessageData {
			return nil
		})

		e.Start()
		e.ExitSignalMonitor()
	}
```



## Context 内容

### 日志

#### 如何开启日志

``` toml
  [WebSocket.log]
    WriteToFile = true
    LogFilePath = "./logs/websocket/"
    ZipTime = 2
    ChannelSize = 25
    
  [http.log]
    WriteToFile = true
    LogFilePath = "./logs/http/"
    ZipTime = 2
    ChannelSize = 25

```

将服务的日志配置项当中的 ```WriteToFile``` 改为 ```true``` ,并且配置剩余选项,即可开启日志服务,配置详情参考**[日志配置](#logconfig)**

#### 如何使用日志

> 日志分为主要的三个级别 一个是 ```info```  和 ```debug```用来存储日常信息, 另个一个是 ```error``` 当系统报错的时候使用

1. Debug

   ``` go
   log.Debugf(format string, v ...interface{})
   ```

2. Info

   ```go
   log.Infof(format string, v ...interface{})
   ```

3. Error

   ```go
   log.Errorf(format string, v ...interface{})
   ```

```format``` 支持所有 ```Golang``` 占位符语法

### 系统监控



### CLI

## Libs 引用库

-   [protobuf](https://github.com/golang/protobuf): 通讯协议库
    安装方式:
    ```$ go get -u github.com/golang/protobuf```
    首次安装ProtoBuff需要安装生成工具[protoc](https://developers.google.com/protocol-buffers/).
    
-   [UUID](http://github.com/satori/go.uuid): UUID 生成
    安装方式:
    ```$ go get -u github.com/satori/go.uuid```
    
-   [x-net](http://golang.org/x/net ): Golang 网络扩展库Websocket
    安装方式:
    ```$ go get -u golang.org/x/net```
    
-  [g-toml](http://github.com/BurntSushi/toml): Toml 配置文件
    安装方式:
    ```go get -u github.com/BurntSushi/toml```
    
-   [gin](http://github.com/BurntSushi/toml): Gin Web服务
    安装方式:
    ```go get -u github.com/gin-gonic/gin```
    
-   [ksuid](http://github.com/segmentio/ksuidl): 不重复的ID生成
    安装方式:
    ```go get -u github.com/segmentio/ksuid```
    
-   [gopsutil](https://github.com/shirou/gopsutil ) 系统信息读取

    安装方式:

    ```go get -u github.com/shirou/gopsutil ```

