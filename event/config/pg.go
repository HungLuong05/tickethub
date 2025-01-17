package pg

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"fmt"
	"os"
	"context"
	"time"
	"log"
)

const (
	defaultMaxConns = int32(8)
	defaultMinConns = int32(2)
	defaultMaxConnLifetime = time.Hour
	defaultMaxConnIdleTime = time.Minute * 40
	defaultHealthCheckPeriod = time.Minute * 2
	defaultConnectTimeout = time.Second * 8
)

var (
	DB *PGPool
)

type PGPool struct {
	ConnPool *pgxpool.Pool
}

func Config() (*pgxpool.Config) {
	DATABASE_URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	fmt.Printf("Database URL: %v\n", DATABASE_URL)

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create config: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
 	dbConfig.MinConns = defaultMinConns
 	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
 	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
 	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
 	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	 dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	 }
	
	 dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	 }
	
	 dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	 }
	
	 return dbConfig
}

func (db *PGPool) ConnectDB() (error) {
	connPool, err := pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		log.Fatal("could not create connection to database", err)
		return err
	}

	db.ConnPool = connPool
	db.DropTables()
	db.CreateTables()
	return nil
}

func (db *PGPool) Query(query string, args ... interface{}) (pgx.Rows, error) {
	connection, err := db.ConnPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("could not acquire connection from connection pool", err)
		return nil, err
	}

	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("could not ping database", err)
		return nil, fmt.Errorf("could not ping database: %v", err)
	}

	if len(args) == 0 {
		rows, err := connection.Query(context.Background(), query)
		if err != nil {
			return nil, fmt.Errorf("query execution failed: %v", err)
		}
		return rows, nil
	}

	rows, err := connection.Query(context.Background(), query, args...)
	fmt.Println("Query: ", query, "Args: ", args)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	return rows, nil
}

func (db *PGPool) Exec(query string, args ... interface{}) (error) {
	connection, err := db.ConnPool.Acquire(context.Background())
	if err != nil {
		fmt.Printf("Could not acquire connection from connection pool: %v\n", err)
		log.Fatal("Could not acquire connection from connection pool")
		return err
	}

	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database", err)
		return fmt.Errorf("Could not ping database: %v", err)
	}

	if len(args) == 0 {
		_, err := connection.Exec(context.Background(), query)
		if err != nil {
			return fmt.Errorf("Query execution failed: %v", err)
		}
		return nil
	}

	_, err = connection.Exec(context.Background(), query, args...)
	fmt.Println("Query: ", query, "Args: ", args)
	if err != nil {
		log.Fatal("Query execution failed", err)
		return fmt.Errorf("Query execution failed: %v", err)
	}
	return nil
}

func (db *PGPool) QueryRow(query string, args ... interface{}) (pgx.Row, error) {
	connection, err := db.ConnPool.Acquire(context.Background())
	if err != nil {
		fmt.Printf("Could not acquire connection from connection pool: %v\n", err)
		log.Fatal("Could not acquire connection from connection pool")
		return nil, err
	}

	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database")
		return nil, fmt.Errorf("Could not ping database: %v", err)
	}

	if len(args) == 0 {
		row := connection.QueryRow(context.Background(), query)
		return row, nil
	}

	row := connection.QueryRow(context.Background(), query, args...)
	return row, nil
}

func (db *PGPool) CreateTables() (error) {
	createEventsTable := `
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			name VARCHAR(200) NOT NULL,
			location VARCHAR(200) NOT NULL,
			date VARCHAR(200) NOT NULL,
			ticket INT NOT NULL
		);
	`

	fmt.Printf("Creating table: %v\n", createEventsTable)

	err := db.Exec(createEventsTable)

	if err != nil {
		log.Fatal("Could not create table")
		return err
	}

	return nil
}

func (db *PGPool) DropTables() (error) {
	dropUsersTable := `DROP TABLE IF EXISTS events`

	fmt.Printf("Dropping table: %v\n", dropUsersTable)
	_, err := db.Query(dropUsersTable)
	fmt.Printf("Dropped table: %v\n", dropUsersTable)

	if err != nil {
		log.Fatal("Could not drop table")
		return err
	}

	return nil
}