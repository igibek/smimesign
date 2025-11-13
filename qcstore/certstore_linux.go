package qcstore

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/spf13/viper"
	"go.step.sm/crypto/kms/apiv1"
	"go.step.sm/crypto/kms/awskms"
)

// Implement this function, just to silence other compiler errors.
func openStore() (Store, error) {
	return &LinuxStore{}, nil
}

// Implementation of Store interface

type LinuxStore struct {
}

func (s *LinuxStore) Identities() ([]Identity, error) {
	// TODO: in reality this should iterate through the well-know folder for the certificates and populate identitfies
	// currently this version only works with a single certificate provided as environment variable. Checkout LinuxIdentity Certificate function
	idents := []Identity{}

	idents = append(idents, &LinuxIdentity{})
	return idents, nil
}

func (s *LinuxStore) Import(data []byte, password string) error {
	return nil
}

func (s *LinuxStore) Close() {

}

// Implemention of Identity interface
type LinuxIdentity struct {
	cert  *x509.Certificate
	chain []*x509.Certificate
}

func (i *LinuxIdentity) Certificate() (*x509.Certificate, error) {
	certPath := viper.GetString("CERT_PATH")
	if certPath == "" {
		return nil, fmt.Errorf("failed CERT_PATH is empty")
	}

	certPem, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the certificate: %w", err)
	}
	rest := certPem
	var certificates []*x509.Certificate

	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}

		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse DER certificate: %w", err)
			}
			certificates = append(certificates, cert)
		}

	}

	i.cert = certificates[0]
	i.chain = certificates

	return i.cert, nil
}

// CertificateChain attempts to get the identity's full certificate chain.
func (i *LinuxIdentity) CertificateChain() ([]*x509.Certificate, error) {
	if i.chain != nil {
		return i.chain, nil
	}
	chain := make([]*x509.Certificate, len(i.chain))
	chain = append(chain, i.cert)
	i.chain = chain
	return i.chain, nil
}

// Signer gets a crypto.Signer that uses the identity's private key.
func (i *LinuxIdentity) Signer() (crypto.Signer, error) {

	// TODO: get the encrypted AWS access key from config file
	kms, err := awskms.New(context.TODO(), apiv1.Options{})

	if err != nil {
		return nil, fmt.Errorf("failed creating KMS %s", err)
	}
	signingKey := viper.GetString("AWS_KEY_ARN")
	if signingKey == "" {
		return nil, fmt.Errorf("failed to read AWS_KEY_ARN value")
	}
	signer, err := kms.CreateSigner(&apiv1.CreateSignerRequest{
		SigningKey: signingKey,
	})

	if err != nil {
		return nil, fmt.Errorf("failed creating signer: %w", err)
	}

	return signer, nil
}

// Delete deletes this identity from the system.
func (i *LinuxIdentity) Delete() error {
	return nil
}

// Close any manually managed memory held by the Identity.
func (i *LinuxIdentity) Close() {

}

/*
Implementation of crypto.Signer interface
*/
type LinuxSigner struct {
}

// Public implements the crypto.Signer interface.
func (signer *LinuxSigner) Public() crypto.PublicKey {
	// return wpk.publicKey
	return nil
}

// Sign implements the crypto.Signer interface.
func (signer *LinuxSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	// if wpk.capiProv != 0 {
	// 	return wpk.capiSignHash(opts.HashFunc(), digest)
	// } else if wpk.cngHandle != 0 {
	// 	return wpk.cngSignHash(opts.HashFunc(), digest)
	// } else {
	// 	return nil, errors.New("bad private key")
	// }
	return nil, nil
}
