package handlers

import (
	"encoding/json"
	"net/http"
	"tmaster/internal/api/maperrors"
	"tmaster/internal/api/utils/helpers"
	"tmaster/internal/api/utils/params"
	"tmaster/internal/dto"
	"tmaster/internal/validation"
)

func (h *Handler) CreateTicketHandler(w http.ResponseWriter, r *http.Request) {
	var ticket dto.TicketDTO

	user_id, user_role, err := params.GetAllParams(r)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validation.Validate(ticket); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateTicketService(r.Context(), ticket, user_id, user_role)
	if err != nil {
		helpers.WriteError(w, maperrors.StatusFromErr(err), err.Error())
		return
	}
	resp := dto.IDResponse{
		ID: id,
	}
	helpers.WriteJSON(w, 201, resp)
}

func (h *Handler) GetOwnTicketsHandler(w http.ResponseWriter, r *http.Request) {
	user_id, user_role, err := params.GetAllParams(r)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	pag_data := params.GetPagData(r)

	res, err := h.service.GetOwnTicketsService(r.Context(), user_id, user_role, pag_data)

	if err != nil {
		helpers.WriteError(w, maperrors.StatusFromErr(err), err.Error())
		return
	}

	helpers.WriteJSON(w, 200, res)
}

func (h *Handler) GetNewTicketsHandler(w http.ResponseWriter, r *http.Request) {
	_, user_role, err := params.GetAllParams(r)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	pag_data := params.GetPagData(r)

	res, err := h.service.GetNewTicketsService(r.Context(), user_role, pag_data)
	if err != nil {
		helpers.WriteError(w, maperrors.StatusFromErr(err), err.Error())
		return
	}
	helpers.WriteJSON(w, 200, res)
}

func (h *Handler) ClaimNextTicketHandler(w http.ResponseWriter, r *http.Request) {
	user_id, user_role, err := params.GetAllParams(r)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.ClaimNextTicketService(r.Context(), user_id, user_role)
	if err != nil {
		helpers.WriteError(w, maperrors.StatusFromErr(err), err.Error())
		return
	}
	helpers.WriteJSON(w, 200, res)
}

func (h *Handler) CloseTicketHandler(w http.ResponseWriter, r *http.Request) {
	user_id, user_role, err := params.GetAllParams(r)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.CloseTicketService(r.Context(), user_id, 0, user_role); err != nil {
		helpers.WriteError(w, maperrors.StatusFromErr(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
