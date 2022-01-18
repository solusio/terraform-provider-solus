package provider

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/solusio/solus-go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testAccProvider          = New()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"solus": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
)

func TestProvider(t *testing.T) {
	t.Run("runs internal validation without error", func(t *testing.T) {
		assert.NoError(t, New().InternalValidate())
	})

	t.Run("configure", func(t *testing.T) {
		p := New()
		d := p.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
		require.Nil(t, d)
		assert.IsType(t, &solus.Client{}, p.Meta())
	})
}

func testAccPreCheck(t *testing.T) func() {
	return func() {
		ee := []string{baseURLEnv, tokenEnv}
		for _, e := range ee {
			assert.NotEmptyf(t, os.Getenv(e), "%q environment variable must be set for acceptance tests", e)
		}
	}
}

func generateResourceName() string {
	const nameLength = 16
	return generateString(nameLength)
}

var (
	runes      = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	genStrRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func generateString(l int) string {
	rr := make([]rune, l)
	for i := range rr {
		rr[i] = runes[genStrRand.Intn(len(runes))]
	}
	return string(rr)
}
