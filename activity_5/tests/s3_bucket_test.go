package test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func TestS3Bucket(t *testing.T) {
	// Allow test to run in parallel with other tests
	t.Parallel()

	// Use the CopyTerraformFolderToTemp function to generate a randomly
	// named directory for holding the root-level module/state.
	//
	// Enables us to deploy multiple clusters from the same root-level
	// configuration with different parameters in parallel.
	workingDir := test_structure.CopyTerraformFolderToTemp(t, "..", "examples/storage")

	defer test_structure.RunTestStage(t, "destroy_infrastructure", func() {
		DestroyInfrastructure(t, workingDir)
	})

	test_structure.RunTestStage(t, "deploy_infrastructure", func() {
		DeployInfrastructure(t, workingDir)
	})

	test_structure.RunTestStage(t, "s3_bucket_check", func() {
		S3BucketCheck(t, workingDir)
	})

	test_structure.RunTestStage(t, "static_website_response", func() {
		StaticWebsiteResponse(t, workingDir)
	})
}

// Stage 1: DeployInfrastructure
// This function copies the root level terraform module to an ephemeral directory
// which becomes the working directory to deploy all Terraform resources from.
// The function then deploys the infrastructure.
func DeployInfrastructure(t *testing.T, workingDir string) {
	// Generate a 6-character random string
	randomID := strings.ToLower(random.UniqueId())
	// Use the random ID and terratest prefix to generate a random name
	name := fmt.Sprintf("terratest-%s", randomID)
	// Get random AWS region
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
	// **IMPORTANT** Save the options and key for later test stages to use
	test_structure.SaveTerraformOptions(t, workingDir, terraformOptions)

	// run "terraform apply"
	fmt.Println("Running 'terraform init' and 'terraform apply'...")
	terraform.InitAndApply(t, terraformOptions)
	fmt.Println("Terraform apply complete.")
}

// Stage 2: S3BucketCheck
// This function grabs the S3 Bucket ID from the module output
// and tests to see if its a non empty string which determines if the
// S3 bucket exists.
func S3BucketCheck(t *testing.T, workingDir string) {
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)
	s3BucketId := terraform.Output(t, terraformOptions, "s3_bucket_id")
	assert.NotEmpty(t, s3BucketId, "S3 bucket ID should not be empty")
}

// Stage 3: StaticWebsiteResponse
// This function grabs the S3 Static Website Endpoint from the module output
// and contructs an HTTP GET request to see if the static content is reachable.
func StaticWebsiteResponse(t *testing.T, workingDir string) {
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)
	s3StaticSiteURL := terraform.Output(t, terraformOptions, "website_endpoint")
	assert.NotEmpty(t, s3StaticSiteURL, "S3 static website endpoint should not be empty")

	// Send GET request
	resp, err := http.Get("http://" + s3StaticSiteURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Printf("Server response code: %d\n", resp.StatusCode)
	fmt.Println("Response from server:")
	fmt.Println(string(body))
	serverResponse := strings.ReplaceAll(string(body), "\n", "")
	serverResponse = strings.ReplaceAll(serverResponse, " ", "")
	expectedHTML := `<html><body><p>OK</p></body></html>`

	assert.Equal(t, resp.StatusCode, http.StatusOK, "The response from the server was not 200.")
	assert.Equal(t, serverResponse, expectedHTML, "The response from the server did not match the expected result.")
}

// Stage 4: DestroyInfrastructure
// This function destroys all of the infrastructure that was deployed to conduct
// the infrastrcture unit tests.
func DestroyInfrastructure(t *testing.T, workingDir string) {
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)
	fmt.Println("Destroying Terraform resources...")
	terraform.Destroy(t, terraformOptions)
	fmt.Println("Terraform destroy complete.")
}
