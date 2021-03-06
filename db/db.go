package db

import "database/sql"

type URLPath struct {
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

func (db *DB) FindURLbyPath(path string) (*URLPath, error) {
	var u URLPath
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

func (db *DB) SaveUrlAndPath(u URLPath) error {
	sqlstmt := `INSERT INTO url_table(path, url) VALUES ($1, $2)`
	_, err := db.db.Exec(sqlstmt, u.Path, u.URL)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUrlAndPath updates an existing path to a new URL. If the path does not exists, it creates a new one
// along with the corresponding URL.
func (db *DB) UpdateUrlAndPath(u URLPath) error {
	sqlstmt := `INSERT INTO url_table(path, url) VALUES ($1, $2) ON CONFLICT (path) DO UPDATE SET url = $2`
	_, err := db.db.Exec(sqlstmt, u.Path, u.URL)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteURLbyPath(path string) error {
	_, err := db.db.Exec("DELETE FROM url_table WHERE path=$1", path)
	if err != nil {
		return err
	}
	return nil
}
