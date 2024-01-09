package handlers

import (
	"context"
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strconv"
)

type proofsLinksByUserDIDRequest struct {
	UserDid string
}

func newProofLinkByUserDIDRequest(r *http.Request) (proofsLinksByUserDIDRequest, error) {
	userDid := chi.URLParam(r, "user_did")
	if userDid == "" {
		return proofsLinksByUserDIDRequest{}, errors.New("user_did is required")
	}

	return proofsLinksByUserDIDRequest{userDid}, validation.Errors{
		"user_did": validation.Validate(userDid, validation.Required),
	}.Filter()
}

func ProofsLinkByUserDID(w http.ResponseWriter, r *http.Request) {
	req, err := newProofLinkByUserDIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	proofsLinks, err := Storage(r).LinkQ().GetProofsLinksByUserID(context.Background(), req.UserDid)
	if err != nil {
		Log(r).WithError(err).Error("failed to get proofs")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if proofsLinks == nil {
		Log(r).WithField("user_did", req.UserDid).Warn("proofs not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var response resources.ProofLinkListResponse

	for _, link := range proofsLinks {
		linkResponse := resources.ProofLinkResponse{
			Data: resources.ProofLink{
				Key: resources.Key{
					ID:   link.ID.String(),
					Type: resources.PROOFS,
				},
				Attributes: resources.ProofLinkAttributes{
					CreatedAt: strconv.FormatInt(link.CreatedAt.Unix(), 10),
					Link:      link.ID.String(),
				},
			},
		}

		response.Data = append(response.Data, linkResponse.Data)
	}

	ape.Render(w, response)
}
