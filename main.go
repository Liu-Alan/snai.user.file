package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegUser(c *gin.Context) {
	account := strings.TrimSpace(c.Query("account"))
	password := strings.TrimSpace(c.Query("password"))

	message := "success"
	status := 200

	if account == "" || password == "" {
		message = "Account or password cannot be empty"
		status = 400
	} else {
		f, err := os.Open("user.log")
		defer f.Close()
		if err != nil {
			message = "Service error" //err.Error()
			status = 400
		} else {
			s := bufio.NewScanner(f)
			for s.Scan() {
				sc := s.Text()
				lastIndex := strings.LastIndex(sc, "|")
				uaccount := sc[:lastIndex]
				fmt.Println(uaccount)
				if account == uaccount {
					message = "Account already exists"
					status = 400
					goto outJson
				}
			}
			f.Close()

			f, err = os.OpenFile("user.log", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				message = "Service error" //err.Error()
				status = 400
				goto outJson
			}

			h := md5.New()
			io.WriteString(h, password)
			pwdmd5 := hex.EncodeToString(h.Sum(nil))
			_, err = fmt.Fprintln(f, account+"|"+pwdmd5)
			if err != nil {
				message = "Service error" //err.Error()
				status = 400
				goto outJson
			}

			message = "success"
			status = 200
		}

	}
outJson:
	c.JSON(200, gin.H{
		"status":  status,
		"message": message,
	})
}

func main() {
	router := gin.Default()
	router.GET("/reguser", RegUser)
	router.Run(":8080")
}
