package rule

type Type string

const (
	BlackList Type = "black"
	WhiteList Type = "white"
)

type Rules []Rule

type Rule struct {
	ID       int
	IP       string
	RuleType Type
}

type IStorage interface {
	Create(rule Rule) (int, error)
	Delete(id int) error

	GetForType(ruleType Type) (*Rules, error)
	Find(ip string, ruleType Type) (*Rules, error)
}

type IService interface {
	InWhiteList(ip string) (bool, error)
	InBlackList(ip string) (bool, error)

	WhiteListAdd(ip string) error
	WhiteListDelete(ip string) error

	BlackListAdd(ip string) error
	BlackListDelete(ip string) error
}
