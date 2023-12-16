package data

import "context"

//go:generate xo schema "postgres://link:link@localhost:15432/link-db?sslmode=disable" -o ./ --single=schema.xo.go --src templates
//go:generate xo schema "postgres://link:link@localhost:15432/link-db?sslmode=disable" -o pg --single=schema.xo.go --src=pg/templates --go-context=both
//go:generate goimports -w ./

type Storage interface {
	ProofQ() ProofQ
}

type ProofQ interface {
	ProofByIDCtx(ctx context.Context, id int, isForUpdate bool) (*Proof, error)
	InsertCtx(ctx context.Context, p *Proof) error
}

type GorpMigrationQ interface {
}
