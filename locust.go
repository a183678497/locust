package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"locust/inet"
	"locust/inetimpl"
	"locust/log"
	"locust/mysql"
	"locust/utils"
)

/**
 * @Author: stydm
 * @Date: 2019-10-30 16:38
 */

//环境数据
type Environment struct {
	sn string
	number uint8
	serialNumber uint16
	temperature float32
	humidity uint8
	illumination float64
	collectorState uint8
	sensorState uint8
}

//解析环境数据路由
type EnvironmentRouter struct {
	inetimpl.BaseRouter
}

func getEnvironment(bytes []byte) Environment {
	var env Environment
	env.sn = utils.BcdToStr(bytes[1:9])
	env.serialNumber = utils.BytesToUint16(bytes[9:11])
	env.number = bytes[13]
	if (utils.BytesToUint16(bytes[14:16])&0x8000)==0 {
		env.temperature = float32(utils.BytesToUint16(bytes[14:16]) & 0x7fff)*0.1
	}else{
		env.temperature = float32(utils.BytesToUint16(bytes[14:16]) & 0x7fff)*-0.1
	}
	env.humidity = bytes[16]
	env.illumination = float64(utils.BytesToUint32(bytes[17:21])) * 0.001
	env.collectorState = bytes[21]
	env.sensorState = bytes[22]
	return env
}

//环境数据处理
func (this *EnvironmentRouter) Handle(request inet.IRequest) {
	fmt.Println("Call EnvironmentRouter Handle")
	sql := "INSERT INTO `caiyuan_data_environment`(`device_no`, `number`, `serial_number`, `temperature`, `humidity`, `illumination`, `collector_state`, `sensor_state`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	env := getEnvironment(request.GetData())
	_,err := mysql.Engine.Exec(sql,env.sn,env.number,env.serialNumber,env.temperature,env.humidity,env.illumination,env.collectorState,env.sensorState)
	if err != nil {
		fmt.Println(err)
	}
}

//氮磷钾数据
type NPK struct {
	sn string
	number uint8
	serialNumber uint16
	nitrogen uint16
	phosphorus uint16
	potassium uint16
	humidity uint8
	waterState uint8
	sensorState uint8
}

func getNPK(bytes []byte) NPK {
	var npk NPK
	npk.sn = utils.BcdToStr(bytes[1:9])
	npk.serialNumber = utils.BytesToUint16(bytes[9:11])
	npk.number = bytes[13]
	npk.nitrogen = utils.BytesToUint16(bytes[14:16])
	npk.phosphorus = utils.BytesToUint16(bytes[16:18])
	npk.potassium = utils.BytesToUint16(bytes[18:20])
	npk.humidity = bytes[20]
	npk.waterState = bytes[21]
	npk.sensorState = bytes[22]
	return npk
}

//解析氮磷钾数据路由
type NPKRouter struct {
	inetimpl.BaseRouter
}

//氮磷钾数据处理
func (this *NPKRouter) Handle(request inet.IRequest) {
	fmt.Println("Call NPKRouter Handle")
	sql := "INSERT INTO `caiyuan_data_npk`(`device_no`, `number`, `serial_number`, `nitrogen`, `phosphorus`, `potassium`, `humidity` ,`water_state` ,`sensor_state`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	npk := getNPK(request.GetData())
	_,err := mysql.Engine.Exec(sql,npk.sn,npk.number,npk.serialNumber,npk.nitrogen,npk.phosphorus,npk.potassium,npk.humidity,npk.waterState,npk.sensorState)
	if err != nil {
		fmt.Println(err)
	}
}

//创建连接的时候执行
func DoConnectionBegin(conn inet.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")

	//=============设置两个链接属性，在连接创建之后===========
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://www.jianshu.com/u/35261429b7f1")
	//===================================================

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

//连接断开的时候执行
func DoConnectionLost(conn inet.IConnection) {
	//============在连接销毁之前，查询conn的Name，Home属性=====
	if name, err:= conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}
	//===================================================

	fmt.Println("DoConneciotnLost is Called ... ")
}

func main() {
	//创建一个server句柄
	s := inetimpl.NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//配置路由
	s.AddRouter(0, &EnvironmentRouter{})
	s.AddRouter(1, &NPKRouter{})

	log.LogFile.WithFields(logrus.Fields{
		"operation": "程序启动中....",
		"data":  "",
	}).Info("A group of walrus emerges from the ocean")

	//开启服务
	s.Serve()


	// AA0000112018100009950001080100007b1600C9E85F30AACC
	// AA010011201810000995000208010304CF03D40012161EB4CC
}