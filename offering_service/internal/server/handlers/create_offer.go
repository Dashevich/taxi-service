package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"offering_service/internal/model"
)

func (handler *Handler) CreateOffer(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("неверный запрос"))
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		err := request.Body.Close()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}()
	var requestOffer model.Request
	err = json.Unmarshal(body, &requestOffer)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("неверный запрос"))
		return
	}
	offerData := handler.offeringService.GetOffer(requestOffer)
	id, err := handler.jwtService.CreateToken(offerData)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ans, err := json.Marshal(model.Response{
		OfferId:  id,
		ClientId: offerData.ClientId,
		From:     offerData.From,
		To:       offerData.To,
		Price:    offerData.Price,
	})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(ans)
}
