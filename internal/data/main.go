package data

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
	// FIXME
	//SelectAll() ([]*Proof, error)
	//GetProofByID(proofID uuid.UUID) (Proof, error)
	//ProofByIDCtx(ctx context.Context, id uuid.UUID, isForUpdate bool) (*Proof, error)
	//ProofsByUserDIDCtx(ctx context.Context, userDID string, isForUpdate bool) ([]Proof, error)
	//InsertCtx(ctx context.Context, p *Proof) error
	//DeleteByID(id uuid.UUID) error
}

type LinkQ interface {
	// FIXME
	//SelectAll() ([]*Link, error)
	//GetProofsByUserDID(userDID string) ([]Link, error)
	//InsertCtx(ctx context.Context, l *Link) error
	//InsertCtxLinkToProof(ctx context.Context, l *LinksToProof) error
	//Transaction(fn func(db LinkQ) error) error
}

type LinksToProofQ interface {
	// FIXME
	//SelectAll() ([]*LinksToProof, error)
	//GetProofsByLinkID(linkID uuid.UUID) ([]LinksToProof, error)
	//InsertCtx(ctx context.Context, l *LinksToProof) error
}

type GorpMigrationQ interface {
}
