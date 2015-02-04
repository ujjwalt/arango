package arango

type Database struct {
	Name     string      // Name of the databse
	Id       string      // Id of the database
	Path     string      // Filesystem path of the database
	IsSystem bool        // Whether or not the current database is the _system database
	conn     *Connection // The connection to which this database belongs
}

var DefaultDatabase *Database

func (c *Connection) NewDatabase() *Database {
	db := &Database{conn: c}
	if c == DefaultConnection && DefaultDatabase == nil {
		DefaultDatabase = db
	}
	return db
}
