var blackList = [
    "baidu.com",
]

var blankList = [

]


var proxy = "SOCKS5 127.0.0.1:8081"
var direct = "DIRECT"

function FindProxyForURL(url, host) {
  for (var i = 0; i < blackList.length; i++) {
      if(dnsDomainIs(host, blackList[i])){
        for (var j = 0; j < blankList.length; i++) {
            if(dnsDomainIs(host, blankList[j])){
                return direct;
            }
        }
        return proxy;
      }
  }
  return direct;
}
