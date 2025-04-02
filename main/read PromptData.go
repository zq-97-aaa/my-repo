package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

var db = InitDB()

type CustomNamingStrategy struct {
	schema.NamingStrategy
}

func (ns CustomNamingStrategy) ColumnName(table, column string) string {
	return strings.ToLower(column[:1]) + column[1:]
}
func (Prompt) TableName() string {
	return "Prompt3"
}

// Prompt 结构体对应 Protobuf 定义
type Prompt struct {
	Id                string          `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt         time.Time       `gorm:"index" json:"createdAt,omitempty"`
	UpdatedAt         time.Time       `json:"updatedAt,omitempty"`
	Title             string          `gorm:"type:text" json:"title,omitempty"`
	Description       string          `gorm:"type:text" json:"description,omitempty"`
	Saves             int64           `gorm:"index" json:"saves,omitempty"`
	UserId            string          `gorm:"index" json:"userId,omitempty"`
	Live              bool            `json:"live,omitempty"`
	Popularity        float64         `gorm:"index" json:"popularity,omitempty"`
	Views             int64           `json:"views,omitempty"`
	Language          string          `json:"language,omitempty"`
	Uses              int64           `gorm:"index" json:"uses,omitempty"`
	Ranking           float64         `json:"ranking,omitempty"`
	Type              string          `json:"type,omitempty"`
	CategoryId        int64           `json:"categoryId,omitempty"`
	SubCategoryId     int64           `json:"subCategoryId,omitempty"`
	UserNsfw          bool            `json:"userNsfw,omitempty"`
	Nsfw              bool            `json:"nsfw,omitempty"`
	NormalizedUses    float64         `json:"normalizedUses,omitempty"`
	TrendingScore     float64         `gorm:"index" json:"trendingScore,omitempty"`
	Fop               int64           `json:"fop,omitempty"`
	Cup               int64           `json:"cup,omitempty"`
	FeatureTags       string          `gorm:"type:text" json:"featureTags,omitempty"`
	Tags              string          `gorm:"type:text" json:"tags,omitempty"`
	Upvotes           int64           `json:"upvotes,omitempty"`
	Downvotes         int64           `json:"downvotes,omitempty"`
	Shares            int64           `json:"shares,omitempty"`
	NewTopScore       float64         `json:"newTopScore,omitempty"`
	NewTrendingScore  float64         `json:"newTrendingScore,omitempty"`
	Comments          int64           `json:"comments,omitempty"`
	Impressions       int64           `json:"impressions,omitempty"`
	PublicChats       int64           `json:"publicChats,omitempty"`
	ThumbnailURL      string          `gorm:"type:text" json:"thumbnailURL,omitempty"`
	URI               string          `json:"uri,omitempty"`
	V1Tags            string          `gorm:"type:text" json:"v1Tags,omitempty"`
	IsMobileRecommend bool            `json:"isMobileRecommend,omitempty"`
	IsWebRecommend    bool            `json:"isWebRecommend,omitempty"`
	AppQualityScore   float64         `json:"appQualityScore,omitempty"`
	TagsMap           map[string]bool `gorm:"-" json:"tagsMap,omitempty"`
}

func InitDB() *gorm.DB {
	// 初始化数据库连接，返回 *gorm.DB
	dataBase, err := gorm.Open(sqlite.Open("/Users/qinjiu/Documents/testDB"), &gorm.Config{
		NamingStrategy: CustomNamingStrategy{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Error,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Fatal(err)
	}
	return dataBase // 请替换为实际数据库初始化代码
}
func main5() {
	startTime := time.Now()
	log.Println("开始导入流程...")
	// 修改批量大小
	const optimalBatchSize = 500 // 根据SQLite限制调整
	// 初始化数据库连接
	db = InitDB()
	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接失败: %v", err)
	}
	sqlDB.SetMaxOpenConns(runtime.NumCPU() * 2)
	sqlDB.SetMaxIdleConns(runtime.NumCPU())
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移表结构
	if err := db.AutoMigrate(&Prompt{}); err != nil {
		log.Fatalf("表结构迁移失败: %v", err)
	}

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
	go resultCollector(db, results, done, optimalBatchSize)

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

	close(jobs)
	wg.Wait()
	close(results)
	<-done

	if err := scanner.Err(); err != nil {
		log.Printf("文件扫描错误: %v", err)
	}

	log.Printf("导入完成! 总行数: %d, 总耗时: %v", totalLines, time.Since(startTime))
}

// worker处理JSON解析
func worker(jobs <-chan []byte, results chan<- *Prompt, wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range jobs {
		var temp struct {
			*Prompt
			//FeatureTags string          `json:"featureTags"`
			Tags      []string        `json:"tags"`
			V1Tags    []string        `json:"v1Tags"`
			TagsMap   map[string]bool `json:"tagsMap"`
			UpdatedAt string          `json:"updatedAt"`
		}

		if err := json.Unmarshal(line, &temp); err != nil {
			log.Printf("JSON解析错误: %v", err)
			continue
		}

		prompt := temp.Prompt
		if prompt == nil {
			continue
		}

		// 处理时间字段
		if temp.UpdatedAt != "" {
			if t, err := time.Parse(time.RFC3339, temp.UpdatedAt); err == nil {
				prompt.UpdatedAt = t
			}
		}

		// 序列化数组和map字段
		if tagsJSON, err := json.Marshal(temp.Tags); err == nil {
			prompt.Tags = string(tagsJSON)
		}
		if v1TagsJSON, err := json.Marshal(temp.V1Tags); err == nil {
			prompt.V1Tags = string(v1TagsJSON)
		}
		if featureTagsJSON, err := json.Marshal(temp.FeatureTags); err == nil {
			prompt.FeatureTags = string(featureTagsJSON)
		}
		if tagsMapJSON, err := json.Marshal(temp.TagsMap); err == nil {
			var tagsMapSlice map[string]bool
			if err := json.Unmarshal(tagsMapJSON, &tagsMapSlice); err == nil {
				prompt.TagsMap = tagsMapSlice
			} else {
				fmt.Println("JSON 解析失败:", err)
			}
		}

		results <- prompt
	}
}

// 结果收集器批量插入数据库
func resultCollector(db *gorm.DB, results <-chan *Prompt, done chan<- struct{}, batchSize int) {
	batch := make([]*Prompt, 0, batchSize)
	totalInserted := 0

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case prompt, ok := <-results:
			if !ok {
				// 插入剩余记录
				if len(batch) > 0 {
					if err := insertBatch(db, batch); err != nil {
						log.Printf("最后批次插入失败: %v", err)
					} else {
						totalInserted += len(batch)
					}
				}
				log.Printf("总共插入 %d 条记录", totalInserted)
				done <- struct{}{}
				return
			}

			batch = append(batch, prompt)
			if len(batch) >= batchSize {
				if err := insertBatch(db, batch); err != nil {
					log.Printf("批量插入失败: %v", err)
				} else {
					totalInserted += len(batch)
				}
				batch = batch[:0] // 清空批次
			}

		case <-ticker.C:
			log.Printf("当前进度: 已插入 %d 条记录", totalInserted)
		}
	}
}

func insertBatch(db *gorm.DB, batch []*Prompt) error {
	const maxVariablesPerBatch = 500                                         // 保守估计的每批最大记录数
	effectiveBatchSize := maxVariablesPerBatch / (getPromptFieldCount() + 1) // 字段数+1为安全边际

	for start := 0; start < len(batch); start += effectiveBatchSize {
		end := start + effectiveBatchSize
		if end > len(batch) {
			end = len(batch)
		}

		tx := db.Begin()
		if tx.Error != nil {
			return tx.Error
		}

		if err := tx.Create(batch[start:end]).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit().Error; err != nil {
			return err
		}
	}
	return nil
}

// 计算Prompt结构体的字段数
func getPromptFieldCount() int {
	// 根据您的Prompt结构体实际字段数返回
	// 这里需要手动更新为您的实际字段数
	return 39 // 根据您提供的Protobuf结构体
}
