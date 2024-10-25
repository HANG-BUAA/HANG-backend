package test

import (
	"HANG-backend/src/utils"
	"fmt"
	"testing"
)

func TestJwt(t *testing.T) {
	token, _ := utils.GenerateToken(1, "22371426")
	claim, err := utils.ParseToken(token)
	fmt.Println(claim)
	fmt.Println(err)
}
