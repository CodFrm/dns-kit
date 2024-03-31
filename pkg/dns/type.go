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
	ID    string
	Type  RecordType
	Name  string
	Value string
	TTL   int
	// 额外字段，各个平台的特殊字段，例如dnspod支持线路、cloudflare支持代理
	Extra map[string]any
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
	Key       string
	Title     string
	FieldType FieldType
	// 选项，只有FieldTypeSelect时有效
	Options []string
	Default any
}

type Domain struct {
	ID     string
	Domain string
}
