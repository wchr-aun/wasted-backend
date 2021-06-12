package models

type User struct {
	Uuid        string `json:"uuid"`
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phonenumber"`
	New         bool   `json:"new"`
}
