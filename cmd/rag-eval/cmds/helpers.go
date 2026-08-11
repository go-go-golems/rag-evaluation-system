package cmds

import (
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
)

// OpenDBAtPath opens and migrates a database at the given path, returning Queries
func OpenDBAtPath(path string) (*db.Queries, error) {
	database, err := db.OpenDB(path)
	if err != nil {
		return nil, err
	}
	if err := db.Migrate(database); err != nil {
		_ = database.Close()
		return nil, err
	}
	return db.NewQueries(database), nil
}
