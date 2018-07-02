#### 浏览器劫持插件

- 参考了github上一个https代理的项目

#### 生成证书
- [linux安装openssl](https://blog.csdn.net/shiyong1949/article/details/78212971?locationNum=10&fps=1)
- [生成自定义证书](https://blog.csdn.net/oldmtn/article/details/52208747)
- [windows信任证书](https://blog.csdn.net/xiuye2015/article/details/54599331)

#### 下载到certmgr.exe文件
- [下载地址](https://go.microsoft.com/fwlink/?LinkID=698771)
- 管理员执行如下命令，可以静默信任证书
>  .\certmgr.exe -add  -c ./ca/ca.cer -r localMachine -s AuthRoot
- [如何在运行bat时获取管理员权限(可参考resources/addCert.bat)](https://zhidao.baidu.com/question/2202641660666565188.html)

- https代理，返回自定义证书，篡改响应实现跳转，参考[该项目](https://github.com/sheepbao/gomitmproxy)

#### 总结
- 本是无意中街道的外包项目，克服了一个个困难。但最后说不能用网卡监听/代理/dns等方法去实现，因为会被网吧禁掉。那就没有别的办法了，
放弃了。