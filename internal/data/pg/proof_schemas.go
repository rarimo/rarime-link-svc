package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-link-svc/internal/data"
)

func (q ProofQ) SelectAllCtx(ctx context.Context) ([]*data.Proof, error) {
	stmt := squirrel.Select("*").From("public.proofs")

	var proofs []*data.Proof

	if err := q.db.SelectContext(ctx, &proofs, stmt); err != nil {
		return nil, err
	}

	return proofs, nil
}

func (q ProofQ) ProofsByCreatorCtx(ctx context.Context, userDID string) ([]data.Proof, error) {
	stmt := squirrel.Select("*").From("public.proofs").Where(squirrel.Eq{"creator": userDID})

	var proofs []data.Proof

	if err := q.db.SelectContext(ctx, &proofs, stmt); err != nil {
		return nil, err
	}

	return proofs, nil
}
