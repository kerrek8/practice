package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"practic/internal/logger/sl"
	"practic/internal/models"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
	CreateUser(name, login string, password []byte) (uid int64, err error)
	User(login string) (models.UserDB, error)
	//CreateListing()
	//UpdateListing()
	//DeleteListing()
	//GetListings()
}

type service struct {
	log *slog.Logger
	db  *sql.DB
}

var (
	dburl      = os.Getenv("BLUEPRINT_DB_URL")
	dbInstance *service
)

func New(log *slog.Logger) Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("sqlite", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Error("error connection", sl.Err(err))
	}

	dbInstance = &service{
		log: log,
		db:  db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		s.log.Error("db down: %v", sl.Err(err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	s.log.Info("Disconnected from database: %s", dburl)
	return s.db.Close()
}

func (s *service) CreateUser(name, login string, password []byte) (uid int64, err error) {
	const op = "sqlite.database.CreateUser"
	const query = `
		INSERT INTO users (username, password, name) VALUES (?, ?, ?) RETURNING id;
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := stmt.Exec(login, password, name)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := resp.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *service) User(login string) (models.UserDB, error) {
	const op = "sqlite.database.User"
	const query = `
		SELECT id, username, password, name
		FROM users
		WHERE username = ?
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return models.UserDB{}, fmt.Errorf("%s: %w", op, err)
	}

	var user models.UserDB

	err = stmt.QueryRow(login).Scan(&user.ID, &user.Login, &user.Password, &user.Name)
	if err != nil {
		return models.UserDB{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil

}

//func (s *service) CreateListing(name, type_l, description, status, city string, price float64, user_id int64) (uid int64, err error) {
//	const op = "sqlite.database.CreateListing"
//	const query = `
//		INSERT INTO listings (name, type, description, status, price, city, user_id) VALUES ($1, $2, $3);
//	`
//	err = s.db.
//}
