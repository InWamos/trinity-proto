package v1

import (
	"github.com/InWamos/trinity-proto/internal/record/presentation/v1/handlers"
	"github.com/go-chi/chi/v5"
)

type RecordMuxV1 struct {
	mux *chi.Mux
}

func NewRecordMuxV1(
	getLatestTelegramRecordsByTelegramID *handlers.GetLatestTelegramRecordsByTelegramIDHandler,
) *RecordMuxV1 {
	mux := chi.NewRouter()
	mux.Get("/telegram/{telegram_id}/records", getLatestTelegramRecordsByTelegramID.ServeHTTP)
	return &RecordMuxV1{
		mux: mux,
	}
}

func (um *RecordMuxV1) GetMux() *chi.Mux {
	return um.mux
}
