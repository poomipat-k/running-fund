package address

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/patrickmn/go-cache"
)

const provinceCachePrefix = "province"
const distByProv = "distByProv"

type store struct {
	db *sql.DB
	c  *cache.Cache
}

func NewStore(db *sql.DB, c *cache.Cache) *store {
	return &store{
		db: db,
		c:  c,
	}
}

func (s *store) GetProvinces() ([]Province, error) {
	// check cache first
	rawCachedData, found := s.c.Get(fmt.Sprintf("%sall", provinceCachePrefix))
	if found {
		log.Println("===PROVINCES found in cache")
		cachedData, ok := rawCachedData.([]Province)
		if ok {
			log.Println("===cached Used")
			return cachedData, nil
		}
	}
	// Fetch data from the db
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
	s.c.Set(fmt.Sprintf("%s__all", provinceCachePrefix), data, cache.NoExpiration)
	return data, nil
}

func (s *store) GetDistrictsByProvince(provinceId int) ([]District, error) {

	rows, err := s.db.Query(getDistrictsSQL, provinceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []District
	for rows.Next() {
		var row District
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
	log.Println("===dist data")
	log.Println(data)
	s.c.Set(fmt.Sprintf("%s__%d", distByProv, provinceId), data, cache.NoExpiration)
	return data, nil
}
