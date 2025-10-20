package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mikrocloud/mikrocloud/internal/domain/users"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*users.Organization, error)
	FindByID(ctx context.Context, id users.OrganizationID) (*users.Organization, error)
	FindByUserID(ctx context.Context, userID users.UserID) ([]*users.Organization, error)
	Save(ctx context.Context, org *users.Organization) error
	Delete(ctx context.Context, id users.OrganizationID) error

	FindMembersByOrganizationID(ctx context.Context, orgID users.OrganizationID) ([]*users.OrganizationMemberWithUser, error)
	FindMemberByID(ctx context.Context, memberID users.OrganizationMemberID) (*users.OrganizationMember, error)
	FindMemberByOrganizationAndUserID(ctx context.Context, orgID users.OrganizationID, userID users.UserID) (*users.OrganizationMember, error)
	SaveMember(ctx context.Context, member *users.OrganizationMember) error
	DeleteMember(ctx context.Context, memberID users.OrganizationMemberID) error
}

type SQLiteOrganizationRepository struct {
	db *sql.DB
}

func NewSQLiteOrganizationRepository(db *sql.DB) Repository {
	return &SQLiteOrganizationRepository{db: db}
}

func (r *SQLiteOrganizationRepository) FindAll(ctx context.Context) ([]*users.Organization, error) {
	query := `
		SELECT id, name, slug, description, owner_id, billing_email, plan, status, created_at, updated_at
		FROM organizations
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query organizations: %w", err)
	}
	defer rows.Close()

	var organizations []*users.Organization
	for rows.Next() {
		var row organizationRow
		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Slug,
			&row.Description,
			&row.OwnerID,
			&row.BillingEmail,
			&row.Plan,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization row: %w", err)
		}

		org, err := r.mapRowToOrganization(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map organization: %w", err)
		}

		organizations = append(organizations, org)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organization rows: %w", err)
	}

	return organizations, nil
}

func (r *SQLiteOrganizationRepository) FindByID(ctx context.Context, id users.OrganizationID) (*users.Organization, error) {
	query := `
		SELECT id, name, slug, description, owner_id, billing_email, plan, status, created_at, updated_at
		FROM organizations
		WHERE id = ?
	`

	var row organizationRow
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&row.ID,
		&row.Name,
		&row.Slug,
		&row.Description,
		&row.OwnerID,
		&row.BillingEmail,
		&row.Plan,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to find organization by ID: %w", err)
	}

	return r.mapRowToOrganization(row)
}

func (r *SQLiteOrganizationRepository) FindByUserID(ctx context.Context, userID users.UserID) ([]*users.Organization, error) {
	query := `
		SELECT o.id, o.name, o.slug, o.description, o.owner_id, o.billing_email, o.plan, o.status, o.created_at, o.updated_at
		FROM organizations o
		LEFT JOIN organization_members om ON o.id = om.organization_id
		WHERE o.owner_id = ? OR om.user_id = ?
		GROUP BY o.id
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID.String(), userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query organizations by user: %w", err)
	}
	defer rows.Close()

	var organizations []*users.Organization
	for rows.Next() {
		var row organizationRow
		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Slug,
			&row.Description,
			&row.OwnerID,
			&row.BillingEmail,
			&row.Plan,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization row: %w", err)
		}

		org, err := r.mapRowToOrganization(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map organization: %w", err)
		}

		organizations = append(organizations, org)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organization rows: %w", err)
	}

	return organizations, nil
}

func (r *SQLiteOrganizationRepository) Save(ctx context.Context, org *users.Organization) error {
	_, err := r.FindByID(ctx, org.ID())

	if err != nil {
		query := `
			INSERT INTO organizations (id, name, slug, description, owner_id, billing_email, plan, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err := r.db.ExecContext(ctx, query,
			org.ID().String(),
			org.Name(),
			org.Slug(),
			org.Description(),
			org.OwnerID().String(),
			org.BillingEmail(),
			org.Plan(),
			org.Status(),
			org.CreatedAt().Format(time.RFC3339),
			org.UpdatedAt().Format(time.RFC3339),
		)
		if err != nil {
			return fmt.Errorf("failed to insert organization: %w", err)
		}
		return nil
	}

	query := `
		UPDATE organizations 
		SET name = ?, slug = ?, description = ?, owner_id = ?, billing_email = ?, plan = ?, status = ?, updated_at = ?
		WHERE id = ?
	`
	_, err = r.db.ExecContext(ctx, query,
		org.Name(),
		org.Slug(),
		org.Description(),
		org.OwnerID().String(),
		org.BillingEmail(),
		org.Plan(),
		org.Status(),
		org.UpdatedAt().Format(time.RFC3339),
		org.ID().String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}

	return nil
}

func (r *SQLiteOrganizationRepository) Delete(ctx context.Context, id users.OrganizationID) error {
	query := `DELETE FROM organizations WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("organization not found: %s", id.String())
	}

	return nil
}

type organizationRow struct {
	ID           string
	Name         string
	Slug         string
	Description  string
	OwnerID      string
	BillingEmail string
	Plan         string
	Status       string
	CreatedAt    string
	UpdatedAt    string
}

func (r *SQLiteOrganizationRepository) mapRowToOrganization(row organizationRow) (*users.Organization, error) {
	orgID, err := users.OrganizationIDFromString(row.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	ownerID, err := users.UserIDFromString(row.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner ID: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at timestamp: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at timestamp: %w", err)
	}

	return users.ReconstructOrganization(
		orgID,
		row.Name,
		row.Slug,
		row.Description,
		ownerID,
		row.BillingEmail,
		users.OrganizationPlan(row.Plan),
		users.OrganizationStatus(row.Status),
		createdAt,
		updatedAt,
	), nil
}

func (r *SQLiteOrganizationRepository) FindMembersByOrganizationID(ctx context.Context, orgID users.OrganizationID) ([]*users.OrganizationMemberWithUser, error) {
	query := `
		SELECT 
			om.id, om.organization_id, om.user_id, om.role, om.invited_by, om.invited_at, om.joined_at, om.status,
			u.id, u.email, u.password_hash, u.name, u.username, u.avatar_url, u.status, u.email_verified_at, u.last_login_at, u.timezone, u.created_at, u.updated_at
		FROM organization_members om
		JOIN users u ON om.user_id = u.id
		WHERE om.organization_id = ?
		ORDER BY om.joined_at DESC, om.invited_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query organization members: %w", err)
	}
	defer rows.Close()

	var members []*users.OrganizationMemberWithUser
	for rows.Next() {
		var memberRow organizationMemberRow
		var userRow userRow

		err := rows.Scan(
			&memberRow.ID,
			&memberRow.OrganizationID,
			&memberRow.UserID,
			&memberRow.Role,
			&memberRow.InvitedBy,
			&memberRow.InvitedAt,
			&memberRow.JoinedAt,
			&memberRow.Status,
			&userRow.ID,
			&userRow.Email,
			&userRow.PasswordHash,
			&userRow.Name,
			&userRow.Username,
			&userRow.AvatarURL,
			&userRow.Status,
			&userRow.EmailVerifiedAt,
			&userRow.LastLoginAt,
			&userRow.Timezone,
			&userRow.CreatedAt,
			&userRow.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization member row: %w", err)
		}

		member, err := r.mapRowToOrganizationMember(memberRow)
		if err != nil {
			return nil, fmt.Errorf("failed to map organization member: %w", err)
		}

		user, err := r.mapRowToUser(userRow)
		if err != nil {
			return nil, fmt.Errorf("failed to map user: %w", err)
		}

		members = append(members, &users.OrganizationMemberWithUser{
			Member: member,
			User:   user,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organization member rows: %w", err)
	}

	return members, nil
}

func (r *SQLiteOrganizationRepository) FindMemberByID(ctx context.Context, memberID users.OrganizationMemberID) (*users.OrganizationMember, error) {
	query := `
		SELECT id, organization_id, user_id, role, invited_by, invited_at, joined_at, status
		FROM organization_members
		WHERE id = ?
	`

	var row organizationMemberRow
	err := r.db.QueryRowContext(ctx, query, memberID.String()).Scan(
		&row.ID,
		&row.OrganizationID,
		&row.UserID,
		&row.Role,
		&row.InvitedBy,
		&row.InvitedAt,
		&row.JoinedAt,
		&row.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization member not found: %s", memberID.String())
		}
		return nil, fmt.Errorf("failed to find organization member by ID: %w", err)
	}

	return r.mapRowToOrganizationMember(row)
}

func (r *SQLiteOrganizationRepository) FindMemberByOrganizationAndUserID(ctx context.Context, orgID users.OrganizationID, userID users.UserID) (*users.OrganizationMember, error) {
	query := `
		SELECT id, organization_id, user_id, role, invited_by, invited_at, joined_at, status
		FROM organization_members
		WHERE organization_id = ? AND user_id = ?
	`

	var row organizationMemberRow
	err := r.db.QueryRowContext(ctx, query, orgID.String(), userID.String()).Scan(
		&row.ID,
		&row.OrganizationID,
		&row.UserID,
		&row.Role,
		&row.InvitedBy,
		&row.InvitedAt,
		&row.JoinedAt,
		&row.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization member not found")
		}
		return nil, fmt.Errorf("failed to find organization member: %w", err)
	}

	return r.mapRowToOrganizationMember(row)
}

func (r *SQLiteOrganizationRepository) SaveMember(ctx context.Context, member *users.OrganizationMember) error {
	var invitedBy *string
	if member.InvitedBy() != nil {
		val := member.InvitedBy().String()
		invitedBy = &val
	}

	var invitedAt *string
	if member.InvitedAt() != nil {
		val := member.InvitedAt().Format(time.RFC3339)
		invitedAt = &val
	}

	var joinedAt *string
	if member.JoinedAt() != nil {
		val := member.JoinedAt().Format(time.RFC3339)
		joinedAt = &val
	}

	_, err := r.FindMemberByID(ctx, member.ID())

	if err != nil {
		query := `
			INSERT INTO organization_members (id, organization_id, user_id, role, invited_by, invited_at, joined_at, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err := r.db.ExecContext(ctx, query,
			member.ID().String(),
			member.OrganizationID().String(),
			member.UserID().String(),
			member.Role(),
			invitedBy,
			invitedAt,
			joinedAt,
			member.Status(),
		)
		if err != nil {
			return fmt.Errorf("failed to insert organization member: %w", err)
		}
		return nil
	}

	query := `
		UPDATE organization_members 
		SET organization_id = ?, user_id = ?, role = ?, invited_by = ?, invited_at = ?, joined_at = ?, status = ?
		WHERE id = ?
	`
	_, err = r.db.ExecContext(ctx, query,
		member.OrganizationID().String(),
		member.UserID().String(),
		member.Role(),
		invitedBy,
		invitedAt,
		joinedAt,
		member.Status(),
		member.ID().String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update organization member: %w", err)
	}

	return nil
}

func (r *SQLiteOrganizationRepository) DeleteMember(ctx context.Context, memberID users.OrganizationMemberID) error {
	query := `DELETE FROM organization_members WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, memberID.String())
	if err != nil {
		return fmt.Errorf("failed to delete organization member: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("organization member not found: %s", memberID.String())
	}

	return nil
}

type organizationMemberRow struct {
	ID             string
	OrganizationID string
	UserID         string
	Role           string
	InvitedBy      sql.NullString
	InvitedAt      sql.NullString
	JoinedAt       sql.NullString
	Status         string
}

type userRow struct {
	ID              string
	Email           string
	PasswordHash    string
	Name            string
	Username        sql.NullString
	AvatarURL       sql.NullString
	Status          string
	EmailVerifiedAt sql.NullString
	LastLoginAt     sql.NullString
	Timezone        string
	CreatedAt       string
	UpdatedAt       string
}

func (r *SQLiteOrganizationRepository) mapRowToOrganizationMember(row organizationMemberRow) (*users.OrganizationMember, error) {
	memberID, err := users.OrganizationMemberIDFromString(row.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization member ID: %w", err)
	}

	orgID, err := users.OrganizationIDFromString(row.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	userID, err := users.UserIDFromString(row.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var invitedBy *users.UserID
	if row.InvitedBy.Valid {
		id, err := users.UserIDFromString(row.InvitedBy.String)
		if err != nil {
			return nil, fmt.Errorf("invalid invited_by user ID: %w", err)
		}
		invitedBy = &id
	}

	var invitedAt *time.Time
	if row.InvitedAt.Valid {
		t, err := time.Parse(time.RFC3339, row.InvitedAt.String)
		if err != nil {
			return nil, fmt.Errorf("invalid invited_at timestamp: %w", err)
		}
		invitedAt = &t
	}

	var joinedAt *time.Time
	if row.JoinedAt.Valid {
		t, err := time.Parse(time.RFC3339, row.JoinedAt.String)
		if err != nil {
			return nil, fmt.Errorf("invalid joined_at timestamp: %w", err)
		}
		joinedAt = &t
	}

	return users.ReconstructOrganizationMember(
		memberID,
		orgID,
		userID,
		users.MemberRole(row.Role),
		invitedBy,
		invitedAt,
		joinedAt,
		users.MemberStatus(row.Status),
	), nil
}

func (r *SQLiteOrganizationRepository) mapRowToUser(row userRow) (*users.User, error) {
	userID, err := users.UserIDFromString(row.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	email, err := users.NewEmail(row.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	var username *users.Username
	if row.Username.Valid {
		username, err = users.NewUsername(row.Username.String)
		if err != nil {
			return nil, fmt.Errorf("invalid username: %w", err)
		}
	}

	var avatarURL *string
	if row.AvatarURL.Valid {
		avatarURL = &row.AvatarURL.String
	}

	var emailVerifiedAt *time.Time
	if row.EmailVerifiedAt.Valid {
		t, err := time.Parse(time.RFC3339, row.EmailVerifiedAt.String)
		if err != nil {
			return nil, fmt.Errorf("invalid email_verified_at timestamp: %w", err)
		}
		emailVerifiedAt = &t
	}

	var lastLoginAt *time.Time
	if row.LastLoginAt.Valid {
		t, err := time.Parse(time.RFC3339, row.LastLoginAt.String)
		if err != nil {
			return nil, fmt.Errorf("invalid last_login_at timestamp: %w", err)
		}
		lastLoginAt = &t
	}

	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at timestamp: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at timestamp: %w", err)
	}

	return users.ReconstructUser(
		userID,
		email,
		row.PasswordHash,
		row.Name,
		username,
		avatarURL,
		users.UserStatus(row.Status),
		emailVerifiedAt,
		lastLoginAt,
		row.Timezone,
		createdAt,
		updatedAt,
	), nil
}
