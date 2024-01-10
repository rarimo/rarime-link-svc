/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ProofCreate struct {
	// The ID of the organization that issued the proof's claim
	OrgId string `json:"org_id"`
	// The proof object in JSON string format
	Proof string `json:"proof"`
	// The type of the proof
	ProofType string `json:"proof_type"`
}
