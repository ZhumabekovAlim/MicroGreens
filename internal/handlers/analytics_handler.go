package handlers

import (
	_ "MicroGreens/internal/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type AnalyticsHandler struct {
	DB *sql.DB
}

func NewAnalyticsHandler(db *sql.DB) *AnalyticsHandler {
	return &AnalyticsHandler{DB: db}
}

type BatchAnalytics struct {
	BatchID      int     `json:"batch_id"`
	TotalDays    int     `json:"total_days"`
	PassedDays   int     `json:"passed_days"`
	Remaining    int     `json:"remaining_days"`
	Progress     float64 `json:"progress_percent"`
	SowingDate   string  `json:"sowing_date"`
	HarvestIn    int     `json:"estimated_harvest_days"`
	Observations int     `json:"observed_days"`
}

type HumidityAnalytics struct {
	Date            string `json:"date"`
	HumidityPercent int    `json:"humidity_percent"`
}

func (h *AnalyticsHandler) GetBatchProgress(w http.ResponseWriter, r *http.Request) {
	batchIDStr := r.URL.Query().Get("batch_id")
	if batchIDStr == "" {
		http.Error(w, "Missing batch_id", http.StatusBadRequest)
		return
	}
	batchID, err := strconv.Atoi(batchIDStr)
	if err != nil {
		http.Error(w, "Invalid batch_id", http.StatusBadRequest)
		return
	}

	query := `
		SELECT b.id, b.sowing_date, b.estimated_harvest_days,
			(SELECT COUNT(*) FROM observations o WHERE o.batch_id = b.id) as observation_count
		FROM batches b
		WHERE b.id = ?`

	row := h.DB.QueryRow(query, batchID)

	today := time.Now().Truncate(24 * time.Hour)

	var id, harvestDays, observations int
	var sowingDate time.Time

	err = row.Scan(&id, &sowingDate, &harvestDays, &observations)
	if err != nil {
		http.Error(w, "Batch not found", http.StatusNotFound)
		return
	}

	totalDays := harvestDays
	passed := int(today.Sub(sowingDate).Hours() / 24)
	remaining := totalDays - passed
	if remaining < 0 {
		remaining = 0
	}
	progress := float64(passed) / float64(totalDays) * 100
	if progress > 100 {
		progress = 100
	}

	result := BatchAnalytics{
		BatchID:      id,
		TotalDays:    totalDays,
		PassedDays:   passed,
		Remaining:    remaining,
		Progress:     progress,
		SowingDate:   sowingDate.Format("2006-01-02"),
		HarvestIn:    harvestDays,
		Observations: observations,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *AnalyticsHandler) GetHumidityLast7Days(w http.ResponseWriter, r *http.Request) {
	batchIDStr := r.URL.Query().Get("batch_id")
	if batchIDStr == "" {
		http.Error(w, "Missing batch_id", http.StatusBadRequest)
		return
	}
	batchID, err := strconv.Atoi(batchIDStr)
	if err != nil {
		http.Error(w, "Invalid batch_id", http.StatusBadRequest)
		return
	}

	query := `
		SELECT date, humidity_percent
		FROM observations
		WHERE batch_id = ? AND date >= CURDATE() - INTERVAL 7 DAY
		ORDER BY date ASC`

	rows, err := h.DB.Query(query, batchID)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []HumidityAnalytics
	for rows.Next() {
		var hEntry HumidityAnalytics
		var date time.Time
		if err := rows.Scan(&date, &hEntry.HumidityPercent); err != nil {
			http.Error(w, "Data parse error", http.StatusInternalServerError)
			return
		}
		hEntry.Date = date.Format("2006-01-02")
		result = append(result, hEntry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *AnalyticsHandler) GetHeightLast7Days(w http.ResponseWriter, r *http.Request) {
	batchIDStr := r.URL.Query().Get("batch_id")
	if batchIDStr == "" {
		http.Error(w, "Missing batch_id", http.StatusBadRequest)
		return
	}
	batchID, err := strconv.Atoi(batchIDStr)
	if err != nil {
		http.Error(w, "Invalid batch_id", http.StatusBadRequest)
		return
	}

	query := `
		SELECT date, height_cm
		FROM observations
		WHERE batch_id = ? AND date >= CURDATE() - INTERVAL 7 DAY
		ORDER BY date ASC`

	rows, err := h.DB.Query(query, batchID)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type HeightEntry struct {
		Date     string  `json:"date"`
		HeightCM float64 `json:"height_cm"`
	}

	var result []HeightEntry
	for rows.Next() {
		var entry HeightEntry
		var date time.Time
		if err := rows.Scan(&date, &entry.HeightCM); err != nil {
			http.Error(w, "Data parse error", http.StatusInternalServerError)
			return
		}
		entry.Date = date.Format("2006-01-02")
		result = append(result, entry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *AnalyticsHandler) GetAIPrediction(w http.ResponseWriter, r *http.Request) {
	batchIDStr := r.URL.Query().Get("batch_id")
	if batchIDStr == "" {
		http.Error(w, "Missing batch_id", http.StatusBadRequest)
		return
	}
	batchID, err := strconv.Atoi(batchIDStr)
	if err != nil {
		http.Error(w, "Invalid batch_id", http.StatusBadRequest)
		return
	}

	query := `
		SELECT date, note, height_cm, water_status, light_type, humidity_percent
		FROM observations
		WHERE batch_id = ?
		ORDER BY date ASC`

	rows, err := h.DB.Query(query, batchID)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Obs struct {
		Date     string
		Note     string
		Height   float64
		Water    string
		Light    string
		Humidity int
	}

	var observations []Obs
	for rows.Next() {
		var o Obs
		var date time.Time
		if err := rows.Scan(&date, &o.Note, &o.Height, &o.Water, &o.Light, &o.Humidity); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		o.Date = date.Format("2006-01-02")
		observations = append(observations, o)
	}
	var microgreenName string

	err = h.DB.QueryRow(`
	SELECT m.name 
	FROM batches b
	JOIN microgreens m ON b.microgreen_id = m.id
	WHERE b.id = ?`, batchID).Scan(&microgreenName)

	if err != nil {
		http.Error(w, "Не удалось найти тип микрозелени", http.StatusNotFound)
		return
	}

	prompt := fmt.Sprintf(`На основе следующих наблюдений по выращиванию микрозелени "%s", строго верни JSON-ответ вида:

{
  "success_rate": "85%s",
  "est_yield": "средний",
  "quality": "хорошее"
} - это просто рандомный пример.

Никаких пояснений. Только JSON. Дай корректный ответ, а не повторяй пример, сделай реальный анализ из Наблюдении .

Наблюдения:
`, microgreenName, "%")

	for _, obs := range observations {
		prompt += fmt.Sprintf("Дата: %s, Высота: %.1fсм, Влажность: %d%%, Полив: %s, Свет: %s, Заметка: %s\n", obs.Date, obs.Height, obs.Humidity, obs.Water, obs.Light, obs.Note)
	}

	// Вызов ChatGPT
	response, err := callChatGPT(prompt)
	if err != nil {
		http.Error(w, "AI prediction error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"prediction": response,
	})
}

func callChatGPT(prompt string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	apiURL := "https://api.openai.com/v1/chat/completions"
	apiKey := os.Getenv("OPENAI_API_KEY")

	body := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "Ты агроном-эксперт по микрозелени."},
			{"role": "user", "content": prompt},
		},
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	message := result["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"]
	return message.(string), nil
}

func (h *AnalyticsHandler) GetTodayObservations(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	today := time.Now().Format("2006-01-02")

	query := `
	SELECT 
		b.id, b.name, b.sowing_date, b.estimated_harvest_days,
		IFNULL(m.name, '') as microgreen_name
	FROM batches b
	LEFT JOIN microgreens m ON b.microgreen_id = m.id
	WHERE b.user_id = ?
	AND b.id NOT IN (
		SELECT o.batch_id
		FROM observations o
		WHERE DATE(o.date) = ?
	)
	`

	rows, err := h.DB.Query(query, userID, today)
	if err != nil {
		http.Error(w, "DB query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type BatchToday struct {
		BatchID    int     `json:"batch_id"`
		Name       string  `json:"name"`
		Microgreen string  `json:"microgreen"`
		SowingDate string  `json:"sowing_date"`
		HarvestIn  int     `json:"estimated_harvest_days"`
		Progress   float64 `json:"progress_percent"`
		PassedDays int     `json:"passed_days"`
		Remaining  int     `json:"remaining_days"`
	}

	var result []BatchToday
	todayTime := time.Now().Truncate(24 * time.Hour)

	for rows.Next() {
		var id, harvestDays int
		var name, microgreen string
		var sowingDate time.Time

		if err := rows.Scan(&id, &name, &sowingDate, &harvestDays, &microgreen); err != nil {
			http.Error(w, "Row scan error", http.StatusInternalServerError)
			return
		}

		passed := int(todayTime.Sub(sowingDate).Hours() / 24)
		remaining := harvestDays - passed
		if remaining < 0 {
			remaining = 0
		}
		progress := float64(passed) / float64(harvestDays) * 100
		if progress > 100 {
			progress = 100
		}

		result = append(result, BatchToday{
			BatchID:    id,
			Name:       name,
			Microgreen: microgreen,
			SowingDate: sowingDate.Format("2006-01-02"),
			HarvestIn:  harvestDays,
			PassedDays: passed,
			Remaining:  remaining,
			Progress:   progress,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
