package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenerateUserFeatures(csvPath string) []*TUserFeature {
	// 读取CSV文件
	file, _ := os.Open(csvPath)
	defer file.Close()
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	var features []*TUserFeature
	rand.Seed(time.Now().UnixNano())

	// 预先生ID池用于随机选择
	var promptIdPool []string
	var userIdPool []string
	for _, record := range records[1:] { // 跳过标题行
		promptIdPool = append(promptIdPool, record[0])
		userIdPool = append(userIdPool, record[6])
	}

	for _, record := range records[1:] {
		// 基础字段
		user := &TUserFeature{
			UserId:       record[6],  // userId列
			Language:     record[10], // language列
			PoolType:     randomPoolType(),
			DeviceType:   randomDeviceType(),
			LoggedIn:     rand.Intn(2) > 0,
			Model:        "gpt-4",
			IsNewUser:    rand.Intn(5) == 0, // 20%概率是新用户
			ReqScene:     randomScene(),
			SelectGender: randomGender(),
		}

		// 数组类字段（随机长度）
		user.UserFollows = randomSelect(userIdPool, 1, 5)
		user.LastRecItemsRecord = randomSelect(promptIdPool, 3, 7)
		user.LastChattedPrompts = randomUUIDs(3, 7)
		user.LastViewedPrompts = randomSelect(promptIdPool, 3, 7)
		user.SavedPrompts = randomSelect(promptIdPool, 3, 7)

		// 标签处理
		if tags := strings.Trim(record[23], "[]"); tags != "" {
			user.SelectedTags = randomSelect(strings.Split(tags, ","), 2, 6)
		}

		// Map类型字段
		user.AB = map[string]string{
			"app_sensitive_image_exp": strconv.Itoa(rand.Intn(3) + 1),
		}

		user.RecPromptMap = make(map[string]bool)
		for _, id := range randomSelect(promptIdPool, 2, 5) {
			user.RecPromptMap[id] = true
		}

		user.ExposedPromptsCount = make(map[string]string)
		for _, id := range randomSelect(promptIdPool, 1, 3) {
			user.ExposedPromptsCount[id] = strconv.Itoa(rand.Intn(20) + 1)
		}

		user.LastRecommendedPrompts = make(map[string]bool)
		for _, id := range randomSelect(promptIdPool, 1, 2) {
			user.LastRecommendedPrompts[id] = true
		}

		features = append(features, user)
	}
	fmt.Println("--->features  Length =", len(features))
	return features
}

// 辅助函数：随机选择数组元素
func randomSelect(source []string, minLen, maxLen int) []string {
	n := rand.Intn(maxLen-minLen+1) + minLen
	shuffled := append([]string{}, source...)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	if n > len(shuffled) {
		n = len(shuffled)
	}
	return shuffled[:n]
}

// 生成随机UUID
func randomUUIDs(minLen, maxLen int) []string {
	n := rand.Intn(maxLen-minLen+1) + minLen
	var uuids []string
	for i := 0; i < n; i++ {
		uuids = append(uuids, generateUUID())
	}
	return uuids
}

// 其他随机生成函数
func randomPoolType() string {
	return []string{"sfw", "nsfw"}[rand.Intn(2)]
}

func randomDeviceType() string {
	return []string{"mobile", "desktop"}[rand.Intn(2)]
}

func randomScene() string {
	return []string{"home", "search", "profile"}[rand.Intn(3)]
}

func randomGender() string {
	return []string{"male", "female", "other"}[rand.Intn(3)]
}

// UUID生成函数（示例实现）
func generateUUID() string {
	// 实际实现应使用标准UUID库
	return fmt.Sprintf("%x", rand.Uint32())
}
