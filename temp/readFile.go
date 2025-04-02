package temp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"strings"
)

// 递归解析 JSON 并提取 botInfo.backData.id
func extractIDs(data interface{}, ids *[]string) {
	switch v := data.(type) {
	case map[string]interface{}:
		if botInfo, ok := v["botInfo"].(map[string]interface{}); ok {
			if backData, ok := botInfo["backData"].(map[string]interface{}); ok {
				if id, ok := backData["id"].(string); ok {
					*ids = append(*ids, id)
				}
			}
		}
		// 递归遍历字典的值
		for _, value := range v {
			extractIDs(value, ids)
		}
	case []interface{}:
		// 遍历列表
		for _, item := range v {
			extractIDs(item, ids)
		}
	}
}

func main2() {
	// JSON 文件路径
	basePath := "/Users/qinjiu/Downloads/"
	files := []string{
		"husband.json",
		"getintrouble.json",
		"nextdoor.json",
		"school.json",
		"warmcaring.json",
		"cold.json",
		"enemy.json",
		"welcomeMessage.json",
		"welcomeMessage(1).json",
	}
	var idList []string

	for _, file := range files {
		// 读取文件内容
		jsonFile, err := os.Open(basePath + file)
		if err != nil {
			fmt.Printf("无法打开文件 %s: %v\n", file, err)
			continue
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		// 解析 JSON
		var data interface{}
		if err := json.Unmarshal(byteValue, &data); err != nil {
			fmt.Printf("解析 JSON 失败: %v\n", err)
			continue
		}

		// 提取 ID
		extractIDs(data, &idList)
	}
	// 输出到文件
	outputFile := "extracted_ids.txt"
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("无法创建输出文件:", err)
		return
	}
	defer file.Close()

	for _, id := range idList {
		file.WriteString(id + "\n")
	}

	fmt.Printf("提取的 ID 已保存到 %s，共 %d 个\n", outputFile, len(idList))

}

func main3() {
	// 输入输出文件路径
	inputFile := "bot_id.txt"
	outputFile := "output2.sql"

	// 固定TagID
	const promptTagID = "82c4635c-c205-4651-9584-3ae199b3395d"

	// 打开输入文件
	inFile, err := os.Open(inputFile)
	if err != nil {
		panic(fmt.Errorf("无法打开输入文件: %v", err))
	}
	defer inFile.Close()

	// 创建输出文件
	outFile, err := os.Create(outputFile)
	if err != nil {
		panic(fmt.Errorf("无法创建输出文件: %v", err))
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(inFile)
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	// 处理每一行
	lineCount := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// 生成UUID
		id := uuid.New().String()

		// 构建SQL语句
		sql := fmt.Sprintf(
			"INSERT INTO public.\"PromptColdStart\" (id, \"promptId\", \"source\", \"createAt\") VALUES ('%s', '%s', '%s');\n",
			id,
			promptTagID,
			strings.ReplaceAll(line, "'", "''"), // 处理单引号
		)

		// 写入文件
		_, err := writer.WriteString(sql)
		if err != nil {
			panic(fmt.Errorf("写入文件失败: %v", err))
		}

		lineCount++
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("读取文件错误: %v", err))
	}

	fmt.Printf("成功生成 %d 条SQL语句到 %s\n", lineCount, outputFile)
}
