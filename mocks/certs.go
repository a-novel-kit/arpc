package arpcmocks

import (
	"crypto/x509"
	"errors"
)

var ErrAppendCert = errors.New("failed to append certificate to default pool")

func ClientCerts(mocked ...[]byte) func() (*x509.CertPool, error) {
	return func() (*x509.CertPool, error) {
		pool := x509.NewCertPool()

		for _, cert := range mocked {
			if !pool.AppendCertsFromPEM(cert) {
				return nil, ErrAppendCert
			}
		}

		return pool, nil
	}
}
