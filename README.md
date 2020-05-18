# NetCore

## Menu 目录

[TOC]
## How to use?  如何使用?

### STEP1 第一步: 创建服务器
```go
  err := common.InitConfig("Config.json")
  serv := socket.Server{}
  serv.Run()
```

#### 配置文件
```toml

```

#### 生成配置文件
```go
  common.GenerateConf("Config.json")
```

### STEP2 第二部: 添加回调函数
```go
  g.ServerHandler.AddHandler(socket.MsgHert,heartHandler)
```
#### 回调函数
```go
  type RecvHandler func(message []byte) *MessageData
```

## How to registry? 如何注册服务
```go
	


```




## Context 内容




### Socket 网络

#### RSA加密

```go
func main() {
  err := RSA.GenRsaKey(1024)
  	if err != nil {
  		fmt.Println(err.Error())
  	}


	// 公钥加密私钥解密
	if err := applyPubEPriD(); err != nil {
		log.Println(err)
	}
	// 公钥解密私钥加密
	if err := applyPriEPubD(); err != nil {
		log.Println(err)
	}
}

// 公钥加密私钥解密
func applyPubEPriD() error {
	pubenctypt, err := gorsa.PublicEncrypt(`hello world`,Pubkey)
	if err != nil {
		return err
	}

	pridecrypt, err := gorsa.PriKeyDecrypt(pubenctypt,Pirvatekey)
	if err != nil {
		return err
	}
	if string(pridecrypt) != `hello world` {
		return errors.New(`解密失败`)
	}
	return nil
}

// 公钥解密私钥加密
func applyPriEPubD() error {
	prienctypt, err := gorsa.PriKeyEncrypt(`hello world`,Pirvatekey)
	if err != nil {
		return err
	}

	pubdecrypt, err := gorsa.PublicDecrypt(prienctypt,Pubkey)
	if err != nil {
		return err
	}
	if string(pubdecrypt) != `hello world` {
		return errors.New(`解密失败`)
	}
	return nil
}
```


### Config 服务配置

### Plugins 插件

### Errors 异常处理

## Libs 引用库

-   [protobuf](https://github.com/golang/protobuf): 通讯协议库
    安装方式:
    ```$ go get github.com/golang/protobuf```
    首次安装ProtoBuff需要安装生成工具[protoc](https://developers.google.com/protocol-buffers/).
-   [UUID](http://github.com/satori/go.uuid): UUID 生成
    安装方式:
    ```$ go get github.com/satori/go.uuid```
-   [go-logging](http://github.com/op/go-logging ): 日志
    安装方式:
    ```$ go get github.com/op/go-logging```

-   [x-net](http://golang.org/x/net ): Golang 网络扩展库Websocket
    安装方式:
    ```$ go get golang.org/x/net```

-  [g-toml](http://github.com/BurntSushi/toml): Toml 配置文件
    安装方式:
    ```go get github.com/BurntSushi/toml```
