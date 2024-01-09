package pg

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-link-svc/internal/data"
)

func (q LinkQ) SelectAllCtx(ctx context.Context) ([]*data.Link, error) {
	stmt := squirrel.Select("*").From("public.links")

	var links []*data.Link

	if err := q.db.SelectContext(ctx, &links, stmt); err != nil {
		return nil, err
	}

	return links, nil
}

func (q LinkQ) Transaction(fn func(db data.LinkQ) error) error {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q LinkQ) InsertCtxLinkToProof(ctx context.Context, linksToProof data.LinksToProof) error {
	stmt := squirrel.Insert("public.links_to_proofs").SetMap(map[string]interface{}{
		"link_id":  linksToProof.LinkID,
		"proof_id": linksToProof.ProofID,
	})

	err := q.db.ExecContext(ctx, stmt)
	return err
}

func (q LinkQ) GetProofsLinksByUserID(ctx context.Context, userID string) ([]*data.Link, error) {
	stmt := squirrel.Select("*").
		From("public.links").
		Where(squirrel.Eq{"user_id": userID})

	var links []*data.Link

	if err := q.db.SelectContext(ctx, &links, stmt); err != nil {
		return nil, err
	}

	return links, nil
}
