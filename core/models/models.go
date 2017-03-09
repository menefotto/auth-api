package models

type User struct {
	Username    string `json:"username"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Photourl    string `json:"photourl"`
	Isactive    string `json:"isactive"`
	Issuperuser string `json:"issuperuser"`
	Isstaff     string `json:"isstaff"`
	Datejoined  string `json:"datejoined"`
	Uuid        string `json:"uuid"`
	Code        string `json:"confirmationcode"`
}
