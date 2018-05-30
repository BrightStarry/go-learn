package util

import "log"

/*工具类*/

/*判断是否有异常，有则打印日志*/
func LogError(err error,v ...interface{}) {
	if err != nil {
		log.Fatalln(v)
	}
}