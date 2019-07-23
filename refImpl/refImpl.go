package refImpl

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-ocf/certificate-authority/service"
	"github.com/go-ocf/kit/log"
	kit "github.com/go-ocf/kit/net/grpc"
	"github.com/go-ocf/kit/security"
	ocfSigner "github.com/go-ocf/kit/security/signer"
	"google.golang.org/grpc"
)

type Config struct {
	Log                 log.Config
	Service             service.Config
	SignerCertificate   string        `envconfig:"SIGNER_CERTIFICATE" required:"True"`
	SignerPrivateKey    string        `envconfig:"SIGNER_PRIVATE_KEY" required:"True"`
	SignerValidDuration time.Duration `envconfig:"SIGNER_VALID_DURATION" default:"87600h"`
}

//String return string representation of Config
func (c Config) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return fmt.Sprintf("config: \n%v\n", string(b))
}

func Init(config Config) (*kit.Server, error) {
	log.Setup(config.Log)

	log.Info(config.String())

	handler, err := NewRequestHandlerFromConfig(config)
	if err != nil {
		return nil, err
	}
	cfg := kit.Config{Addr: config.Service.Addr, TLSConfig: config.Service.TLSConfig}
	svr, err := kit.NewServer(cfg, func(s *grpc.Server) {
		service.Register(s, handler)
	})
	if err != nil {
		return nil, err
	}
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
