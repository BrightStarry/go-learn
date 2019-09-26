package util

import (
	"os"
	"github.com/grafov/m3u8"
	"bufio"
	"strings"
)

/**
处理m3u8文件
 */


 /**
 get uris and keyUri
  */
func GetUrisAndKeyURIByM3u8File(m3u8FilePath string) (uris []string,keyURI string,err error) {
	m3u8File, err := os.Open(m3u8FilePath)
	if err != nil {
		return
	}
	result, _, err := m3u8.DecodeFrom(bufio.NewReader(m3u8File), true)
	if err != nil {
		return
	}
	data := result.(*m3u8.MediaPlaylist)


	/**
	获取ts 的uri列表
	 */
	uris = make([]string,0)
	for _,x := range data.Segments {
		if x == nil {
			continue
		}
		uris = append(uris, x.URI)
	}
	/**
	获取keyURI
	 */
	keyURI = strings.TrimPrefix(data.Key.URI,"file@")
	return
}