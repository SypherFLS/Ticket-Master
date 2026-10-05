package service

import (
	"context"
	"tmaster/internal/apperrors"
	"tmaster/internal/auth"
	"tmaster/internal/constants"
	"tmaster/internal/constants/params"
	"tmaster/internal/dto"
)

func (s *Service) CreateTicketService(ctx context.Context, ticketDTO dto.TicketDTO, user_id int, user_Role constants.UserRole) (int, error) {
	if !auth.HasPermission(user_Role, auth.PermissionTicketCreate) {
		return 0, apperrors.NoPermission
	}
	ticket := dto.TicketToModel(ticketDTO, user_id)
	if err := s.repo.CreateTicketRepo(ctx, &ticket); err != nil {
		return -1, err
	}
	logger := params.GetLogger(ctx)
	logger.Info(
		"ticket created",
		"user_id", user_id,
		"ticket_id", ticket.ID,
	)
	return ticket.ID, nil
}

func (s *Service) GetOwnTicketsService(ctx context.Context, user_id int, user_Role constants.UserRole, pag_data dto.PaginationData) ([]dto.UserTicketResponse, error) {
	tickets := []dto.UserTicketResponse{}
	if !auth.HasPermission(user_Role, auth.PermissionTicketReadOwn) {
		return tickets, apperrors.NoPermission
	}

	raw_data, err := s.repo.GetOwnTicketsRepo(ctx, user_id, pag_data)

	if err != nil {
		return []dto.UserTicketResponse{}, err
	}

	tickets = dto.ManyMTUT(raw_data)

	return tickets, nil
}

func (s *Service) GetNewTicketsService(ctx context.Context, user_Role constants.UserRole, pag_data dto.PaginationData) ([]dto.OperatorTicketResponse, error) {
	if !auth.HasPermission(user_Role, auth.PermissionTicketReadQueue) {
		return nil, apperrors.NoPermission
	}

	raw_data, err := s.repo.GetNewTickets(ctx, pag_data.Limit)

	if err != nil {
		return nil, err
	}

	res := dto.ManyMTOT(raw_data)
	return res, nil
}

func (s *Service) ClaimNextTicketService(ctx context.Context, operator_id int, user_Role constants.UserRole) (dto.OperatorTicketResponse, error) {

	if !auth.HasPermission(user_Role, auth.PermissionTicketAssign) {
		return dto.OperatorTicketResponse{}, apperrors.NoPermission
	}

	raw_data, err := s.repo.ClaimNextTicketRepo(ctx, operator_id)
	if err != nil {
		return dto.OperatorTicketResponse{}, err
	}

	logger := params.GetLogger(ctx)
	duration := raw_data.ClaimedAt.Sub(raw_data.CreatedAt)
	logger.Info(
		"ticket claimed",
		"operator_id", operator_id,
		"ticket_id", raw_data.ID,
		"created_at", raw_data.CreatedAt,
		"claimed_at", raw_data.ClaimedAt,
		"awaiting_time", duration,
	)

	resp := dto.ModelToOperatorTicket(*raw_data)

	return resp, nil
}

func (s *Service) CloseTicketService(ctx context.Context, operator_id int, ticket_id int, user_Role constants.UserRole) error {
	if !auth.HasPermission(user_Role, auth.PermissionTicketUpdate) {
		return apperrors.NoPermission
	}

	if err := s.repo.CloseTicketRepo(ctx, ticket_id, operator_id); err != nil {
		return err
	}

	logger := params.GetLogger(ctx)
	logger.Info(
		"ticket closed",
		"operator_id", operator_id,
		"ticket_id", ticket_id,
	)

	return nil
}

func (s *Service) GetOwnClaimedTicketsService(ctx context.Context, operator_id int, user_Role constants.UserRole) (dto.OperatorTicketResponse, error) {
	if !auth.HasPermission(user_Role, auth.PermissionTicketOwnClaimed) {
		return dto.OperatorTicketResponse{}, apperrors.NoPermission
	}

	raw_data, err := s.repo.GetClaimedTicket(ctx, operator_id)

	if err != nil {
		return dto.OperatorTicketResponse{}, err
	}

	resp := dto.ModelToOperatorTicket(raw_data)
	return resp, nil
}
