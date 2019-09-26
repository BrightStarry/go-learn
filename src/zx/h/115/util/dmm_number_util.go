package util

import (
	"strings"
	"regexp"
)

/**
番号处理
 */
type NO struct {
	Pre string
	Suf string
}

/*no对象比较*/
func (s *NO) Equals(other *NO)bool{
	if strings.EqualFold(s.Pre,other.Pre) && strings.EqualFold(s.Suf,other.Suf) {
		return true
	}
	return false
}
/*no对象判断是否为空*/
func (s *NO) IsNull() bool{
	if s == nil || s.Pre == ""  || s.Suf == "" {
		return true
	}
	return false
}

const (
	FC2 = "FC2"
	fc2 = "fc2"
	ZERO = "0"
)

/**
提取番号
 */
var getNOReg = regexp.MustCompile("^([A-Za-z\\d]+|[\\d]+)[-_\\s]?([\\d]+)")
var getNORegFC2 = regexp.MustCompile("[\\d]{4,}")
func GetNO(name string)(n NO){
	// 处理fc2番号
	if strings.HasPrefix(name,FC2) || strings.HasPrefix(name,fc2){
		temp := getNORegFC2.FindAllString(name,1)
		// 格式错误，直接返回空对象
		if len(temp) < 1{
			return
		}
		return NO{FC2, strings.TrimLeft(temp[0],ZERO)}
	}

	// 处理其他番号
	temp := getNOReg.FindStringSubmatch(name)
	// 格式错误，直接返回空对象
	if len(temp) < 3 {
		return
	}
	return NO{temp[1], strings.TrimLeft(temp[2],ZERO)}
}