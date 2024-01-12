/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Link struct {
	Key
	Attributes LinkAttributes `json:"attributes"`
}
type LinkResponse struct {
	Data     Link     `json:"data"`
	Included Included `json:"included"`
}

type LinkListResponse struct {
	Data     []Link          `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *LinkListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *LinkListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustLink - returns Link from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLink(key Key) *Link {
	var link Link
	if c.tryFindEntry(key, &link) {
		return &link
	}
	return nil
}
