package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rarimo/rarime-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-link-svc/internal/data"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"golang.org/x/exp/utf8string"
)

const maxLinkNameLen = 50

type ProofLink struct {
	UserDID   string      `json:"user_did"`
	ProofsIds []uuid.UUID `json:"proofs_ids"`
	LinkName  *string     `json:"link_name,omitempty"`
}

type ProofLinkRequest struct {
	Data ProofLink `json:"data"`
}

func newProofLinkCreateRequest(r *http.Request) (*ProofLinkRequest, error) {
	var req ProofLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "failed to decode body")
	}

	if req.Data.LinkName != nil {
		if utf8string.NewString(*req.Data.LinkName).IsASCII() {
			if len([]rune(*req.Data.LinkName)) > maxLinkNameLen {
				return nil, errors.New("link name length is too big")
			}
		}
	}

	return &req, nil
}

func CreateProofLink(w http.ResponseWriter, r *http.Request) {
	req, err := newProofLinkCreateRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !auth.Authenticates(UserClaim(r), auth.UserGrant(req.Data.UserDID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	var (
		timestamp      = time.Now().UTC()
		linkID         = uuid.New()
		proofs         []data.Proof
		proofNotFound  = errors.New("proof not found")
		invalidCreator = errors.New("invalid proof creator")
		linkNameStr    *string
	)

	err = Storage(r).LinkQ().Transaction(func(q data.LinkQ) error {
		linkName := sql.NullString{}
		if req.Data.LinkName != nil {
			linkName.String = *req.Data.LinkName
			linkName.Valid = true
			linkNameStr = req.Data.LinkName
		}
		err = q.Insert(&data.Link{
			ID:        linkID,
			UserID:    req.Data.UserDID,
			CreatedAt: timestamp,
			Name:      linkName,
		})
		if err != nil {
			ape.RenderErr(w, problems.InternalError())
			return err
		}

		for _, proofID := range req.Data.ProofsIds {
			p, err := Storage(r).ProofQ().ProofByID(proofID, false)
			if err != nil {
				ape.RenderErr(w, problems.InternalError())
				return err
			}
			if p == nil {
				ape.RenderErr(w, problems.NotFound())
				return proofNotFound
			}

			if p.Creator != req.Data.UserDID {
				ape.RenderErr(w, problems.Unauthorized())
				return invalidCreator
			}

			proofs = append(proofs, *p)

			err = q.InsertCtxLinkToProof(r.Context(), data.LinksToProof{
				LinkID:  linkID,
				ProofID: proofID,
			})
			if err != nil {
				ape.RenderErr(w, problems.InternalError())
				return err
			}
		}

		return nil
	})

	if err != nil {
		Log(r).WithError(err).Error("failed to create proof link entry")
		// Response error rendered before
		return
	}

	ape.Render(w, resources.LinkResponse{
		Data: resources.Link{
			Key: resources.Key{
				Type: resources.LINKS,
			},
			Attributes: resources.LinkAttributes{
				Link:      linkID.String(),
				CreatedAt: timestamp,
				LinkName:  linkNameStr,
			},
		},
		Included: resources.Included{},
	})
}
