package gcr

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/datagravity-ai/keel/extension/credentialshelper"
	"github.com/datagravity-ai/keel/types"
)

func init() {
	credentialshelper.RegisterCredentialsHelper("gcr", New())
}

type CredentialsHelper struct {
	enabled     bool
	credentials string
}

func New() *CredentialsHelper {
	ch := &CredentialsHelper{}

	credentialsFile, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !ok {
		return ch
	}

	credentials, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		return ch
	}

	ch.enabled = true
	ch.credentials = string(credentials)
	return ch
}

func (h *CredentialsHelper) IsEnabled() bool {
	return h.enabled
}

func (h *CredentialsHelper) GetCredentials(image *types.TrackedImage) (*types.Credentials, error) {
	if !h.enabled {
		return nil, errors.New("not initialised")
	}

	if image.Image.Registry() != "gcr.io" {
		return nil, credentialshelper.ErrUnsupportedRegistry
	}

	return &types.Credentials{
		Username: "_json_key",
		Password: h.credentials,
	}, nil
}
