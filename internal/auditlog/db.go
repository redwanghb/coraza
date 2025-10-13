package auditlog

import (
	"fmt"
	"reflect"

	"waap/experimental/plugins/plugintypes"
)

// TODO 使用其他方式定位数据库Writer，避免出现单词拼写错误的情况
const (
	POSTGRESQL = "postgresql"
)

const (
	DEFAULTTABLENAME = "eventlog"
	DEFAULTDBNAME    = "insights"
)

// 定义要支持的表结构列表
var tableMap map[string]any = map[string]any{
	"eventlog": &EventLog{},
}

// 通用数据库Writer，可以具体指定到数据库类型的Writer，目前仅支持postgresql，待补充其他数据库类型
type dbWriter struct {
	database plugintypes.AuditLogWriter
}

func (d *dbWriter) Init(c plugintypes.AuditLogConfig) error {
	switch c.DB.DBType() {
	case POSTGRESQL:
		d.database = &pgWriter{}
		return d.database.Init(c)
	default:
		return fmt.Errorf("database type %s unsupported yet", c.DB.DBType())

	}
}

func (d *dbWriter) Write(al plugintypes.AuditLog) error {
	return d.database.Write(al)
}

func (d *dbWriter) Close() error {
	return d.database.Close()
}

var _ plugintypes.AuditLogWriter = (*dbWriter)(nil)

// TODO 定义表结构
type EventLog struct {
	// 发生时间戳
	Timestamp int64 `json:"timestamp"`
	// Event ID
	EventID string `json:"event_id"`
	// 威胁等级
	Severity string `json:"severity"`
	// 威胁类别
	Category string `json:"category"`
	// 威胁指标
	IOC string `json:"ioc"`
	// 规则ID
	RuleID string `json:"rule_id"`
	// 威胁描述
	Message string `json:"message"`
	// 处置方式
	Action string `json:"action"`
	// TODO 攻击结果
	Status string `json:"status"`
	// TODO 载荷信息(存储命中规则的载荷列表)，待确认是否需要base64编码解决特殊字符问题？
	Payloads []string `json:"payloads"`

	// 访问数据信息
	// 访问源，也就是ClientIP和ClientPort
	SrcIp   string `json:"src_ip"`
	SrcPort int    `json:"src_port"`

	// 访问目的，也就是HostIP和HostPort
	DstIp   string `json:"dst_ip"`
	DstPort int    `json:"dst_port"`

	// HTTP请求
	// 请求头信息
	Host           string `json:"host"`
	URI            string `json:"uri"`
	Method         string `json:"method"`
	AppProtocol    string `json:"app_protocol"`
	Protocol       string `json:"protocol"`
	Version        string `json:"version"`
	RequestHeaders string `json:"request_headers"`
	// 请求体信息
	RequestBody string `json:"request_body"`
	Question    string `json:"question"`

	// HTTP应答
	// 响应头
	StatusCode      string `json:"status_code"`
	ResponseHeaders string `json:"response_headers"`
	// 响应体
	ResponseBody string `json:"response_body"`
	Answers      string `json:"answers"`
}

func GetTableFields(tablename string) []string {
	if tablename == "" {
		return nil
	}
	st, ok := tableMap[tablename]
	if !ok {
		return nil
	}

	var fieldList []string

	t := reflect.TypeOf(st).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get("json")
		if tagValue != "" {
			fieldList = append(fieldList, tagValue)
		} else {
			fieldList = append(fieldList, field.Name)
		}
	}
	return fieldList
}
