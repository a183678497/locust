package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"locust/common"
	"log"
)

/**
 * @Author: stydm
 * @Date: 2019-11-04 17:15
 */


/*
   定义一个全局的对象
*/

var Engine *xorm.Engine

/*
   初始化数据库连接
*/
func newEngine() *xorm.Engine {
	engine ,err := xorm.NewEngine(common.ConfigObject.MysqlDriverName, common.ConfigObject.MysqlDataSourceName)
	if err != nil {
		log.Fatal(newEngine, err)
		return nil
	}

	engine.SetMaxIdleConns(common.ConfigObject.MysqlMaxIdleConns)
	engine.SetMaxOpenConns(common.ConfigObject.MysqlMaxOpenConns)

	return engine
}

/*
   提供init方法，默认加载
*/
func init() {
	Engine = newEngine()
}