package refImpl

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	dir, err := ioutil.TempDir("", "gotesttmp")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	var cfg Config
	cfg = testSignerCerts(t, dir, cfg)

	type args struct {
		config Config
	}
	tests := []struct {
		name       string
		args       args
		wantServer bool
		wantErr    bool
	}{
		{
			name:    "invalid",
			wantErr: true,
		},
		{
			name: "valid",
			args: args{
				config: Config{
					SignerCertificate:   cfg.SignerCertificate,
					SignerPrivateKey:    cfg.SignerPrivateKey,
					SignerValidDuration: time.Hour * 10,
					Addr:                ":1234",
				},
			},
			wantServer: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Init(tt.args.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.wantServer {
				assert.NotEmpty(t, got)
			} else {
				assert.Empty(t, got)

			}
		})
	}
}

func testSignerCerts(t *testing.T, dir string, cfg Config) Config {
	crt := filepath.Join(dir, "cert.crt")
	if err := ioutil.WriteFile(crt, IdentityIntermediateCA, 0600); err != nil {
		assert.NoError(t, err)
	}
	crtKey := filepath.Join(dir, "cert.key")
	if err := ioutil.WriteFile(crtKey, IdentityIntermediateCAKey, 0600); err != nil {
		assert.NoError(t, err)
	}
	cfg.SignerCertificate = crt
	cfg.SignerPrivateKey = crtKey
	return cfg
}

var (
	IdentityIntermediateCA = []byte(`-----BEGIN CERTIFICATE-----
MIIBczCCARmgAwIBAgIRANntjEpzu9krzL0EG6fcqqgwCgYIKoZIzj0EAwIwETEP
MA0GA1UEAxMGUm9vdENBMCAXDTE5MDcxOTIwMzczOVoYDzIxMTkwNjI1MjAzNzM5
WjAZMRcwFQYDVQQDEw5JbnRlcm1lZGlhdGVDQTBZMBMGByqGSM49AgEGCCqGSM49
AwEHA0IABKw1/6WHFcWtw67hH5DzoZvHgA0suC6IYLKms4IP/pds9wU320eDaENo
5860TOyKrGn7vW/cj/OVe2Dzr4KSFVijSDBGMA4GA1UdDwEB/wQEAwIBBjATBgNV
HSUEDDAKBggrBgEFBQcDATASBgNVHRMBAf8ECDAGAQH/AgEAMAsGA1UdEQQEMAKC
ADAKBggqhkjOPQQDAgNIADBFAiEAgPtnYpgwxmPhN0Mo8VX582RORnhcdSHMzFjh
P/li1WwCIFVVWBOrfBnTt7A6UfjP3ljAyHrJERlMauQR+tkD/aqm
-----END CERTIFICATE-----
`)
	IdentityIntermediateCAKey = []byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPF4DPvFeiRL1G0ROd6MosoUGnvIG/2YxH0CbHwnLKxqoAoGCCqGSM49
AwEHoUQDQgAErDX/pYcVxa3DruEfkPOhm8eADSy4Lohgsqazgg/+l2z3BTfbR4No
Q2jnzrRM7Iqsafu9b9yP85V7YPOvgpIVWA==
-----END EC PRIVATE KEY-----
`)
)
