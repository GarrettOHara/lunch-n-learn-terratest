package test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/aws"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkEc2Instance(awsRegion string, instanceId string) string {
	fmt.Println("Inside checkEc2Instance")
	session, err := aws.NewAuthenticatedSession(awsRegion)
	check(err)
	client := ec2.New(session)

	request := ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{&instanceId},
	}
	result, err := client.DescribeInstanceStatus(&request)
	check(err)
	for _, instanceStatus := range result.InstanceStatuses {
		fmt.Println("Instance ID:", *instanceStatus.InstanceId)
		fmt.Println("Instance Status:", *instanceStatus.InstanceState.Name)
		return *instanceStatus.InstanceState.Name
	}
	return ""
}

func TestWebServer(t *testing.T) {
	// Allow test to run in parrallel with other tests
	t.Parallel()
	// Generate a 6-character random string
	randomID := strings.ToLower(random.UniqueId())
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

	// Wait for EC2 instance to become available
	// for {
	//     status := checkEc2Instance("us-west-1", instance_id)
	//     if status == "running" {
	//         break
	//     } else {
	//         time.Sleep(10 * time.Second)
	//         fmt.Println("Waiting for web server to come online...")
	//     }
	// }
	fmt.Println("Waiting 5 minutes for web server to come online...")
	time.Sleep(300 * time.Second)

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

	// Print response body
	fmt.Println("Response body:")
	fmt.Println(string(body))
}
