package domain

type User struct {
	ID      int64
	Name    string
	Email   string
	PwdHash string
}

func NewUser(name, email, pwdHash string) *User {
	return &User{
		Email:   email,
		Name:    name,
		PwdHash: pwdHash,
	}
}

type UserModel struct {
	ID      int64
	Name    string
	Email   string
	PwdHash string
}
