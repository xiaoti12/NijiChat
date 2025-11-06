package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/syumai/workers/cloudflare/d1" // 注册D1驱动
)

// D1Client Cloudflare D1 数据库客户端
type D1Client struct {
	db *sql.DB
}

// NewD1Client 创建D1客户端实例
// 在Cloudflare Worker环境中，通过binding名称连接D1数据库
func NewD1Client(bindingName string) (*D1Client, error) {
	// 使用标准的 sql.Open 方法，D1驱动会自动从环境中获取绑定
	db, err := sql.Open("d1", bindingName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to D1 database: %w", err)
	}

	return &D1Client{db: db}, nil
}

// Close 关闭数据库连接
func (c *D1Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Query 执行查询语句
func (c *D1Client) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

// QueryRow 查询单行记录
func (c *D1Client) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return c.db.QueryRowContext(ctx, query, args...)
}

// Exec 执行SQL语句（INSERT, UPDATE, DELETE等）
func (c *D1Client) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

// Transaction 执行事务
func (c *D1Client) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // 重新抛出panic
		}
	}()

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// JSONArrayToString 将字符串数组转换为JSON字符串（用于存储tags等字段）
func JSONArrayToString(arr []string) string {
	if arr == nil {
		return "[]"
	}
	bytes, err := json.Marshal(arr)
	if err != nil {
		return "[]"
	}
	return string(bytes)
}

// StringToJSONArray 将JSON字符串转换为字符串数组
func StringToJSONArray(str string) ([]string, error) {
	if str == "" || str == "[]" {
		return []string{}, nil
	}

	var arr []string
	err := json.Unmarshal([]byte(str), &arr)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON array: %w", err)
	}

	return arr, nil
}

// BoolToInt SQLite布尔值转换（0=false, 1=true）
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// IntToBool SQLite整数转布尔值
func IntToBool(i int) bool {
	return i != 0
}
