package hackatticgo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const baseUrl = "https://hackattic.com/challenges"

var accessToken string

func init() {
	accessToken = os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatalf("ACCESS_TOKEN is required")
	}
}

func GetProblem[T any](problemName string) (T, error) {
	var problem T

	url := fmt.Sprintf("%s/%s/problem", baseUrl, problemName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return problem, err
	}

	q := req.URL.Query()
	q.Add("access_token", accessToken)
	req.URL.RawQuery = q.Encode()

	log.Println("[Request] Requesting to Get Problem:", problemName)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return problem, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return problem, err
	}

	if err := json.Unmarshal(resBody, &problem); err != nil {
		return problem, err
	}

	log.Println("[Response] Get Problem response:", problem)

	return problem, nil
}
