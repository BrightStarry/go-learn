### Go语言学习

#### 打包运行方式
- 在项目根目录对需要导入的go包执行 go install xxx,在pkg目录中生成.a文件
- 再对需要运行的go文件的包名执行 go install xxx,在bin目录生成exe文件(注意，运行的主文件必须package main)
- 也可直接执行 go build xxx.go ，生成exe或sh可执行文件

#### bug
- runnerw.exe: CreateProcess failed with error 216:    修改package main即可
- 函数名和 类库名相同，会有bug
- testList := make([]interface{},3) 和 testList := []string{"1"},前者是[]interface{}类型，但后者不是
- 用io.Open方式打开文件,写入数据时会报 拒绝访问异常. 需要用 os.OpenFile


#### 要点
- 有些方法接收的参数是指针，则表示该方法中，可以改变该参数的引用的值。否则是值传递。
- 当方法参数接收的是指针类型时，出现方法都无法调用，是因为该指针类型的原类型是接口，将其改为实现类即可（例如net.Conn改为net.TcpConn）
- 测试文件必须以_test结尾，然后在方法中添加参数t *testing.T
- go中所有函数均为值传递,包括数组.但用分片是引用传递.分片其实就是一个数组的视图.
    - 如果改变分片的值,原数组的值也会被改变.
    - 如果用数组A构建的分片B的容量超过了数组B,那么系统会重新分配一个更大的底层数组
    - 如果分片长度递增,那么每次长度等于容量,进行扩容时,会按照2的幂等进行扩容,即 2,4,6,8,16,32
    - 如果想要删除分片中的某个元素,则 s = append(s[:3],s[4:]),则删除了分片s的下标为3的元素
- rune类型相当于char    
#### 跨平台编译
~~~
在根目录执行
set GOOS=linux
set GOARCH=amd64
然后再执行其他命令

要恢复的话执行
set GOOS=windows
set GOARCH=amd64

也可以直接如下，构建出linux的执行文件
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build HTTPServer.go
~~~


#### 多main编译
```
对于如下目录结构
-GOPATH
    -src
        -bt
            -main
            -util
        -web
            -main
            -util
如果想要install bt中的main，需要在根目录执行
go install bt/main
```

#### 奇淫巧技
- goland快捷键： C+A+v ,快速生成方法返回对象