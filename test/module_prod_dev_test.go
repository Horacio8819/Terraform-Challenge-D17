package test

import (
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestWebserverCluster(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/services/webserver-cluster",
		Vars: map[string]interface{}{
			"cluster_name":            "webservers-production",
			"environment":             "production",
			"create_dns_record":       false,
			"aws_region":              "eu-central-1",
			"project_name":            "devops-project",
			"team_name":               "devOps",
			"alert_email":             "horace.djousse@yahoo.com",
			"server_template_version": "latest",

			"public_subnets": map[string]interface{}{
				"a": "172.31.108.0/24",
				"b": "172.31.109.0/24",
				"c": "172.31.110.0/24",
			},
		},
	})

	//defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	albDnsName := terraform.Output(t, terraformOptions, "alb_dns_name")
	url := "http://" + albDnsName

	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		url,
		nil,
		30,
		10*time.Second,
		func(status int, body string) bool {
			return status == 200 && len(body) > 0
		},
	)
}
