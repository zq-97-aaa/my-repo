package common

const (
	SFW         = "sfw"
	NSFW        = "nsfw"
	parallelNum = 8
	Male        = "Male"
	Female      = "Female"
)

const (
	TopReqType    = "top"
	TrendyReqType = "trendy"
	NewReqType    = "new"
	ForYouReqType = "for_you"
	VideoReqType  = "video"
	PostReqType   = "post"
	StoryReqType  = "story"
)

const (
	MobileDeviceType = "mobile"
	DestopDeviceType = "desktop"
)

const (
	RulePathRecallSourcePreFix = "recall_source_"
	RulePathUserTagPrefix      = "user_tag_"
)

// EXP Key
const (
	AppSpecialBotExpKey  = "app_rs_special_bot_exp"
	AppSensitiveImageExp = "app_sensitive_image_exp"

	ExpGroup0 = "0"
	ExpGroup1 = "1"
	ExpGroup2 = "2"
	ExpGroup3 = "3"
	ExpGroup4 = "4"
)

const (
	InterfaceKey = "tag"
	OrRelation   = "or"
	AndRelation  = "and"
)

// Scene
const (
	FeatureHashTag     = "HashTag"
	AppTagsFilterScene = "AppTagsFilterScene"
	AppExplore         = "_explore"
)

// bool
const (
	True  = "true"
	False = "false"
)

const (
	AdminIllegal    = "ADMIN_ILLEGAL"
	AdminSensitive  = "ADMIN_SENSITIVE"
	AdminNormal     = "ADMIN_NORMAL"
	OpenaiIllegal   = "OPENAI_ILLEGAL"
	OpenaiSensitive = "OPENAI_SENSITIVE"
	MachineIllegal  = "MACHINE_ILLEGAL"
)

var IllegalModerationStatus = []string{AdminIllegal, OpenaiIllegal, MachineIllegal}

// es const
const (
	EsIndexName = "prompt_index"
)

const (
	ENLanguage = "en"
)

// tags
const (
	SensitiveImage2 = "sensitive:SENSITIVE_IMAGE2"
	SensitiveImage  = "sensitive:SENSITIVE_IMAGE"
	NvJianBeiYaoFu  = "sexy:luolu:nvjianbeiyaofu:REVIEW"
	SeQing          = "porn:seqing:seqing:REVIEW"
)

type PromptV2ColdStartStatus int32

const (
	ACTIVE  PromptV2ColdStartStatus = 0
	REMOVED PromptV2ColdStartStatus = 1
)

const Event_Type_Show_Prompt_Card = "show_prompt_card"

// kafka
// TOPIC NAMES
const KAFKA_EVENT_TRACKING_TOPIC = "event-tracking"
