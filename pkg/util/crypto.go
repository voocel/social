package util

import "golang.org/x/crypto/bcrypt"

func VerifyPassword(loginPass, userPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPass), []byte(loginPass))
	return err == nil
}

func EncryptPassword(password string) (string, error) {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashPwd), err
}
