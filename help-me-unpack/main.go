package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
)

var accessToken string

const baseUrl = "https://hackattic.com/challenges/help_me_unpack"

func init() {
	accessToken = os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatalf("ACCESS_TOKEN is required")
	}
}

type GetBytesRequest struct {
	Bytes string `json:"bytes"`
}

func getBytesString() (string, error) {
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

	var data GetBytesRequest
	if err := json.Unmarshal(resBody, &data); err != nil {
		return "", err
	}

	return data.Bytes, nil
}

type SubmitAnswerRequest struct {
	Myint               int32   `json:"int"`
	Myuint              uint32  `json:"uint"`
	Myshort             int16   `json:"short"`
	Myfloat             float32 `json:"float"`
	Mydouble            float64 `json:"double"`
	Mybig_endian_double float64 `json:"big_endian_double"`
}

func unpack(requestBytes []byte) SubmitAnswerRequest {
	// Little-endian signed integer
	myint := int32(binary.LittleEndian.Uint32(requestBytes[:4]))
	// Little-endian unsigned integer
	uint := binary.LittleEndian.Uint32(requestBytes[4:8])
	// Little-endian signed short
	short := int16(binary.LittleEndian.Uint16(requestBytes[8:10]))
	// Little-endian float
	myfloat := math.Float32frombits(binary.LittleEndian.Uint32(requestBytes[12:16]))
	// Little-endian double
	double := math.Float64frombits(binary.LittleEndian.Uint64(requestBytes[16:24]))
	// Big-endian double
	bigeDouble := math.Float64frombits(binary.BigEndian.Uint64(requestBytes[24:32]))

	return SubmitAnswerRequest{
		Myint:               myint,
		Myuint:              uint,
		Myshort:             short,
		Myfloat:             myfloat,
		Mydouble:            double,
		Mybig_endian_double: bigeDouble,
	}
}

func submitAnswer(reqBody SubmitAnswerRequest) error {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/solve?access_token=%s", baseUrl, accessToken)

	log.Println("Submitting answer: ", body, string(body))

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Println("Getting submit answer response", string(resBody))

	return nil
}

func main() {
	bytesString, err := getBytesString()
	if err != nil {
		log.Fatal("Error get Bytes string", err.Error())
	}

	requestBytes, err := base64.StdEncoding.DecodeString(bytesString)
	if err != nil {
		log.Fatal("Error decode", err.Error())
	}

	submitReq := unpack(requestBytes)
	if err = submitAnswer(submitReq); err != nil {
		log.Fatal("Error submit request", err.Error())
	}
}
