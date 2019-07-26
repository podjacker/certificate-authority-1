package local

import (
	"context"

	"github.com/go-ocf/certificate-authority/pb"
)

type IdentityCertificateSigner struct {
	accessToken string
	client      pb.CertificateAuthorityClient
}

func NewIdentityCertificateSigner(client pb.CertificateAuthorityClient, accessToken string) *IdentityCertificateSigner {
	return &IdentityCertificateSigner{client: client, accessToken: accessToken}
}

func (s *IdentityCertificateSigner) Sign(ctx context.Context, csr []byte) (signedCsr []byte, err error) {
	req := pb.SignCertificateRequest{
		CertificateSigningRequest: csr,
		AuthorizationContext: &pb.AuthorizationContext{
			AccessToken: s.accessToken,
		},
	}
	resp, err := s.client.SignIdentityCertificate(ctx, &req)
	return resp.GetCertificate(), err
}
