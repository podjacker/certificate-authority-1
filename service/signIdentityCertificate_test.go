package service

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	authTest "github.com/go-ocf/authorization/oauth/test"
	"github.com/go-ocf/certificate-authority/pb"
	ocfSigner "github.com/go-ocf/kit/security/signer"
	"github.com/stretchr/testify/require"
)

func newIdentitySigner(t *testing.T) CertificateSigner {
	identityIntermediateCABlock, _ := pem.Decode(IdentityIntermediateCA)
	require.NotEmpty(t, identityIntermediateCABlock)
	identityIntermediateCA, err := x509.ParseCertificates(identityIntermediateCABlock.Bytes)
	require.NoError(t, err)
	identityIntermediateCAKeyBlock, _ := pem.Decode(IdentityIntermediateCAKey)
	require.NotEmpty(t, identityIntermediateCAKeyBlock)
	identityIntermediateCAKey, err := x509.ParseECPrivateKey(identityIntermediateCAKeyBlock.Bytes)
	require.NoError(t, err)
	return ocfSigner.NewIdentityCertificateSigner(identityIntermediateCA, identityIntermediateCAKey, time.Hour*86400)
}

func TestRequestHandler_SignIdentityCertificate(t *testing.T) {
	type args struct {
		req *pb.SignCertificateRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.SignCertificateResponse
		wantErr bool
	}{
		{
			name: "invalid auth",
			args: args{
				req: &pb.SignCertificateRequest{},
			},
			wantErr: true,
		},
		{
			name: "valid",
			args: args{
				req: &pb.SignCertificateRequest{
					AuthorizationContext: &pb.AuthorizationContext{
						AccessToken: authTest.UserToken,
					},
					CertificateSigningRequest: testCSR,
				},
			},
			wantErr: false,
		},
	}

	r := NewRequestHandler(nil, newIdentitySigner(t))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.SignIdentityCertificate(context.Background(), tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, got)
		})
	}
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
	testCSR = []byte(`-----BEGIN CERTIFICATE REQUEST-----
MIIBRjCB7QIBADA0MTIwMAYDVQQDEyl1dWlkOjAwMDAwMDAwLTAwMDAtMDAwMC0w
MDAwLTAwMDAwMDAwMDAwMTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABLiT0onX
Dw9JpJR9L1+SfyvILLZfluLTuxC7DNa0CdAhrGU2f6SCv+7VJQiQ02wlCt4iFCMx
u1XoaoEZuwcGKaSgVzBVBgkqhkiG9w0BCQ4xSDBGMAwGA1UdEwQFMAMBAQAwCwYD
VR0PBAQDAgGIMCkGA1UdJQQiMCAGCCsGAQUFBwMBBggrBgEFBQcDAgYKKwYBBAGC
3nwBBjAKBggqhkjOPQQDAgNIADBFAiAl/msC2XmurMvieTSOGt9aEgjZ197rchKL
IpK9P9vnXgIhAJ64cyN2X2uWu+x4NqpRkcneK0L3o0yOR4+DxF683pQ2
-----END CERTIFICATE REQUEST-----
`)
)
