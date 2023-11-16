package library

import (
	"encoding/json"
	"golibrary/utils"
	"net/http"
)

func (lf *LibraryFacade) GetAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	authors, err := utils.GetAuthors(lf.DB)
	if err != nil {
		http.Error(w, "Ошибка при получении информации об авторах", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}
