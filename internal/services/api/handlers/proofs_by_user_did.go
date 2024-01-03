package handlers

import (
	"encoding/base64"
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strconv"
	"strings"
)

type proofsByUserDIDRequest struct {
	UserDid string
}

type proofsByUserDIDResponse struct {
	Data         []resources.Proof  `json:"data"`
	Include      resources.Included `json:"included"`
	Base64Proofs string             `json:"base64_proofs"`
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

func ProofsByUserDID(w http.ResponseWriter, r *http.Request) {
	req, err := newProofsByUserDIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	proofs, err := Storage(r).ProofQ().ProofsByUserDIDCtx(r.Context(), req.UserDid, false)
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

	resourceProofs := make([]resources.Proof, len(proofs))
	proofIDs := make([]string, len(proofs))

	for i, proof := range proofs {
		resourceProofs[i] = resources.Proof{
			Key: resources.Key{
				ID:   strconv.FormatInt(int64(proof.ID), 10),
				Type: resources.PROOFS,
			},
			Attributes: resources.ProofAttributes{
				CreatedAt: strconv.FormatInt(proof.CreatedAt.Unix(), 10),
				Creator:   proof.Creator,
				Proof:     string(proof.Proof),
			},
		}
		proofIDs[i] = strconv.FormatInt(int64(proof.ID), 10)
	}

	base64Data := req.UserDid + ":" + strings.Join(proofIDs, ",")
	base64Proofs := base64.StdEncoding.EncodeToString([]byte(base64Data))

	response := proofsByUserDIDResponse{
		Data:         resourceProofs,
		Base64Proofs: base64Proofs,
		Include:      resources.Included{},
	}

	ape.Render(w, response)
}
