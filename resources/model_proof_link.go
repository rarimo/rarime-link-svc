/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type ProofLink struct {
	Key
	Attributes ProofLinkAttributes `json:"attributes"`
}
type ProofLinkResponse struct {
	Data     ProofLink `json:"data"`
	Included Included  `json:"included"`
}

type ProofLinkListResponse struct {
	Data     []ProofLink     `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *ProofLinkListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *ProofLinkListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustProofLink - returns ProofLink from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustProofLink(key Key) *ProofLink {
	var proofLink ProofLink
	if c.tryFindEntry(key, &proofLink) {
		return &proofLink
	}
	return nil
}
