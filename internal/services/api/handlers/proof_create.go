package handlers

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-link-svc/internal/data"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strconv"
	"time"
)

func newProofCreateRequest(r *http.Request) (*resources.ProofCreate, error) {
	var req resources.ProofCreate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "failed to decode body")
	}

	if valid := json.Valid([]byte(req.Proof)); !valid {
		return nil, validation.Errors{
			"proof": errors.New("invalid json"),
		}
	}

	if req.Type == "" {
		return nil, validation.Errors{
			"type": errors.New("type is required"),
		}
	}

	return &req, nil
}

func ProofCreate(w http.ResponseWriter, r *http.Request) {
	req, err := newProofCreateRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	proof := data.Proof{
		Creator:   UserID(r),
		CreatedAt: time.Now().UTC(),
		Proof:     []byte(req.Proof),
		Type:      req.Type,
	}

	err = Storage(r).ProofQ().InsertCtx(r.Context(), &proof)
	if err != nil {
		Log(r).WithError(err).Error("failed to create proof")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, resources.ProofResponse{
		Data: resources.Proof{
			Key: resources.Key{
				ID:   strconv.FormatInt(int64(proof.ID), 10),
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
