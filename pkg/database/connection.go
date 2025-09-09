// Package database provides database connection management and abstraction layer.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// safeInt32 safely converts an int to int32, capping at max value to prevent overflow
func safeInt32(value int) int32 {
	if value > math.MaxInt32 {
		return math.MaxInt32
	}
	if value < 0 {
		return 0
	}
	return int32(value)
}

// Connection represents a database connection that abstracts different database drivers.
type Connection struct {
	db      *sqlx.DB
	pgxPool *pgxpool.Pool
	driver  types.Driver
	config  types.Config
}

// NewConnection creates a new database connection based on the provided configuration.
func NewConnection(config types.Config) (*Connection, error) {
	conn := &Connection{
		driver: config.Driver,
		config: config,
	}

	switch config.Driver {
	case types.MySQL:
		return conn.connectMySQL()
	case types.PostgreSQL:
		return conn.connectPostgreSQL()
	default:
		return nil, fmt.Errorf("unsupported driver: %s", config.Driver)
	}
}

func (c *Connection) connectMySQL() (*Connection, error) {
	dsn := c.buildMySQLDSN()

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	c.configurePool(db)
	c.db = db
	return c, nil
}

func (c *Connection) connectPostgreSQL() (*Connection, error) {
	dsn := c.buildPostgreSQLDSN()

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL with sqlx: %w", err)
	}

	c.configurePool(db)
	c.db = db

	if err := c.setupPgxPool(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connection) buildMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		c.config.Username,
		c.config.Password,
		c.config.Host,
		c.config.Port,
		c.config.Database,
		c.getCharset(),
		c.getTimezone(),
	)
}

func (c *Connection) buildPostgreSQLDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.config.Host,
		c.config.Port,
		c.config.Username,
		c.config.Password,
		c.config.Database,
		c.getSSLMode(),
		c.getTimezone(),
	)
}

func (c *Connection) setupPgxPool() error {
	pgxDSN := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s&timezone=%s",
		c.config.Username,
		c.config.Password,
		c.config.Host,
		c.config.Port,
		c.config.Database,
		c.getSSLMode(),
		c.getTimezone(),
	)

	poolConfig, err := pgxpool.ParseConfig(pgxDSN)
	if err != nil {
		return fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	maxOpenConns := c.getMaxOpenConns()
	poolConfig.MaxConns = safeInt32(maxOpenConns)

	maxIdleConns := c.getMaxIdleConns()
	poolConfig.MinConns = safeInt32(maxIdleConns)
	poolConfig.MaxConnLifetime = c.config.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = c.config.ConnMaxIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return fmt.Errorf("failed to create pgx pool: %w", err)
	}

	c.pgxPool = pool
	return nil
}

func (c *Connection) configurePool(db *sqlx.DB) {
	db.SetMaxOpenConns(c.getMaxOpenConns())
	db.SetMaxIdleConns(c.getMaxIdleConns())
	db.SetConnMaxLifetime(c.getConnMaxLifetime())
	db.SetConnMaxIdleTime(c.getConnMaxIdleTime())
}

func (c *Connection) getCharset() string {
	if c.config.Charset == "" {
		return "utf8mb4"
	}
	return c.config.Charset
}

func (c *Connection) getSSLMode() string {
	if c.config.SSLMode == "" {
		return "disable"
	}
	return c.config.SSLMode
}

func (c *Connection) getTimezone() string {
	if c.config.Timezone == "" {
		return "UTC"
	}
	return c.config.Timezone
}

func (c *Connection) getMaxOpenConns() int {
	if c.config.MaxOpenConns == 0 {
		return 25
	}
	return c.config.MaxOpenConns
}

func (c *Connection) getMaxIdleConns() int {
	if c.config.MaxIdleConns == 0 {
		return 25
	}
	return c.config.MaxIdleConns
}

func (c *Connection) getConnMaxLifetime() time.Duration {
	if c.config.ConnMaxLifetime == 0 {
		return 5 * time.Minute
	}
	return c.config.ConnMaxLifetime
}

func (c *Connection) getConnMaxIdleTime() time.Duration {
	if c.config.ConnMaxIdleTime == 0 {
		return 5 * time.Minute
	}
	return c.config.ConnMaxIdleTime
}

// QueryContext executes a query that returns rows, typically a SELECT.
func (c *Connection) QueryContext(ctx context.Context, query string, args ...interface{}) (types.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
func (c *Connection) QueryRowContext(ctx context.Context, query string, args ...interface{}) types.Row {
	return c.db.QueryRowContext(ctx, query, args...)
}

// ExecContext executes a query without returning any rows.
func (c *Connection) ExecContext(ctx context.Context, query string, args ...interface{}) (types.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

// Begin starts a transaction with default options.
func (c *Connection) Begin() (types.Tx, error) {
	tx, err := c.db.Beginx()
	if err != nil {
		return nil, err
	}
	return NewTransaction(tx, c.driver), nil
}

// BeginTx starts a transaction with the specified context and options.
func (c *Connection) BeginTx(ctx context.Context, opts *types.TxOptions) (types.Tx, error) {
	sqlOpts := &sql.TxOptions{}
	if opts != nil {
		sqlOpts.Isolation = sql.IsolationLevel(opts.Isolation)
		sqlOpts.ReadOnly = opts.ReadOnly
	}

	tx, err := c.db.BeginTxx(ctx, sqlOpts)
	if err != nil {
		return nil, err
	}
	return NewTransaction(tx, c.driver), nil
}

// Driver returns the database driver type.
func (c *Connection) Driver() types.Driver {
	return c.driver
}

// Close closes the database connection and releases resources.
func (c *Connection) Close() error {
	var err error
	if c.db != nil {
		err = c.db.Close()
	}
	if c.pgxPool != nil {
		c.pgxPool.Close()
	}
	return err
}

// Ping verifies a connection to the database is still alive.
func (c *Connection) Ping() error {
	if c.db != nil {
		return c.db.Ping()
	}
	if c.pgxPool != nil {
		return c.pgxPool.Ping(context.Background())
	}
	return fmt.Errorf("no connection available")
}

// Stats returns database statistics including connection pool information.
func (c *Connection) Stats() types.DBStats {
	if c.db != nil {
		stats := c.db.Stats()
		return types.DBStats{
			OpenConnections: stats.OpenConnections,
			InUse:           stats.InUse,
			Idle:            stats.Idle,
		}
	}
	if c.pgxPool != nil {
		stat := c.pgxPool.Stat()
		return types.DBStats{
			OpenConnections: int(stat.TotalConns()),
			InUse:           int(stat.AcquiredConns()),
			Idle:            int(stat.IdleConns()),
		}
	}
	return types.DBStats{}
}

// PgxPool returns the underlying pgx connection pool for PostgreSQL connections.
func (c *Connection) PgxPool() *pgxpool.Pool {
	return c.pgxPool
}

// SqlxDB returns the underlying sqlx database connection.
func (c *Connection) SqlxDB() *sqlx.DB {
	return c.db
}
