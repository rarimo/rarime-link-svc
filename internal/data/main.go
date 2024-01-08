package data

import "context"

//go:generate xo schema "postgres://link:link@localhost:15432/link-db?sslmode=disable" -o ./ --single=schema.xo.go --src templates
//go:generate xo schema "postgres://link:link@localhost:15432/link-db?sslmode=disable" -o pg --single=schema.xo.go --src=pg/templates --go-context=both
//go:generate goimports -w ./

type Storage interface {
	ProofQ() ProofQ
	ProofLinkQ() ProofLinkQ
	LinkToProofQ() LinkToProofQ
}

type ProofQ interface {
	SelectAll() ([]*Proof, error)
	GetProofByID(proofID int) (Proof, error)
	ProofByIDCtx(ctx context.Context, id int, isForUpdate bool) (*Proof, error)
	ProofsByUserDIDCtx(ctx context.Context, userDID string, isForUpdate bool) ([]Proof, error)
	InsertCtx(ctx context.Context, p *Proof) error
	DeleteByID(id int) error
}

type ProofLinkQ interface {
	SelectAll() ([]*Link, error)
	GetProofsByIndex(index int) ([]Link, error)
	InsertCtx(ctx context.Context, l *Link) error
	GetLastIndex() (int, error)
	Transaction(fn func(db ProofLinkQ) error) error
}

type LinkToProofQ interface {
	SelectAll() ([]*LinkToProof, error)
	InsertCtx(ctx context.Context, l *LinkToProof) error
}

type GorpMigrationQ interface {
}
