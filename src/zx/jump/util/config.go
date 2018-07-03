package util

import (
	"bytes"
)

/*通用配置类*/

const(
	// 未知异常
	Fail = iota
	// 成功
	Success
	// 密码错误
	PwdError
	// 连接到目标服务器失败
	ConnectPeerError
)



/*可转为数组接口*/
type Byteable interface {
	ToBytes() []byte
}
/**
	代理请求对象
 */
type JumpRequest struct{
	PwdLen byte
	// 密码
	Pwd []byte
	// 目标长度
	TargetLen byte
	// 目标 domain或ip:port
	Target []byte
}


func (this JumpRequest)ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(byte(this.PwdLen))
	buf.Write( this.Pwd)
	buf.WriteByte(this.TargetLen)
	buf.Write(this.Target)
	return buf.Bytes()
}

/**
	代理请求响应
 */
 type JumpResponse struct{
 	// 1成功，0失败
 	Status byte
 }
