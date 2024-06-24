package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// TokenExp - время жизни токена
const TokenExp = time.Hour * 3

// SecretKey - секретный ключ для шифрования
const SecretKey = "tratatata"

// Claims заявления при создании куки
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// CookieMiddleware создает куки если её не было, и добавляет к запросу и к ответу.
func CookieMiddleware(h http.Handler) http.Handler {
	cookieFn := func(w http.ResponseWriter, r *http.Request) {
		if !cookieIsValid(r) {
			newID := uuid.New().String()
			cookie, err := CreateNewCookie(newID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.AddCookie(&cookie)
			http.SetCookie(w, &cookie)
		}
		h.ServeHTTP(w, r)
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
	_, ok := GetUserID(token)
	return ok
}

// GetUserID получает id пользователя из куки
func GetUserID(tokenString string) (string, bool) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
	if err != nil {
		return "", false
	}

	if !token.Valid {
		log.Printf("Token is not valid")
		return "", false
	}

	log.Printf("Token is valid.")
	return claims.UserID, true
}

// CreateNewCookie - создание новой куки для юзера если такой куки не существует или она не проходит проверку подлинности.
func CreateNewCookie(newID string) (http.Cookie, error) {
	if len(newID) == 0 {
		return http.Cookie{}, errors.New("no id to put into cookie")
	}

	tokenString, err := buildJWTString(newID)
	if err != nil {
		return http.Cookie{}, err
	}

	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   3600 * 3,
		HttpOnly: true,
		Secure:   false,
	}

	return cookie, nil
}

func buildJWTString(newID string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		// собственное утверждение
		UserID: newID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}
