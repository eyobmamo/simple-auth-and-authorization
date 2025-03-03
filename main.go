package main

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID uint `json:"id"`
	Email  string `json:"email"`
	Password string `json:"-"`
}

var users = make(map[string]*User)
var JwtSecret = []byte("your_jwt_secret")

func main() {
	router := gin.Default()

	router.GET("/",func (c *gin.Context){
		c.JSON(200,gin.H{
			"message":"Welecome to the Go Authentiations  ans authorization",
		})
	})
	router.POST("/Register",userRegister)
	router.POST("/Login",userLogin)
	router.GET("/profile",AuthMiddleware(),ChackAuth)
	router.Run(":8080")
}

func userRegister(c *gin.Context){
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400,gin.H{"error":"Invalid request payload"})
		return
	}

	//user registration logic
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500,gin.H{"error":"Internal server error"})
		return
	}

	user.Password =string(hashedPassword)
	users[user.Email] = &user

	c.JSON(200,gin.H{"message" : "User registered successful"})

}

func userLogin(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400,gin.H{"error": "Invalid request payload"})
		return
	}

	var jwtSecret = []byte("your_jwt_secret")

	//User login logic
	existingUser ,ok := users[user.Email]
	if !ok || bcrypt.CompareHashAndPassword([]byte(existingUser.Password),[]byte(user.Password)) != nil {
		c.JSON(401,gin.H{"error":"Invalid email or passoword"})
		return

	}

	//Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"user_id" : existingUser.ID,
		"email" : existingUser.Email,
	})

	jwtToken,err := token.SignedString(jwtSecret)

	if err != nil {
		c.JSON(500,gin.H{"error":"Internal server error"})
		return
	}


	c.JSON(200,gin.H{"message" : " User logged in successful","token":jwtToken})

}



func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//jwt validation logic
		authHeader := c.GetHeader("Authorization")
		if authHeader == ""{
			c.JSON(401,gin.H{"error":"Authorization header is required"})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader," ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(401,gin.H{"error" : "Invalid authorization header"})
			c.Abort()
			return
		}

		token,err := jwt.Parse(authParts[1],func(token *jwt.Token) (interface{},error) {
			if _,ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
				return nil, fmt.Errorf("Unexpected signing method : %v",token.Header["alg"])
			}
			return JwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(401,gin.H{"error":"Invalid JWT"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ChackAuth( c *gin.Context) {
	c.JSON(200,gin.H{"message" : "This is a secure router"})
}