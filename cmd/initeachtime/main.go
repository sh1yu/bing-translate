package main

import (
	"flag"
	"fmt"
	"github.com/psy-core/bing-translate/trans"
	"go.uber.org/atomic"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func main() {

	concurrent := flag.Int("c", 1, "concurrent number")
	sleepTime := flag.Int("sleep", 1000, "每次请求后睡眠毫秒数")
	timeout := flag.Int("timeout", 5000, "http请求超时毫秒数")
	flag.Parse()

	count := atomic.NewInt64(0)
	errCount := atomic.NewInt64(0)

	go func() {
		for {
			if count.Load() == 0 {
				time.Sleep(2 * time.Second)
				continue
			}
			perccentage := float64(errCount.Load()) / float64(count.Load())
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"),
				"totalCount:", count.Load(),
				"errCount:", errCount.Load(),
				"err percentage:", strconv.FormatFloat(perccentage, 'f', 2, 32))
			time.Sleep(10 * time.Second)
		}
	}()

	for i := 0; i < *concurrent; i++ {
		go func() {
			for {
				t, err := trans.NewBingTranslator(*timeout)
				if err != nil {
					fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "trans init error:", err)
					continue
				}
				count.Inc()
				result, err := t.Translate("zh-CHS", "en", "今天中午吃啥?")
				if err != nil {
					fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "trans err:", err)
					errCount.Inc()
					continue
				}

				if result != "What do you have for lunch today?" {
					fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "trans err: result is not correct:", result)
					errCount.Inc()
					continue
				}

				content, err := t.TSpeak("zh-CHS", "今天中午吃啥？")
				if err != nil {
					fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "speak err:", err)
					errCount.Inc()
					continue
				}

				ioutil.WriteFile("a.mp3", content, 0644)

				time.Sleep(time.Duration(*sleepTime) * time.Millisecond)
				os.Remove("a.mp3")
			}
		}()
	}

	select {}
}
