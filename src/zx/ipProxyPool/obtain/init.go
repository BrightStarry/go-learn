package obtain

import (
	"time"
)

func init(){
	InitWebObtainers()
}

/**
	初始化 Obtinaer
 */
func InitWebObtainers() {
	var xiciHttp Obtainer = &XiciHttpObtainer{NewDefaultWebInfo("http://www.xicidaili.com/nt/", "西刺Http", 0, 20*time.Minute)}
	var xiciHttps Obtainer = &XiciHttpsObtainer{NewDefaultWebInfo("http://www.xicidaili.com/wn/", "西刺Https", 1, 30*time.Minute)}
	var ip66Common Obtainer = &Ip66CommonObtainer{NewDefaultWebInfo("http://www.66ip.cn/mo.php?tqsl=", "66ip普通", 2, 2*time.Minute)}
	var ip66Https Obtainer = &Ip66HttpsObtainer{NewDefaultWebInfo("http://www.66ip.cn/nmtq.php?isp=0&anonymoustype=0&area=0&proxytype=2&api=66ip&getnum=", "66ipHttps", 3, 5*time.Minute)}
	var kuaidailiNn Obtainer = &KuaidailiNnObtainer{NewDefaultWebInfo("https://www.kuaidaili.com/free/inha/", "快代理高匿", 4, 6 * time.Hour)}
	var kuaidailiCommon Obtainer = &KuaidailiCommonObtainer{NewDefaultWebInfo("https://www.kuaidaili.com/free/intr/", "快代理普通", 2, 6 * time.Hour)}
	var ip3366O Obtainer = &Ip3366Obtainer{NewDefaultWebInfo("http://www.ip3366.net/?stype=%d&page=%d", "ip3366", 2, 15 * time.Hour)}
	WebObtainers = append(WebObtainers,xiciHttp)
	WebObtainers = append(WebObtainers,xiciHttps)
	WebObtainers = append(WebObtainers,ip66Common)
	WebObtainers = append(WebObtainers,ip66Https)
	WebObtainers = append(WebObtainers,kuaidailiNn)
	WebObtainers = append(WebObtainers,kuaidailiCommon)
	WebObtainers = append(WebObtainers,ip3366O)
}