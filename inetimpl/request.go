package inetimpl

import (
	"locust/inet"
	"locust/message"
)

/**
 * @Author: stydm
 * @Date: 2019-10-30 16:05
 */

type Request struct {
	conn inet.IConnection //已经和客户端建立好的 链接
	msg  message.IMessage //客户端请求的数据

}
//获取请求连接信息
func(r *Request) GetConnection() inet.IConnection {
	return r.conn
}
//获取请求消息的数据
func(r *Request) GetData() []byte {
	return r.msg.GetData()
}

//获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
