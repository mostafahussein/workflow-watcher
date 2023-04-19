package main

type workflowStatus string

const (
	envVarRepoName        string = "INPUT_REPOSITORY-NAME"
	envVarHeadSha         string = "INPUT_HEAD-SHA"
	envVarBaseBranch      string = "INPUT_BASE-BRANCH"
	envVarRepoOwner       string = "INPUT_REPOSITORY-OWNER"
	envVarPollingInterval string = "INPUT_POLLING-INTERVAL"
	envVarToken           string = "INPUT_SECRET"

	workflowStatusFailed        workflowStatus = "failure"
	workflowStatusCompleted     workflowStatus = "completed"
	workflowConclusionSuccess   workflowStatus = "success"
	workflowConclusionFailed    workflowStatus = "failure"
	workflowConclusionCancelled workflowStatus = "cancelled"
	workflowConclusionSkipped   workflowStatus = "skipped"
	workflowConclusionTimeOut   workflowStatus = "timed_out"
)
