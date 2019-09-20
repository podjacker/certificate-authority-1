package service

import (
	"context"

	"github.com/go-ocf/certificate-authority/pb"
	"github.com/go-ocf/kit/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *RequestHandler) SignCertificate(ctx context.Context, req *pb.SignCertificateRequest) (*pb.SignCertificateResponse, error) {
	log.Debugf("RequestHandler.SignCertificate: %v", string(req.CertificateSigningRequest))
	cert, err := r.signer.Sign(ctx, req.CertificateSigningRequest)
	if err != nil {
		return nil, logAndReturnError(status.Errorf(codes.InvalidArgument, "cannot sign certificate: %v", err))
	}
	return &pb.SignCertificateResponse{
		Certificate: cert,
	}, nil
}
