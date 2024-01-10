package handlers

import (
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strconv"
)

type proofsByUserDIDRequest struct {
	UserDid string
}

func newProofsByUserDIDRequest(r *http.Request) (proofsByUserDIDRequest, error) {
	userDid := chi.URLParam(r, "user_did")
	if userDid == "" {
		return proofsByUserDIDRequest{}, errors.New("user_did is required")
	}

	return proofsByUserDIDRequest{userDid}, validation.Errors{
		"user_did": validation.Validate(userDid, validation.Required),
	}.Filter()
}

func GetProofsByUserDID(w http.ResponseWriter, r *http.Request) {
	req, err := newProofsByUserDIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	proofs, err := Storage(r).ProofQ().ProofsByUserDIDCtx(r.Context(), req.UserDid)
	if err != nil {
		Log(r).WithError(err).Error("failed to get proofs")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if proofs == nil {
		Log(r).WithField("user_did", req.UserDid).Warn("proofs not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var response resources.ProofListResponse

	for _, proof := range proofs {
		proofResponse := resources.ProofResponse{
			Data: resources.Proof{
				Key: resources.Key{
					ID:   proof.ID.String(),
					Type: resources.PROOFS,
				},
				Attributes: resources.ProofAttributes{
					CreatedAt: strconv.FormatInt(proof.CreatedAt.Unix(), 10),
					Creator:   proof.Creator,
					Proof:     string(proof.Proof),
					ProofType: proof.Type,
					OrgId:     proof.OrgID.String(),
				},
			},
		}
		response.Data = append(response.Data, proofResponse.Data)
	}

	ape.Render(w, response)
}
