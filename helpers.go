package main

import (
	"context"
	"fmt"
	"os"
)

func handleInterrupt(ctx context.Context) {
	closeComment := "Workflow cancelled."
	fmt.Println(closeComment)
}

func validateInput() error {
	missingEnvVars := []string{}
	if os.Getenv(envVarRepoName) == "" {
		missingEnvVars = append(missingEnvVars, envVarRepoName)
	}

	if os.Getenv(envVarHeadSha) == "" {
		missingEnvVars = append(missingEnvVars, envVarHeadSha)
	}

	if os.Getenv(envVarBaseBranch) == "" {
		missingEnvVars = append(missingEnvVars, envVarBaseBranch)
	}

	if os.Getenv(envVarRepoOwner) == "" {
		missingEnvVars = append(missingEnvVars, envVarRepoOwner)
	}

	if os.Getenv(envVarPollingInterval) == "" {
		missingEnvVars = append(missingEnvVars, envVarPollingInterval)
	}

	if os.Getenv(envVarToken) == "" {
		missingEnvVars = append(missingEnvVars, envVarToken)
	}

	if len(missingEnvVars) > 0 {
		return fmt.Errorf("missing env vars: %v", missingEnvVars)
	}
	return nil
}
