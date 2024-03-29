// Package data contains generated code for schema 'public'.
package data

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"database/sql/driver"
	"encoding/csv"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rarimo/xo/types/xo"

	"github.com/google/uuid"
)

// StringSlice is a slice of strings.
type StringSlice []string

// quoteEscapeRegex is the regex to match escaped characters in a string.
var quoteEscapeRegex = regexp.MustCompile(`([^\\]([\\]{2})*)\\"`)

// Scan satisfies the sql.Scanner interface for StringSlice.
func (ss *StringSlice) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid StringSlice")
	}

	// change quote escapes for csv parser
	str := quoteEscapeRegex.ReplaceAllString(string(buf), `$1""`)
	str = strings.Replace(str, `\\`, `\`, -1)

	// remove braces
	str = str[1 : len(str)-1]

	// bail if only one
	if len(str) == 0 {
		*ss = StringSlice([]string{})
		return nil
	}

	// parse with csv reader
	cr := csv.NewReader(strings.NewReader(str))
	slice, err := cr.Read()
	if err != nil {
		fmt.Printf("exiting!: %v\n", err)
		return err
	}

	*ss = StringSlice(slice)

	return nil
}

// Value satisfies the driver.Valuer interface for StringSlice.
func (ss StringSlice) Value() (driver.Value, error) {
	v := make([]string, len(ss))
	for i, s := range ss {
		v[i] = `"` + strings.Replace(strings.Replace(s, `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
	}
	return "{" + strings.Join(v, ",") + "}", nil
} // GorpMigration represents a row from 'public.gorp_migrations'.
type GorpMigration struct {
	ID        string       `db:"id" json:"id" structs:"-"`                          // id
	AppliedAt sql.NullTime `db:"applied_at" json:"applied_at" structs:"applied_at"` // applied_at

}

// Link represents a row from 'public.links'.
type Link struct {
	ID        string    `db:"id" json:"id" structs:"-"`                          // id
	UserID    string    `db:"user_id" json:"user_id" structs:"user_id"`          // user_id
	CreatedAt time.Time `db:"created_at" json:"created_at" structs:"created_at"` // created_at

}

// LinksToProof represents a row from 'public.links_to_proofs'.
type LinksToProof struct {
	LinkID  string    `db:"link_id" json:"link_id" structs:"-"`   // link_id
	ProofID uuid.UUID `db:"proof_id" json:"proof_id" structs:"-"` // proof_id

}

// Proof represents a row from 'public.proofs'.
type Proof struct {
	ID        uuid.UUID     `db:"id" json:"id" structs:"-"`                          // id
	Creator   string        `db:"creator" json:"creator" structs:"creator"`          // creator
	CreatedAt time.Time     `db:"created_at" json:"created_at" structs:"created_at"` // created_at
	Proof     xo.Jsonb      `db:"proof" json:"proof" structs:"proof"`                // proof
	OrgID     uuid.UUID     `db:"org_id" json:"org_id" structs:"org_id"`             // org_id
	Type      string        `db:"type" json:"type" structs:"type"`                   // type
	SchemaURL string        `db:"schema_url" json:"schema_url" structs:"schema_url"` // schema_url
	Operator  ProofOperator `db:"operator" json:"operator" structs:"operator"`       // operator
	Field     string        `db:"field" json:"field" structs:"field"`                // field

}
