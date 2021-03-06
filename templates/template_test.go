package templates

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/gobuffalo/packr"
	"github.com/stelligent/mu/common"
	"github.com/stretchr/testify/assert"
)

func TestNewTemplate(t *testing.T) {
	assert := assert.New(t)

	templates := []string{common.TemplatePolicyDefault, common.TemplatePolicyAllowAll,
		common.TemplateApp, common.TemplateBucket, common.TemplateBuildspec,
		common.TemplateCommonIAM, common.TemplateDatabase, common.TemplateELB,
		common.TemplateEnvEC2, common.TemplateEnvECS, common.TemplatePipelineIAM,
		common.TemplatePipeline, common.TemplateRepo, common.TemplateSchedule,
		common.TemplateServiceEC2, common.TemplateServiceECS, common.TemplateServiceIAM,
		common.TemplateVPCTarget, common.TemplateVPC}

	for _, templateName := range templates {
		templateBody, err := GetAsset(templateName)

		assert.Nil(err)
		assert.NotNil(templateBody)
		assert.NotEmpty(templateBody)
	}
}

func TestNewTemplate_invalid(t *testing.T) {
	assert := assert.New(t)

	templateBodyReader, err := GetAsset("invalid-template-name.yml")
	assert.Empty(templateBodyReader)
	assert.NotNil(err)
}

func TestNewTemplate_assets(t *testing.T) {
	assert := assert.New(t)

	sessOptions := session.Options{SharedConfigState: session.SharedConfigEnable}
	sess, err := session.NewSessionWithOptions(sessOptions)
	assert.Nil(err)

	svc := cloudformation.New(sess)

	box := packr.NewBox("./assets")
	templates := box.List()
	assert.NotZero(len(templates))
	for _, templateName := range templates {
		if !strings.HasPrefix(templateName, "cloudformation/") {
			continue
		}

		templateBody, err := GetAsset(templateName, ExecuteTemplate(nil))

		assert.Nil(err, templateName)
		assert.NotNil(templateBody, templateName)

		if templateBody != "" {

			assert.NotEmpty(templateBody, templateName)

			params := &cloudformation.ValidateTemplateInput{
				TemplateBody: aws.String(templateBody),
			}

			_, err := svc.ValidateTemplate(params)
			if err != nil {
				if awsErr, ok := err.(awserr.Error); ok {
					if awsErr.Code() == "InvalidClientTokenId" && awsErr.Message() == "The security token included in the request is invalid." {
						return
					}
					assert.Fail(awsErr.Code(), awsErr.Message(), templateName)
				}
				assert.Fail(err.Error(), templateName)
			}
		}
	}
}
