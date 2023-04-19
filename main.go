package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/google/go-github/v51/github"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
)

func newWorkflowStatusLoopChannel(ctx context.Context, wrkflw *workflowEnvironment, client *github.Client) chan int {
	channel := make(chan int)
	go func() {
		for {
			workflowRuns, _, err := client.Actions.ListRepositoryWorkflowRuns(ctx, wrkflw.repoOwner, wrkflw.repo, &github.ListWorkflowRunsOptions{
				HeadSHA: wrkflw.headSha,
			})
			if err != nil {
				fmt.Printf("error listing workflow for a repository: %v\n", err)
				channel <- 1
				close(channel)
			}

			if workflowRuns.GetTotalCount() == 0 {
				fmt.Printf("No workflows found for the specified commit sha")
				channel <- 0
				close(channel)
			}
			parsedOutput, err := json.Marshal(workflowRuns)
			if err != nil {
				fmt.Printf("error parsing output to json: %v\n", err)
				channel <- 1
				close(channel)
			}
			headBranchWorkflowExists := gjson.Get(string(parsedOutput), "workflow_runs.#(head_branch==\""+wrkflw.baseBranch+"\")").Exists()
			if headBranchWorkflowExists {
				headBranchWorkflow := gjson.Get(string(parsedOutput), "workflow_runs.#(head_branch==\""+wrkflw.baseBranch+"\")")
				headBranchWorkflowStatus := gjson.Get(headBranchWorkflow.String(), "status").String()
				fmt.Printf("Base branch workflow status: %s\n", headBranchWorkflowStatus)
				switch headBranchWorkflowStatus {
				case string(workflowStatusCompleted):
					fmt.Println("Base branch workflow status is completed, verifying base branch workflow conclusion...")
					headBranchWorkflowConclusion := gjson.Get(headBranchWorkflow.String(), "conclusion").String()
					switch headBranchWorkflowConclusion {
					case string(workflowConclusionSuccess):
						fmt.Printf("Base branch workflow status is success")
						channel <- 0
						close(channel)
					case string(workflowConclusionFailed):
						fmt.Printf("Base branch workflow status is failed")
						channel <- 1
						close(channel)
					case string(workflowConclusionSkipped):
						fmt.Printf("Base branch workflow status is skipped")
						channel <- 1
						close(channel)
					case string(workflowConclusionTimeOut):
						fmt.Printf("Base branch workflow status is timeout")
						channel <- 1
						close(channel)
					case string(workflowConclusionCancelled):
						fmt.Printf("Base branch workflow status is cancelled")
						channel <- 1
						close(channel)
					}
				case string(workflowStatusFailed):
					fmt.Printf("Base branch workflow status is failed")
					channel <- 1
					close(channel)
				}

				time.Sleep(wrkflw.pollingInterval)
			}
		}
	}()
	return channel
}

func newGithubClient(ctx context.Context) (*github.Client, error) {
	token := os.Getenv(envVarToken)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	serverUrl, serverUrlPresent := os.LookupEnv("GITHUB_SERVER_URL")
	apiUrl, apiUrlPresent := os.LookupEnv("GITHUB_API_URL")

	if serverUrlPresent {
		if !apiUrlPresent {
			apiUrl = serverUrl
		}
		return github.NewEnterpriseClient(apiUrl, serverUrl, tc)
	}
	return github.NewClient(tc), nil
}

func main() {
	if err := validateInput(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	repo := os.Getenv(envVarRepoName)
	repoOwner := os.Getenv(envVarRepoOwner)

	ctx := context.Background()
	client, err := newGithubClient(ctx)
	if err != nil {
		fmt.Printf("error connecting to server: %v\n", err)
		os.Exit(1)
	}

	headSha := os.Getenv(envVarHeadSha)
	baseBranch := os.Getenv(envVarBaseBranch)
	pollingInterval := os.Getenv(envVarPollingInterval)

	wrkflw, err := newWorkflowEnvironment(repo, repoOwner, headSha, baseBranch, pollingInterval)
	if err != nil {
		fmt.Printf("error creating workflow environment: %v\n", err)
		os.Exit(1)
	}

	killSignalChannel := make(chan os.Signal, 1)
	signal.Notify(killSignalChannel, os.Interrupt)

	workflowStatusLoopChannel := newWorkflowStatusLoopChannel(ctx, wrkflw, client)

	select {
	case exitCode := <-workflowStatusLoopChannel:
		os.Exit(exitCode)
	case <-killSignalChannel:
		handleInterrupt(ctx)
		os.Exit(1)
	}
}
