package service

import (
	"context"
	"fmt"

	"github.com/mikrocloud/mikrocloud/internal/domain/organizations/repository"
	"github.com/mikrocloud/mikrocloud/internal/domain/users"
)

type OrganizationService struct {
	orgRepo repository.Repository
}

func NewOrganizationService(orgRepo repository.Repository) *OrganizationService {
	return &OrganizationService{
		orgRepo: orgRepo,
	}
}

func (s *OrganizationService) ListOrganizations(ctx context.Context) ([]*users.Organization, error) {
	return s.orgRepo.FindAll(ctx)
}

func (s *OrganizationService) GetOrganization(ctx context.Context, id string) (*users.Organization, error) {
	orgID, err := users.OrganizationIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	return s.orgRepo.FindByID(ctx, orgID)
}

func (s *OrganizationService) GetUserOrganizations(ctx context.Context, userID string) ([]*users.Organization, error) {
	uid, err := users.UserIDFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return s.orgRepo.FindByUserID(ctx, uid)
}

func (s *OrganizationService) CreateOrganization(ctx context.Context, name, slug, description, billingEmail, ownerIDStr string) (*users.Organization, error) {
	ownerID, err := users.UserIDFromString(ownerIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid owner ID: %w", err)
	}

	org := users.NewOrganization(name, slug, ownerID)
	org.UpdateDescription(description)
	if billingEmail != "" {
		org.SetBillingEmail(billingEmail)
	}

	if err := s.orgRepo.Save(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	return org, nil
}

func (s *OrganizationService) UpdateOrganization(ctx context.Context, id, name, description, billingEmail string) (*users.Organization, error) {
	orgID, err := users.OrganizationIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	if name != "" {
		org.UpdateName(name)
	}
	if description != "" {
		org.UpdateDescription(description)
	}
	if billingEmail != "" {
		org.SetBillingEmail(billingEmail)
	}

	if err := s.orgRepo.Save(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	return org, nil
}

func (s *OrganizationService) DeleteOrganization(ctx context.Context, id string) error {
	orgID, err := users.OrganizationIDFromString(id)
	if err != nil {
		return fmt.Errorf("invalid organization ID: %w", err)
	}

	allOrgs, err := s.orgRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to check organization count: %w", err)
	}

	if len(allOrgs) <= 1 {
		return fmt.Errorf("cannot delete the last organization")
	}

	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return fmt.Errorf("organization not found: %w", err)
	}

	if org.Slug() == "default" {
		return fmt.Errorf("cannot delete the default organization")
	}

	if err := s.orgRepo.Delete(ctx, orgID); err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil
}

func (s *OrganizationService) ListOrganizationMembers(ctx context.Context, orgID string) ([]*users.OrganizationMemberWithUser, error) {
	id, err := users.OrganizationIDFromString(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	return s.orgRepo.FindMembersByOrganizationID(ctx, id)
}

func (s *OrganizationService) InviteMember(ctx context.Context, orgID, email, roleStr, invitedByID string) error {
	id, err := users.OrganizationIDFromString(orgID)
	if err != nil {
		return fmt.Errorf("invalid organization ID: %w", err)
	}

	inviterID, err := users.UserIDFromString(invitedByID)
	if err != nil {
		return fmt.Errorf("invalid inviter ID: %w", err)
	}

	userID, err := users.UserIDFromString(email)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	role := users.MemberRole(roleStr)

	existing, err := s.orgRepo.FindMemberByOrganizationAndUserID(ctx, id, userID)
	if err == nil && existing != nil {
		return fmt.Errorf("user is already a member of this organization")
	}

	member := users.NewOrganizationMember(id, userID, role, &inviterID)

	if err := s.orgRepo.SaveMember(ctx, member); err != nil {
		return fmt.Errorf("failed to invite member: %w", err)
	}

	return nil
}

func (s *OrganizationService) UpdateMemberRole(ctx context.Context, memberID, roleStr string) error {
	id, err := users.OrganizationMemberIDFromString(memberID)
	if err != nil {
		return fmt.Errorf("invalid member ID: %w", err)
	}

	member, err := s.orgRepo.FindMemberByID(ctx, id)
	if err != nil {
		return fmt.Errorf("member not found: %w", err)
	}

	role := users.MemberRole(roleStr)
	member.UpdateRole(role)

	if err := s.orgRepo.SaveMember(ctx, member); err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}

	return nil
}

func (s *OrganizationService) RemoveMember(ctx context.Context, memberID string) error {
	id, err := users.OrganizationMemberIDFromString(memberID)
	if err != nil {
		return fmt.Errorf("invalid member ID: %w", err)
	}

	member, err := s.orgRepo.FindMemberByID(ctx, id)
	if err != nil {
		return fmt.Errorf("member not found: %w", err)
	}

	if member.Role() == users.RoleOwner {
		return fmt.Errorf("cannot remove the organization owner")
	}

	if err := s.orgRepo.DeleteMember(ctx, id); err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}

func (s *OrganizationService) TransferOwnership(ctx context.Context, orgID, newOwnerID string) error {
	id, err := users.OrganizationIDFromString(orgID)
	if err != nil {
		return fmt.Errorf("invalid organization ID: %w", err)
	}

	newOwner, err := users.UserIDFromString(newOwnerID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("organization not found: %w", err)
	}

	member, err := s.orgRepo.FindMemberByOrganizationAndUserID(ctx, id, newOwner)
	if err != nil {
		return fmt.Errorf("new owner must be a member of the organization: %w", err)
	}

	oldOwnerMember, err := s.orgRepo.FindMemberByOrganizationAndUserID(ctx, id, org.OwnerID())
	if err == nil {
		oldOwnerMember.UpdateRole(users.RoleAdmin)
		if err := s.orgRepo.SaveMember(ctx, oldOwnerMember); err != nil {
			return fmt.Errorf("failed to update old owner role: %w", err)
		}
	}

	member.UpdateRole(users.RoleOwner)
	if err := s.orgRepo.SaveMember(ctx, member); err != nil {
		return fmt.Errorf("failed to update new owner role: %w", err)
	}

	org.TransferOwnership(newOwner)
	if err := s.orgRepo.Save(ctx, org); err != nil {
		return fmt.Errorf("failed to transfer ownership: %w", err)
	}

	return nil
}
