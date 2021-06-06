package models

type Waste struct {
	Type        string `json:"type"`
	SubType     string `json:"subType"`
	Description string `json:"description"`
	Disposal    string `json:"disposal"`
	ImgUrl      string `json:"imgUrl"`
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	TrashTag    int64  `json:"trashTag"`
}
