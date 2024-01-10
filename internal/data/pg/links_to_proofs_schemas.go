package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-link-svc/internal/data"
)

func (q LinksToProofQ) SelectAllCtx(ctx context.Context) ([]*data.LinksToProofQ, error) {
	stmt := squirrel.Select("*").From("public.links_to_proofs")

	var links []*data.LinksToProofQ

	if err := q.db.SelectContext(ctx, &links, stmt); err != nil {
		return nil, err
	}

	return links, nil
}

func (q LinksToProofQ) GetLinksToProofsByLinkID(ctx context.Context, linkID uuid.UUID) ([]*data.LinksToProof, error) {
	stmt := squirrel.Select("*").
		From("public.links_to_proofs").
		Where(squirrel.Eq{"link_id": linkID})

	var links []*data.LinksToProof

	if err := q.db.SelectContext(ctx, &links, stmt); err != nil {
		return nil, err
	}

	return links, nil
}

func (q LinksToProofQ) GetLinksToProofsByProofID(ctx context.Context, proofID uuid.UUID) ([]*data.LinksToProof, error) {
	stmt := squirrel.Select("*").
		From("public.links_to_proofs").
		Where(squirrel.Eq{"proof_id": proofID})

	var links []*data.LinksToProof

	if err := q.db.SelectContext(ctx, &links, stmt); err != nil {
		return nil, err
	}

	return links, nil
}
