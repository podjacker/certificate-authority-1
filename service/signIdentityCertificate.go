package service

import (
	"context"

	"github.com/go-ocf/certificate-authority/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *RequestHandler) SignIdentityCertificate(ctx context.Context, req *pb.SignCertificateRequest) (*pb.SignCertificateResponse, error) {
	if err := r.auth(req.GetAuthorizationContext().GetAccessToken()); err != nil {
		return nil, logAndReturnError(status.Errorf(codes.InvalidArgument, "cannot sign identity certificate: %v", err))
	}

	cert, err := r.identitySigner.Sign(ctx, req.CertificateSigningRequest)
	if err != nil {
		return nil, logAndReturnError(status.Errorf(codes.InvalidArgument, "cannot sign identity certificate: %v", err))
	}
	return &pb.SignCertificateResponse{
		Certificate: cert,
	}, nil
}
