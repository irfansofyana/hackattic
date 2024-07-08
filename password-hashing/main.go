package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
	"hackatticgo"
	"log"
)

const problemName string = "password_hashing"

type GetPasswordHashingProblem struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Pbkdf2   struct {
		Hash   string `json:"hash"`
		Rounds int    `json:"rounds"`
	} `json:"pbkdf2"`
	Scrypt struct {
		N       int    `json:"N"`
		R       int    `json:"r"`
		P       int    `json:"p"`
		Buflen  int    `json:"buflen"`
		Control string `json:"_control"`
	} `json:"scrypt"`
}

func computesScrypt(problem GetPasswordHashingProblem) string {
	dk, err := scrypt.Key(
		[]byte(problem.Password),
		[]byte(problem.Salt),
		problem.Scrypt.N, problem.Scrypt.R, problem.Scrypt.P, problem.Scrypt.Buflen)

	if err != nil {
		log.Fatal(err.Error())
	}

	return hex.EncodeToString(dk)
}

func computeSHA256(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedPasswd := hash.Sum(nil)
	return hex.EncodeToString(hashedPasswd)
}

func computeHMACSHA256(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func computePBKDF2(password, salt string, iterations, keyLength int) string {
	return hex.EncodeToString(pbkdf2.Key([]byte(password), []byte(salt), iterations, keyLength, sha256.New))
}

type SubmitAnswerRequest struct {
	Sha256 string `json:"sha256"`
	Hmac   string `json:"hmac"`
	Pbkdf2 string `json:"pbkdf2"`
	Scrypt string `json:"scrypt"`
}

func main() {
	problem, err := hackatticgo.GetProblem[GetPasswordHashingProblem](problemName)
	if err != nil {
		log.Fatal("Error get problem", err.Error())
	}

	d, _ := base64.StdEncoding.DecodeString(problem.Salt)
	problem.Salt = string(d)

	answer := SubmitAnswerRequest{
		Sha256: computeSHA256(problem.Password),
		Hmac:   computeHMACSHA256(problem.Salt, problem.Password),
		Pbkdf2: computePBKDF2(problem.Password, problem.Salt, problem.Pbkdf2.Rounds, 32),
		Scrypt: computesScrypt(problem),
	}

	log.Println(answer)

	if err := hackatticgo.SubmitAnswer[SubmitAnswerRequest](problemName, answer); err != nil {
		log.Fatal("Error submitting answer", err.Error())
	}

	log.Println("Done")
}
