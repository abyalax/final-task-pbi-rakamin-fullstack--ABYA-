package authcontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Backend/go-jwt-gin-gorm/config"
	"github.com/Backend/go-jwt-gin-gorm/database"
	"github.com/Backend/go-jwt-gin-gorm/helpers"
	"github.com/Backend/go-jwt-gin-gorm/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// mengambil inputan dari json
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()
	fmt.Printf("User input: %+v\n", userInput.Username)

	// ambil data user berdasarkan username
	var user models.User
	// fmt.Println(user)
	if err := database.DB.Where("username = ?", userInput.Username).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": "Username atau Password salah"}
			helpers.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helpers.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// cek password apakah valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{"message": "Username atau Password salah"}
		helpers.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// proses pembuatan token jwt
	expTime := time.Now().Add(time.Minute * 10)
	claims := &config.JWTClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// mendeklarasikan algoritma yang akan digunakan untuk signing
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//signed token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// set token ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	response := map[string]string{"message": "login berhasil"}
	helpers.ResponseJSON(w, http.StatusOK, response)

}
func Register(w http.ResponseWriter, r *http.Request) {
	// mengambil inputan dari json
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	//hash pass menggunakan bcrypt
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	//Insert ke database
	if err := database.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
	}

	response := map[string]string{"message": "succes"}
	helpers.ResponseJSON(w, http.StatusOK, response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// hapus token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "logout berhasil"}
	helpers.ResponseJSON(w, http.StatusOK, response)
}
