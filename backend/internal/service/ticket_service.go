package service

import (
	"context"
	"tmaster/internal/apperrors"
	"tmaster/internal/auth"
	"tmaster/internal/constants"
	"tmaster/internal/dto"
)

func (s *Service) CreateTicketService(ctx context.Context, ticketDTO dto.TicketDTO, user_id int, user_Role constants.UserRole) (int, error) {
	if !auth.HasPermission(user_Role, auth.PermissionTicketCreate) {
		return 0, apperrors.NoPermission
	}
	ticket := dto.TicketToModel(ticketDTO, user_id)
	return ticket.ID, s.repo.CreateTicketRepo(ctx, &ticket)
}

func (s *Service) GetOwnTicketsService(ctx context.Context, user_id int, user_Role constants.UserRole, pag_data dto.PaginationData) ([]dto.UserTicketResponse, error) {
	tickets := []dto.UserTicketResponse{}
	if !auth.HasPermission(user_Role, auth.PermissionTicketReadOwn) {
		return tickets, apperrors.NoPermission
	}

	raw_data, err := s.repo.GetOwnTicketsRepo(ctx, user_id, pag_data)

	if err != nil {
		return tickets, err
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

func (s *Service) ClaimNextTicketService(ctx context.Context, user_Role constants.UserRole) error {
	if !auth.HasPermission(user_Role, auth.PermissionTicketAssign) {
		return apperrors.NoPermission
	}

	return nil
}

func (s *Service) CloseTicketService(ctx context.Context, user_Role constants.UserRole) error {
	if !auth.HasPermission(user_Role, auth.PermissionTicketUpdate) {
		return apperrors.NoPermission
	}

	return nil
}

func (s *Service) GetOwnClaimedTicketsService(ctx context.Context, operator_id int, user_Role constants.UserRole) error {
	if !auth.HasPermission(user_Role, auth.PermissionTicketOwnClaimed) {
		return apperrors.NoPermission
	}

	return nil
}
