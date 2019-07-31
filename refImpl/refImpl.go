package refImpl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-ocf/certificate-authority/service"
	"github.com/go-ocf/kit/log"
	kit "github.com/go-ocf/kit/net/grpc"
	"github.com/go-ocf/kit/security"
	"github.com/go-ocf/kit/security/jwt"
	ocfSigner "github.com/go-ocf/kit/security/signer"
)

type Config struct {
	Log                 log.Config
	Service             service.Config
	SignerCertificate   string        `envconfig:"SIGNER_CERTIFICATE" required:"True"`
	SignerPrivateKey    string        `envconfig:"SIGNER_PRIVATE_KEY" required:"True"`
	SignerValidDuration time.Duration `envconfig:"SIGNER_VALID_DURATION" default:"87600h"`
	JwksUrl             string        `envconfig:"JWKS_URL" required:"True"`
}

//String return string representation of Config
func (c Config) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return fmt.Sprintf("config: \n%v\n", string(b))
}

func Init(config Config) (*kit.Server, error) {
	log.Setup(config.Log)
	log.Info(config.String())

	cfg := kit.Config{Addr: config.Service.Addr, TLSConfig: config.Service.TLSConfig}
	auth := NewAuth(config.JwksUrl, "ocf.cert.sign")

	svr, err := kit.NewServerWithoutPeerVerification(cfg, auth.Stream(), auth.Unary())
	if err != nil {
		return nil, err
	}

	handler, err := NewRequestHandlerFromConfig(config)
	if err != nil {
		return nil, err
	}
	service.Register(svr.Server, handler)

	return svr, nil
}

// NewRequestHandlerFromConfig creates RegisterGrpcGatewayServer with all dependencies.
func NewRequestHandlerFromConfig(config Config) (*service.RequestHandler, error) {
	chainCerts, err := security.LoadX509(config.SignerCertificate)
	if err != nil {
		return nil, err
	}
	privateKey, err := security.LoadX509PrivateKey(config.SignerPrivateKey)
	if err != nil {
		return nil, err
	}

	identitySigner := ocfSigner.NewIdentityCertificateSigner(chainCerts, privateKey, config.SignerValidDuration)
	signer := ocfSigner.NewBasicCertificateSigner(chainCerts, privateKey, config.SignerValidDuration)
	return service.NewRequestHandler(signer, identitySigner), nil
}

func NewAuth(jwksUrl string, scope string) kit.AuthInterceptors {
	return kit.MakeJWTInterceptors(jwksUrl, func(context.Context) kit.Claims {
		return jwt.NewScopeClaims(scope)
	})
}
