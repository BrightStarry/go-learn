package test

import (
	"testing"
	"os"
	"github.com/grafov/m3u8"
	"bufio"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"io/ioutil"
	"zx/h/m3u8Merge/util"
	"os/exec"
	"time"
	"log"
	"github.com/alvatar/multipart-downloader"
	"path"
	"sync"
	"path/filepath"
	"zx/h/m3u8Merge/myLog"
	"errors"
	"strings"
)
type Data struct {
	filename     string
	filesize     uint
	dirname      string
	fullfilename string
}

/**
下载 ts
 */
func TestDownloadTs(t *testing.T) {
	uris,keyURI,err := util.GetUrisAndKeyURIByM3u8File("e:/0snis00584.m3u8")
	if err != nil {
		panic("解析m3u8失败:" + err.Error())
	}
	if len(uris) <= 0 {
		panic("解析m3u8失败，ts链接为空")
	}
	if keyURI == "" {
		panic("解析m3u8失败，key链接为空")
	}
	nConns := 32
	timeout := time.Duration(60000) * time.Second
	waitGroup := sync.WaitGroup{}
	multipartdownloader.SetVerbose(true)
	for _,uri := range uris{
		waitGroup.Add(1)
		uriCopy := uri
		go func() {
			defer waitGroup.Done()

			dldr := multipartdownloader.NewMultiDownloader([]string{uriCopy}, nConns, timeout)

			// Gather info from all sources
			for {
				_, err =dldr.GatherInfo()
				if err == nil {
					break
				}
				fmt.Println(err)
			}

			// Prepare the file to write downloaded blocks on it
			_, err = dldr.SetupFile("E:\\新建文件夹\\" + path.Base(uriCopy))
			if err != nil {
				panic(err)
			}


			// Perform download
			err = dldr.Download(func(feedback []multipartdownloader.ConnectionProgress) {
				log.Println(feedback)
			})

			if err != nil {
				panic(err)
			}
		}()
	}
	waitGroup.Wait()
	fmt.Println("success")






}

/**
解析m3u8 文件
 */
func TestParseM3u8(t *testing.T) {
	//f, err := os.Open("E:/0snis00584.m3u8")
	f, err := os.Open("C:\\h\\av2\\m3u8\\all\\ipx00180.m3u8")
	if err != nil {
		panic(err)
	}
	p, _, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	data := p.(*m3u8.MediaPlaylist)


	/**
	获取ts 的uri列表
	 */
	uris := make([]string,0)
	for _,x := range data.Segments {
		if x == nil {
			continue
		}
		uris = append(uris, x.URI)
	}
	for _,x := range uris {
		fmt.Println(x)
	}
	fmt.Println("success")

	/**
	获取key
	 */
	 keyURI := data.Key.URI
	 // 加密方法
	 keyMethod := data.Key.Method
	 fmt.Println(keyURI)
	 fmt.Println(keyMethod)

}

/**
测试下载
 */
func TestDownload(t *testing.T) {
	request := gorequest.New()
	resp, body, errs := request.Get("https://str.dmm.com:443/digital/st1:M60Dzvs0lJVTYRZXdBQpE9ECwfPaXiZF3joWn5vd-n4XPEf9RttYx3WxCmMpWMAw/3Gc6MVxQvftqVNSW6fuZ3Ul/-/media_b3000000_3.ts").EndBytes()
	if len(errs) > 0 {
		for _,x := range errs {
			fmt.Println(x)
		}
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("err:" + string(resp.StatusCode))
	}

	ioutil.WriteFile("E:/test.zx",body,0666)

}

/**
批量处理
解密 + 合并
 */
func TestBatchProcess(t *testing.T) {
	//rootDir := "E:\\新建文件夹\\"
	keyRootDir := "c:\\h\\av2\\m3u8\\all\\"
	dir := util.GetParam("ts目录:")
	dir = "E:\\新建文件夹\\118abp00108"



	// 读取key
	keyFilePath := keyRootDir + filepath.Base(dir) + ".key"
	keyBytes,err := ioutil.ReadFile(keyFilePath)
	if err!= nil {
		panic("读取key异常：" + err.Error())
	}

	// 获取ts文件名
	filePaths := util.GetAllFileName(dir)
	// 解密后输出路径
	outPath := dir+string(os.PathSeparator)+"out"+string(os.PathSeparator)
	err = util.CreateFileByFileIsExist(outPath)
	if err!= nil {
		panic("创建文件异常：" + err.Error())
	}

	// 设置线程池
	threadPool := util.ThreadPool{}
	threadPool.Init(32, func(args []interface{}) error {
		filePath := args[0].(string)
		// 读取加密视频
		videoBytes,err := ioutil.ReadFile(filePath)
		if err!= nil {
			return errors.New("读取加密视频异常：" + err.Error())
		}
		decryptBytes := util.AESDecrypt(videoBytes,keyBytes)
		err = ioutil.WriteFile(outPath + filepath.Base(filePath) + ".out.ts",decryptBytes,0666)
		if err!= nil {
			return errors.New("写入解密后视频异常：" + err.Error())
		}
		return nil
	})
	// 启动
	threadPool.Start()
	// 获取结果,必须放在入队前，否则会死锁
	errResults := make([]util.Result,0)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		// 获取结果
		for i:=0;i< len(filePaths);i++{
			result := threadPool.Take()
			if !result.Success {
				errResults = append(errResults,result)
			}
		}
	}()

	//任务入队
	for i,item := range filePaths{
		threadPool.Put(i,[]interface{}{item})
	}
	// 关闭
	threadPool.Close()
	// 等待获取结果
	waitGroup.Wait()
	// 打结果
	if len(errResults)> 0 {
		myLog.Error("失败结果如下:%v" ,errResults)
	}else{
		myLog.Info("视频解密成功" )
	}

	/**
	合并
	  */
	var allBytes []byte
	for _,item := range filePaths{
		tempBytes,_ :=ioutil.ReadFile(outPath + filepath.Base(item)+ ".out.ts")
		allBytes = append(allBytes,tempBytes...)
	}
	ioutil.WriteFile(outPath+"all.ts",allBytes,0666)



}


/**
测试解密
 */
func TestAes(t *testing.T) {
	/**
	加密视频
	 */
	bytes,err := ioutil.ReadFile("E:/test.zx")
	if err!= nil {
		panic(err)
	}
	/**
	key
	 */
	keyBytes,err := ioutil.ReadFile("E:/0snis00584.key")
	if err!= nil {
		panic(err)
	}
	decryptBytes := util.AESDecrypt(bytes,keyBytes)
	ioutil.WriteFile("e:/test_out.ts",decryptBytes,0666)
}

/**
测试合并
 */
func TestMerge(t *testing.T) {
	bytes1,_ :=ioutil.ReadFile("E:/test_out.ts")
	bytes2,_ :=ioutil.ReadFile("E:/test_out2.ts")
	allBytes := append(bytes1,bytes2...)
	ioutil.WriteFile("e:/test_out3.ts",allBytes,0666)
}

/**
测试转码
 */
func TestVideo(t *testing.T) {
	cmd := exec.Command("ffmpeg.exe",
		"-i","E:\\新建文件夹\\118abp00593-out\\result.ts",
		"-c","copy",
		//"-bsf:a","aac_adtstoasc",
		"-y", "E:\\新建文件夹\\result\\118abp00593.mp4")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("stdout=[%s]\n", string(out))
}

/**
测试日志
 */
func TestLog(t *testing.T) {

	myLog.Error("dfdfa")

	myLog.Info("dfdfa:%v",util.Result{true,"xxx",[]interface{}{""}})

	log.Printf("dfdfa:%s","xx")
	myLog.Error("dfdfa")

}

/**
测试传参
 */
 func Test1(t *testing.T) {
	 // 输入执行的命令
	 cmd := exec.Command("copy","str.dmm.com")
	 buf, err := cmd.Output()
	 fmt.Printf("%s\n%s",buf,err)
	 //
	 //// 获取子进程标准输出
	 //stdout, _ := cmd.StdoutPipe()
	 //
	 //// 执行命令
	 //cmd.Start()
	 //
	 //// 读取子进程
	 //reader := bufio.NewReader(stdout)
	 //for {
		// line, err2 := reader.ReadString('\n')
		// if err2 != nil || io.EOF == err2 {
		//	 break
		// }
		// // 转换CMD的编码为GBK
		// reader := transform.NewReader(
		//	 bytes.NewReader([]byte(line)),
		//	 simplifiedchinese.GBK.NewDecoder(),
		// )
		// d, _ := ioutil.ReadAll(reader)
	 //
		// // 将子进程的内容输出
		// print(string(d))
	 //}
 }

 /**
 temp
  */
func TestTemp(t *testing.T) {

	bytes,err := ioutil.ReadFile(`C:\Users\97038\Desktop\新建文本文档 (2).txt`)
	if err != nil {
		panic(err)
	}
	str := string(bytes)
	arr := strings.Split(str,"\r\n")
	for _, v := range arr {
		fmt.Println(v + `,file@C:\h\av2\m3u8\all\` + v+".m3u8")
	}
}

