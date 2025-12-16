package rule

import (
	"errors"
	"net"
)

var (
	ErrRuleNotFound   = errors.New("rule not found")
	ErrInvalidInputIP = errors.New("incorrect IP passed")
)

type Service struct {
	ruleStorage IStorage
}

func NewService(ruleStorage IStorage) *Service {
	return &Service{ruleStorage: ruleStorage}
}

func (s Service) InWhiteList(ip string) (bool, error) {
	return s.inList(ip, WhiteList)
}

func (s Service) InBlackList(ip string) (bool, error) {
	return s.inList(ip, BlackList)
}

func (s Service) WhiteListAdd(ip string) error {
	return s.listAdd(ip, WhiteList)
}

func (s Service) WhiteListDelete(ip string) error {
	return s.listDelete(ip, WhiteList)
}

func (s Service) BlackListAdd(ip string) error {
	return s.listAdd(ip, BlackList)
}

func (s Service) BlackListDelete(ip string) error {
	return s.listDelete(ip, BlackList)
}

func (s Service) listAdd(ip string, listType Type) error {
	_, err := s.ruleStorage.Create(Rule{
		IP:       ip,
		RuleType: listType,
	})

	return err
}

func (s Service) listDelete(ip string, listType Type) error {
	rules, err := s.ruleStorage.Find(ip, listType)
	if err != nil {
		return err
	}

	if len(*rules) == 0 {
		return ErrRuleNotFound
	}

	for _, rule := range *rules {
		deleteErr := s.ruleStorage.Delete(rule.ID)
		if deleteErr != nil {
			return deleteErr
		}
	}

	return nil
}

func (s Service) inList(ip string, listType Type) (bool, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false, ErrInvalidInputIP
	}

	rules, err := s.ruleStorage.GetForType(listType)
	if err != nil {
		return false, err
	}

	for _, rule := range *rules {
		if parsedIP.Equal(net.ParseIP(rule.IP)) {
			return true, nil
		}

		if _, netIP, err := net.ParseCIDR(rule.IP); err == nil && netIP.Contains(parsedIP) {
			return true, nil
		}
	}

	return false, nil
}
