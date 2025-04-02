package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"helloworld/common"
	"strconv"
	"time"
)

type TUserFeature struct {
	UserId                 string
	Language               string
	PoolType               string
	DeviceType             string
	UserFollows            []string
	LoggedIn               bool
	AB                     map[string]string
	RecPromptMap           map[string]bool
	LastRecItemsRecord     []string
	ExposedPromptsCount    map[string]string
	LastRecommendedPrompts map[string]bool
	LastChattedPrompts     []string
	LastViewedPrompts      []string
	SavedPrompts           []string
	SelectedTags           []string
	Model                  string
	IsNewUser              bool
	ReqScene               string
	SelectGender           string
}
type PRecall struct {
	TypeName          string // redis key prefix, as recall source marker
	RealName          string // name of recall
	NeedSize          int64  // how many items this recall needs
	RealSize          int64  // how many items this recall actually got (excludes the added ones)
	AddSize           int64  // how many items were added to this recall
	OrderLevel        float64
	DbType            int64
	NoAdd             bool // whether to add more items to this recall if the number is not enough
	Tags              [][]string
	FeatureTags       []string
	ReqType           string
	UserType          string // new or old, logged in or not
	CutoffDate        string
	DeviceType        string
	AdditionalFilters map[string][]string
	TagVersion        string
	TagsRelation      string
	Scene             string
}

var logger2 = log.GetLogger()

func FilterPromptList(user *TUserFeature, allPromptList []*Prompt, skipPreFilters bool) []*Prompt {
	fmt.Println("Start FilterPromptList...")
	//var cutoffDateTime time.Time
	//var err error
	//if len(pRecall.CutoffDate) > 0 {
	//	cutoffDateTime, err = time.Parse(time.RFC3339, pRecall.CutoffDate)
	//	if err != nil {
	//		logger2.Log(log.LevelError, "Error parsing cutoff date", zap.Error(err))
	//	}
	//}

	start := time.Now()
	defer func() {
		logger2.Log(log.LevelInfo, "FilterPromptList complete, costTime", time.Since(start).String())
	}()

	// Initialize filtered prompt list
	var filteredPromptList []*Prompt

	// Start filtering in a single loop
	for _, prompt := range allPromptList {
		if !skipPreFilters {
			if !RecallPoolTypeFilter(user.PoolType, prompt.Nsfw) {
				continue
			}

			// Step 2: Filter by exposure limits
			if !RecallExposureLimit(user, prompt) {
				continue
			}

			// Step 3: Filter by language
			if user.Language != "" && prompt.Language != user.Language {
				continue
			}
		}

		// Step 4: Filter by tags and feature tags
		//if !applyTagMatching(pRecall, prompt) {
		//	continue
		//}

		// todo: for sensitive image filter exp;  will full or delete later
		if user.DeviceType == common.MobileDeviceType && user.ReqScene == common.AppExplore {
			if user.AB[common.AppSensitiveImageExp] == common.ExpGroup1 && (prompt.TagsMap[common.SensitiveImage] || prompt.TagsMap[common.NvJianBeiYaoFu] || prompt.TagsMap[common.SeQing]) {
				continue
			} else if user.AB[common.AppSensitiveImageExp] == common.ExpGroup2 && prompt.TagsMap[common.SensitiveImage2] {
				continue
			}
		}

		// Step 5: Filter by cutoff date
		//if len(pRecall.CutoffDate) > 0 && prompt.CreatedAt.Before(cutoffDateTime) {
		//	continue
		//}

		// Step 6: Filter by additional filters (category, author, etc.)
		//if !RecallAdditionalFilters(pRecall.AdditionalFilters, prompt) {
		//	continue
		//}

		// Step 7: Filter by device type and tag criteria
		if user.DeviceType == common.MobileDeviceType {
			// Filter by thumbnail URL
			if isDefaultThumbnail(prompt.ThumbnailURL) {
				continue
			}
		}

		// If passed all filters, append to result list
		filteredPromptList = append(filteredPromptList, prompt)

	}

	logger2.Log(log.LevelInfo, "filtered  prompts:", len(filteredPromptList))
	// Return filtered prompt list
	return filteredPromptList
}

func RecallExposureLimit(user *TUserFeature, p *Prompt) bool {
	if user.LastRecommendedPrompts[p.Id] {
		return false
	}
	countStr := user.ExposedPromptsCount[p.Id]
	if countStr != "" {
		count, err := strconv.Atoi(countStr)
		if err != nil {
			logger2.Log(log.LevelError, "Error converting exposed prompts count to int", zap.Error(err))
		}
		if count > 15 {
			return false
		}
	}
	return true
}

var (
	AuthorRecallFilter   = "author"
	CategoryRecallFilter = "category"
)

func RecallAdditionalFilters(filters map[string][]string, p *Prompt) bool {
	for filterK, filterV := range filters {
		switch filterK {
		case CategoryRecallFilter:
			if len(filterV) > 0 {
				categoryIdInt, err := strconv.ParseInt(filterV[0], 10, 64)
				if err != nil || p.CategoryId != categoryIdInt {
					return false
				}
			}
		case AuthorRecallFilter:
			if len(filterV) > 0 && p.UserId != filterV[0] {
				return false
			}
		}
	}
	return true
}

func RecallPoolTypeFilter(poolType string, isNsfw bool) bool {
	switch poolType {
	case "all":
		return true
	case "sfw":
		return !isNsfw
	case "nsfw":
		return isNsfw
	default:
		return true
	}
}

func applyTagMatching(pRecall *PRecall, p *Prompt) bool {
	if len(pRecall.Tags) == 0 {
		return true
	}
	tagVersion := pRecall.TagVersion
	var v1Tags []string
	if err := json.Unmarshal([]byte(p.V1Tags), &v1Tags); err != nil {
		logger2.Log(log.LevelError, "Error parsing v1 tags", zap.Error(err))
	}
	for _, tags := range pRecall.Tags {
		if tagVersion == "v1" && !tagMatchOr(tags, v1Tags) {
			return false
		} else if tagVersion == "v2" && pRecall.TagsRelation == "or" && !tagMatchOrByMap(tags, p.TagsMap) {
			return false
		} else if tagVersion == "v2" && !tagMatchOrByMap(tags, p.TagsMap) {
			return false
		}
	}
	return true
}

func tagMatchOrByMap(tags []string, promptTagsMap map[string]bool) bool {
	// check if at least one tag in tags is in promptTags
	for _, tag := range tags {
		if promptTagsMap[tag] {
			return true
		}
	}
	return false
}

func tagMatchOr(tags []string, promptTags []string) bool {
	// check if at least one tag in tags is in promptTags
	for _, tag := range tags {
		for _, promptTag := range promptTags {
			if tag == promptTag {
				return true
			}
		}
	}
	return false
}

func isDefaultThumbnail(url string) bool {
	switch url {
	case "https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i10.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i11.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i12.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i7.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i8.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i9.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f10.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f11.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f12.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f13.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f14.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f15.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f16.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f17.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f18.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f19.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f1.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f20.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f21.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f22.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f23.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f24.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f25.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f26.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f27.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f28.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f29.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f2.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f30.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f31.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f32.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f33.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f34.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f35.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f36.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f37.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f38.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f3.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f4.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f5.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f6.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f7.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f8.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f9.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m10.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m11.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m12.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m13.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m14.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m15.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m16.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m17.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m18.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m19.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m1.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m20.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m21.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m22.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m23.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m24.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m25.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m26.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m27.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m28.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m29.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m2.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m30.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m31.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m32.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m33.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m34.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m35.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m36.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m37.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m38.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m39.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m3.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m40.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m41.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m4.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m5.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m6.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m7.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m8.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m9.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f10.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f11.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f12.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f13.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f14.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f15.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f16.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f17.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f18.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f19.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f1.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f2.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f3.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f4.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f5.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f6.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f7.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f8.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f9.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m10.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m11.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m12.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m13.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m14.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m15.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m16.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m1.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m2.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m3.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m4.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m5.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m6.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m7.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m8.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m9.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i10.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i11.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i12.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i13.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i14.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i15.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i16.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i17.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i18.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i19.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i20.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i21.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i22.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i7.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i8.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i9.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_10.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_7.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_8.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_9.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_7.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i10.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i7.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i8.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i9.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i10.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i11.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i12.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i13.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i14.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i15.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i16.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i17.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i18.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i19.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i20.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i21.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i22.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i23.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i24.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i7.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i8.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i9.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i10.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i11.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i12.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i13.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i14.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i15.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i16.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i17.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i18.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i19.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i7.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i8.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i9.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_10.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_11.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_12.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_13.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_14.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_7.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_8.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_9.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n10.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n11.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n12.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n13.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n14.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n15.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n16.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n17.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n18.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n19.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n1.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n20.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n21.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n22.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n23.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n24.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n25.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n26.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n27.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n28.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n29.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n2.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n30.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n31.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n32.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n33.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n34.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n35.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n36.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n37.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n38.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n39.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n3.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n40.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n41.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n42.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n43.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n44.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n45.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n46.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n4.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n5.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n6.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n7.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n8.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n9.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f10.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f11.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f12.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f13.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f14.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f15.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f16.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f17.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f18.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f19.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f1.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f20.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f21.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f22.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f23.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f24.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f25.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f26.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f27.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f28.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f29.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f2.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f30.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f31.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f32.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f33.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f34.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f35.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f36.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f37.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f3.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f4.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f5.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f6.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f7.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f8.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f9.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m10.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m11.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m12.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m13.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m14.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m15.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m16.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m17.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m18.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m19.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m1.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m20.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m21.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m22.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m23.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m24.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m25.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m26.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m27.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m28.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m29.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m2.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m30.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m31.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m32.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m33.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m34.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m35.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m36.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m37.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m38.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m39.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m3.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m40.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m4.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m5.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m6.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m7.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m8.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m9.jpg",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_7.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_8.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_9.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_10.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_11.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_1.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_2.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_3.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_4.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_5.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_6.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_7.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_8.png",
		"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_9.png":
		return true
	default:
		return false
	}
}
