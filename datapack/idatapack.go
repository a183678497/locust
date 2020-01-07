package datapack

import (
	"locust/message"
)

/**
 * @Author: stydm
 * @Date: 2019-10-30 19:11
 */

/*
   封包数据和拆包数据
   直接面向TCP连接中的数据流,为传输数据添加头部信息，用于处理TCP粘包问题。
*/
type IDataPack interface{
	Pack(msg message.IMessage)([]byte, error) //封包方法
	UnPack([]byte)(message.IMessage, int8)    //拆包方法
}