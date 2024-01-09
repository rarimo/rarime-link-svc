package handlers

import (
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type proofsByLinkID struct {
	LinkID uuid.UUID
}

func newProofLinkByIDRequest(r *http.Request) (proofsByLinkID, error) {
	linkID := chi.URLParam(r, "link_id")
	if linkID == "" {
		return proofsByLinkID{}, errors.New("user_did is required")
	}

	uuidLinkID, err := uuid.Parse(linkID)
	if err != nil {
		return proofsByLinkID{}, errors.New("invalid link_id")
	}

	return proofsByLinkID{uuidLinkID}, nil
}

func ProofLinkByID(w http.ResponseWriter, r *http.Request) {
	req, err := newProofLinkByIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	proofs, err := Storage(r).LinkToProofQ().GetProofsByLinkID(req.LinkID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get link to proofs")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if proofs == nil {
		Log(r).WithField("link_id", req.LinkID).Warn("link not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var response resources.ProofLinkByIdListResponse
	for _, proof := range proofs {
		proof, err := Storage(r).ProofQ().GetProofByID(proof.ProofID)
		if err != nil {
			Log(r).WithError(err).Error("failed to get proof")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		proofsResponse := resources.ProofLinkByIdResponse{
			Data: resources.ProofLinkById{
				Key: resources.Key{
					ID:   proof.ID.String(),
					Type: resources.PROOFS,
				},
				Attributes: resources.ProofLinkByIdAttributes{
					CreatedAt: proof.CreatedAt.String(),
					Creator:   proof.Creator,
					Link:      req.LinkID.String(),
					Proof:     string(proof.Proof),
					Type:      proof.Type,
				},
			},
		}

		response.Data = append(response.Data, proofsResponse.Data)
	}

	ape.Render(w, response)
}
