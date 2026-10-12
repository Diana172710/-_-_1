package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/Diana172710/-_-_1/mint/internal/service"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func (h *Handler) statsPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "stats", 200, pageData{Title: "Статистика", Year: time.Now().Year()})
}
func (h *Handler) byCategory(w http.ResponseWriter, r *http.Request) {
	f, err := service.ParseFilter(r.URL.Query())
	if err != nil {
		h.fail(w, r, 400, err.Error())
		return
	}
	if f.CategoryID != 0 {
		cats, ok := h.categories(w, r)
		if !ok {
			return
		}
		if !hasCategory(cats, f.CategoryID) {
			h.fail(w, r, 400, "Категория не найдена.")
			return
		}
	}
	data, err := h.repo.ByCategory(r.Context(), currentUser(r).ID, f)
	if err != nil {
		h.dbError(w, r, err)
		return
	}
	h.json(w, data)
}
func (h *Handler) byMonth(w http.ResponseWriter, r *http.Request) {
	year := time.Now().Year()
	if v := r.URL.Query().Get("year"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 9999 {
			h.fail(w, r, 400, "Год должен быть от 1 до 9999.")
			return
		}
		year = n
	}
	data, err := h.repo.ByMonth(r.Context(), currentUser(r).ID, year)
	if err != nil {
		h.dbError(w, r, err)
		return
	}
	h.json(w, data)
}
func (h *Handler) json(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

// Prevent spreadsheet software from interpreting user text as a formula.
func csvText(s string) string {
	trimmed := strings.TrimLeftFunc(s, unicode.IsSpace)
	if trimmed != "" && strings.ContainsRune("=+-@", rune(trimmed[0])) {
		return "'" + s
	}
	if strings.ContainsAny(s, "\t\r\n") {
		return "'" + s
	}
	return s
}
func (h *Handler) export(w http.ResponseWriter, r *http.Request) {
	f, err := service.ParseFilter(r.URL.Query())
	if err != nil {
		h.fail(w, r, 400, err.Error())
		return
	}
	if f.CategoryID != 0 {
		cats, ok := h.categories(w, r)
		if !ok {
			return
		}
		if !hasCategory(cats, f.CategoryID) {
			h.fail(w, r, 400, "Категория не найдена.")
			return
		}
	}
	items, err := h.repo.GetAll(r.Context(), currentUser(r).ID, f)
	if err != nil {
		h.dbError(w, r, err)
		return
	}
	var b bytes.Buffer
	b.WriteString("\ufeff")
	cw := csv.NewWriter(&b)
	cw.Comma = ';'
	cw.UseCRLF = true
	_ = cw.Write([]string{"Дата", "Категория", "Описание", "Сумма (RUB)"})
	for _, e := range items {
		_ = cw.Write([]string{e.Date, csvText(e.Category), csvText(e.Description), money(e.Amount)})
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		h.dbError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="expenses-`+time.Now().Format("2006-01-02")+`.csv"`)
	_, _ = w.Write(b.Bytes())
}
