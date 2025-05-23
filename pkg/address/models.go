package address

type Province struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type District struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ProvinceId int    `json:"provinceId,omitempty"`
}

type Subdistrict struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	DistrictId int    `json:"districtId,omitempty"`
}

type Postcode struct {
	Id            int `json:"id"`
	Code          int `json:"code"`
	SubdistrictId int `json:"subdistrictId,omitempty"`
}
