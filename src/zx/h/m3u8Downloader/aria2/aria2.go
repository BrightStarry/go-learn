package aria2

import (
	"time"
	"github.com/zyxar/argo/rpc"
	"context"
	"strconv"
	"os"
	"io"
	"strings"
	"zx/h/m3u8Downloader/util"
)

/**
aria2相关
 */


/**
运行单个下载
 */
func RunAria2(loaPath, m3u8Path, outDir,aria2Path string,aria2Args []string) (error) {
	args := []string{"/c",aria2Path,"-i", m3u8Path,"-d", outDir,}
	args = append(args,aria2Args...)
	return util.StartCMDToLog(loaPath, "cmd", args...)
}

/**
根据日志判断任务执行是否成功
 */
func IsSuccessByLog(loaPath string) (success bool, err error) {
	success = false
	// 文件不存在 失败
	if !util.FileIsExist(loaPath) {
		return
	}
	logFile, err := os.OpenFile(loaPath, os.O_RDONLY, 0777)
	if err != nil {
		return
	}
	defer logFile.Close()
	// 读取末尾的20个字节，判断是否包含成功字符串
	_, err = logFile.Seek(-20, io.SeekEnd)
	if err != nil {
		return
	}
	bytes := make([]byte, 20)
	logFile.Read(bytes)
	s := string(bytes)
	return strings.Contains(s, "download completed"), nil
}

/**
 启动aria2 rpc服务
传入各个配置根目录
 */
func StartRPC(rootDir string, port int) (err error) {
	return util.StartCMD("cmd", "/c",
		rootDir+"aria2c.exe",
		"--conf-path="+rootDir+"aria2.conf",
		" --log="+rootDir+"task1.log",
		"--input-file="+rootDir+"aria2.session",
		"--save-session="+rootDir+"aria2.session",
		"--rpc-listen-port="+strconv.Itoa(port))

}

/**
连接到rpc，根据端口
! 需要关闭
 */
func ConnRPC(port int) (aria2Conn rpc.Client, err error) {
	aria2Conn, err = rpc.New(context.Background(), "ws://localhost:"+strconv.Itoa(port)+"/jsonrpc", "", time.Second, &rpc.DummyNotifier{})
	return
}
