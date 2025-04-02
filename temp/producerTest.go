package temp

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"
)

const msg1 = `{"appId":"flowgpt_app","compression":"gzip","data":"H4sIAAAAAAAAA4VVXY+cNhR931+xQn2Md21jwOxbFUVtmq8+VEqrKLKMMYM1YBPbDDOp8t9rGJiB3YmKePI5596L772HL//e3YcnkgepPfOnTkZP95GrzcA6a9rOM8FtGb1as1Q5clJaVAgSCiRBBSCiFKCIEwzipMJlkiYiFnLRlfKghJyFOKcwzWAOCkyCkMIM0DKpQJVkKOY4LyREi7DX6lu/CKWsRIGyElRlHoRcIJCXGQGIxgWsaEbiPF6EvOsWVWtErdggiwXzqpXO87YLMMoIzlOc5yRN8xmvjfNM81Ze5Q/CtIu8MYI3E/a9Bq8/Lsfu5LxsWcP1rue7GV9AYXrt7YkJU07IVfbdaMlMVTnpA4ApheMzgy3fT/TfpNibRdE13FfGtiPw+fpVTlgpNaul2tVjqJziLTKo0tfjJydoQQ5KDp2x/qqiefwcW3QUpUvAcn+5Hz444EzTe2U0EI0Se+et5C0IpGjFP0jrAmWUwAf6AC8ddtKGw/MYTidsPGKVsqEL3vSiZuuGXYnnKnnTy6mPcYYwxUmKyFzmhROulm17vuJeqD9ePS9hmp+f5Ivid+lr9Al9en8YxDFmJ/f+r28k+t/MBGVZHBqMEnjNfLfKH3HvrSp6L932WkRvbdi+P+fZqk0rV+mizoYdM737Ge558XbaByaPXWPsc/Dj3M83N9DZFcaJ2KYcDeIctfn7wx4cft/BTuvP5Nf2jwH9AzZBlD8vzTvuajXwnbp/65xUa47SpTxOE7I+LZXjRSM/qGPHtWxG3Nt+U+GRVU1wrHamsCDxSovFqn45+89T8BtJOUlBSUsMSExiwBMkAU1EmWBCRZWSG2GDl7iqOYWJ6M/xbnDCuMqCu2BVOjRbX8xu+3kBapoQp+NiH/p0WaJe77UZ9JrMpLXGLoaBIcQvwTBYbm73jTzT5I0YyjMIIArvPcyekjy8D3meb7K5ECqs51z14qYAQ5zAGKcAEhzDnKbJTVXIZ19M+uSrdNTckpS95f5sCOEXQG5RdN8WkzngDdyNN3cZpzeTP0cvCL2dJqX2vnNPj49XG3/ccIOXhpqZlVW40ilZFJ3X8e7H1/8AYjPt6R0HAAA=","date":"2025-03-27T09:02:25+00:00","fakeIp":null,"ingest_time":1743066145000,"ip":"37.201.199.106","method":"POST","path":"/collect","platform":"Android","rid":"b49c17dda548c4c788c1bed81d54e4f3","server_ingest_time":1743066145000,"source_type":"http_server","timestamp":"2025-03-27T09:02:25.108828633Z","ua":"amplify-android/1.38.8 (Android 13; Xiaomi 23090RA98G; ar_EG)","uri":"/collect?platform=Android&appId=flowgpt_app&hashCode=296be5d1&event_bundle_sequence_id=26824&upload_timestamp=1743066144388&compression=gzip"}`

func main22() {
	// Kafka 配置
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true            // 需要 Successes 回调
	config.Producer.Return.Errors = true               // 需要 Errors 回调
	config.Producer.RequiredAcks = sarama.WaitForLocal // 只需 Leader 确认

	// Kafka Broker 地址
	brokers := []string{"localhost:9092"}

	// 创建异步生产者
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// 监听中断信号（Ctrl+C）
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// 使用 WaitGroup 确保协程完成
	var wg sync.WaitGroup
	wg.Add(1)

	// 启动协程处理 Successes 和 Errors
	go func() {
		defer wg.Done()
		for {
			select {
			case success := <-producer.Successes():
				log.Printf("Message sent! Partition: %d, Offset: %d", success.Partition, success.Offset)
			case err := <-producer.Errors():
				log.Printf("Failed to send message: %v", err.Err)
			case <-signals:
				return
			}
		}
	}()

	// 发送消息
	topic := "event-tracking"
	for i := 0; i < 10000; i++ {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(msg1),
		}
		producer.Input() <- msg
	}

	// 等待中断信号
	<-signals
	log.Println("Producer shutting down...")
	wg.Wait()
}
