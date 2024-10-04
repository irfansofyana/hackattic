package main

import (
	"encoding/base64"
	"fmt"
	"hackatticgo"
	"log"
	"os"
	"os/exec"
)

const (
	problemName = "collision_course"
)

type GetProblem struct {
	Include string `json:"include"`
}

type SubmitSolution struct {
	Files []string `json:"files"`
}

func getCommonString() (string, error) {
	problem, err := hackatticgo.GetProblem[GetProblem](problemName)
	if err != nil {
		return "", err
	}

	return problem.Include, nil
}

func getBase64FromFile(filePath string) (string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(fileBytes)

	return base64String, nil
}

// reference: https://github.com/brimstone/fastcoll
func main() {
	commonString, err := getCommonString()
	if err != nil {
		log.Fatal("error get common string", err.Error())
	}

	err = os.WriteFile("input", []byte(commonString), 0o644)
	if err != nil {
		log.Fatal("error writing input file:", err)
	}

	if _, err := exec.LookPath("docker"); err != nil {
		log.Fatal("Docker is not installed or not in PATH:", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("error getting current working directory:", err)
	}

	cmd := exec.Command("docker", "run", "--rm", "-v", pwd+":/work", "-w", "/work", "-u", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()), "brimstone/fastcoll", "--prefixfile", "input", "-o", "out1.bin", "out2.bin")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("error executing Docker command: %v\nOutput: %s", err, output)
	}

	base64Out1, err := getBase64FromFile("out1.bin")
	if err != nil {
		log.Fatal("error reading out1.bin:", err)
	}

	base64Out2, err := getBase64FromFile("out2.bin")
	if err != nil {
		log.Fatal("error reading out2.bin:", err)
	}

	solution := SubmitSolution{
		Files: []string{base64Out1, base64Out2},
	}

	if err := hackatticgo.SubmitAnswer[SubmitSolution](problemName, solution); err != nil {
		log.Fatalf("error submitting answer: %v", err)
	}
}
