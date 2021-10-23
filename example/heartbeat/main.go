package main

import (
	bg "github.com/iyear/biligo"
	"log"
	"time"
)

func main() {
	b, err := bg.NewBiliClient(&bg.BiliSetting{
		Auth: &bg.CookieAuth{
			DedeUserID:      "",
			DedeUserIDCkMd5: "",
			SESSDATA:        "",
			BiliJCT:         "",
		},
		DebugMode: true,
	})
	if err != nil {
		log.Fatal("failed to make new bili client; error: ", err)
		return
	}

	sig := make(chan struct{})

	go func(client *bg.BiliClient, sig chan struct{}) {
		start := time.Now()
		// 每十秒上报一次
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				err := client.VideoHeartBeat(242531611, 173439442, int64(time.Since(start).Seconds()))
				if err != nil {
					log.Println("failed to send heartbeat; error:", err)
					continue
				}
				log.Println("heartbeat...")
			case <-sig:
				ticker.Stop()
				log.Println("stop...")
				return
			}
		}
	}(b, sig)

	// 模拟一分钟
	m := time.NewTimer(time.Minute)
	<-m.C
	m.Stop()
	// 发出信号终止心跳包
	close(sig)
}
