package data

import (
	"context"
	"github.com/google/uuid"
)

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
	GetProofByID(proofID uuid.UUID) (Proof, error)
	ProofByIDCtx(ctx context.Context, id uuid.UUID, isForUpdate bool) (*Proof, error)
	ProofsByUserDIDCtx(ctx context.Context, userDID string, isForUpdate bool) ([]Proof, error)
	InsertCtx(ctx context.Context, p *Proof) error
	DeleteByID(id uuid.UUID) error
}

type ProofLinkQ interface {
	SelectAll() ([]*Link, error)
	GetProofsByUserDID(userDID string) ([]Link, error)
	InsertCtx(ctx context.Context, l *Link) error
	InsertCtxLinkToProof(ctx context.Context, l *LinkToProof) error
	Transaction(fn func(db ProofLinkQ) error) error
}

type LinkToProofQ interface {
	SelectAll() ([]*LinkToProof, error)
	GetProofsByLinkID(linkID uuid.UUID) ([]LinkToProof, error)
	InsertCtx(ctx context.Context, l *LinkToProof) error
}

type GorpMigrationQ interface {
}
