package models

type User struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhotoUrl    string `json:"photourl"`
	IsActive    bool   `json:"isactive"`
	IsSuperUser bool   `json:"issuperuser"`
	IsStaff     bool   `json:"isstaff"`
	DateJoined  string `json:"datejoined"`
	Uuid        string `json:"uuid"`
	Code        string `json:"confirmationcode"`
}
