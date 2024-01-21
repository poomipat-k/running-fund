package address

const getProvincesSQL = "SELECT id, name FROM province;"

const getDistrictsSQL = "SELECT id, name, province_id FROM district;"

const getSubdistrictsSQL = "SELECT id, name, district_id FROM subdistrict;"
