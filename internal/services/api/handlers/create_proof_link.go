package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-link-svc/internal/data"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"time"
)

type ProofLink struct {
	ProofsIds []uuid.UUID `json:"proofs_ids"`
}

type ProofLinkRequest struct {
	Data ProofLink `json:"data"`
}

func newProofLinkCreateRequest(r *http.Request) (*ProofLinkRequest, error) {
	var req ProofLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "failed to decode body")
	}

	return &req, nil
}

func CreateProofLink(w http.ResponseWriter, r *http.Request) {
	req, err := newProofLinkCreateRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	timestamp := time.Now().UTC()
	linkID := uuid.New()

	var proofs []data.Proof
	err = Storage(r).LinkQ().Transaction(func(q data.LinkQ) error {

		err = q.Insert(&data.Link{
			ID:        linkID,
			UserID:    UserID(r),
			CreatedAt: timestamp,
		})
		if err != nil {
			ape.RenderErr(w, problems.InternalError())
			return err
		}

		for _, proofID := range req.Data.ProofsIds {
			p, err := Storage(r).ProofQ().ProofByID(proofID, false)
			if err != nil {
				ape.RenderErr(w, problems.InternalError())
				return err
			}

			proofs = append(proofs, *p)

			err = q.InsertCtxLinkToProof(r.Context(), data.LinksToProof{
				LinkID:  linkID,
				ProofID: proofID,
			})
			if err != nil {
				ape.RenderErr(w, problems.InternalError())
				return err
			}
		}

		if len(proofs) != len(req.Data.ProofsIds) {
			ape.RenderErr(w, problems.BadRequest(errors.New("proofs not found"))...)
			return errors.New("proofs not found")
		}

		return nil
	})

	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, resources.LinkResponse{
		Data: resources.Link{
			Key: resources.Key{
				Type: resources.LINKS,
			},
			Attributes: resources.LinkAttributes{
				Link:      linkID.String(),
				CreatedAt: timestamp,
			},
		},
		Included: resources.Included{},
	})
}
