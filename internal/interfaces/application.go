package interfaces

type Application interface {
	LimitCheck(ip, login, password string) (bool, error)
	LimitReset(ip, login string) error

	WhiteListAdd(ip string) error
	WhiteListDelete(ip string) error

	BlackListAdd(ip string) error
	BlackListDelete(ip string) error
}
