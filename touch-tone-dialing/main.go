package main

import (
	"io"
	"log"
	"net/http"

	"hackatticgo"

	"github.com/Hallicopter/go-dtmf/dtmf"
)

const (
	problemName = "touch_tone_dialing"
)

type GetProblem struct {
	WavURL string `json:"wav_url"`
}

type SubmitSolution struct {
	Sequence string `json:"sequence"`
}

func grabWavContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
	problem, err := hackatticgo.GetProblem[GetProblem](problemName)
	if err != nil {
		log.Fatalf("Failed to get problem: %v", err)
	}

	content, err := grabWavContent(problem.WavURL)
	if err != nil {
		log.Fatalf("Failed to grab WAV content: %v", err)
	}

	sequence, err := dtmf.DecodeDTMFFromBytes([]byte(content), 4096, 5)
	if err != nil {
		log.Fatalf("Failed to decode DTMF: %v", err)
	}

	log.Printf("Sequence: %s", sequence)

	submitSolution := SubmitSolution{
		Sequence: sequence,
	}

	if err := hackatticgo.SubmitAnswer(problemName, submitSolution); err != nil {
		log.Fatalf("Failed to submit solution: %v", err)
	}

	log.Println("Successfully submitted solution")
}
