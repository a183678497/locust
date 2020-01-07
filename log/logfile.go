package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

/**
 * @Author: stydm
 * @Date: 2019-11-04 19:47
 */

/*
   定义一个全局的对象
*/


var LogFile *logrus.Logger



/*
   提供init方法，默认加载
*/
func init() {
	LogFile = logrus.New()

	file, err := os.OpenFile("locust.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		LogFile.Out = file
	} else {
		LogFile.Info("Failed to log to file, using default stderr")
	}
}