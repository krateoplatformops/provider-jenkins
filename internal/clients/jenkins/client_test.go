//go:build integration
// +build integration

package jenkins

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	httphelper "github.com/krateoplatformops/provider-jenkins/internal/helpers/http"
	"github.com/lucasepe/dotenv"
	"github.com/stretchr/testify/assert"
)

func TestUseCrumbs(t *testing.T) {
	cli := createClient(t)
	ok, err := cli.UseCrumbs(context.TODO())
	assert.True(t, ok, "expecting not using crumbs")
	assert.Nil(t, err, "expecting nil error")
}

func TestGetCrumbs(t *testing.T) {
	cli := createClient(t)
	key, val, err := cli.GetCrumbs(context.TODO())
	assert.Nil(t, err, "expecting nil error")

	t.Logf("Key: %s\n", key)
	t.Logf("Val: %s\n", val)
}

func TestGetJobConfigKO(t *testing.T) {
	var expErr *HTTPStatusError

	cli := createClient(t)
	dat, err := cli.GetJobConfig(context.TODO(), "oonode-jwt")
	assert.Nil(t, dat, "expecting nil data")
	assert.NotNil(t, err, "expecting not nil error")
	assert.ErrorAs(t, err, &expErr, "expecting not found error")
}

func TestGetJobConfigOk(t *testing.T) {
	cli := createClient(t)
	dat, err := cli.GetJobConfig(context.TODO(), "node-jwt")
	assert.NotNil(t, dat, "expecting not nil data")
	assert.Nil(t, err, "expecting nil error")

	ioutil.WriteFile("config.xml", dat, 0644)
	//t.Logf("\n%s\n", dat)
}

func TestCreateJob(t *testing.T) {
	dat, err := ioutil.ReadFile("config.xml")
	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, dat, "expecting not nil data")

	cli := createClient(t)
	err = cli.CreateJob(context.TODO(), "node-jwt-2", dat)
	assert.Nil(t, err, "expecting nil error creating job")
}

func TestDeleteJob(t *testing.T) {
	cli := createClient(t)
	err := cli.DeleteJob(context.TODO(), "node-jwt-2")
	assert.Nil(t, err, "expecting nil error deleting job")
}

func createClient(t *testing.T) *Client {
	envMap, err := dotenv.FromFile("../../../.env")
	if err != nil {
		t.Fatal(err)
	}

	dotenv.PutInEnv(envMap, false)

	return &Client{
		baseUrl: os.Getenv("JENKINS_URL"),
		auth: &Auth{
			Username: os.Getenv("JENKINS_USERNAME"),
			ApiToken: os.Getenv("JENKINS_API_TOKEN"),
		},
		httpClient: httphelper.ClientFromOpts(httphelper.ClientOpts{
			Verbose:  true,
			Insecure: false,
		}),
		crumbData: make(map[string]string),
	}
}
