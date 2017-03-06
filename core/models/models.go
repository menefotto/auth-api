package models

type User struct {
	Username    string `json:"username"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Photourl    string `json:"photourl"`
	Isactive    bool   `json:"isactive"`
	Issuperuser bool   `json:"issuperuser"`
	Isstaff     bool   `json:"isstaff"`
	Datejoined  string `json:"datejoined"`
	Uuid        string `json:"uuid"`
	Code        string `json:"confirmationcode"`
}
