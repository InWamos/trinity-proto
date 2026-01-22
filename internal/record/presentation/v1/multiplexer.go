package v1

import "github.com/go-chi/chi/v5"

type RecordMuxV1 struct {
	mux *chi.Mux
}

func NewRecordMuxV1() *RecordMuxV1 {
	mux := chi.NewRouter()
	return &RecordMuxV1{
		mux: mux,
	}
}

func (um *RecordMuxV1) GetMux() *chi.Mux {
	return um.mux
}
