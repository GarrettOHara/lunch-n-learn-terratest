package test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func TestWebServer(t *testing.T) {
	// Allow test to run in parrallel with other tests
	t.Parallel()
	// Generate a 6-character random string
	randomID := strings.ToLower(random.UniqueId())
	// Use the random ID and terratest prefix to generate a random name
	name := fmt.Sprintf("terratest-%s", randomID)
    // Get random AWS region
    awsRegion := aws.GetRandomStableRegion(t, []string{"us-west-1", "us-west-2", "eu-west-1"}, nil)

	// Use the CopyTerraformFolderToTemp function to generate a randomly
	// named directory for holding the root-level module/state.
	//
	// Enables us to deploy multiple clusters from the same root-level
	// configuration with different parameters in parallel.
	workingDir := test_structure.CopyTerraformFolderToTemp(t, "..", "examples/web_server")

	testName := t.Name()
	bucketName := fmt.Sprintf("lunch-n-learn-terratest/%s/%s/%s.tfstate", testName, awsRegion, name)
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workingDir,
		BackendConfig: map[string]interface{}{
			"key": bucketName,
		},
		Vars: map[string]interface{}{
			"name": name,
            "region": awsRegion,
		},
		NoColor: true,
	})
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	instance_id := terraform.Output(t, terraformOptions, "instance_id")
	public_ipv4 := terraform.Output(t, terraformOptions, "public_ipv4_addr")
	public_dns := terraform.Output(t, terraformOptions, "public_dns")

	// Ensure outputs are all non-empty
	assert.NotEmpty(t, instance_id, "instance_id should not be empty")
	assert.NotEmpty(t, public_ipv4, "public_ipv4 should not be empty")
	assert.NotEmpty(t, public_dns, "public_dns should not be empty")

	fmt.Println("Waiting 2 minutes for web server to come online...")
	time.Sleep(120 * time.Second)

	// Send GET request to API
	resp, err := http.Get("http://" + public_dns)
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

	serverResponse := strings.ReplaceAll(string(body), "\n", "")
	serverResponse = strings.ReplaceAll(serverResponse, " ", "")
    fmt.Printf("Response status code: %d\n", resp.StatusCode)
	fmt.Printf("Response body:\n%s\n", serverResponse)

    // Test server response
    assert.Equal(t, serverResponse, "OK", "The server sent an unexpected response.")
}
