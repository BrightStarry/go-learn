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
	"math/rand"
	"github.com/influxdata/influxdb/client/v2"
	"strings"
	"flag"
	"net/http"
	"encoding/json"
)

/*日志处理*/

/*日志结构*/
type Message struct {
	// 日志中的时间
	TimeLocal time.Time
	// 流量
	BytesSent int
	// 路径/方法/状态等
	Path, Method, Schema, Status string
	//
	UpstreamTime, RequestTime float64
}

/*系统状态监控*/
type SystemInfo struct {
	// 总处理日志行数
	HandleLine int `json:"handleLine"`
	// 系统吞吐量
	Tps float64 `json:"tps"`
	// read channel 长度
	ReadChanLen int `json:"readChanLen"`
	// write channel 长度
	WriteChanLen int `json:"writeChanLen"`
	// 运行总时间
	RunTime string `json:"tunTime"`
	// 错误数
	ErrNum int `json:"errNum"`
}

const (
	TypeHandleLine = 0
	TypeErrNum     = 1
)

var TypeMonitorChan = make(chan int, 200)

/*监控器*/
type Monitor struct {
	// 系统开始运行时间
	startTime time.Time
	// 暂存定时器每次计算的tps
	tpsTemp []int
	data    SystemInfo
}

func (this *Monitor) start(logProcess *LogProcess) {
	// 一旦监控通道收到数据,就对 异常,或者行数统计进行 ++
	go func() {
		for i := range TypeMonitorChan {
			switch i {
			case TypeHandleLine:
				this.data.HandleLine += 1
			case TypeErrNum:
				this.data.ErrNum += 1
			}
		}
	}()

	// 定时记录当前 处理日志总长
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			//
			<-ticker.C
			// 存储每次定时任务运行时的处理日志总长
			this.tpsTemp = append(this.tpsTemp, this.data.HandleLine)
			// 这个数组实际上只要存储最近一次和最近第二次的 日志总长就可以了,所以如果长度过大,就删除之前的
			if len(this.tpsTemp) > 2 {
				this.tpsTemp = this.tpsTemp[1:]
			}
		}
	}()

	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		// 系统运行时间
		this.data.RunTime = time.Now().Sub(this.startTime).String()
		this.data.ReadChanLen = len(logProcess.readChannel)
		this.data.WriteChanLen = len(logProcess.writeChannel)
		// 防止刚运行,就查看监控,
		if len(this.tpsTemp) >= 2 {
			// 用存储的最近一次的日志总行数 -  最近第二次的日志总行数, 除以定时器的运行间隔5,即可得出近乎实时的tps
			this.data.Tps = float64(this.tpsTemp[1]-this.tpsTemp[0]) / 5
		}

		// 对象转json, 对象/前缀无/用\t作为缩进
		result, _ := json.MarshalIndent(this.data, "", "\t")
		io.WriteString(writer, string(result))
	})
	http.ListenAndServe(":9193", nil)
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

	configArr := strings.Split(this.influxDBDsn, "@")

	//创建http连接
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: configArr[0],
		//Username:
		//Password:
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// 创建客户端
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		// 数据库名
		Database: configArr[1],
		// 精度, "s"表示秒
		Precision: configArr[2],
	})
	if err != nil {
		log.Fatal(err)
	}

	for v := range writeChannel {
		// 构造字段(需要作为选择条件的),标签
		// Tags: Path,Method,,Schema,Status
		tags := map[string]string{"Path": v.Path, "Method": v.Method, "Schema": v.Schema, "Status": v.Status}
		fields := map[string]interface{}{
			"BytesSent":    v.BytesSent,
			"RequestTime":  v.RequestTime,
			"UpstreamTime": v.UpstreamTime,
		}

		// 表
		pt, err := client.NewPoint("nginx_log", tags, fields, v.TimeLocal)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)

		// 写入
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}
		log.Println("写入成功!")
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}
	}()

}

/*文件读取对象-读取方法*/
func (this *FileReader) Read(readChannel chan []byte) {
	/**
		打开文件
	 */
	file, err := os.Open(this.path)
	if err != nil {
		panic(fmt.Sprintf("打开文件异常:%s", err.Error()))
	}
	defer file.Close()

	/**
		从文件末尾逐行读取数据
	 */
	// 设置文件下次读取的位置,指 相对于文件末尾(参数2)的第0个字节(0),也就是文件末尾
	file.Seek(0, 2)
	reader := bufio.NewReader(file)
	for {
		// 读取文件内容,直到 换行
		line, err := reader.ReadBytes('\n')
		// 读到文件末尾,则等待
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			// 其他异常
			panic(fmt.Sprintf("读取文件异常:%s", err.Error()))
		}
		// 每读取到一行,往监控通道传入标记
		TypeMonitorChan <- TypeHandleLine

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
	location, _ := time.LoadLocation("Asia/Hangzhou")

	// 正则
	r := regexp.MustCompile(`([\d\.]{7,16}) - (\S+) \[(.*)\] \"(\S+) (/\S*) (\S*)\" (\d+) (\d+) \"(\S+)\" \"(.+)\" \"(.*)\"`)
	for v := range this.readChannel {
		// 返回匹配到的内容 []string
		ret := r.FindStringSubmatch(string(v))
		if len(ret) != 12 {
			// 出现异常,往监控通道传入标记
			TypeMonitorChan <- TypeErrNum
			log.Println("解析格式异常:", string(v))
			continue
		}

		var err error
		message := &Message{}
		// 时间, 时间可能为 "-"
		if ret[3] != "-" {
			TypeMonitorChan <- TypeErrNum
			message.TimeLocal, err = time.ParseInLocation("02/Jan/2006:15:04:05 -0700", ret[3], location)
			if err != nil {
				log.Println("时间格式异常:", err.Error(), )
			}
			// TODO 此处暂时将时间修改为当前,因为我的nginx日志中的数据都是旧数据,在展示时不好看
			message.TimeLocal = time.Now()
		}
		// 字节数
		message.BytesSent, _ = strconv.Atoi(ret[8])

		// 请求路径  例如 /list?pageNo=2,需要解析出路径,不要参数
		u, err := url.Parse(ret[5])
		if err != nil {
			TypeMonitorChan <- TypeErrNum
			log.Println("解析请求路径异常:", ret[5])
			continue
		}
		message.Path = u.Path

		// 协议, 例如http
		message.Schema = ret[6]

		// 状态
		message.Status = ret[7]

		// 因为我自己的nginx文件中没有记录响应时间等,所以随机生成
		message.UpstreamTime = rand.Float64()
		message.RequestTime = rand.Float64()

		this.writeChannel <- message
	}
}

func main() {
	// 读取传入的参数
	var path, influxDBDsn string
	flag.StringVar(&path, "path", "./access.log", "读取文件路径")
	flag.StringVar(&influxDBDsn, "influxDBDsn", "http://106.14.7.29:8086@zx@s", "influxDB data source")
	flag.Parse()

	// 构建读取接口
	read := &FileReader{
		path: path,
	}
	// 构建写入接口
	write := &InfluxDBWriter{
		influxDBDsn: influxDBDsn,
	}

	// 构建日志处理对象
	logProcess := &LogProcess{
		readChannel:  make(chan []byte, 200),
		writeChannel: make(chan *Message, 200),
		read:         read,
		write:        write,
	}

	// 异步进行操作
	go logProcess.read.Read(logProcess.readChannel)
	// 由于读取是最快的,所以处理模块和写入模块可以多开几个协程同时处理
	for i := 0; i < 2; i++ {
		go logProcess.Process()
	}
	for i := 0; i < 4; i++ {
		go logProcess.write.Write(logProcess.writeChannel)

	}

	// 设置监控
	monitor := &Monitor{
		startTime: time.Now(),
		data:      SystemInfo{},
	}
	monitor.start(logProcess)

}
