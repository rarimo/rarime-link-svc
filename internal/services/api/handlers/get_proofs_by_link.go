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

func newLinkByIDRequest(r *http.Request) (proofsByLinkID, error) {
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

func GetLinkByID(w http.ResponseWriter, r *http.Request) {
	req, err := newLinkByIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	link, err := Storage(r).LinkQ().LinkByID(req.LinkID, false)
	if err != nil {
		Log(r).WithError(err).Error("failed to get link by UUID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if link == nil {
		Log(r).WithField("link_id", req.LinkID).Warn("link not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	links, err := Storage(r).LinksToProofQ().GetLinksToProofsByLinkID(r.Context(), req.LinkID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get link to proofs")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if len(links) == 0 {
		Log(r).WithField("link_id", req.LinkID).Warn("links not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	response := resources.LinkResponse{
		Data: resources.Link{
			Key: resources.Key{
				ID:   link.ID.String(),
				Type: resources.LINKS,
			},
			Attributes: resources.LinkAttributes{
				CreatedAt: link.CreatedAt.UTC().String(),
				Link:      link.ID.String(),
			},
		},
	}

	included := resources.Included{}
	for _, linkToProof := range links {
		proof, err := Storage(r).ProofQ().ProofByID(linkToProof.ProofID, false)
		if err != nil {
			Log(r).WithError(err).Error("failed to get proof")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		included.Add(&resources.Proof{
			Key: resources.Key{
				ID:   proof.ID.String(),
				Type: resources.PROOFS,
			},
			Attributes: resources.ProofAttributes{
				CreatedAt: proof.CreatedAt.String(),
				Creator:   proof.Creator,
				Proof:     string(proof.Proof),
				ProofType: proof.Type,
				OrgId:     proof.OrgID.String(),
			},
		})
	}
	response.Included = included

	ape.Render(w, response)
}
