package common

import (
	"encoding/json"
	"io/ioutil"
	"locust/inet"
)

/**
 * @Author: stydm
 * @Date: 2019-10-30 17:08
 */

/*
   存储一切有关的全局参数，供其他模块使用
   一些参数也可以通过 用户根据 locust.json来配置
*/
type Config struct {
	/*
		server
	*/
	TcpServer inet.IServer   //当前的全局Server对象
	Host      string         //当前服务器主机IP
	TcpPort   int            //当前服务器主机监听端口号
	Name      string         //当前服务器名称
	Version   string         //当前版本号

	/*
		worker
	*/
	MaxPacketSize uint32 //数据包的最大值
	MaxConn       int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize uint32   //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen uint32  //SendBuffMsg发送消息的缓冲最大长度

	/*
		mysql
	*/
	MysqlMaxOpenConns   int    //设置最大打开的连接数
	MysqlMaxIdleConns   int    //设置闲置的连接数
	MysqlDriverName     string //驱动类型 mysql
	MysqlDataSourceName string //连接串 dev:Dev!123456@tcp(114.67.97.22:3306)/stydm?charset=utf8mb4

	/*
		log
	 */
	LogName string
	MaxRemainCnt uint
}

/*
   定义一个全局的对象
*/
var ConfigObject *Config

//读取用户的配置文件
func (c *Config) Reload() {
	data, err := ioutil.ReadFile("locust.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &ConfigObject)
	if err != nil {
		panic(err)
	}
}

/*
   提供init方法，默认加载
*/
func init() {
	//初始化GlobalObject变量，设置一些默认值
	ConfigObject = &Config{
		Name:    "LocustApp",
		Version: "tcp",
		TcpPort: 1437,
		Host:    "0.0.0.0",
		MaxConn: 12000,
		MaxPacketSize:4096,
		WorkerPoolSize:1024,
		MaxWorkerTaskLen:4096,
		MysqlMaxOpenConns:   100,
		MysqlMaxIdleConns:   10,
		MysqlDriverName:     "",
		MysqlDataSourceName: "",
		LogName:"locust",
		MaxRemainCnt : 240,
	}

	//从配置文件中加载一些用户配置的参数
	ConfigObject.Reload()
}