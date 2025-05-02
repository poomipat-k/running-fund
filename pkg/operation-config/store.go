package operationConfig

import (
	"database/sql"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{
		db: db,
	}
}

func (s *store) GetLatestConfig() (OperationConfig, error) {
	// Query count from database
	row := s.db.QueryRow(getOperationConfigSQL)
	var conf OperationConfig
	err := row.Scan(&conf.AllowNewProject)
	if err != nil {
		return OperationConfig{}, err
	}
	return conf, nil
}
