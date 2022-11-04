package auth

import (
    "net/http"
    db "example.com/chatapp/db/config"
    "example.com/chatapp/db/models"
    "time"
    _"os"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
    "example.com/chatapp/utils"
)

type RefreshToken struct {
    Token string `json:"refresh_token"`
}

var jwtKey = []byte(utils.GetDotEnvVariable("ACCESS_SECRET"))

// GenerateTokens returns the access and refresh tokens
func GenerateTokens(username string) (string, string) {
    start := time.Now()
    claim, accessToken := GenerateAccessClaims(username)
    end := time.Since(start)
    log.Println("GenerateAccessClaims",end)
    start = time.Now()
    refreshToken := GenerateRefreshClaims(claim)
    end = time.Since(start)
    log.Println("GenerateRefreshClaims",end)

    return accessToken, refreshToken
}

// GenerateAccessClaims returns a claim and a acess_token string
func GenerateAccessClaims(username string) (*models.Claims, string) {

    t := time.Now()
    claim := &models.Claims{
        StandardClaims: jwt.StandardClaims{
            Issuer:    username,
            ExpiresAt: t.Add(15 * time.Minute).Unix(),
            Subject:   "access_token",
            IssuedAt:  t.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        panic(err)
    }

    return claim, tokenString
}

// GenerateRefreshClaims returns refresh_token
func GenerateRefreshClaims(cl *models.Claims) string {

    start := time.Now()
    result := db.DB.Where(&models.Claims{
        StandardClaims: jwt.StandardClaims{
            Issuer: cl.Issuer,
        },
    }).Find(&models.Claims{})
    end := time.Since(start)
    log.Println("Finding proccess",end)

    // checking the number of refresh tokens stored.
    // If the number is higher than 0, remove all the refresh tokens and leave only new one.
    log.Println(result.RowsAffected,cl.Issuer)
    start = time.Now()
    if result.RowsAffected > 0 {
        db.DB.Where(&models.Claims{
            StandardClaims: jwt.StandardClaims{Issuer: cl.Issuer},
        }).Delete(&models.Claims{})
    }
    end = time.Since(start)
    log.Println("Deleting proccess",end)

    t := time.Now()
    refreshClaim := &models.Claims{
        StandardClaims: jwt.StandardClaims{
            Issuer:    cl.Issuer,
            ExpiresAt: t.Add(30 * 24 * time.Hour).Unix(),
            Subject:   "refresh_token",
            IssuedAt:  t.Unix(),
        },
    }

    // create a claim on DB
    start = time.Now()
    db.DB.Create(&refreshClaim)
    end = time.Since(start)
    log.Println("Creating proccess",end)

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
    refreshTokenString, err := refreshToken.SignedString(jwtKey)
    if err != nil {
        panic(err)
    }

    return refreshTokenString
}


func GetAccessToken(c *gin.Context) {
    var refresh_token RefreshToken
    // refreshToken := c.Cookies("refresh_token")

    if err := c.ShouldBindJSON(&refresh_token); err != nil {
        log.Fatal(err)
        c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"error":"Invalid json provided."})
        return
    }

    refreshClaims := new(models.Claims)
    token, _ := jwt.ParseWithClaims(refresh_token.Token, refreshClaims,
        func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

    if res := db.DB.Where(
        "expires_at = ? AND issued_at = ? AND issuer = ?",
        refreshClaims.ExpiresAt, refreshClaims.IssuedAt, refreshClaims.Issuer,
    ).First(&models.Claims{}); res.RowsAffected <= 0 {
        // no such refresh token exist in the database
        // c.ClearCookie("access_token", "refresh_token")
        log.Println("Not found from db")
        c.IndentedJSON(http.StatusForbidden,gin.H{"error":"Forbidden."})
        return
    }

    if token.Valid {
        if refreshClaims.ExpiresAt < time.Now().Unix() {
            // refresh token is expired
            // c.ClearCookie("access_token", "refresh_token")
            c.IndentedJSON(http.StatusForbidden,gin.H{"error":"Forbidden."})
        }
    } else {
        // malformed refresh token
        // c.ClearCookie("access_token", "refresh_token")
        c.IndentedJSON(http.StatusForbidden,gin.H{"error":"Forbidden."})
        return
    }

    _, accessToken := GenerateAccessClaims(refreshClaims.Issuer)

    c.IndentedJSON(http.StatusOK,gin.H{"access_token": accessToken})
    return
}
