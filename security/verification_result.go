package security

// VerificationResult represents the outcome of a signature verification process.
type VerificationResult struct {
	success bool
	errors  []string
}

// Success returns true if the verification was successful.
func (v *VerificationResult) Success() bool {
	return v.success
}

// Errors returns the list of errors (if any) encountered during verification.
func (v *VerificationResult) Errors() []string {
	return v.errors
}

// NewSuccess returns a VerificationResult representing a successful verification.
func NewSuccess() VerificationResult {
	return VerificationResult{
		success: true,
		errors:  []string{},
	}
}

// NewFailure returns a VerificationResult representing a failed verification.
func NewFailure(errors ...string) VerificationResult {
	return VerificationResult{
		success: false,
		errors:  errors,
	}
}
