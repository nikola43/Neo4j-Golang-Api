package test

import (
	"github.com/nikola43/ecodadys_api/models"
	"github.com/nikola43/ecodadys_api/utils"
	"log"
	"testing"
	"time"
)

func TestGetAll(t *testing.T) {

	url := "http://localhost:8080/api/user"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InNzQGdtYWlsLmNvbSIsImV4cCI6MTU3MDc1MjcyNX0.5SiPPXjU1iSnW_BWFtvzm1yyB6pCPUSBhpQV29Fs8Rw"

	start := time.Now()
	utils.GetRequest(url, token, nil)
	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
}

func TestInsert(t *testing.T) {

	url := "http://localhost:8080/api/user/new"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InNzQGdtYWlsLmNvbSIsImV4cCI6MTU3MDc1MjcyNX0.5SiPPXjU1iSnW_BWFtvzm1yyB6pCPUSBhpQV29Fs8Rw"

	i := 0
	start := time.Now()
	for i = 0; i < 1000; i++ {
		user := &models.User{Name: utils.GenerateRandomString(10), Username: utils.GenerateRandomString(10), Password: utils.GenerateRandomString(10)}
		utils.PostRequest(url, token, user)
		time.Sleep(50 * time.Microsecond)
	}
	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
}

func TestGetUser(t *testing.T) {

	url := "http://localhost:8080/api/user/85"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InNzQGdtYWlsLmNvbSIsImV4cCI6MTU3MDc1MjcyNX0.5SiPPXjU1iSnW_BWFtvzm1yyB6pCPUSBhpQV29Fs8Rw"

	i := 0
	start := time.Now()
	for i = 0; i < 10000; i++ {
		utils.GetRequest(url, token, nil)

	}
	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
}
