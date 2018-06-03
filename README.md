### Go语言学习

#### 打包运行方式
- 在项目根目录对需要导入的go包执行 go install xxx,在pkg目录中生成.a文件
- 再对需要运行的go文件的包名执行 go install xxx,在bin目录生成exe文件(注意，运行的主文件必须package main)
- 也可直接执行 go build xxx.go ，生成exe或sh可执行文件

#### bug
- runnerw.exe: CreateProcess failed with error 216:    修改package main即可
- 函数名和 类库名相同，会有bug