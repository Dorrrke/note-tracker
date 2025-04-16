package models

type UserCreds struct {
	Login    string
	Password string
}

type User struct {
	UID      string
	Name     string
	Login    string
	Password string
}
