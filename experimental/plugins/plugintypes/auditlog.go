// Copyright 2024 Juan Pablo Tosso and the OWASP Coraza contributors
// SPDX-License-Identifier: Apache-2.0

package plugintypes

import (
	"io/fs"
	"strconv"

	"github.com/redwanghb/coraza/v3/internal/collections"
	"github.com/redwanghb/coraza/v3/types"
)

// AuditLog represents the main struct for audit log data
type AuditLog interface {
	Parts() types.AuditLogParts
	Transaction() AuditLogTransaction
	Messages() []AuditLogMessage
}

// AuditLogTransaction contains transaction specific information
type AuditLogTransaction interface {
	Timestamp() string
	UnixTimestamp() int64
	ID() string
	ClientIP() string
	ClientPort() int
	HostIP() string
	HostPort() int
	ServerID() string
	Request() AuditLogTransactionRequest
	HasRequest() bool
	Response() AuditLogTransactionResponse
	HasResponse() bool
	Producer() AuditLogTransactionProducer
	HighestSeverity() string // The highest severity of the matched rules for the transaction
	IsInterrupted() bool     // True if the transaction was interrupted
	// 新增LastRID() string 方法，用于获取最后匹配的规则ID
	LastRID() string
	// 新增LastMessage() string 方法，用于获取最后匹配的规则的描述信息
	LastMessage() string
	// 新增Payload() string方法，用于获取最后匹配规则的数据特征，对应输出告警的Payload
	Payload() string
	// 新增RequestHeader() string方法，用于获取请求头字符串
	RequestHeader() string
	// 新增ResponseHeader() string方法，用于获取响应头字符串
	ResponseHeader() string
	// 新增LLMQuestion() string方法，用于获取问题
	LLMQuestion() string
	// 新增LLMAnswer() string方法，用于获取答案
	LLMAnswer() string
	// 新增Action() string方法，用于获取动作
	Action() string
}

// AuditLogTransactionResponse contains response specific information
type AuditLogTransactionResponse interface {
	Protocol() string
	Status() int
	Headers() map[string][]string
	Body() string
}

// AuditLogTransactionProducer contains producer specific information
// for debugging
type AuditLogTransactionProducer interface {
	Connector() string
	Version() string
	Server() string
	RuleEngine() string
	Stopwatch() string
	Rulesets() []string
}

// AuditLogTransactionRequest contains request specific information
type AuditLogTransactionRequest interface {
	Method() string
	Protocol() string
	URI() string
	HTTPVersion() string
	Headers() map[string][]string
	Body() string
	Files() []AuditLogTransactionRequestFiles
	Args() *collections.ConcatKeyed // A string representation of all request agruments in the format 'k=v,'
	Length() int32                  // The total size of the request in bytes
}

// AuditLogTransactionRequestFiles contains information for the
// uploaded files using multipart forms
type AuditLogTransactionRequestFiles interface {
	Name() string
	Size() int64
	Mime() string
}

// AuditLogMessage contains information about the triggered rules
type AuditLogMessage interface {
	Actionset() string
	Message() string
	Data() AuditLogMessageData
}

// AuditLogMessageData contains information about the triggered rules
// in detail
type AuditLogMessageData interface {
	File() string
	Line() int
	ID() int
	Rev() string
	Msg() string
	Data() string
	Severity() types.RuleSeverity
	Ver() string
	Maturity() int
	Accuracy() int
	Tags() []string
	Raw() string
}

// AuditLogConfig is the configuration of a Writer.
type AuditLogConfig struct {
	// Target is the path to the file to write the raw audit log to.
	Target string

	// FileMode is the mode to use when creating File.
	FileMode fs.FileMode

	// Dir is the path to the directory to write formatted audit logs to.
	Dir string

	// DirMode is the mode to use when creating Dir.
	DirMode fs.FileMode

	// Formatter is the formatter to use when writing formatted audit logs.
	Formatter AuditLogFormatter

	// DB是与数据库有关的信息
	DB *DB
}

// 新增了修改AuditLogConfig.Target的方法
func (a *AuditLogConfig) WriteTarget(s string) *AuditLogConfig {
	a.Target = s
	return a
}

// 新增了修改AuditLogConfig.Formatter的方法
func (a *AuditLogConfig) WriteFormatter(f AuditLogFormatter) *AuditLogConfig {
	a.Formatter = f
	return a
}

// 新增返回DB结构体的方法
func (a *AuditLogConfig) GetDBType() string {
	return a.DB.dbType
}

// 新增DB结构体，用于存储数据库信息，将auditlog写入数据库
type DB struct {
	name     string
	user     string //数据库用户名
	password string //数据库密码
	address  string //数据库地址
	port     int    //数据库服务端口号
	tls      bool   //是否开启tls
	dbType   string //数据库类型
}

func (d *DB) Name() string {
	return d.name
}

func (d *DB) User() string {
	return d.user
}

func (d *DB) Password() string {
	return d.password
}

func (d *DB) Address() string {
	return d.address
}

func (d *DB) Port() string {
	return strconv.Itoa(d.port)
}

func (d *DB) EnableTLS() bool {
	return d.tls
}

func (d *DB) DBType() string {
	return d.dbType
}

func NewDB(opts ...DBOption) *DB {
	db := new(DB)
	for _, opt := range opts {
		opt(db)
	}
	return db
}

type DBOption func(*DB)

func WithName(name string) DBOption {
	return func(d *DB) {
		d.name = name
	}
}

func WithUser(user string) DBOption {
	return func(d *DB) {
		d.user = user
	}
}

func WithPassword(password string) DBOption {
	return func(d *DB) {
		d.password = password
	}
}

func WithAddress(address string) DBOption {
	return func(d *DB) {
		d.address = address
	}
}

func WithPort(port int) DBOption {
	return func(d *DB) {
		d.port = port
	}
}

func WithTLS(tls bool) DBOption {
	return func(d *DB) {
		d.tls = tls
	}
}

func WithDBType(dbtype string) DBOption {
	return func(d *DB) {
		d.dbType = dbtype
	}
}

// AuditLogWriter is the interface for all log writers.
// It receives an auditlog and writes it to the output stream
// An output stream may be a file, a socket, an URL, etc
type AuditLogWriter interface {
	// Init the writer requires previous preparations
	Init(AuditLogConfig) error
	// Write the audit log to the output destination.
	// Using the Formatter is mandatory to generate a "readable" audit log
	// It is not sent as a bslice because some writers may require some Audit
	// metadata.
	Write(AuditLog) error
	// Close the writer if required
	Close() error
}

// AuditLogFormatter serializes an AuditLog into a byte slice.
// It is used to construct the formatted audit log.
type AuditLogFormatter interface {
	Format(AuditLog) ([]byte, error)
	MIME() string
}
