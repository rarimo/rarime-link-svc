package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-link-svc/internal/data"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type proofCreateRequest struct {
	Data resources.ProofCreate `json:"data"`
}

func newProofCreateRequest(r *http.Request) (*proofCreateRequest, error) {
	var req proofCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "failed to decode body")
	}

	if valid := json.Valid([]byte(req.Data.Proof)); !valid {
		return nil, validation.Errors{
			"proof": errors.New("invalid json"),
		}
	}

	if req.Data.ProofType == "" {
		return nil, validation.Errors{
			"type": errors.New("type is required"),
		}
	}

	return &req, nil
}

func CreateProof(w http.ResponseWriter, r *http.Request) {
	req, err := newProofCreateRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !auth.Authenticates(UserClaim(r), auth.UserGrant(req.Data.UserDid)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	orgID, err := uuid.Parse(req.Data.OrgId)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	proof := data.Proof{
		ID:        uuid.New(),
		Creator:   req.Data.UserDid,
		CreatedAt: time.Now().UTC(),
		Proof:     []byte(req.Data.Proof),
		Type:      req.Data.ProofType,
		OrgID:     orgID,
		SchemaURL: req.Data.SchemaUrl,
	}

	err = Storage(r).ProofQ().Insert(&proof)
	if err != nil {
		Log(r).WithError(err).Error("failed to create proof")
		ape.RenderErr(w, problems.InternalError())
		return
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
				SchemaUrl: req.Data.SchemaUrl,
			},
		},
		Included: resources.Included{},
	})
}
