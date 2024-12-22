package auditlog

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redwanghb/coraza/v3/experimental/plugins/plugintypes"
)

// TODO 待设计结构体内容
type pgWriter struct {
	Pool      *pgxpool.Pool
	DBName    string
	TableName string
	// 表的列名
	Fields []string
	// 写入表的sql字符串
	// sql: "INSERT INTO eventlog (timestamp, event_id, serverity) VALUES ($1, $2, $2)"
	sql string
}

// 初始化pgWriter
func (pg *pgWriter) Init(c plugintypes.AuditLogConfig) error {
	if c.GetDBType() != "postgresql" {
		return fmt.Errorf("unsupported database type %s", c.GetDBType())
	}
	return pg.init(c)
}

func (pg *pgWriter) init(c plugintypes.AuditLogConfig) error {
	pg.DBName = DBNAME
	pg.TableName = TABLENAME
	pool, err := pg.initPool(c)
	if err != nil {
		return err
	}
	pg.Pool = pool

	pg.Fields = GetTableFields(pg.TableName)
	if pg.Fields == nil {
		return fmt.Errorf("there is no fields get from the table %s", pg.TableName)
	}
	tablefields, fieldMarks := func() (string, string) {
		var tablefields, fieldMarks string
		for i := 0; i < len(pg.Fields); i++ {
			if i == 0 {
				tablefields = pg.Fields[0]
				fieldMarks = fmt.Sprintf("$%s", strconv.Itoa(i+1))
			} else {
				tablefields = fmt.Sprintf("%s, %s", tablefields, pg.Fields[i])
				fieldMarks = fmt.Sprintf("%s, $%s", fieldMarks, strconv.Itoa(i+1))
			}
		}
		return tablefields, fieldMarks
	}()

	pg.sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", pg.TableName, tablefields, fieldMarks)

	return nil
}

func (pg *pgWriter) initPool(c plugintypes.AuditLogConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.DB.User(), c.DB.Password(), c.DB.Address(), c.DB.Port(), pg.DBName)
	if !c.DB.EnableTLS() {
		dsn = dsn + "?sslmode=disable"
	}
	pgconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pgconfig.MaxConns = 100
	pgconfig.MinConns = 2
	pgconfig.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), pgconfig)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// TODO 待设计写入方法
func (pg *pgWriter) Write(al plugintypes.AuditLog) error {
	var (
		uri         string
		host        string
		method      string
		version     string
		status_code string
		req_body    string
		res_body    string
		req_header  string
		res_header  string
	)
	if al.Transaction().HasRequest() {
		uri = al.Transaction().Request().URI()
		method = al.Transaction().Request().Method()
		host = al.Transaction().ServerID()
		version = al.Transaction().Request().Protocol()
		req_body = al.Transaction().Request().Body()
		req_header = al.Transaction().RequestHeader()
	}

	if al.Transaction().HasResponse() {
		status_code = strconv.Itoa(al.Transaction().Response().Status())
		res_body = al.Transaction().Response().Body()
		res_header = al.Transaction().ResponseHeader()
	}
	_, err := pg.Pool.Exec(
		context.Background(),
		pg.sql,
		al.Transaction().UnixTimestamp()/1000000,
		uuid.New().String(),
		al.Transaction().HighestSeverity(),
		// TODO Category，暂时为空，等后续处理，当前规则结构中不支持Category定义，考虑使用Category与规则ID匹配的方式处理
		"",
		// TODO IOC，目前先使用Message的内容作为IOC，待后续完善
		al.Transaction().LastMessage(),
		al.Transaction().LastRID(),
		al.Transaction().LastMessage(),
		al.Transaction().Action(),
		// TODO Status，暂时固定为attempt，后续再考虑如何处理
		"attempt",
		al.Transaction().Payload(),
		al.Transaction().ClientIP(),
		al.Transaction().ClientPort(),
		al.Transaction().HostIP(),
		al.Transaction().HostPort(),
		host,
		uri,
		method,
		"http",
		"tcp",
		version,
		req_header,
		req_body,
		al.Transaction().LLMQuestion(),
		status_code,
		res_header,
		res_body,
		al.Transaction().LLMAnswer(),
	)
	if err != nil {
		return err
	}
	return nil
}

// 主要是关闭pgx的Pool
func (pg *pgWriter) Close() error {
	pg.Pool.Close()
	return nil
}

var _ plugintypes.AuditLogWriter = (*pgWriter)(nil)
