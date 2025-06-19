package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"time"
)

type Config struct {
	Host     string
	Port     uint32
	Database string
	Username string
	Password string
}

type Client struct {
	conn clickhouse.Conn
}

type AuditLog struct {
	User     string    `json:"user"`
	Action   string    `json:"action"`
	Type     string    `json:"type"`
	Metadata string    `json:"metadata"`
	Service  string    `json:"service"`
	Created  time.Time `json:"created"`
}

func New(cfg Config) (*Client, error) {
	ctx := context.Background()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateAuditLogTable() error {
	ctx := context.Background()
	err := c.conn.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS audit_logs (
            user String,
            action String,
            type String,
            metadata String,
            service String,
            created DateTime
        ) ENGINE = MergeTree()
        ORDER BY created
    `)

	return err
}

func (c *Client) AddAuditLog(auditLog AuditLog) error {
	ctx := context.Background()
	query := `
        INSERT INTO audit_logs (user, action, type, metadata, service, created)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	err := c.conn.Exec(ctx, query,
		auditLog.User,
		auditLog.Action,
		auditLog.Type,
		auditLog.Metadata,
		auditLog.Service,
		auditLog.Created,
	)

	return err
}
