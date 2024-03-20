package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-link-svc/internal/data"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"golang.org/x/exp/utf8string"
)

type proofsByLinkID struct {
	linkID string
}

func newLinkByIDRequest(r *http.Request) (proofsByLinkID, error) {
	linkID := chi.URLParam(r, "link_id")
	if linkID == "" {
		return proofsByLinkID{}, errors.New("link_id is required")
	}

	uuidLinkID, err := uuid.Parse(linkID)
	if err != nil {
		if utf8string.NewString(linkID).IsASCII() {
			return proofsByLinkID{
				linkID: linkID,
			}, nil
		}
		return proofsByLinkID{}, errors.New("invalid link_id")
	}

	return proofsByLinkID{
		linkID: uuidLinkID.String(),
	}, nil
}

func GetLinkByID(w http.ResponseWriter, r *http.Request) {
	req, err := newLinkByIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	link, err := Storage(r).LinkQ().LinkByID(req.linkID, false)
	if err != nil {
		Log(r).WithError(err).Error("failed to get link by UUID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if link == nil {
		Log(r).WithField("link_id", req.linkID).Warn("link not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	links, err := Storage(r).LinksToProofQ().GetLinksToProofsByLinkID(r.Context(), link.ID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get link to proofs")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if len(links) == 0 {
		Log(r).WithField("link_id", req.linkID).Warn("links not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	authorized := true
	verify := false
	who := "Authorized(owner) proofs verification"
	if len(UserClaim(r)) == 0 {
		authorized = false
		who = "Public proofs verification"
	}

	if authorized && !auth.Authenticates(UserClaim(r), auth.UserGrant(link.UserID)) {
		who = "Authorized(not owner) proofs verification"
		verify = true
	}
	Log(r).Debug(who)

	var proofs []data.Proof

	response := resources.LinkResponse{
		Data: resources.Link{
			Key: resources.Key{
				ID:   link.ID,
				Type: resources.LINKS,
			},
			Attributes: resources.LinkAttributes{
				CreatedAt: link.CreatedAt,
				Link:      link.ID,
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

		proofs = append(proofs, *proof)

		included.Add(&resources.Proof{
			Key: resources.Key{
				ID:   proof.ID.String(),
				Type: resources.PROOFS,
			},
			Attributes: resources.ProofAttributes{
				CreatedAt: proof.CreatedAt,
				Creator:   proof.Creator,
				Proof:     string(proof.Proof),
				ProofType: proof.Type,
				OrgId:     proof.OrgID.String(),
				SchemaUrl: proof.SchemaURL,
				Operator:  proof.Operator.String(),
				Field:     proof.Field,
			},
		})
	}
	response.Included = included

	if verify {
		getPointsForVerifyProofs(r, proofs, link.UserID, UserClaim(r)[0].User)
	}

	ape.Render(w, response)
}
