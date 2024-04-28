package dns

// RecordType 记录类型
type RecordType string

const (
	A     RecordType = "A"
	AAAA  RecordType = "AAAA"
	CNAME RecordType = "CNAME"
	TXT   RecordType = "TXT"
	MX    RecordType = "MX"
	NS    RecordType = "NS"
)

type Record struct {
	ID    string     `json:"id"`
	Type  RecordType `json:"type"`
	Name  string     `json:"name"`
	Value string     `json:"value"`
	TTL   int        `json:"ttl"`
	// 额外字段，各个平台的特殊字段，例如dnspod支持线路、cloudflare支持代理
	Extra map[string]any `json:"extra"`
}

type FieldType string

const (
	FieldTypeText   FieldType = "text"   // 文本
	FieldTypeNumber FieldType = "number" // 数字
	FieldTypeSwitch FieldType = "switch" // 开关
	FieldTypeSelect FieldType = "select" // 下拉框
)

// Extra 额外字段定义描述
type Extra struct {
	Key       string    `json:"key"`
	Title     string    `json:"title"`
	FieldType FieldType `json:"field_type"`
	// 选项，只有FieldTypeSelect时有效
	Options []string `json:"options"`
	Default any      `json:"default"`
}

type Domain struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
