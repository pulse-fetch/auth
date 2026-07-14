package service

import "fmt"

func (s *Service) Register(username, email, password string) (string, error) {
	fmt.Println(username)
	return "", nil
}
