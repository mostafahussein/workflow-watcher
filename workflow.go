package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type workflowEnvironment struct {
	repo            string
	repoOwner       string
	headSha         string
	baseBranch      string
	pollingInterval time.Duration
}

func newWorkflowEnvironment(repo string, repoOwner string, headSha string, baseBranch string, pollingIntervalInput string) (*workflowEnvironment, error) {

	pollingInterval, err := strconv.Atoi(pollingIntervalInput)
	if err != nil {
		fmt.Printf("error converting to int: %v\n", err)
		os.Exit(1)
	}
	duration := time.Duration(pollingInterval) * time.Second

	repoOwnerAndName := strings.Split(repo, "/")
	var repoName string
	if len(repoOwnerAndName) == 2 {
		repoName = repoOwnerAndName[1]
	} else {
		repoName = repoOwnerAndName[0]
	}

	return &workflowEnvironment{
		repo:            repoName,
		repoOwner:       repoOwner,
		headSha:         headSha,
		baseBranch:      baseBranch,
		pollingInterval: duration,
	}, nil
}
