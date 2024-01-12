package handlers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type proofsByUserDIDRequest struct {
	UserDid string
}

func newProofsRequest(r *http.Request) (proofsByUserDIDRequest, error) {
	request := proofsByUserDIDRequest{}
	if err := urlval.DecodeSilently(r.URL.Query(), &request); err != nil {
		return request, err
	}
	request.UserDid = r.URL.Query().Get("filter[did]")

	return proofsByUserDIDRequest{request.UserDid}, validation.Errors{
		"did": validation.Validate(request.UserDid, validation.Required),
	}.Filter()
}

func GetProofs(w http.ResponseWriter, r *http.Request) {
	req, err := newProofsRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !auth.Authenticates(UserClaim(r), auth.UserGrant(req.UserDid)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	proofs, err := Storage(r).ProofQ().ProofsByCreatorCtx(r.Context(), req.UserDid)
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
					CreatedAt: proof.CreatedAt,
					Creator:   proof.Creator,
					Proof:     string(proof.Proof),
					ProofType: proof.Type,
					OrgId:     proof.OrgID.String(),
					SchemaUrl: proof.SchemaURL,
				},
			},
		}
		response.Data = append(response.Data, proofResponse.Data)
	}

	ape.Render(w, response)
}
