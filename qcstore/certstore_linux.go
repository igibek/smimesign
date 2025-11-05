package qcstore

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

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
	// TODO: read the certificate
	certPem, err := os.ReadFile("./real-smime-certificate.crt")
	if err != nil {
		return nil, fmt.Errorf("failed to read the certificate: %w", err)
	}

	block, _ := pem.Decode(certPem)

	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to find or decode PEM block of type certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)

	if err != nil {
		return nil, fmt.Errorf("failed to parse DER certificate: %w", err)
	}

	i.cert = cert

	return i.cert, nil
}

// CertificateChain attempts to get the identity's full certificate chain.
func (i *LinuxIdentity) CertificateChain() ([]*x509.Certificate, error) {
	chain := make([]*x509.Certificate, len(i.chain))
	chain = append(chain, i.cert)
	i.chain = chain
	return i.chain, nil
}

// Signer gets a crypto.Signer that uses the identity's private key.
func (i *LinuxIdentity) Signer() (crypto.Signer, error) {
	// TODO: configure and return the KMS signer
	kms, err := awskms.New(context.TODO(), apiv1.Options{})

	if err != nil {
		return nil, fmt.Errorf("failed creating KMS %s", err)
	}

	signer, err := kms.CreateSigner(&apiv1.CreateSignerRequest{
		SigningKey: "arn:aws:kms:eu-central-1:523676012530:key/5b0fcd89-7688-425f-9c46-219727925d89", // TODO: get this from environment variable
	})

	if err != nil {
		return nil, fmt.Errorf("failed creating signer: %w", err)
	}
	// pub, _ := x509.MarshalPKIXPublicKey(signer.Public())
	// fmt.Printf("debug kms signer %s", pub)
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
