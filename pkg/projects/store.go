package projects

import (
	"database/sql"
	"log"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{
		db: db,
	}
}

func (s *store) GetAll() ([]Project, error) {
	rows, err := s.db.Query("SELECT id, project_code, project_version, created_at FROM project")
	if err != nil {
		log.Println("Error on Query: ", err)
		return nil, err
	}
	defer rows.Close()

	var row Project
	var data []Project
	for rows.Next() {
		err = rows.Scan(&row.Id, &row.ProjectCode, &row.ProjectVersion, &row.CreatedAt)
		if err != nil {
			log.Println("Error on Scan: ", err)
			return nil, err
		}
		data = append(data, row)
	}
	// get any error cncountered during iteration
	err = rows.Err()
	if err != nil {
		log.Println("Error on rows.Err: ", err)
		return nil, err
	}
	return data, nil
}
