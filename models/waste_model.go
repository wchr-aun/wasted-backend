package models

type MasterWasteType struct {
	WasteId     int    `json:"wasteId"`
	Type        string `json:"type"`
	SubType     string `json:"subType"`
	Description string `json:"description"`
	Disposal    string `json:"disposal"`
	ImgUrl      string `json:"imgUrl"`
	Name        string `json:"name"`
	TrashTag    string `json:"trashTag"`
}

type SellerWaste struct {
	SellerId string `json:"sellerId" form:"sellerId"`
	WasteId  int    `json:"wasteId" form:"wasteId"`
	Amount   int    `json:"amount" form:"amount"`
}
