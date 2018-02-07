package main

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"strings"
)

func main() {
	r := gin.Default()
	r.StaticFile("/", "./index.html")
	r.GET("/profile", profile)
	r.Run()
}

func handleError(err error) {
	log.Printf("%+v", err)
}

func extractToken(r *http.Request) string {
	ss := strings.Split(r.Header.Get("Authorization"), " ")
	if len(ss) > 1 {
		return ss[1]
	}
	return ""
}

func createAppAuth(ctx context.Context) (*auth.Client, error) {
	opt := option.WithCredentialsFile("firebase-secret.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	appAuth, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return appAuth, nil
}

func getUser(ctx context.Context, r *http.Request) (*auth.UserRecord, error) {
	requestToken := extractToken(r)
	if requestToken == "" {
		return nil, nil
	}

	appAuth, err := createAppAuth(ctx)
	token, err := appAuth.VerifyIDToken(requestToken)
	if err != nil {
		return nil, err
	}

	user, err := appAuth.GetUser(ctx, token.UID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func profile(c *gin.Context) {
	user, err := getUser(c, c.Request)
	if err != nil {
		handleError(err)
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "user": user.UserInfo})
}
