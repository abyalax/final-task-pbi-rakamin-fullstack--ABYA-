package photocontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Backend/go-jwt-gin-gorm/database"
	"github.com/Backend/go-jwt-gin-gorm/helpers"
	"github.com/Backend/go-jwt-gin-gorm/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Index(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := database.DB.Preload("Photos").Find(&users).Error; err != nil {
		helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	helpers.ResponseJSON(w, http.StatusOK, users)
}

func Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
		return
	}

	var user models.User
	if err := database.DB.Preload("Photos").First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ResponseJSON(w, http.StatusNotFound, map[string]string{"message": "User not found"})
		} else {
			helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return
	}

	helpers.ResponseJSON(w, http.StatusOK, user)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	var userInput models.Photo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	if err := database.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Success"}
	helpers.ResponseJSON(w, http.StatusOK, response)
}

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid photo ID"})
		return
	}

	var userInput models.Photo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	defer r.Body.Close()

	var photo models.Photo
	if err := database.DB.First(&photo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ResponseJSON(w, http.StatusNotFound, map[string]string{"message": "Photo not found"})
		} else {
			helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return
	}

	photo.Title = userInput.Title
	photo.Caption = userInput.Caption
	photo.PhotoURL = userInput.PhotoURL

	if err := database.DB.Save(&photo).Error; err != nil {
		helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	helpers.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Success"})
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid photo ID"})
		return
	}

	if err := database.DB.Delete(&models.Photo{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ResponseJSON(w, http.StatusNotFound, map[string]string{"message": "Photo not found"})
		} else {
			helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return
	}

	helpers.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Success"})
}
