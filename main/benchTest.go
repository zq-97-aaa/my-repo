package main

import (
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var user1 = &TUserFeature{
	UserId:      "_w71hZTqVqX_rBfgSw9QU",
	Language:    "en",
	PoolType:    "sfw",
	DeviceType:  "mobile",
	UserFollows: []string{"W5kl4av3gh5n5pBxymDtA", "f8G4gYpTTRcx9c2-mEXwm"},
	LoggedIn:    true,
	AB: map[string]string{
		"app_sensitive_image_exp": "2",
	},
	RecPromptMap: map[string]bool{
		"iS24myZtdmdvNvIa4rXWr": true,
		"WoDcPaApNzynQD7-a3_ku": true,
	},
	LastRecItemsRecord: []string{"NdVJpd17haqdXxOYZaFyI", "UEdvgwcvP6RFGTC4TMWGo", "pH1T6ZqO7N3vKNljo5jac"},
	ExposedPromptsCount: map[string]string{
		"7OuOXdcLyTL5QGFYNkT_l": "5",
		"uL7UT0kiu41yz4pdvrF9v": "12",
	},
	LastRecommendedPrompts: map[string]bool{
		"eo3GGoI3RQfOAwQP0vL3G": true,
	},
	LastChattedPrompts: []string{"f3cefc12-caed-4f36-a668-accbefd0a270", "542e3f2d-851a-477a-aa67-9c1ec0986e2a", "437d8716-b680-472e-91e0-1b4b9e1afacd", "4287367f-e04c-459d-8763-073fd6b7322e", "f54d5e18-88c9-4f22-bd9a-ce6101abf771", "10432f5e-870f-428c-a821-d75a29a35a00", "13468fe0-6fc3-47cd-a924-62c44b27b1b6"},
	LastViewedPrompts:  []string{"f3cefc12-caed-4f36-a668-accbefd0a270", "542e3f2d-851a-477a-aa67-9c1ec0986e2a", "437d8716-b680-472e-91e0-1b4b9e1afacd", "13468fe0-6fc3-47cd-a924-62c44b27b1b6"},
	SavedPrompts:       []string{"f3cefc12-caed-4f36-a668-accbefd0a270", "542e3f2d-851a-477a-aa67-9c1ec0986e2a", "437d8716-b680-472e-91e0-1b4b9e1afacd", "13468fe0-6fc3-47cd-a924-62c44b27b1b6"},
	SelectedTags:       []string{"V4", "Character", "Gender", "Male"},
	Model:              "gpt-4",
	IsNewUser:          false,
	ReqScene:           "home",
	SelectGender:       "female",
}
var usage syscall.Rusage
var m runtime.MemStats
var users = GenerateUserFeatures("/Users/qinjiu/Desktop/Prompt3.csv")
var wg sync.WaitGroup

func main() {

	syscall.Getrusage(syscall.RUSAGE_SELF, &usage)

	fmt.Printf("用户态 CPU 时间：%v 秒 %v 微秒\n", usage.Utime.Sec, usage.Utime.Usec)
	fmt.Printf("内核态 CPU 时间：%v 秒 %v 微秒\n", usage.Stime.Sec, usage.Stime.Usec)
	start := time.Now()

	//TestForLoop()
	TestQuery()

	wg.Wait()
	end := time.Now()

	syscall.Getrusage(syscall.RUSAGE_SELF, &usage)
	fmt.Printf("执行过滤后的 CPU 时间：%v 秒 %v 微秒\n", usage.Utime.Sec, usage.Utime.Usec)
	fmt.Printf("程序运行时间：%v\n", end.Sub(start))

	runtime.ReadMemStats(&m)
	fmt.Printf("当前内存分配：%v KB\n", m.Alloc/1024)
	fmt.Printf("累计分配内存：%v KB\n", m.TotalAlloc/1024)
	fmt.Printf("系统占用内存：%v KB\n", m.Sys/1024)
	fmt.Printf("垃圾回收次数：%v 次\n", m.NumGC)
}

func TestForLoop() {
	//LoadPromptDataAsync()

	//user := users[rand.Intn(1000)]
	for _, user := range users[:100] {
		wg.Add(1)
		tmp := user
		go func() {
			defer wg.Done()
			FilterPromptList(tmp, promptDataList, false)
			fmt.Println("--->promptListLength:", len(promptDataList))
		}()

	}

}

func TestQuery() {
	//user := users[rand.Intn(1000)]
	for _, user := range users[:100] {
		wg.Add(1)
		tmp := user
		go func() {
			defer wg.Done()
			_, err := FilterPromptList2(tmp, nil)
			if err != nil {
				return
			}
		}()
	}
}
