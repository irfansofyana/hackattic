package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/liyue201/goqr"
)

var accessToken string

const baseUrl = "https://hackattic.com/challenges/reading_qr"

func init() {
	accessToken = os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatalf("ACCESS_TOKEN is required")
	}
}

type GetQrImageResponse struct {
	Url string `json:"image_url"`
}

type SubmitAnswerRequest struct {
	Code string `json:"code"`
}

func getQrImageUrl() (string, error) {
	url := fmt.Sprintf("%s/%s", baseUrl, "problem")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("access_token", accessToken)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	var getQrResponse GetQrImageResponse
	if err := json.Unmarshal(resBody, &getQrResponse); err != nil {
		return "", err
	}

	return getQrResponse.Url, nil
}

func getQr(qrImageUrl string) (string, error) {
	res, err := http.Get(qrImageUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	img, _, err := image.Decode(bytes.NewReader(resBody))
	if err != nil {
		return "", err
	}

	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		return "", err
	}

	if len(qrCodes) == 0 {
		return "", fmt.Errorf("no QR code found")
	}

	return string(qrCodes[0].Payload), nil
}

func submitAnswer(qrString string) error {
	reqBody := SubmitAnswerRequest{
		Code: qrString,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/solve?access_token=%s", baseUrl, accessToken)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Println("Submit answer response", string(resBody))

	return nil
}

func main() {
	url, err := getQrImageUrl()
	if err != nil {
		log.Fatal("Error get QR Image URL:", err.Error())
	}

	qrString, err := getQr(url)
	if err != nil {
		log.Fatal("error get QR string:", err.Error())
	}

	log.Println("QR String:", qrString)

	if err := submitAnswer(qrString); err != nil {
		log.Fatal("Error submitting answer:", err.Error())
	}

	log.Println("Successfully run")
}
