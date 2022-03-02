package db

import (
	"log"
	"testing"
)

func TestUserDao_GetById(t *testing.T) {
	u := UserDao.GetById(1)
	log.Println(u)
}

func TestUserDao_GetByEmail(t *testing.T) {
	u := UserDao.GetByEmail("974875956@qq.com")
	log.Println(u)
}

func TestUserDao_GetByPhone(t *testing.T) {
	u := UserDao.GetByPhone("15652087361")
	log.Println(u)
}
