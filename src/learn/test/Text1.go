package main

import (
	"fmt"
	"strings"
	"strconv"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"bytes"
	"os"
	"bufio"
	"io"
)

/*
	文本操作  / struct 互转 xml / os.Args取参  / 读取文件列表（解析xml）/
	模拟cp操作
*/
func main() {
	//base()
	//structOfXml()
	osOfArgs()
	//parseXml()
	//cp()
}

/*字符串基本操作*/
func base() {
	s := "hello world"
	// 是否包含
	fmt.Println(strings.Contains(s, "hello"))
	// 索引
	fmt.Println(strings.Index(s, "o"))

	ss := "1#2#345"
	// 分割字符串
	ssArr := strings.Split(ss, "#")
	fmt.Println(ssArr)
	// 合并字符串数组 为 字符串
	fmt.Println(strings.Join(ssArr, "#"))

	// 是否以该字符开头/结尾
	fmt.Println(strings.HasPrefix(s, "he"), strings.HasSuffix(s, "ld"))
}

/*字符串转换操作*/
func convert() {
	// int 转 string
	strconv.Itoa(10)
	// string 转 int
	strconv.Atoi("25")

	// string 转 bool  case "1", "t", "T", "true", "TRUE", "True":  case "0", "f", "F", "false", "FALSE", "False":
	strconv.ParseBool("false")
	// bool 转 string
	strconv.FormatBool(true)
	// string转float,32表示精度
	strconv.ParseFloat("3.14", 32)

	// int 转 string，可设置进制
	strconv.FormatInt(123, 10)
}

/*结构体序列化和反序列化*/
type person struct {
	Name string `xml:"name1,attr"` // 在属性后加这个，可以将该成员变量在xml转换时，转为xml标签的属性     `xml:",attr"`,逗号前的字符，则为别名
	Age  int
}

func structOfXml() {
	p := person{
		Name: "zx",
		Age:  18,
	}

	var xmlBytes []byte

	// 注意，这样写，这个bytes的作用域就在这个if中,
	// MarshaIndent()可以加前缀和标签间的缩进
	if bytes, err := xml.Marshal(p); err != nil {
		fmt.Println("序列化失败", err)
		return
	} else {
		xmlBytes = bytes
		fmt.Println(string(bytes))
	}

	p2 := new(person)
	// 反序列化，可能抛出error
	xml.Unmarshal(xmlBytes, p2)
	fmt.Println(p2)
}

/*os.Args取参*/
func osOfArgs() {
	// 默认有个index：0的参数（命令自己的路径名字），所以后面加3个参数的时候，参数总数为4
	//fmt.Println(len(os.Args))
	//fmt.Println(os.Args[1])

	/*
		以下两种方法，如果输入命令时，输入错误，会自动进行提示，也就是显示最后一个参数中的内容
		他是通过 flag.PrintDefaults() 函数打印出来的
	*/

	// 或
	// 如果追加  "-method 123" 这样的key/value参数
	// 使用该方法读取，参数1:key值，2：默认值，3：说明，  会返回一个*string
	// 下面这句话会读取到追加的  -method xxx 中的xxx
	//methodPtr := flag.String("method","defaultValue","method desc")
	// 下面会读取到追加的 -value 123 的123
	//valuePtr := flag.Int("value",-1,"value desc")
	// 解析
	//flag.Parse()
	//fmt.Println("xxx",*methodPtr,*valuePtr)

	// 或
	var method string
	var value int
	flag.StringVar(&method, "method", "default", "method of sample")
	flag.IntVar(&value, "value", 12355, "value of sample")
	flag.Parse()
	fmt.Println("xxx", method, value)
}

/*解析xml*/
func parseXml() {
	// 读取文件为字节
	content, _ := ioutil.ReadFile("H:\\goSpace\\go-learn\\.idea\\misc.xml")
	// 解码器
	decoder := xml.NewDecoder(bytes.NewBuffer(content))

	// 匿名函数，用于从xml属性数组中，获取对应属性名的属性值
	var getAttrFunc = func(attr []xml.Attr,name string) string {
		for _,item := range attr {
			if item.Name.Local == name {
				return item.Value
			}
		}
		return ""
	}

	// 状态机的一个标志，是否在 <component>标签中
	var inComponent bool

	// 循环，每次调用Token()都是获取xml的一个节点
	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		// 根据token的类型进行不同处理
		switch token := t.(type) {
			// 一个开始元素(标签)
			case xml.StartElement:
				name := token.Name.Local

				// 如果在<component>中
				if inComponent {
					if name == "option" {
						// 获取该标签中的 name属性
						value := getAttrFunc(token.Attr,"name")
						fmt.Println(value)
					}
				}else{
					if name == "component" {
						// 如果进入了该标签
						inComponent = true
					}
				}

			// 一个结束元素(标签)
			case xml.EndElement:
				// 如果离开了<component>
				if inComponent && token.Name.Local == "component" {
					inComponent = false
				}
		}
	}
}

/*模拟拷贝命令*/
func cp() {
	// 两个参数：  和强制
	var showProgress,force bool

	// 此处有个bug，即使加了 -f -v，也用的默认值false
	flag.BoolVar(&force,"f",false,"是否强制")
	flag.BoolVar(&showProgress,"v",false,"解释正在做的事情")
	flag.Parse()

	// 参数个数小于两个，则有误
	if flag.NArg() < 2 {
		// 打印用法说明
		flag.Usage()
		return
	}

	// 第0个参数为 源文件，第1个参数为目标文件
	cp1(flag.Arg(0),flag.Arg(1),showProgress,force)
}
/*模拟拷贝命令,真正操作方法*/
func cp1(src,dest string,showProgress,force bool) {
	// 如果不强制，需要判断目标文件是否存在
	if !force {
		// 文件状态
		_,err := os.Stat(dest)
		// 没有异常表示，存在
		isExist := err ==nil || os.IsExist(err)
		if isExist {
			fmt.Println(dest,"已存在，是否重写?y/n")
			// 读取标准输入
			reader := bufio.NewReader(os.Stdin)
			// 读取一行
			data,_,_ := reader.ReadLine()
			if strings.TrimSpace(string(data)) != "y" {
				return
			}
		}
	}

	// 进行拷贝
	srcFile,err := os.Open(src)
	if err != nil {
		fmt.Println("打开源文件有误:",err)
		return
	}
	defer srcFile.Close()

	destFile,err := os.Create(dest)
	if err != nil {
		fmt.Println("创建目标文件有误:",err)
		return
	}
	defer destFile.Close()

	// 返回拷贝的字节数n
	n,err :=io.Copy(destFile,srcFile)

	// 如果 -v ，则输出说明
	if showProgress {
		fmt.Println(src,"拷贝到",dest,"字节大小：",n)
	}
}
