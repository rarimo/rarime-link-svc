/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Proof struct {
	Key
	Attributes ProofAttributes `json:"attributes"`
}
type ProofResponse struct {
	Data     Proof    `json:"data"`
	Included Included `json:"included"`
}

type ProofListResponse struct {
	Data     []Proof         `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *ProofListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *ProofListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustProof - returns Proof from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustProof(key Key) *Proof {
	var proof Proof
	if c.tryFindEntry(key, &proof) {
		return &proof
	}
	return nil
}
