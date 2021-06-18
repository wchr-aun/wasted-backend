package models

type UserResponse struct {
	Fullname    string `json:"fullname,omitempty"`
	PhoneNumber string `json:"phonenumber,omitempty"`
	PhotoUrl    string `json:"photourl,omitempty"`
	New         bool   `json:"new"`
}

type UserTable struct {
	Uuid        string `json:"uuid,omitempty"`
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phonenumber"`
	PhotoUrl    string `json:"photourl"`
}
