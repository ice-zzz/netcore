# NetCore_TODOs



网络通讯

package

 ```0xFF|0xFF|token(32bit)|包总数|当前数|len(高)|len(低)|MsgType|Message(zlib)|MD5(32bit)|0xFF|0xFE```

1. 0xFF|0xFF 起始标识符
2. token(32bit) 为随机字符,客户端需要返回 tempSid(27bit)+uid 并且md5(32bit)
3. 其中len为Message的长度，实际长度为len(高)*256+len(低)
4. MD5，为Message的MD5码,防止被串改数据


socket&websocket

1. 客户端连接
2. 服务器发送欢迎消息,并且要求客户端进行身份验证
3. 客户端发出身份信息, 用户名
4. 服务器验证身份信息,并且生成临时唯一ID
5. 客户端接收ID并将 ID+用户名称 进行md5加密.作为token
6. 每次通讯服务器都验证 token是否正确
7. 如果连接断开则销毁临时ID, 如果没有断开的时候再次请求身份验证的话则向全体同用户名的连接发送入侵消息



http&https

短连接 JWT

