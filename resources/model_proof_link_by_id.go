/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type ProofLinkById struct {
	Key
	Attributes ProofLinkByIdAttributes `json:"attributes"`
}
type ProofLinkByIdResponse struct {
	Data     ProofLinkById `json:"data"`
	Included Included      `json:"included"`
}

type ProofLinkByIdListResponse struct {
	Data     []ProofLinkById `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *ProofLinkByIdListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *ProofLinkByIdListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustProofLinkById - returns ProofLinkById from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustProofLinkById(key Key) *ProofLinkById {
	var proofLinkByID ProofLinkById
	if c.tryFindEntry(key, &proofLinkByID) {
		return &proofLinkByID
	}
	return nil
}
