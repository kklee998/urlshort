package db

import "database/sql"

type URL struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

type DB struct {
	db *sql.DB
}

func (db *DB) Close() error {
	return db.db.Close()
}

// StartDB will create the db file and table if it does not exist
func (db *DB) StartDB() error {
	err := db.db.Ping()
	if err != nil {
		return err
	}
	sqlstmt := `CREATE TABLE IF NOT EXISTS url_table (path TEXT UNIQUE, url TEXT)`

	_, err = db.db.Exec(sqlstmt)
	if err != nil {
		return err
	}
	return nil

}

func (db *DB) FindURLbyPath(path string) (*URL, error) {
	var u URL
	row := db.db.QueryRow("SELECT url from url_table WHERE path=$1", path)

	err := row.Scan(&u.URL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &u, nil
}

func (db *DB) SaveUrlAndPath(u URL) error {
	sqlstmt := `INSERT INTO url_table VALUES ($1, $2)`
	_, err := db.db.Exec(sqlstmt, u.Path, u.URL)
	if err != nil {
		return err
	}
	return nil
}
