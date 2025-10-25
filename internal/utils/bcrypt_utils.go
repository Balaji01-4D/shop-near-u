package utils

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func ParseUintParam(param string) (uint, error) {
	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func ParseFloatParam(param string) (float64, error) {
	value, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func ParseIntParam(param string) (int, error) {
	value, err := strconv.Atoi(param)

	if err != nil {
		return 0, err
	}
	return value, nil
}