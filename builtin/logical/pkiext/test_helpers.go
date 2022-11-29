package pkiext

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/stretchr/testify/require"
)

func requireFieldsSetInResp(t *testing.T, resp *logical.Response, fields ...string) {
	var missingFields []string
	for _, field := range fields {
		value, ok := resp.Data[field]
		if !ok || value == nil {
			missingFields = append(missingFields, field)
		}
	}

	require.Empty(t, missingFields, "The following fields were required but missing from response:\n%v", resp.Data)
}

func requireSuccessNonNilResponse(t *testing.T, resp *logical.Response, err error, msgAndArgs ...interface{}) {
	require.NoError(t, err, msgAndArgs...)
	if resp.IsError() {
		errContext := fmt.Sprintf("Expected successful response but got error: %v", resp.Error())
		require.Falsef(t, resp.IsError(), errContext, msgAndArgs...)
	}
	require.NotNil(t, resp, msgAndArgs...)
}

func requireSuccessNilResponse(t *testing.T, resp *logical.Response, err error, msgAndArgs ...interface{}) {
	require.NoError(t, err, msgAndArgs...)
	if resp.IsError() {
		errContext := fmt.Sprintf("Expected successful response but got error: %v", resp.Error())
		require.Falsef(t, resp.IsError(), errContext, msgAndArgs...)
	}
	if resp != nil {
		msg := fmt.Sprintf("expected nil response but got: %v", resp)
		require.Nilf(t, resp, msg, msgAndArgs...)
	}
}

func parseCert(t *testing.T, pemCert string) *x509.Certificate {
	block, _ := pem.Decode([]byte(pemCert))
	require.NotNil(t, block, "failed to decode PEM block")

	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)
	return cert
}

func parseKey(t *testing.T, pemKey string) crypto.Signer {
	block, _ := pem.Decode([]byte(pemKey))
	require.NotNil(t, block, "failed to decode PEM block")

	key, _, err := certutil.ParseDERKey(block.Bytes)
	require.NoError(t, err)
	return key
}

func requireSerialNumberInCRL(t *testing.T, revokeList pkix.TBSCertificateList, cert *x509.Certificate) bool {
	for _, revokeEntry := range revokeList.RevokedCertificates {
		if revokeEntry.SerialNumber.Cmp(cert.SerialNumber) == 0 {
			return true
		}
	}

	if t != nil {
		t.Fatalf("the serial number for cert %v, was not found in the CRL list containing: %v", cert, revokeList)
	}

	return false
}

func parseCRL(t *testing.T, crlPem string) pkix.TBSCertificateList {
	certList, err := x509.ParseCRL([]byte(crlPem))
	require.NoError(t, err)
	return certList.TBSCertList
}