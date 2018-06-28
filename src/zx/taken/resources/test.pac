var proxy = "SOCKS5 127.0.0.1:8081"

function FindProxyForURL(url, host) {
  return proxy;
}
