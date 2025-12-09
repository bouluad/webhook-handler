package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
)

// ValidateSignature checks the GitHub X-Hub-Signature-256 header against the request body.
// It returns true if the signature is valid, false otherwise.
func ValidateSignature(signatureHeader string, payloadBody []byte, secret string) bool {
	// 1. Check for prefix and split the signature
	if !strings.HasPrefix(signatureHeader, "sha256=") {
		log.Println("Signature header missing 'sha256=' prefix.")
		return false
	}
	
	actualSignatureHex := strings.TrimPrefix(signatureHeader, "sha256=")
	
	// 2. Compute expected signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payloadBody)
	expectedSignature := mac.Sum(nil)
	
	expectedSignatureHex := hex.EncodeToString(expectedSignature)
	
	// 3. Compare using constant time comparison (CRITICAL for security)
	// We use the standard Go library function for constant-time comparison.
	// We compare the raw bytes of the computed signature and the decoded header signature.
	
	actualSignature, err := hex.DecodeString(actualSignatureHex)
	if err != nil {
		log.Printf("Failed to decode signature header: %v", err)
		return false
	}

	return hmac.Equal(actualSignature, expectedSignature)
}
