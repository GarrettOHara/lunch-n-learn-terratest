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

func TestWebServer(t *testing.T) {
	// Generate a 6-character random string
	randomID := random.UniqueId()
	// Use the random ID and terratest prefix to generate a random name
	name := fmt.Sprintf("terratest-%s", randomID)

	// Use the CopyTerraformFolderToTemp function to generate a randomly
	// named directory for holding the root-level module/state.
	//
	// Enables us to deploy multiple clusters from the same root-level
	// configuration with different parameters in parallel.
	workingDir := test_structure.CopyTerraformFolderToTemp(t, "..", "examples/web_server")

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

	instance_id := terraform.Output(t, terraformOptions, "instance_id")
	public_ipv4 := terraform.Output(t, terraformOptions, "public_ipv4_addr")
	public_dns := terraform.Output(t, terraformOptions, "public_dns")

    // Ensure outputs are all non-empty
	assert.NotEmpty(t, instance_id, "instance_id should not be empty")
	assert.NotEmpty(t, public_ipv4, "public_ipv4 should not be empty")
	assert.NotEmpty(t, public_dns, "public_dns should not be empty")
}
