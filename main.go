package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MiddlewaresCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	})
}

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

func Login(c *gin.Context) {
	account := strings.TrimSpace(c.Query("account"))
	password := strings.TrimSpace(c.Query("password"))

	message := ""
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
				if account == uaccount {
					h := md5.New()
					io.WriteString(h, password)
					pwdmd5 := hex.EncodeToString(h.Sum(nil))
					upassword := sc[lastIndex+1:]
					if pwdmd5 == upassword {
						message = "Welcome " + account
						status = 200
						goto outJson
					} else {
						message = "Account or password is wrong"
						status = 400
						goto outJson
					}
				}
			}
			message = "Account or password is wrong"
			status = 400
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
	router.Use(MiddlewaresCors())
	router.GET("/reguser", RegUser)
	router.GET("/login", Login)
	router.Run(":8080")
}
