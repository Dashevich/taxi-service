package handlers

import (
	"encoding/json"
	"net/http"
	"offering_service/internal/model"
)

func (handler *Handler) ParseOffer(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("неверный запрос"))
		return
	}
	offerId := request.URL.Query().Get("offer_id")
	data, err := handler.jwtService.ExtractClaims(offerId)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ans, err := json.Marshal(model.Response{
		OfferId:  offerId,
		ClientId: data.ClientId,
		From:     data.From,
		To:       data.To,
		Price:    data.Price,
	})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(ans)
}
