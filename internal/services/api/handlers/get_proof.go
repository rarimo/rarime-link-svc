package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-link-svc/internal/data"
	"github.com/rarimo/rarime-link-svc/resources"
	points "github.com/rarimo/rarime-points-svc/pkg/connector"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type proofByIDRequest struct {
	ID uuid.UUID
}

func newProofByIDRequest(r *http.Request) (*proofByIDRequest, error) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to int")
	}

	return &proofByIDRequest{id}, validation.Errors{
		"id": validation.Validate(id, validation.Required),
	}.Filter()
}

func ProofByID(w http.ResponseWriter, r *http.Request) {
	req, err := newProofByIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	proof, err := Storage(r).ProofQ().ProofByID(req.ID, false)
	if err != nil {
		Log(r).WithError(err).Error("failed to get proof")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if proof == nil {
		Log(r).WithField("id", req.ID).Warn("proof not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	// User is not eligible for rewards for their own proof verification
	if !auth.Authenticates(UserClaim(r), auth.UserGrant(proof.Creator)) {
		Log(r).Debug("Proof verification")
		// while auth is mandatory on this endpoint, it is guaranteed that we have at
		// least one claim with valid user DID
		verifier := UserClaim(r)[0].User
		ok := getPointsForVerifyProofs(r, []data.Proof{*proof}, proof.Creator, verifier)
		if !ok {
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	ape.Render(w, resources.ProofResponse{
		Data: resources.Proof{
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
		},
		Included: resources.Included{},
	})

}

func getPointsForVerifyProofs(r *http.Request, proofs []data.Proof, creator, verifier string) (ok bool) {
	req := points.FulfillVerifyProofEventRequest{
		UserDID:     creator,
		ProofTypes:  make([]string, len(proofs)),
		VerifierDID: verifier,
	}

	for i, proof := range proofs {
		req.ProofTypes[i] = proof.Type
	}

	pointsError := Points(r).FulfillVerifyProofEvent(context.Background(), req)
	if !isNormalFlowError(pointsError) {
		Log(r).WithError(pointsError).Errorf("Error occurred while event fulfillment for proofs %s", proofs)
		return false
	}

	return true
}
