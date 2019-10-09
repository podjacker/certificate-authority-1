package refImpl

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-ocf/kit/security/acme"

	"github.com/go-ocf/certificate-authority/pb"
	"github.com/go-ocf/certificate-authority/service"
	"github.com/go-ocf/kit/log"
	kitNetGrpc "github.com/go-ocf/kit/net/grpc"
	"github.com/go-ocf/kit/security"
	"github.com/go-ocf/kit/security/jwt"
	ocfSigner "github.com/go-ocf/kit/security/signer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	Log                 log.Config
	Service             service.Config
	ListenAcme          acme.Config   `envconfig:"LISTEN_ACME"`
	SignerCertificate   string        `envconfig:"SIGNER_CERTIFICATE" required:"True"`
	SignerPrivateKey    string        `envconfig:"SIGNER_PRIVATE_KEY" required:"True"`
	SignerValidDuration time.Duration `envconfig:"SIGNER_VALID_DURATION" default:"87600h"`
	JwksURL             string        `envconfig:"JWKS_URL" required:"True"`
}

type RefImpl struct {
	handle            *service.RequestHandler
	server            *kitNetGrpc.Server
	listenCertManager *acme.CertManager
}

// NewRequestHandlerFromConfig creates RegisterGrpcGatewayServer with all dependencies.
func NewRefImplFromConfig(config Config, auth kitNetGrpc.AuthInterceptors) (*RefImpl, error) {
	chainCerts, err := security.LoadX509(config.SignerCertificate)
	if err != nil {
		return nil, err
	}
	privateKey, err := security.LoadX509PrivateKey(config.SignerPrivateKey)
	if err != nil {
		return nil, err
	}

	listenCertManager, err := acme.NewCertManagerFromConfiguration(config.ListenAcme)
	if err != nil {
		return nil, fmt.Errorf("cannot create server cert manager %v", err)
	}

	serverTLSConfig := listenCertManager.GetServerTLSConfig()
	serverTLSConfig.ClientAuth = tls.NoClientCert

	svr, err := kitNetGrpc.NewServer(config.Service.Addr, grpc.Creds(credentials.NewTLS(&serverTLSConfig)), auth.Stream(), auth.Unary())
	if err != nil {
		listenCertManager.Close()
		return nil, err
	}

	identitySigner := ocfSigner.NewIdentityCertificateSigner(chainCerts, privateKey, config.SignerValidDuration)
	signer := ocfSigner.NewBasicCertificateSigner(chainCerts, privateKey, config.SignerValidDuration)
	return &RefImpl{
		handle:            service.NewRequestHandler(signer, identitySigner),
		listenCertManager: listenCertManager,
		server:            svr,
	}, nil
}

//String return string representation of Config
func (c Config) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return fmt.Sprintf("config: \n%v\n", string(b))
}

func Init(config Config) (*RefImpl, error) {
	//auth := NewAuth(config.JwksURL, "ocf.cert.sign")
	auth := kitNetGrpc.MakeAuthInterceptors(func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	})
	return InitWithAuth(config, auth)
}

func InitWithAuth(config Config, auth kitNetGrpc.AuthInterceptors) (*RefImpl, error) {
	log.Setup(config.Log)
	log.Info(config.String())

	impl, err := NewRefImplFromConfig(config, auth)
	if err != nil {
		return nil, err
	}

	pb.RegisterCertificateAuthorityServer(impl.server.Server, impl.handle)

	return impl, nil
}

func (r *RefImpl) Serve() error {
	return r.server.Serve()
}

func (r *RefImpl) Shutdown() {
	r.server.Stop()
	r.listenCertManager.Close()
}

func NewAuth(jwksUrl string, scope string) kitNetGrpc.AuthInterceptors {
	return kitNetGrpc.MakeJWTInterceptors(jwksUrl, func(context.Context) kitNetGrpc.Claims {
		return jwt.NewScopeClaims(scope)
	})
}
