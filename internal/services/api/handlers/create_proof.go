package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	if err := validation.Validate(req.Data.Operator, ValidationOperator); err != nil {
		return nil, validation.Errors{
			"data/operator": errors.Wrap(err, "invalid operator value"),
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
		Operator:  data.MustProofOperatorFromString(req.Data.Operator),
		Field:     req.Data.Field,
	}

	err = Storage(r).ProofQ().Insert(&proof)
	if err != nil {
		Log(r).WithError(err).Error("failed to create proof")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	pointsError := Points(r).FulfillEvent(context.Background(), points.FulfillEventRequest{
		UserDID:   req.Data.UserDid,
		EventType: fmt.Sprintf("generate_proof_%s", req.Data.ProofType),
	})
	if !isNormalFlowError(pointsError) {
		Log(r).WithError(pointsError).Errorf("error occurred while fulfilling event")
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
				Operator:  proof.Operator.String(),
				Field:     req.Data.Field,
			},
		},
		Included: resources.Included{},
	})
}

// There are cases when we can safely ignore error and continue the flow, because
// those codes can occur in normal flow
func isNormalFlowError(err *points.Error) bool {
	if err == nil {
		return true
	}
	switch err.Code {
	case points.CodeEventExpired, points.CodeEventDisabled, points.CodeEventNotFound:
		return true
	}
	return false
}
