package receipt

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/tagptroll1/receipt-processor/lib/jsonlib"
)

type IdResponse struct {
	Id string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"point"`
}

func ProcessReceipts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var rec Receipt
	err := jsonlib.UnmarshalReader(r.Body, &rec)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	id, errs := SaveNewReceipt(ctx, rec)
	if len(errs) > 0 {
		sb := strings.Builder{}
		for _, err := range errs {
			sb.WriteString(err.Error())
			sb.WriteString("\n")
		}
		_, _ = w.Write([]byte(sb.String()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(IdResponse{Id: id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	writeLogErr(w, jsonResponse)
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	rec, err := GetReceipt(ctx, id)
	switch {
	case errors.Is(err, ErrReceiptNotFound):
		writeStringLogErr(w, "receipt not found")
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	points, err := CalculatePoints(ctx, rec)
	if err != nil {
		writeStringLogErr(w, "failed to calculate points")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(PointsResponse{Points: points})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeLogErr(w, jsonResponse)
}

func writeLogErr(w http.ResponseWriter, b []byte) {
	if _, err := w.Write(b); err != nil {
		log.Printf("failed to write: %s\n", b)
	}
}

func writeStringLogErr(w http.ResponseWriter, s string) {
	if _, err := w.Write([]byte(s)); err != nil {
		log.Printf("failed to write: %s\n", s)
	}
}
