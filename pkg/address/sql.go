package address

const getProvincesSQL = "SELECT id, name FROM province;"

const getDistrictsSQL = "SELECT id, name FROM district WHERE province_id = $1;"

const getSubdistrictsSQL = "SELECT id, name FROM subdistrict WHERE district_id = $1;"

const getPostcodeSQL = "SELECT id, code FROM postcode WHERE subdistrict_id = $1;"
