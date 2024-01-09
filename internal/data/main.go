package data

import (
	"context"
	"github.com/google/uuid"
)

//go:generate xo schema "postgres://link:link@localhost:15432/link-db?sslmode=disable" -o ./ --single=schema.xo.go --src templates
//go:generate xo schema "postgres://link:link@localhost:15432/link-db?sslmode=disable" -o pg --single=schema.xo.go --src=pg/templates --go-context=both
//go:generate goimports -w ./

type Storage interface {
	Transaction(func() error) error
	ProofQ() ProofQ
	LinkQ() LinkQ
	LinksToProofQ() LinksToProofQ
}

type ProofQ interface {
	SelectAllCtx(ctx context.Context) ([]*Proof, error)
	ProofByID(proofID uuid.UUID, isForUpdate bool) (*Proof, error)
	Insert(p *Proof) error
	Update(p *Proof) error
	Upsert(p *Proof) error
	Delete(p *Proof) error
	ProofsByUserDIDCtx(ctx context.Context, userDID string) ([]Proof, error)
}

type LinkQ interface {
	Insert(l *Link) error
	Update(l *Link) error
	Upsert(l *Link) error
	Delete(l *Link) error
	SelectAllCtx(ctx context.Context) ([]*Link, error)
	LinkByID(id uuid.UUID, isForUpdate bool) (*Link, error)
	Transaction(fn func(db LinkQ) error) error
	InsertCtxLinkToProof(ctx context.Context, linksToProof LinksToProof) error
	GetProofsLinksByUserID(ctx context.Context, userID string) ([]*Link, error)
}

type LinksToProofQ interface {
	Insert(l *LinksToProof) error
	Delete(l *LinksToProof) error
	LinksToProofByLinkIDProofIDCtx(ctx context.Context, linkID, proofID uuid.UUID, isForUpdate bool) (*LinksToProof, error)
	SelectAllCtx(ctx context.Context) ([]*LinksToProofQ, error)
	GetLinksToProofsByLinkID(ctx context.Context, linkID uuid.UUID) ([]*LinksToProof, error)
	GetLinksToProofsByProofID(ctx context.Context, proofID uuid.UUID) ([]*LinksToProof, error)
}

type GorpMigrationQ interface {
}
