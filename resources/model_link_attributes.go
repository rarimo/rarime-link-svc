/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type LinkAttributes struct {
	// The date and time when the proof was created in the RFC3339 format
	CreatedAt time.Time `json:"created_at"`
	// Link to proofs ID (UUID or custom ASCII string)
	Link string `json:"link"`
}
