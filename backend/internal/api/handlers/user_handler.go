package handlers

import (
	"encoding/json"
	"net/http"
	"tmaster/internal/api/maperrors"
	"tmaster/internal/api/utils/helpers"
	"tmaster/internal/api/utils/params"
	"tmaster/internal/dto"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var logData dto.LoginDTO

	if err := json.NewDecoder(r.Body).Decode(&logData); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.LoginService(r.Context(), logData)

	if err != nil {
		helpers.WriteError(w, maperrors.StatusFromErr(err), err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(token); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user dto.RegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.WriteError(w, maperrors.StatusFromErr(err), err.Error())
		return
	}

	err := h.service.RegisterService(r.Context(), user)
	if err != nil {
		helpers.WriteError(w, maperrors.StatusFromErr(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) SetRoleHandler(w http.ResponseWriter, r *http.Request) {
	var ru dto.RequestUser

	_, user_role, err := params.GetAllParams(r)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&ru); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.SetRoleService(r.Context(), user_role, ru.Email, ru.Role); err != nil {
		helpers.WriteError(w, maperrors.StatusFromErr(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
