package handlers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/rarimo/rarime-link-svc/internal/data"
	"github.com/rarimo/rarime-link-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strconv"
	"time"
)

type ProofLink struct {
	ProofsIds []int `json:"proofs_ids"`
}

type ProofLinkRequest struct {
	Data ProofLink `json:"data"`
}

func newProofLinkCreateRequest(r *http.Request) (*ProofLinkRequest, error) {
	var req ProofLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "failed to decode body")
	}

	return &req, nil
}

func ProofLinkCreate(w http.ResponseWriter, r *http.Request) {
	req, err := newProofLinkCreateRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	lastIndex, err := Storage(r).ProofLinkQ().GetLastIndex()
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	index := lastIndex + 1
	timestamp := time.Now().UTC()

	var proofs []data.Proof
	for _, proof := range req.Data.ProofsIds {
		p, err := Storage(r).ProofQ().GetProofByID(proof)
		if err != nil {
			ape.RenderErr(w, problems.InternalError())
			return
		}
		proofs = append(proofs, p)
	}

	if len(proofs) != len(req.Data.ProofsIds) {
		ape.RenderErr(w, problems.BadRequest(errors.New("proofs not found"))...)
		return
	}

	var proofsIDs string

	err = Storage(r).ProofLinkQ().Transaction(func(q data.ProofLinkQ) error {
		for _, proof := range req.Data.ProofsIds {
			err = q.InsertCtx(r.Context(), &data.Link{
				ID:        proof,
				Index:     index,
				CreatedAt: timestamp,
			})
			if err != nil {
				return err
			}
			proofsIDs += strconv.Itoa(proof) + ","
		}
		// remove last comma
		proofsIDs = proofsIDs[:len(proofsIDs)-1]

		return nil
	})
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	base64Link := base64.StdEncoding.EncodeToString([]byte(proofsIDs))

	err = Storage(r).LinkToProofQ().InsertCtx(r.Context(), &data.LinkToProof{
		LinkID: index,
	})
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, resources.ProofLinkResponse{
		Data: resources.ProofLink{
			Key: resources.Key{
				ID:   strconv.Itoa(index),
				Type: resources.PROOFS,
			},
			Attributes: resources.ProofLinkAttributes{
				Base64Link: base64Link,
				CreatedAt:  timestamp.String(),
			},
		},
		Included: resources.Included{},
	})
}
