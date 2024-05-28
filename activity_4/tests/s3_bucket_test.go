package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func TestS3Bucket(t *testing.T) {
	// Generate a 6-character random string
	randomID := strings.ToLower(random.UniqueId())

	// Use the random ID and terratest prefix to generate a random name
	name := fmt.Sprintf("terratest-%s", randomID)

	// Use the CopyTerraformFolderToTemp function to generate a randomly
	// named directory for holding the root-level module/state.
	//
	// Enables us to deploy multiple clusters from the same root-level
	// configuration with different parameters in parallel.
	workingDir := test_structure.CopyTerraformFolderToTemp(t, "..", "examples/complete")

	// t.Name() returns the name of the currently running test
	testName := t.Name()
	// Interpolate S3 bucket name for BackendConfig
	bucketName := fmt.Sprintf("terratest-lunch-n-learn/%s/us-west-1/%s.tfstate", testName, name)

	// Full terraform.Options struct:
	// https://github.com/gruntwork-io/terratest/blob/64a1856f2695fe1c24658fe8fc66090e83c7a530/modules/terraform/options.go#L39-L74
	//
	// Options struct for running Terraform commands
	options := terraform.Options{
		TerraformDir: workingDir,
		BackendConfig: map[string]interface{}{
			"key": bucketName,
		},
		Vars: map[string]interface{}{
			"name": name,
		},
		NoColor: true,
	}

	// WithDefaultRetryableErrors function source:
	// https://github.com/gruntwork-io/terratest/blob/64a1856f2695fe1c24658fe8fc66090e83c7a530/modules/terraform/options.go#L110-L127
	//
	// This function makes a copy of the Options object and returns an updated object with sensible defaults
	// for retryable errors. The included retryable errors are typical errors that most terraform modules encounter during
	// testing, and are known to self resolve upon retrying.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &options)

	// This will run `terraform destroy` when function is about to exit
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply`
	terraform.InitAndApply(t, terraformOptions)

	// This will run `terraform output <VALUE>` to get the value of the specified output
	s3BucketId := terraform.Output(t, terraformOptions, "s3_bucket_id")

	// Assert that the ARN is not an empty string
	assert.NotEmpty(t, s3BucketId, "S3 bucket ARN should not be empty")
}
