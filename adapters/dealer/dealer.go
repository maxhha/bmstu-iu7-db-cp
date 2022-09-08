package dealer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"encoding/json"
)

type DealerAdapter struct {
	baseUrl string
}

func New(url string) DealerAdapter {
	return DealerAdapter{
		url,
	}
}

func (d *DealerAdapter) request(url string, data []byte) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", d.baseUrl, url), bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("response(status=%s): ioutil.ReadAll: %w", response.Status, err)
	}

	return fmt.Errorf("response(status=%s): %s", response.Status, string(body))
}

func (d *DealerAdapter) OwnerAccept(offerId string, userId *string, comment *string) error {
	data, err := json.Marshal(struct {
		OfferId string  `json:"offerId"`
		UserId  *string `json:"userId"`
		Comment *string `json:"comment"`
	}{
		OfferId: offerId,
		UserId:  userId,
		Comment: comment,
	})

	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	return d.request("owner-accept", data)
}

func (d *DealerAdapter) BuyerAccept(offerId string, userId *string) error {
	data, err := json.Marshal(struct {
		OfferId string  `json:"offerId"`
		UserId  *string `json:"userId"`
		Comment *string `json:"comment"`
	}{
		OfferId: offerId,
		UserId:  userId,
	})

	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	return d.request("consumer-accept", data)
}
