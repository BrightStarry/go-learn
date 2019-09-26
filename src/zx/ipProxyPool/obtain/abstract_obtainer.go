package obtain

import (
	"time"
	"fmt"
)

func init(){
	InitWebObtainers()
}

/**
	初始化 Obtinaer
 */
func InitWebObtainers() {
	WebObtainers =  []Obtainer{
		&XiciHttpObtainer{NewDefaultWebInfo("http://www.xicidaili.com/nt/", "西刺Http", 0, 20*time.Minute)},
		&XiciHttpsObtainer{NewDefaultWebInfo("http://www.xicidaili.com/wn/", "西刺Https", 1, 30*time.Minute)},
		&Ip66CommonObtainer{NewDefaultWebInfo("http://www.66ip.cn/mo.php?tqsl=", "66ip免费", 2, 2*time.Minute)},
		&Ip66AnonymityObtainer{NewDefaultWebInfo("http://www.66ip.cn/nmtq.php?proxytype=2&api=66ip&area=1&getnum=", "66ip匿名（国外）", 3, 1*time.Minute)},
		&KuaidailiNnObtainer{NewDefaultWebInfo("https://www.kuaidaili.com/free/inha/", "快代理高匿", 4, 6 * time.Hour)},
		&KuaidailiCommonObtainer{NewDefaultWebInfo("https://www.kuaidaili.com/free/intr/", "快代理普通", 2, 6 * time.Hour)},
		&Ip3366Obtainer{NewDefaultWebInfo("http://www.ip3366.net/?stype=%d&page=%d", "ip3366", 2, 15 * time.Minute)},
		&Ip89Obtainer{NewDefaultWebInfo("http://www.89ip.cn/tqdl.html?api=1&num=", "ip89", 2, 3 * time.Minute)},
	}
}

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