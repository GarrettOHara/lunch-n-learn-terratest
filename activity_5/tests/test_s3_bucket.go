package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func TestS3Bucket(t *testing.T) {
    // Allow test to run in parrallel with other tests
    t.parrallel()

	// Generate a 6-character random string
	randomID := random.UniqueId()
	// Use the random ID and terratest prefix to generate a random name
	name := fmt.Sprintf("terratest-%s", randomID)

	// Use the CopyTerraformFolderToTemp function to generate a randomly
	// named directory for holding the root-level module/state.
	//
	// Enables us to deploy multiple clusters from the same root-level
	// configuration with different parameters in parallel.
	workingDir := test_structure.CopyTerraformFolderToTemp(t, "..", "examples/complete")

	testName := t.Name()
	bucketName := fmt.Sprintf("terratest-lunch-n-learn/%s/us-west-1/%s.tfstate", testName, name)
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workingDir,
		BackendConfig: map[string]interface{}{
			"key": bucketName,
		},
		Vars: map[string]interface{}{
			"name": name,
		},
		NoColor: true,
	})
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	s3BucketARN := terraform.Output(t, terraformOptions, "s3_bucket_arn")
	assert.NotEmpty(t, s3BucketARN, "S3 bucket ARN should not be empty")
}
