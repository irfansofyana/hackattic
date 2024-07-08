package hackatticgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func SubmitAnswer[T any](problemName string, submission T) error {
	body, err := json.Marshal(submission)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s/solve?access_token=%s", baseUrl, problemName, accessToken)

	log.Println("[Request] Submitting answer:", submission)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Println("[Response] Submit answer response:", string(resBody))

	return nil
}
