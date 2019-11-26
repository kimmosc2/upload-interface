TanPan 接口服务
===
负责响应用户的请求,并向数据服务发送定位请求  
本项目为《分布式对象存储——原理、架构及 GO 语言实现》一书的实践

## REST接口
GET /locate/<object_name> -- 响应:定位结果  
客户端通过GET方法发起对象定位请求，接口服务节点收到该请求后会向数据服务层发送一个定位消息，
然后等待数据服务节点的反馈，如果有数据服务节点发回确认消息，则返回该数据服务节点的地址，
如果超过一段时间没有任何反馈，就返回HTTP错误代码 404

## 主要流程
TODO

## Auther
[BuTn](https://github.com/kimmosc2)

## 环境变量
RabbitMQ地址:RABBITMQ_SERVER  
本地监听地址:LISTEN_ADDRESS  

## 其他组件
[数据服务](https://github.com/kimmosc2/upload-data)