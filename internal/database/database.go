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
	CreateListing(name, type_l, description, status, city string, price int64, user_id int64) (uid int64, err error)
	GetListings(userID int64, offset int64, filter string) ([]models.ListingDB, error)
	GetCities(userID int64) ([]string, error)
	UpdateListing(name, typel, description, status, city string, price int64, id int64) error
	DeleteListing(id int64) error
	GetAnalytics(userID int64) (map[string]any, error)
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

func (s *service) CreateListing(name, type_l, description, status, city string, price int64, user_id int64) (uid int64, err error) {
	const op = "sqlite.database.CreateListing"
	const query = `
		INSERT INTO listings (name, type, description, status, price, city, user_id) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id;
	`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	resp, err := stmt.Exec(name, type_l, description, status, price, city, user_id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := resp.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *service) GetListings(userID int64, offset int64, filter string) ([]models.ListingDB, error) {
	const op = "sqlite.database.GetListings"
	const queryWithFilter = `
		SELECT id, name, type, description, status, price, city, user_id, date_created
		FROM listings where user_id = ? and city = ? limit 10 offset ?
	`
	const queryWithoutFilter = `
		SELECT id, name, type, description, status, price, city, user_id, date_created
		from listings where user_id = ? limit 10 offset ?
	`
	var rows *sql.Rows // Prepare the query based on whether a filter is provided
	if filter != "" {
		//filter = "%" + filter + "%"
		stmt, err := s.db.Prepare(queryWithFilter)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		rows, err = stmt.Query(userID, filter, offset)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()
	} else {
		stmt, err := s.db.Prepare(queryWithoutFilter)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		rows, err = stmt.Query(userID, offset)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()
	}

	var listings []models.ListingDB
	for rows.Next() {
		var l models.ListingDB
		if err := rows.Scan(&l.ID, &l.Name, &l.Typel, &l.Description, &l.Status, &l.Price, &l.City, &l.UserID, &l.Date_created); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return listings, nil
}

func (s *service) GetCities(userID int64) ([]string, error) {
	const op = "sqlite.database.GetCities"
	const query = `
		SELECT DISTINCT city FROM listings WHERE user_id = ? ORDER BY city
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cities, nil

}

func (s *service) UpdateListing(name, typel, description, status, city string, price int64, id int64) error {
	const op = "sqlite.database.UpdateListing"
	const query = `
		UPDATE listings SET name = ?, type = ?, description = ?, status = ?, price = ?, city = ? WHERE id = ?;
	`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(name, typel, description, status, price, city, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *service) DeleteListing(id int64) error {
	const op = "sqlite.database.DeleteListing"
	const query = `
		DELETE FROM listings WHERE id = ?;
	`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *service) GetAnalytics(userID int64) (map[string]any, error) {
	const op = "sqlite.database.AnalyticsHandler"
	const query = `
		SELECT COUNT(*), AVG(price)
		FROM listings where user_id = ?;
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var count int64
	var avgPrice float64
	err = stmt.QueryRow(userID).Scan(&count, &avgPrice)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	const query2 = `
		SELECT city, COUNT(*) as count
		FROM listings
		WHERE user_id = ?
		GROUP BY city
		ORDER BY count DESC
		LIMIT 3
	`
	stmt, err = s.db.Prepare(query2)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	type CityStat struct {
		City  string `json:"city"`
		Count int    `json:"count"`
	}
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var topCities []CityStat
	for rows.Next() {
		var cs CityStat
		if err := rows.Scan(&cs.City, &cs.Count); err == nil {
			topCities = append(topCities, cs)
		}
	}

	return map[string]any{
		"total_listings": count,
		"avg_price":      avgPrice,
		"top_cities":     topCities,
	}, nil
}
