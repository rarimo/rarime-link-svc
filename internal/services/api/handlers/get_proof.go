package handlers

import (
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strconv"
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

	ape.Render(w, resources.ProofResponse{
		Data: resources.Proof{
			Key: resources.Key{
				ID:   proof.ID.String(),
				Type: resources.PROOFS,
			},
			Attributes: resources.ProofAttributes{
				CreatedAt: strconv.FormatInt(proof.CreatedAt.Unix(), 10),
				Creator:   proof.Creator,
				Proof:     string(proof.Proof),
				Type:      proof.Type,
			},
		},
		Included: resources.Included{},
	})

}
