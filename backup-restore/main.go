package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/lib/pq"
	"hackatticgo"
	"io"
	"log"
	"os"
	"os/exec"
)

var accessToken string

const problemName string = "backup_restore"

const dbFile string = "db.sql"

const dbName string = "backup_restore"

const dbUser string = "irfanputra"

func init() {
	accessToken = os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatalf("ACCESS_TOKEN is required")
	}
}

type GetBackupRestoreProblem struct {
	Dump string `json:"dump"`
}

type SubmitBackupRestoreRequest struct {
	AliveSsns []string `json:"alive_ssns"`
}

func downloadDumpedDb() {
	problem, err := hackatticgo.GetProblem[GetBackupRestoreProblem](problemName)
	if err != nil {
		log.Fatal("Error get problem", err.Error())
	}

	decodedData, err := base64.StdEncoding.DecodeString(problem.Dump)
	if err != nil {
		log.Fatal("Error decode the string", err.Error())
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(decodedData))
	if err != nil {
		log.Fatal("Error creating gzip reader: ", err.Error())
	}
	defer gzipReader.Close()

	outputFile, err := os.Create(dbFile)
	if err != nil {
		log.Fatal("Error when creating output file: ", err.Error())
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, gzipReader)
	if err != nil {
		log.Fatal("Error copying decompressed data: ", err.Error())
	}
}

func restoreDb() {
	cmd := exec.Command("psql", "-d", dbName, "-f", dbFile)
	if err := cmd.Run(); err != nil {
		log.Fatal("Error when restore the db", err.Error())
	}
}

func getAliveSsns() []string {
	ssns := make([]string, 0)

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName))

	if err != nil {
		log.Fatal("Error connect to db", err.Error())
	}

	defer db.Close()

	rows, err := db.Query(`SELECT ssn FROM criminal_records WHERE status='alive';`)
	if err != nil {
		log.Fatal("Error query to db", err.Error())
	}
	defer rows.Close()

	var ssn string
	for rows.Next() {
		err := rows.Scan(&ssn)
		if err != nil {
			log.Fatal(err)
		}
		ssns = append(ssns, ssn)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal("Error get data", err.Error())
	}

	return ssns
}

func submitSolution(aliveSsns []string) {
	answer := SubmitBackupRestoreRequest{
		AliveSsns: aliveSsns,
	}
	if err := hackatticgo.SubmitAnswer[SubmitBackupRestoreRequest](problemName, answer); err != nil {
		log.Fatal("Error submitting answer", err.Error())
	}
}

func main() {
	downloadDumpedDb()

	restoreDb()

	aliveSssns := getAliveSsns()

	submitSolution(aliveSssns)
}
