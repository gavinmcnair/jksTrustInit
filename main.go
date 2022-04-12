package main

import (
//	"crypto/x509"
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
	Password    string `env:"PASSWORD" envDefault:"password"`
	Mode	    bool   `env:"FILE_MODE" envDefault: false`
	Key         string `env:"KEY"`
	Certificate string `env:"CERTIFICATE"`
	KeyFile         string `env:"KEY_FILE,file" envDefault:"key.pem"`
	CertificateFile string `env:"CERTIFICATE_FILE,file" envDefault:"cert.pem"`
	OutputJKS   string `env:"OUTPUT_JKS_FILE" envDefault:""`
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
	password := []byte(cfg.Password)

	ks1 := keystore.New()

	privateKey, errRPK := readPem("PRIVATE KEY", cfg.KeyFile)
	if errRPK != nil {
		return errRPK
	}

	certificate, errRS := readPem("CERTIFICATE", cfg.CertificateFile)
	if errRS != nil {
		return errRS
	}

	pkeIn := keystore.PrivateKeyEntry{
		CreationTime: time.Now(),
		PrivateKey:   privateKey,
		CertificateChain: []keystore.Certificate{
			{
				Type:    "X509",
				Content: certificate,
			},
		},
	}

	if errSPKE := ks1.SetPrivateKeyEntry("alias", pkeIn, password); errSPKE != nil {
		return errSPKE
	}

	if errWKS := writeKeyStore(ks1, "keystore.jks", password); errWKS != nil {
		return errWKS
	}

	// ks2, errRKS := readKeyStore("keystore.jks", password)
	// if errRKS != nil {
	// 	return errRKS
	// }

	// pkeOut, errGPKE := ks2.GetPrivateKeyEntry("alias", password)
	// if errGPKE != nil {
	// 	return errGPKE
	// }

	// _, errPK := x509.ParsePKCS8PrivateKey(pkeOut.PrivateKey)
	// if errPK != nil {
	// 	return errPK
	// }

	return nil
}

// func readKeyStore(filename string, password []byte) (keystore.KeyStore, error) {

// 	ks := keystore.New()

// 	f, err := os.Open(filename)
// 	if err != nil {
// 		return ks, err
// 	}

// 	defer func() {
// 		if err := f.Close(); err != nil {
// 			return
// 		}
// 	}()

// 	if err := ks.Load(f, password); err != nil {
// 		return ks, err
// 	}

// 	return ks, nil
// }

func writeKeyStore(ks keystore.KeyStore, filename string, password []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()

	err = ks.Store(f, password)
	if err != nil {
		return err
	}
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
