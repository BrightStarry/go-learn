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
	var xiciHttp Obtainer = &XiciHttpObtainer{NewDefaultWebInfo("http://www.xicidaili.com/nt/", "西刺Http", 0, 30*time.Minute)}
	var xiciHttps Obtainer = &XiciHttpsObtainer{NewDefaultWebInfo("http://www.xicidaili.com/wn/", "西刺Https", 1, 35*time.Minute)}
	var ip66Common Obtainer = &Ip66CommonObtainer{NewDefaultWebInfo("http://www.66ip.cn/mo.php?tqsl=", "66ip普通", 2, 5*time.Minute)}
	var ip66Https Obtainer = &Ip66HttpsObtainer{NewDefaultWebInfo("http://www.66ip.cn/nmtq.php?isp=0&anonymoustype=0&area=0&proxytype=2&api=66ip&getnum=", "66ipHttps", 3, 10*time.Minute)}
	var kuaidailiNn Obtainer = &KuaidailiNnObtainer{NewDefaultWebInfo("https://www.kuaidaili.com/free/inha/", "快代理高匿", 4, 6 * time.Hour)}
	var kuaidailiCommon Obtainer = &KuaidailiCommonObtainer{NewDefaultWebInfo("https://www.kuaidaili.com/free/intr/", "快代理普通", 2, 6 * time.Hour)}
	WebObtainers[0] = xiciHttp
	WebObtainers[1] = xiciHttps
	WebObtainers[2] = ip66Common
	WebObtainers[3] = ip66Https
	WebObtainers[4] = kuaidailiNn
	WebObtainers[5] = kuaidailiCommon
}