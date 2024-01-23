package handlers

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/urlval"
)

type proofsLinksByUserDIDRequest struct {
	UserDid string
}

func newGetLinksRequest(r *http.Request) (proofsLinksByUserDIDRequest, error) {
	request := proofsLinksByUserDIDRequest{}
	if err := urlval.DecodeSilently(r.URL.Query(), &request); err != nil {
		return request, err
	}
	request.UserDid = r.URL.Query().Get("filter[did]")

	return proofsLinksByUserDIDRequest{request.UserDid}, validation.Errors{
		"did": validation.Validate(request.UserDid, validation.Required),
	}.Filter()
}

func GetLinks(w http.ResponseWriter, r *http.Request) {
	req, err := newGetLinksRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !auth.Authenticates(UserClaim(r), auth.UserGrant(req.UserDid)) {
		ape.RenderErr(w, problems.Unauthorized())
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

	var response resources.LinkListResponse
	for _, link := range proofsLinks {
		var linkName *string
		if link.Name.Valid {
			linkName = &link.Name.String
		}
		linkResponse := resources.LinkResponse{
			Data: resources.Link{
				Key: resources.Key{
					ID:   link.ID.String(),
					Type: resources.LINKS,
				},
				Attributes: resources.LinkAttributes{
					CreatedAt: link.CreatedAt,
					Link:      link.ID.String(),
					LinkName:  linkName,
				},
			},
		}
		response.Data = append(response.Data, linkResponse.Data)
	}

	ape.Render(w, response)
}
