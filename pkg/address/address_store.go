package address

import (
	"database/sql"
	"fmt"

	"github.com/patrickmn/go-cache"
)

const provinceCachePrefix = "province"
const distByProv = "distByProv"
const subDistrictByDist = "subByDist"

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
	raw, found := s.c.Get(fmt.Sprintf("%sall", provinceCachePrefix))
	if found {
		cachedData, ok := raw.([]Province)
		if ok {
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
	if len(data) > 0 {
		s.c.Set(fmt.Sprintf("%s__all", provinceCachePrefix), data, cache.NoExpiration)
	}
	return data, nil
}

func (s *store) GetDistrictsByProvince(provinceId int) ([]District, error) {
	// Check cache
	raw, found := s.c.Get(fmt.Sprintf("%s__%d", distByProv, provinceId))
	if found {
		cachedData, ok := raw.([]District)
		if ok {
			return cachedData, nil
		}
	}
	// If not found check in db
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
	if len(data) > 0 {
		s.c.Set(fmt.Sprintf("%s__%d", distByProv, provinceId), data, cache.NoExpiration)
	}
	return data, nil
}

func (s *store) GetSubdistrictsByDistrict(districtId int) ([]Subdistrict, error) {
	// Check cache
	raw, found := s.c.Get(fmt.Sprintf("%s__%d", subDistrictByDist, districtId))
	if found {
		cachedData, ok := raw.([]Subdistrict)
		if ok {
			return cachedData, nil
		}
	}
	// If not found check in db
	rows, err := s.db.Query(getSubdistrictsSQL, districtId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []Subdistrict
	for rows.Next() {
		var row Subdistrict
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
	if len(data) > 0 {
		s.c.Set(fmt.Sprintf("%s__%d", subDistrictByDist, districtId), data, cache.NoExpiration)
	}
	return data, nil
}

func (s *store) GetPostcodeBySubdistrict(subdistrictId int) ([]Postcode, error) {
	rows, err := s.db.Query(getPostcodeSQL, subdistrictId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []Postcode
	for rows.Next() {
		var row Postcode
		err := rows.Scan(&row.Id, &row.Code)
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
}
