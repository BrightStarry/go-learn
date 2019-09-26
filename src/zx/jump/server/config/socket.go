package config

import (
	"net"
	"zx/jump/util"
	"bytes"
	"log"
	"io"
)

/*socket相关*/

/**
	处理方法
 */
func Handler(conn *net.TCPConn) {
	defer conn.Close()
	// 校验请求
	errFlag, request := VerifyRequest(conn)
	if err := util.SendMessage(conn,util.JumpResponse{Status: errFlag},WriteTimeout);err != nil{
		log.Println("异常：",err)
		return
	}
	if errFlag != util.Success{
		return
	}
	// 连接到目标服务器
	targetConn, err := util.ConnectToTarget(string(request.Target),ReadTimeout)
	if err != nil{
		util.SendMessage(conn,util.JumpResponse{util.ConnectPeerError},WriteTimeout)
		return
	}
	defer targetConn.Close()

	go io.Copy(targetConn, conn)
	io.Copy(conn, targetConn)
}


/**
	校验请求
 */
func VerifyRequest(conn *net.TCPConn) (errFlag uint8,request util.JumpRequest) {
	errFlag = util.Fail
	buf, err := util.ReadMessage(conn, 512, ReadTimeout)
	if err != nil {
		return
	}
	reader := bytes.NewReader(buf)
	switch Config.Encrypt {
		case NotEncrypt:
		//... 什么都不做
	default:
		log.Fatalln("加密方式有误，当前加密方式:",Config.Encrypt)
	}

	if  request.PwdLen,request.Pwd,err = util.ReadByLen(reader);err != nil { // 读取密码
		return
	}
	if bytes.Compare(Config.PwdByte,request.Pwd) != 0{
		errFlag = util.PwdError
		return
	}

	if  request.TargetLen,request.Target,err = util.ReadByLen(reader);err != nil {// 读取目标服务器
		return
	}
	errFlag = util.Success
	return
}

