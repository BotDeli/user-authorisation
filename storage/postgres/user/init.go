package user

type Display interface {
	newUser(login, password string) error
	changePassword(login, password, newPassword string) error
	getUserHash(login string) string
}
