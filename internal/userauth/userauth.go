package userauth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "tratatata"

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func CookieMiddleware(h http.Handler) http.Handler {
	cookieFn := func(w http.ResponseWriter, r *http.Request) {
		if !cookieIsValid(r) {
			cookie := createNewCookie()
			http.SetCookie(w, &cookie)
		}
	}
	return http.HandlerFunc(cookieFn)
}

func cookieIsValid(r *http.Request) bool {
	cookie, err := r.Cookie("auth_token")
	// проверяем есть ли кука
	if err != nil {
		return false
	}

	// в случае если кука есть проверяем что она проходит проверку подлинности
	token := cookie.Value
	_, ok := GetUserId(token)
	return ok
}



func GetUserId(tokenString string) (string, bool) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return "", false
	}

	if !token.Valid {
		log.Printf("Token is not valid")
		return "", false
	}

	log.Printf("Token is valid")
	return claims.UserID, true
}

func createNewCookie() http.Cookie {
	tokenString, err := buildJWTString()
	if err != nil {
		log.Fatal(err)
	}
	// создание новой куки для юзера если такой куки не существует или она не проходит проверку подлинности
	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		MaxAge:   3600 * 3,
		HttpOnly: true,
		Secure:   true,
	}

	return cookie
}

func buildJWTString() (string, error) {
	newId := uuid.New().String()
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		// собственное утверждение

		UserID: newId,
		// TODO: тут добавить данные по сокращенным урл
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}
