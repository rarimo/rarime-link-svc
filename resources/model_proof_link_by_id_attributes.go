/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ProofLinkByIdAttributes struct {
	// The date and time when the proof was created in the timestamp format
	CreatedAt string `json:"created_at"`
	// The ID of the user who created the proof
	Creator string `json:"creator"`
	// UUID Link to proofs
	Link string `json:"link"`
	// The proof object in JSON string format
	Proof string `json:"proof"`
	// The type of the proof
	Type string `json:"type"`
}
