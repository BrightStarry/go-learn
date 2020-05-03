40B7791AC6787268A01B62FAD4DB233C
下载
```
./aria2c -i C:\h\av2\m3u8\all\118abp00108.m3u8 -j300 -x256 -m0 -s1024 -d ./118abp00108 
./aria2c -i C:\h\av2\m3u8\all\118abp00635.m3u8 -j300 -x256 -m0 -s1024 -d ./118abp00635
E:/process/aria2c -i C:\h\av2\m3u8\all\118abp00637.m3u8 -j256 -x1024 -m0 -s1024 -k 128k -d E:/process/tsTODO/118abp00637
-i example.m3u8	下载列出的 example.m3u8
-j100	每个队列项的最大并行下载数 100
-x32	每次下载到一台服务器的最大连接数 32
-m0	尝试次数，0 意味着无限
-s1024	使用 1024 个连接下载文件
–all-proxy=http://127.0.0.1:1080	使用代理服务器 127.0.0.1 进行 HTTP，端口 1080
```

```
调用idm
.\IDMan.exe -d "https://str.dmm.com:443/digital/st1:M60Dzvs0lJVTYRZXdBQpE9ECwfPaXiZF3joWn5vd-n4XPEf9RttYx3WxCmMpWMAw/3Gc6MVxQvftqVNSW6fuZ3Ul/-/media_b3000000_0.ts" /s /p "e:\test"  /f "media_b3000000_0.ts"
```


解密
```
加密媒体使用 AES-128-CBC 加密，这里借助 Openssl 解密
在解密之前需要得知密钥，即 KEY 文件，我们尽量使用简单的程序完成这一操作，在这里我使用常见的 Notepad++ 打开 KEY，看见的会是一行乱码文件，将其全选中，在菜单栏依次选择 插件-Converter-‘ASCII -> HEX’，会得到一行十六位字符的密码，复制。
openssl aes-128-cbc -d -K 40B7791AC6787268A01B62FAD4DB233C -iv 00000000000000000000000000000000 -nosalt -in test.zx -out test.out.ts
-d 解密
-K 密钥
-IV 00000000000000000000000000000000 没有 IV
-nosalt 不加盐
```

合并
```
使用 DOS 的 copy 合并三个 ts 文件命令如下
copy /B media_b6000000_1.out.ts+media_b6000000_2.out.ts+media_b6000000_3.out.ts all.ts
/B	表示一个二进位文件

合并目录下下所有ts
copy /B  E:\新建文件夹\118abp00108\out\*.ts  E:\新建文件夹\118abp00108\out\all.ts
```

转码
```
ffmpeg -i all.ts -c copy -bsf:a aac_adtstoasc -y all.mp4
ffmpeg -i E:\新建文件夹\pcotta071\out\all.ts -c copy  -y E:\新建文件夹\all1.mp4
a aac_adtstoasc：（官方注释）将 MPEG-2/4 AAC ADTS 转换为 MPEG-4 音频特定配置比特流。此过滤器从 MPEG-2/4 AAC ADTS 标头创建 MPEG-4 AudioSpecificConfig 并删除 ADTS 标头。例如，当将 AAC 流从原始 ADTS AAC 或 MPEG-TS 容器复制到 MP4A-LATM、FLV 文件或 MOV/MP4 文件以及相关格式（如 3GP 或 M4A ）时，需要此过滤器。请注意，它是为 MP4A-LATM 和 MOV/MP4 以及相关格式自动插入的。
TS 文件还是 MP4 文件呢？我作为外行个人猜测，TS 是可以直接首位合并的格式，且单个 TS 分段可以播放，加之 TS copy 编码转化为 MP4 时码率会降低，或许因为 TS 合并后有视频标头等大量冗余信息，且考虑到完全是更换“容器”，所以选择 MP4 格式。
还有哪些可以优化之处？下载速度（队列、连接、线程），考虑到下载完全不影响最后的视频文件，所以只能是 ffmpeg，而在目前看来，除了 copy 编码+aac_adtstoasc，暂时没有值得使用的参数了。
```

bug
```
合并ts时遇到一个bug。
无论是自己写代码，还是使用dos下的copy命令，只要不是一次性将所有ts文件合并成一个。而是批量合并
（即使想同时合并也无法做到，内存不够），就会导致合并后的文件大小比正确的文件大32K。
```