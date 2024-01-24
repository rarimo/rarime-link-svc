/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type ProofAttributes struct {
	// The date and time when the proof was created in the RFC3339 format
	CreatedAt time.Time `json:"created_at"`
	// The ID of the user who created the proof
	Creator string `json:"creator"`
	// The operator that will be used to check the proof
	Operator string `json:"operator"`
	// The ID of the organization that issued the proof's claim
	OrgId string `json:"org_id"`
	// The proof object in JSON string format
	Proof string `json:"proof"`
	// The type of the proof
	ProofType string `json:"proof_type"`
	// The schema URL of the claim the proof was created based on
	SchemaUrl string `json:"schema_url"`
}
