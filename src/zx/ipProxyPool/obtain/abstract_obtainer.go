package obtain

import (
	"time"
	"fmt"
)

// 保存所有网站信息
var WebObtainers  []Obtainer

/**
   获取者接口
*/
type Obtainer interface {
	// 初始获取全部ip方法
	InitObtain() int
	// 增量获取ip方法
	IncrementObtain() int
	// 获取WebObtainer
	GetWebObtainer() *WebObtainer
}


/**
	目标网站获取者
 */
type WebObtainer struct {
	// 网站名-作日志打印
	Name string
	// 网址
	Url string
	// 爬取间隔
	Interval time.Duration
	// 最后标记-作增量使用
	LastLabel interface{}
	// 权重
	Weight uint8
}

func (this *WebObtainer) String() string {
	return fmt.Sprintf("权重:%v,名称:%v,网址:%v,间隔:%v,最后标记:%v \n", this.Weight, this.Name, this.Url, this.Interval, this.LastLabel)
}

/**
	创建默认网站对象
	url:网址
	name:网站名
	interval:爬取间隔
	weight: 权重
 */
func NewDefaultWebInfo(url string, name string, weight uint8, interval time.Duration) *WebObtainer {
	return &WebObtainer{
		Url:      url,
		Name:     name,
		Interval: interval,
		Weight:   weight,
	}
}