package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

var promptDataList []*Prompt

func LoadPromptDataAsync() []*Prompt {
	startTime := time.Now()

	// 打开数据文件
	filePath := "/Users/qinjiu/Downloads/prompt_info_with_tag.txt"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 创建处理管道
	jobs := make(chan []byte, 10000)
	results := make(chan *Prompt, 10000)
	done := make(chan struct{})

	// 启动工作池
	var wg sync.WaitGroup
	workerCount := runtime.NumCPU() * 2
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// 启动结果收集器
	go resultCollector2(results, done)

	// 逐行读取文件
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024) // 支持最大10MB的单行数据

	totalLines := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		lineCopy := make([]byte, len(line))
		copy(lineCopy, line)
		jobs <- lineCopy
		totalLines++

		if totalLines%10000 == 0 {
			log.Printf("已读取 %d 行", totalLines)
		}
	}

	log.Printf("close job chan")
	close(jobs)
	wg.Wait()
	log.Printf("close results chan")
	close(results)
	log.Printf("DONE!!!")
	<-done

	if err := scanner.Err(); err != nil {
		log.Printf("文件扫描错误: %v", err)
	}

	log.Printf("导入完成! 总行数: %d, 总耗时: %v", totalLines, time.Since(startTime))
	return promptDataList
}

func LoadPromptData() []*Prompt {
	filePath := "/Users/qinjiu/Downloads/prompt_info_with_tag.txt"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	var batch []*Prompt
	totalLines := 0
	startTime := time.Now()

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		prompt := parseLine(line)
		if prompt != nil {
			batch = append(batch, prompt)
		}

		totalLines++
		if totalLines%10000 == 0 {
			log.Printf("已处理 %d 行", totalLines)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Printf("文件扫描错误: %v", err)
	}

	log.Printf("导入完成! 总行数: %d, 总耗时: %v", totalLines, time.Since(startTime))
	return batch

}

func parseLine(line []byte) *Prompt {
	var temp struct {
		*Prompt
		Tags      []string        `json:"tags"`
		V1Tags    []string        `json:"v1Tags"`
		TagsMap   map[string]bool `json:"tagsMap"`
		UpdatedAt string          `json:"updatedAt"`
	}

	if err := json.Unmarshal(line, &temp); err != nil {
		log.Printf("JSON解析错误: %v", err)
		return nil
	}

	if temp.Prompt == nil {
		return nil
	}

	prompt := temp.Prompt

	if temp.UpdatedAt != "" {
		if t, err := time.Parse(time.RFC3339, temp.UpdatedAt); err == nil {
			prompt.UpdatedAt = t
		}
	}

	if tagsJSON, err := json.Marshal(temp.Tags); err == nil {
		prompt.Tags = string(tagsJSON)
	}
	if v1TagsJSON, err := json.Marshal(temp.V1Tags); err == nil {
		prompt.V1Tags = string(v1TagsJSON)
	}
	if tagsMapJSON, err := json.Marshal(temp.TagsMap); err == nil {
		var tagsMapSlice map[string]bool
		if err := json.Unmarshal(tagsMapJSON, &tagsMapSlice); err == nil {
			prompt.TagsMap = tagsMapSlice
		} else {
			log.Println("JSON 解析失败:", err)
		}
	}

	return prompt
}

func resultCollector2(results <-chan *Prompt, done chan struct{}) {
	var batch []*Prompt
	totalInserted := 0
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case prompt, ok := <-results:
			if !ok {
				log.Println("数据收集完成")
				promptDataList = batch
				close(done)
				return
			}
			batch = append(batch, prompt)
			totalInserted++

		case <-ticker.C:
			log.Printf("当前进度: 已插入 %d 条记录", totalInserted)
		}
	}
}
