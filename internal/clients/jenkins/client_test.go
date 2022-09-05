package jenkins

import (
	"context"
	"errors"
	"testing"

	"github.com/krateoplatformops/provider-jenkins/internal/helpers"
	httphelper "github.com/krateoplatformops/provider-jenkins/internal/helpers/http"
)

func TestGet(t *testing.T) {
	cli := Client{
		baseUrl:  "https://jenkins.zipem.io/",
		username: helpers.StringPtr("Admin"),
		password: helpers.StringPtr("XXXXXXXXX"),
		httpClient: httphelper.ClientFromOpts(httphelper.ClientOpts{
			Verbose:  false,
			Insecure: false,
		}),
		crumbData: make(map[string]string),
	}

	dat, err := cli.GetJobConfig(context.TODO(), "oonode-jwt")
	if err != nil {
		if errors.Is(err, httphelper.ErrResourceNotFound) {
			t.Fatal("Eccolo")
		}
		t.Fatal(err)
	}

	t.Logf("%s\n", string(dat))
}
