/*
This file is part of configNexus.

configNexus is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

configNexus is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with configNexus.  If not, see <https://www.gnu.org/licenses/>.

Copyright (C) 2023 Operistech Inc.
*/
package utils

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerateSelfSignedCert(t *testing.T) {
	// Define temporary paths for storing the certificate and key
	certPath := "test_cert.pem"
	keyPath := "test_key.pem"

	t.Run("successfully generate self-signed certificate and key", func(t *testing.T) {
		err := GenerateSelfSignedCert(certPath, keyPath)
		if err != nil {
			t.Fatalf("Expected to successfully generate cert and key, got error: %v", err)
		}

		// Clean up
		defer os.Remove(certPath)
		defer os.Remove(keyPath)

		// Verify if certificate file is generated
		_, err = os.Stat(certPath)
		if err != nil {
			t.Fatalf("Failed to generate certificate file: %v", err)
		}

		// Verify if key file is generated
		_, err = os.Stat(keyPath)
		if err != nil {
			t.Fatalf("Failed to generate key file: %v", err)
		}

		// Read and parse certificate
		certData, err := ioutil.ReadFile(certPath)
		if err != nil {
			t.Fatalf("Failed to read cert file: %v", err)
		}

		block, _ := pem.Decode(certData)
		if block == nil {
			t.Fatal("Failed to decode PEM block")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			t.Fatalf("Failed to parse certificate: %v", err)
		}

		// Verify basic certificate attributes
		if len(cert.Subject.Organization) == 0 || cert.Subject.Organization[0] != "ConfigNexus" {
			t.Errorf("Unexpected organization in certificate. Got: %v", cert.Subject.Organization)
		}

		if len(cert.DNSNames) != 1 || cert.DNSNames[0] != "localhost" {
			t.Errorf("Unexpected DNS names in certificate. Got: %v", cert.DNSNames)
		}
	})
}
