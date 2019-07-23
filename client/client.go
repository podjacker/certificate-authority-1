package client

import (
	"fmt"

	"github.com/go-ocf/certificate-authority/pb"
	"github.com/go-ocf/kit/net/grpc"
)

// NewClient creates grpc gateway client
func NewClient(host string, opts ...grpc.ClientConnOption) (pb.CertificateAuthorityClient, error) {
	conn, err := grpc.NewClientConnWithOptions(host, opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot create resource directory client: %v", err)
	}
	return pb.NewCertificateAuthorityClient(conn), nil
}
