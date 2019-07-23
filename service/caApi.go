package service

import (
	"context"

	"github.com/go-ocf/certificate-authority/pb"
	"github.com/go-ocf/kit/log"
	"google.golang.org/grpc"
)

type CertificateSigner interface {
	//csr is encoded by PEM and returns PEM
	Sign(ctx context.Context, csr []byte) ([]byte, error)
}

// RequestHandler handles incoming requests.
type RequestHandler struct {
	identitySigner CertificateSigner
	signer         CertificateSigner
}

// Register registers the handler instance with a gRPC server.
func Register(server *grpc.Server, handler *RequestHandler) {
	pb.RegisterCertificateAuthorityServer(server, handler)
}

// NewRequestHandler factory for new RequestHandler.
func NewRequestHandler(signer, identitySigner CertificateSigner) *RequestHandler {
	return &RequestHandler{
		signer:         signer,
		identitySigner: identitySigner,
	}
}

func logAndReturnError(err error) error {
	log.Errorf("%v", err)
	return err
}
