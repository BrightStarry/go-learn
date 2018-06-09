package main

import (
	"fmt"
	"time"
	"os"
	"bufio"
	"io"
	"regexp"
	"log"
	"strconv"
	"net/url"
)

/*日志处理*/

/*日志结构*/
type Message struct{
	// 日志中的时间
	TimeLocal time.Time
	// 流量
	BytesSent int
	// 路径/方法/状态等
	Path,Method,Schema,Status string

}


/*读取接口*/
type Reader interface {
	Read(readChannel chan []byte)
}

/*写入接口*/
type Writer interface {
	Write(writeChannel chan *Message)
}


/*日志处理类*/
type LogProcess struct {
	// 读取通道
	readChannel chan []byte
	// 写入通道
	writeChannel chan *Message
	// 读取接口
	read Reader
	// 写入接口
	write Writer
}

/*文件读取对象*/
type FileReader struct {
	// 读取文件路径
	path string
}

/*influxDB写入对象*/
type InfluxDBWriter struct {
	// influxDB data source
	influxDBDsn string
}

/*influxDB写入模块-写入方法*/
func (this *InfluxDBWriter) Write(writeChannel chan *Message) {
	for v := range writeChannel {
		fmt.Println(v)
	}
}

/*文件读取对象-读取方法*/
func (this *FileReader) Read(readChannel chan []byte) {
	/**
		打开文件
	 */
	file,err := os.Open(this.path)
	if err != nil {
		panic(fmt.Sprintf("打开文件异常:%s",err.Error()))
	}
	defer file.Close()

	/**
		从文件末尾逐行读取数据
	 */
	// 设置文件下次读取的位置,指 相对于文件末尾(参数2)的第0个字节(0),也就是文件末尾
	file.Seek(0,2)
	reader :=bufio.NewReader(file)
	for{
		// 读取文件内容,直到 换行
		line,err := reader.ReadBytes('\n')
		// 读到文件末尾,则等待
		if err == io.EOF{
			time.Sleep(500 * time.Millisecond)
			continue
		}else if err != nil {
			// 其他异常
			panic(fmt.Sprintf("读取文件异常:%s",err.Error()))
		}
		readChannel <- line
	}
}


/**
解析模块
 */
func (this *LogProcess) Process() {

	/**
	我自己福利球网站的nginx日志格式
 	'$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"'
	例如
	101.41.60.216 - -
	[09/Jun/2018:12:37:29 +0800]
	"GET / HTTP/1.1" 502 575 "-"
	"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36" "-"
	 */

	 // 时区
	 location,_ := time.LoadLocation("Asia/Hangzhou")

	 // 正则
	r := regexp.MustCompile(`([\d\.]{7,16}) - (\S+) \[(.*)\] \"(\S+) (/\S*) (\S*)\" (\d+) (\d+) \"(\S+)\" \"(.+)\" \"(.*)\"`)
	for v := range this.readChannel {
		// 返回匹配到的内容 []string
		ret := r.FindStringSubmatch(string(v))
		if len(ret) != 12 {
			log.Println("解析格式异常:",string(v))
			continue
		}

		var err error
		message := &Message{}
		// 时间, 时间可能为 "-"
		if ret[3] != "-" {
			message.TimeLocal,err = time.ParseInLocation("02/Jan/2006:15:04:05 -0700",ret[3],location)
			if err != nil {
				log.Println("时间格式异常:",err.Error(),)
			}
		}
		// 字节数
		message.BytesSent,_ = strconv.Atoi(ret[8])

		// 请求路径  例如 /list?pageNo=2,需要解析出路径,不要参数
		u,err :=url.Parse(ret[5])
		if err != nil {
			log.Println("解析请求路径异常:",ret[5])
			continue
		}
		message.Path = u.Path

		// 协议, 例如http
		message.Schema = ret[6]

		// 状态
		message.Status = ret[7]

		this.writeChannel <- message
	}
}


func main() {
	read := &FileReader{
		path:"./access.log",
	}
	write := &InfluxDBWriter{
		influxDBDsn: "username&password..",
	}

	logProcess := &LogProcess{
		readChannel:make(chan []byte),
		writeChannel:make(chan *Message),
		read:read,
		write:write,
	}

	go logProcess.read.Read(logProcess.readChannel)
	go logProcess.Process()
	go logProcess.write.Write(logProcess.writeChannel)

	time.Sleep(100 * time.Second)
}
