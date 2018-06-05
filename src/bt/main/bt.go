package main

import (
	"net"
	"errors"
	"fmt"
	"encoding/hex"
	"crypto/rand"
	"time"
	"bytes"
	"io"
	"bt/util"
	"encoding/binary"
	"log"
	"io/ioutil"
	"crypto/sha1"
)

/*仅测试 从peer获取metadata信息*/

//获取种子元信息时,第一条握手信息的前缀, 28位byte. 第2-20位,是ASCII码的BitTorrent protocol,
// 第一位19,是固定的,表示这个字符串的长度.后面八位是BT协议的版本.可以全为0,某些软件对协议进行了扩展,协议号不全为0,不必理会.
var handshakePrefix = []byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114,
	111, 116, 111, 99, 111, 108, 0, 0, 0, 0, 0, 16, 0, 1}

// BEP-009定义了三种不同的消息类型, 值分别为 0，1，2
const (
	// 表示请求消息类型
	REQUEST = iota // 0
	// 表示数据消息类型
	DATA // 1
)

const (
	// 表示是一个扩展消息
	EXTENDED = 20
	// 表示握手的一位
	HANDSHAKE = 0
	// dht中metadata数据默认分块大小，也是扩容大小，给buffer扩容时使用 2 ^ 14
	BLOCK = 16384
	// 最大可接受的metadata_size
	MaxMetadataSize = BLOCK * 10
)

func main() {
	//addr := ""
	//conn := connect(addr)

	//bytes := hexInfoHashToByte("eb8abb5d2b4711b4d545b9d0ebb05f22b63f5ca3")
	//fmt.Println(bytes)

	r := Request{"7ac05fc84bdf9f25fed2df8485a8160194972dd0","154.45.216.216:2308"}
	metadata,err := fetchMetadata(r)
	if err != nil {
		fmt.Println("异常：",err)
		return
	}

	d,err := util.Decode(metadata)
	if err != nil {
		fmt.Println(err)
	}
	dict := d.(map[string]interface{})
	fmt.Println(dict)
}

/*请求对象*/
type Request struct {
	InfoHash string
	Address  string
}

/*获取metadata*/
func fetchMetadata(r Request)(metadata []byte,err error){
	infoHash := r.InfoHash
	address := r.Address

	var (
		// 存储分块数据
		blocks [][]byte
		// 分块数
		blockNum int
		// 元数据总长
		metadataSize int
		// 标识
		utMetadata int
		// 标识是否读取完成
		isDone bool
		dial net.Conn
		msgType byte
		length int
		extendID byte
		payload []byte
	)

	// 进行拨号连接
	dial, err = net.DialTimeout("tcp", address, time.Second*15)
	if err != nil {
		return
	}

	// 转换
	conn := dial.(*net.TCPConn)
	// 设置一个连接的关闭行为，当sec<0(默认)，在后台完成发送操作； ==0，丢弃未发送或未确认的消息; >0,和<0时类似，在某些操作系统中，支持sec秒后，将未发送消息丢弃
	conn.SetLinger(0)
	defer conn.Close()

	// 创建缓冲
	data := bytes.NewBuffer(nil)
	// 对其直接扩容
	data.Grow(BLOCK)

	// 依次进行： 发送握手，读取握手响应数据，处理握手响应数据，发送扩展握手
	//if sendHandshakeMessage(conn, infoHash, randomBytes(20)) != nil ||
	//	read(conn, 68, data) != nil ||
	//	onHandshake(data.Next(68)) != nil ||
	//	sendExtHandshake(conn) != nil {
	//	return nil,errors.New("握手失败")
	//}

	if err = sendHandshakeMessage(conn,infoHash,randomBytes(20)); err != nil {
			return
	}
	if err = read(conn, 68, data); err != nil {
		return
	}
	if err = onHandshake(data.Next(68)); err != nil {
		return
	}
	if err = sendExtHandshake(conn); err != nil {
		return
	}

	for {
		// 读取一个字节的消息
		length, err = readMessage(conn, data)
		if err != nil {
			return
		}
		if length == 0 {
			continue
		}

		// 再读取下个字节
		msgType, err = data.ReadByte()
		if err != nil {
			return
		}

		// 根据该字节判断消息类型
		switch msgType {
		// 如果是扩展消息
		case EXTENDED:
			// 继续读取一个字节
			extendID, err = data.ReadByte()
			if err != nil {
				return
			}

			// 读取所有字节
			payload, err = ioutil.ReadAll(data)
			if err != nil {
				return
			}


			// 处理还未请求分块数据的情况
			if extendID == 0 {
				// 如果此时分块数据不为空，表示出错了
				if blocks != nil {
					return
				}

				// 获取握手响应中，对方传过来的一些数据
				utMetadata, metadataSize, err = getUtMetadataSize(payload)
				if err != nil {
					return
				}

				// 计算出分块数，此处/ 是整除
				blockNum = metadataSize / BLOCK
				// 取余，如果有余数，需要+1
				if metadataSize % BLOCK != 0 {
					blockNum++
				}

				// 构造出block
				blocks = make([][]byte,blockNum)
				// 进行分块请求
				go sendRequestBlock(conn,utMetadata,metadataSize,blockNum)
				// 跳出该循环，进行下一轮读取
				continue
			}

			// 进入到该判断时，应该是第二次循环，第一次循环执行了if extendID == 0中的语句，已经构造出了block
			// 所以此时如果还是空，则表明之前未进入过if extendID == 0，表示出错
			if blocks == nil {
				return
			}

			// 对响应进行解码
			d,index,err1 :=util.DictDecode(payload,0)
			if err1 != nil {
				err = err1
				return
			}
			dict := d.(map[string]interface{})

			// 检查key是否存在
			if err = parseKeys(dict,[][]string{{"msg_type", "int"},{"piece", "int"}}); err != nil {
				return
			}

			// 如果不是数据类型，也就是值为1，则跳过
			if dict["msg_type"].(int) != DATA {
				continue
			}

			// 获取当前是第几个分块
			blockIndex := dict["piece"].(int)
			// 该分块长度为，当前读取的数据总长 - 之前读取出的2个字节 - 进行dict解码后，当前索引的位置（也就是减去该map的长度）
			blockLen := length - 2 - index
			// 如果当前分块不是最后一块，大小却不是默认值，或者是最后一块，大小却不等于余数，则出错
			if blockIndex != blockNum-1 && blockLen != BLOCK ||
				blockIndex == blockNum-1 && blockLen != metadataSize%BLOCK{
					return
			}
			// 给对应分块赋值
			blocks[blockIndex] = payload[index:]

			// 判断分块是否读取完成
			for _,b := range blocks{
				if len(b) == 0 {
					isDone = false
				}
			}
			isDone = true
			// 如果读取完成
			if isDone {
				// 拼接所有分块，此处传入的就是二维数组，根据数组的第一维索引进行拼接
				metadata = bytes.Join(blocks,nil)
				// 校验哈希
				hash := sha1.Sum(metadata)
				// 之前返回的hash对象是[20]byte类型，判断需要[]byte类型，直接写hash[:]将数组转为分片
				if !bytes.Equal(hexInfoHashToByte(infoHash), hash[:]) {
					return
				}
				return
			}
		default:
			data.Reset()
		}
	}
}

/*连接到peer*/
func connect(addr string) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(errors.New(fmt.Sprintln("连接到", addr, ",异常:"+err.Error())))
	}
	return conn
}

/*发送握手消息*/
func sendHandshakeMessage(conn *net.TCPConn, infoHash string, peerId []byte) error {
	data := make([]byte, 68)

	copy(data[:28], handshakePrefix)
	copy(data[28:48], hexInfoHashToByte(infoHash))
	copy(data[48:], peerId)

	// 设置超时时间10s
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	_, err := conn.Write(data)
	return err
}

/*处理握手响应*/
func onHandshake(data []byte) (err error) {
	if !bytes.Equal(handshakePrefix[:20], data[:20]) && data[25]&0x10 != 0 {
		err = errors.New("无效握手响应")
	}
	return
}

/*发送扩展握手协议,获取ut_metadata  和 metadata_size*/
func sendExtHandshake(conn *net.TCPConn) error {
	// 构造扩展协议
	data := append(
		[]byte{EXTENDED, HANDSHAKE},
		util.Encode(
			map[string]interface{}{
				"m": map[string]interface{}{
					"ut_metadata": 1,
				},
			},
		)...,
	)
	return sendMessage(conn, data)
}

/*发送请求分块请求*/
func sendRequestBlock(conn *net.TCPConn, utMetadata int, metadataSize int, blockNum int) {
	buffer := make([]byte, 1024)
	// 依次请求每个分块
	for i := 0; i < blockNum; i++ {
		buffer[0] = EXTENDED
		buffer[1] = byte(utMetadata)
		msg := util.Encode(map[string]interface{}{
			"msg_type": REQUEST,
			"piece":    i,
		})

		// 拷贝到buffer
		length := len(msg) + 2
		copy(buffer[2:length], msg)

		// 发送请求
		sendMessage(conn, buffer[:length])
	}
	// ? gc
	buffer = nil
}

/*从数据中获取utMetadata 和 metadataSize*/
func getUtMetadataSize(data []byte) (utMetadata, metadataSize int, err error) {
	// 解码
	v, err := util.Decode(data)
	if err != nil {
		return
	}

	// 转为map类型
	dict, ok := v.(map[string]interface{})
	if !ok {
		err = errors.New("无效的dict类型")
		return
	}

	// 判断map中对应的key是否存在和key类型
	if err = parseKeys(dict, [][]string{{"metadata_size", "int"}, {"m", "map"}}); err != nil {
		log.Println("key有误:", err)
		return
	}

	// 取出m
	m := dict["m"].(map[string]interface{})
	if err = parseKey(m, "ut_metadata", "int"); err != nil {
		log.Println("key有误:", err)
		return
	}

	utMetadata = m["ut_metadata"].(int)
	metadataSize = dict["metadata_size"].(int)

	if metadataSize > MaxMetadataSize {
		err = errors.New(fmt.Sprintf("metadata_size过长，当前长度:%v", metadataSize))
	}
	return

}

/*解析keys，传入二维数组，一维表示key，二维表示keyType，判断整个数组是否都存在对应类型的key*/
func parseKeys(data map[string]interface{}, pairs [][]string) error {
	for _, pair := range pairs {
		key, keyType := pair[0], pair[1]
		if err := parseKey(data, key, keyType); err != nil {
			return err
		}
	}
	return nil
}

/*解析key,判断是否有指定类型的key*/
func parseKey(data map[string]interface{}, key string, keyType string) error {
	// 获取值
	val, ok := data[key]
	if !ok {
		return errors.New("key不存在")
	}

	// 转换值类型
	switch keyType {
	case "string":
		_, ok = val.(string)
	case "int":
		_, ok = val.(int)
	case "list":
		_, ok = val.([]interface{})
	case "map":
		_, ok = val.(map[string]interface{})
	default:
		panic("无效类型")
	}

	// 转化是否成功
	if !ok {
		return errors.New("key类型不匹配")
	}
	return nil
}

/*发送消息, 前4个字节表示消息长度*/
func sendMessage(conn *net.TCPConn, data []byte) error {
	// 创建一个buffer
	lenBytes := bytes.NewBuffer(nil)
	// 将data的长度转为 [4]byte,写入buffer,
	// 看了下源码，如下方法实际调用了binary.BigEndian.PutUint32(),然后再将[]byte，用buffer的writer方法写入buffer
	binary.Write(lenBytes, binary.BigEndian, int32(len(data)))

	// 设置写入超时时间
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	// 追加上实际数据写入
	_, err := conn.Write(append(lenBytes.Bytes(), data...))
	return err
}

/*从连接中读取指定长度的消息到data*/
func read(conn *net.TCPConn, len int, data *bytes.Buffer) error {
	conn.SetReadDeadline(time.Now().Add(15 * time.Second))

	n, err := io.CopyN(data, conn, int64(len))
	// 如果读取有异常，或者字节数不等
	if err != nil || n != int64(len) {
		return errors.New("读取异常" + err.Error())
	}
	return nil
}

/*从连接中读取消息*/
func readMessage(conn *net.TCPConn, data *bytes.Buffer) (length int, err error) {
	// 读取4个字节到 data
	if err = read(conn, 4, data); err != nil {
		return
	}
	// 将读取到的字节转为int
	length = int(bytes2int(data.Next(4)))
	if length == 0 {
		return
	}

	// 读取后续所有消息到data
	if err = read(conn, length, data); err != nil {
		return
	}
	return

}

/*16进制infoHash转[]byte*/
func hexInfoHashToByte(infoHash string) []byte {
	result, _ := hex.DecodeString(infoHash)
	return result
}

/*生成指定长度的随机字符*/
func randomString(size int) string {
	return string(randomBytes(size))
}

/*生成指定长度的随机字节*/
func randomBytes(size int) []byte {
	buff := make([]byte, size)
	rand.Read(buff)
	return buff
}

/*[]byte 转 int64*/
func bytes2int(data []byte) uint64 {
	n, val := len(data), uint64(0)
	if n > 8 {
		panic("数据过长")
	}

	// 每次进行左位移，因为是大端序，第一次将data[0]左位移8-0-1 * 8个距离，也就是最左边
	for i, b := range data {
		val += uint64(b) << uint64((n-i-1)*8)
	}
	return val
}

/*int 转[]byte*/
func int2Bytes(data uint64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, data)
	return result
}
