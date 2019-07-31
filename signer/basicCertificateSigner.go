package local

import (
	"context"

	"github.com/go-ocf/certificate-authority/pb"
	"github.com/go-ocf/kit/net/grpc"
)

type BasicCertificateSigner struct {
	accessToken string
	client      pb.CertificateAuthorityClient
}

func NewBasicCertificateSigner(client pb.CertificateAuthorityClient, accessToken string) *BasicCertificateSigner {
	return &BasicCertificateSigner{client: client, accessToken: accessToken}
}

func (s *BasicCertificateSigner) Sign(ctx context.Context, csr []byte) (signedCsr []byte, err error) {
	req := pb.SignCertificateRequest{CertificateSigningRequest: csr}
	resp, err := s.client.SignCertificate(grpc.CtxWithToken(ctx, s.accessToken), &req)
	return resp.GetCertificate(), err
}
