package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	// We need to import all used modules for Go to compile
	// Are we missing one?
)

func TestBucket(t *testing.T) {
	t.Parallel()
	workingDir := test_structure.CopyTerraformFolderToTemp(t, "..", "bucket")

	defer test_structure.RunTestStage(t, "destroy_infrastructure", func() {
		DestroyInfrastructure(t, workingDir)
	})

	test_structure.RunTestStage(t, "deploy_infrastructure", func() {
		DeployInfrastructure(t, workingDir)
	})

	test_structure.RunTestStage(t, "s3_bucket_check", func() {
		S3BucketCheck(t, workingDir)
	})
}

// Stage 1: DeployInfrastructure
// This function copies the root level terraform module to an ephemeral directory
// which becomes the working directory to deploy all Terraform resources from.
// The function then deploys the infrastructure.
func DeployInfrastructure(t *testing.T, workingDir string) {
	randomID := strings.ToLower(random.UniqueId())
	name := fmt.Sprintf("terratest-%s", randomID)
	awsRegion := aws.GetRandomStableRegion(t, []string{"us-east-2", "us-west-1", "us-west-2", "eu-west-1"}, nil)

	testName := t.Name()
	bucketName := fmt.Sprintf("lunch-n-learn-terratest/%s/%s/%s.tfstate", testName, awsRegion, name)
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workingDir,
		BackendConfig: map[string]interface{}{
			"key": bucketName,
		},
		Vars: map[string]interface{}{
			"name":   name,
			"region": awsRegion,
		},
		NoColor: true,
	})
	// Save our Terraform Options struct to use in another stage
	test_structure.SaveTerraformOptions(t, workingDir, terraformOptions)

	fmt.Println("Running 'terraform init' and 'terraform apply'...")
	terraform.InitAndApply(t, terraformOptions)
	fmt.Println("Terraform apply complete.")
}

// Stage 2: S3BucketCheck
// This function grabs the S3 Bucket ID from the module output
// and tests to see if its a non empty string which determines if the
// S3 bucket exists.
func S3BucketCheck(t *testing.T, workingDir string) {
	// Load Terraform Options struct from Stage 1
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)
	s3BucketId := terraform.Output(t, terraformOptions, "aws_s3_bucket_id")
	assert.NotEmpty(t, s3BucketId, "S3 bucket ID should not be empty")
}

// Stage 3: DestroyInfrastructure
// This function destroys all of the infrastructure that was deployed to conduct
// the infrastrcture unit tests.
func DestroyInfrastructure(t *testing.T, workingDir string) {
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)
	fmt.Println("Destroying Terraform resources...")
	terraform.Destroy(t, terraformOptions)
	fmt.Println("Terraform destroy complete.")
}
