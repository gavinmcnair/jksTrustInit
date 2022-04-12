package main

import (
	"encoding/pem"
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/pavel-v-chernykh/keystore-go/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type config struct {
	Password        string `env:"PASSWORD" envDefault:"password"`
	Mode            bool   `env:"FILE_MODE" envDefault: "false"`
	Key             string `env:"KEY"`
	Certificate     string `env:"CERTIFICATE"`
	KeyFile         string `env:"KEY_FILE,file"`
	CertificateFile string `env:"CERTIFICATE_FILE,file"`
	OutputJKS       string `env:"OUTPUT_FILE" envDefault:"/var/run/secrets/truststore.jks"`
}

func main() {
	zerolog.DurationFieldUnit = time.Second
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run")
	}
	log.Info().Msg("Gracefully exiting")
}

func run() error {
	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		return err
	}

	var key,certificate string

	if cfg.Mode == true {
		key = cfg.KeyFile
		certificate = cfg.CertificateFile
	} else {
		key = cfg.Key
		certificate = cfg.Certificate
	}

	password := []byte(cfg.Password)

	ks := keystore.New()

	pemPrivateKey, err := readPem("PRIVATE KEY", key)
	if err != nil {
		return err
	}

	pemCertificate, err := readPem("CERTIFICATE", certificate)
	if err != nil {
		return err
	}

	pkeIn := keystore.PrivateKeyEntry{
		CreationTime: time.Now(),
		PrivateKey:   pemPrivateKey,
		CertificateChain: []keystore.Certificate{
			{
				Type:    "X509",
				Content: pemCertificate,
			},
		},
	}

	if err := ks.SetPrivateKeyEntry("alias", pkeIn, password); err != nil {
		return err
	}

	if err := writeKeyStore(ks, cfg.OutputJKS, password); err != nil {
		return err
	}

	return nil
}

func writeKeyStore(ks keystore.KeyStore, filename string, password []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	if err = ks.Store(f, password); err != nil {
		f.Close()
		return err
	}

	f.Close()
	return nil
}

func readPem(expectedType string, data string) ([]byte, error) {
	b, _ := pem.Decode([]byte(data))
	if b == nil {
		return nil, errors.New("should have at least one pem block")
	}

	if b.Type != expectedType {
		return nil, errors.New("should be a " + expectedType)
	}

	return b.Bytes, nil
}
