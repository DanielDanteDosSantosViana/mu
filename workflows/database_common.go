package workflows

import (
	"errors"

	"github.com/stelligent/mu/common"
)

type databaseWorkflow struct {
	serviceName           string
	codeRevision          string
	appRevisionBucket     string
	repoName              string
	cloudFormationRoleArn string
	databaseKeyArn        string
}

func (workflow *databaseWorkflow) databaseInput(ctx *common.Context, serviceName string, environmentName string) Executor {
	return func() error {
		// Repo Name
		if serviceName != "" {
			workflow.serviceName = serviceName
		} else if ctx.Config.Service.Name != "" {
			workflow.serviceName = ctx.Config.Service.Name
		} else if ctx.Config.Repo.Name != "" {
			workflow.serviceName = ctx.Config.Repo.Name
		} else {
			return errors.New("Service name must be provided")
		}

		workflow.appRevisionBucket = ctx.Config.Service.Pipeline.Build.Bucket

		return nil
	}
}
