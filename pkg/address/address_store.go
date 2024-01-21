package address

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

func (s *store) GetProvinces() ([]Province, error) {
	rows, err := s.db.Query(getProvincesSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []Province
	for rows.Next() {
		var row Province
		err := rows.Scan(&row.Id, &row.Name)
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	log.Println("===store data", data)
	return data, nil
}

// func (s *store) GetDistricts() ([]District, error) {
// 	rows, err := s.db.Query(getDistrictsSQL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var data []Province
// 	for rows.Next() {
// 		var row Province
// 		err := rows.Scan(&row.Id, &row.Name)
// 		if err != nil {
// 			return nil, err
// 		}
// 		data = append(data, row)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Println("===store data", data)
// 	return data, nil
// }
