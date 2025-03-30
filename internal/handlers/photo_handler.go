package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"MicroGreens/internal/models"
	"MicroGreens/internal/services"
)

type PhotoHandler struct {
	Service *services.PhotoService
}

func (h *PhotoHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Cannot read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Генерируем имя файла
	filename := handler.Filename
	filePath := "./uploads/" + filename

	dst, err := createFile(filePath)
	if err != nil {
		http.Error(w, "Cannot save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Сохраняем содержимое
	_, err = dst.ReadFrom(file)
	if err != nil {
		http.Error(w, "Cannot write file", http.StatusInternalServerError)
		return
	}

	// Сохраняем данные в БД
	observationID, _ := strconv.Atoi(r.FormValue("observation_id"))
	label := r.FormValue("label")

	newPhoto := models.ObservationPhoto{
		ObservationID: observationID,
		PhotoURL:      "/uploads/" + filename,
		Label:         label,
	}

	created, err := h.Service.Create(r.Context(), newPhoto)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(created)
}

func createFile(path string) (*os.File, error) {
	return os.Create(path)
}

func (h *PhotoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.Service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (h *PhotoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	item, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func (h *PhotoHandler) Update(w http.ResponseWriter, r *http.Request) {
	var p models.ObservationPhoto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.Service.Update(r.Context(), p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *PhotoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if err := h.Service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
