package model

import (
	"github.com/aws/aws-sdk-go/service/iot"
)

type ThingConfig struct {
	// The ARN of the certificate.
	CertificateArn string

	// The ID of the certificate. AWS IoT issues a default subject name for the
	// certificate (e.g., AWS IoT Certificate).
	CertificateID string

	// The certificate data, in PEM format.
	CertificatePem string
}

type CsrConfig struct {
	CsrText string
}

// NewThingConfig create a new thing configuration using the response
func NewThingConfig(resp *iot.CreateCertificateFromCsrOutput) *ThingConfig {
	return &ThingConfig{
		CertificateArn: *resp.CertificateArn,
		CertificateID:  *resp.CertificateId,
		CertificatePem: *resp.CertificatePem,
	}
}
