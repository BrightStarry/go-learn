package main

import (
	"zx/ipProxyPool/verify"
	"zx/ipProxyPool/config"
	"log"
	"time"
	"zx/ipProxyPool/obtain"
)

/**
	启动
 */
func main() {
	go verify.StartVerifier()

	go func() {
		for _,v := range obtain.WebObtainers{
			v.IncrementObtain()
		}
	}()

	go func() {
		for v:= range config.VerifiedChan{
			log.Println("入库:",v)
		}
	}()
	time.Sleep(time.Hour)
}
