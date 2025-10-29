package repository

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/activities"
)

type ActivitiesRepository struct {
	db *sql.DB
}

func NewActivitiesRepository(db *sql.DB) *ActivitiesRepository {
	return &ActivitiesRepository{db: db}
}

func (r *ActivitiesRepository) Create(activity *activities.Activity) error {
	_, err := r.db.Exec(`
		INSERT INTO activities (id, event_type, level, description, initiator_id, 
			resource_type, resource_id, resource_name, metadata, organization_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		activity.ID().String(),
		string(activity.EventType()),
		string(activity.Level()),
		activity.Description(),
		uuidToNullString(activity.InitiatorID()),
		strPtrToNullString(activity.ResourceType()),
		uuidToNullString(activity.ResourceID()),
		strPtrToNullString(activity.ResourceName()),
		activity.Metadata(),
		activity.OrganizationID().String(),
		activity.CreatedAt(),
	)
	return err
}

func (r *ActivitiesRepository) ListByOrganization(organizationID uuid.UUID, limit int, offset int) ([]*activities.Activity, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.event_type, a.level, a.description, a.initiator_id,
			a.resource_type, a.resource_id, a.resource_name, a.metadata, a.organization_id, a.created_at
		FROM activities a
		WHERE a.organization_id = ?
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?
	`, organizationID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActivities(rows)
}

func (r *ActivitiesRepository) ListByResource(resourceType string, resourceID uuid.UUID, limit int) ([]*activities.Activity, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.event_type, a.level, a.description, a.initiator_id,
			a.resource_type, a.resource_id, a.resource_name, a.metadata, a.organization_id, a.created_at
		FROM activities a
		WHERE a.resource_type = ? AND a.resource_id = ?
		ORDER BY a.created_at DESC
		LIMIT ?
	`, resourceType, resourceID.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActivities(rows)
}

func (r *ActivitiesRepository) DeleteOlderThan(days int) error {
	_, err := r.db.Exec(`
		DELETE FROM activities
		WHERE created_at < datetime('now', '-' || ? || ' days')
	`, days)
	return err
}

func (r *ActivitiesRepository) scanActivities(rows *sql.Rows) ([]*activities.Activity, error) {
	var result []*activities.Activity

	for rows.Next() {
		var (
			id             string
			eventType      string
			level          string
			description    string
			initiatorID    sql.NullString
			resourceType   sql.NullString
			resourceID     sql.NullString
			resourceName   sql.NullString
			metadata       string
			organizationID string
			createdAt      time.Time
		)

		if err := rows.Scan(
			&id, &eventType, &level, &description, &initiatorID,
			&resourceType, &resourceID, &resourceName, &metadata, &organizationID, &createdAt,
		); err != nil {
			return nil, err
		}

		activityID, _ := activities.ActivityIDFromString(id)
		orgID, _ := uuid.Parse(organizationID)

		activity := activities.ReconstructActivity(
			activityID,
			activities.EventType(eventType),
			activities.ActivityLevel(level),
			description,
			nullStringToUUIDPtr(initiatorID),
			nullStringToStrPtr(resourceType),
			nullStringToUUIDPtr(resourceID),
			nullStringToStrPtr(resourceName),
			metadata,
			orgID,
			createdAt,
		)

		result = append(result, activity)
	}

	return result, rows.Err()
}

func uuidToNullString(u *uuid.UUID) sql.NullString {
	if u == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: u.String(), Valid: true}
}

func strPtrToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func nullStringToUUIDPtr(ns sql.NullString) *uuid.UUID {
	if !ns.Valid {
		return nil
	}
	u, err := uuid.Parse(ns.String)
	if err != nil {
		return nil
	}
	return &u
}

func nullStringToStrPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

type ActivityDTO struct {
	ID          string                 `json:"id"`
	EventType   string                 `json:"event_type"`
	Level       string                 `json:"level"`
	Description string                 `json:"description"`
	Initiator   *InitiatorDTO          `json:"initiator"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
}

type InitiatorDTO struct {
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
	Initials string `json:"initials"`
}

func (r *ActivitiesRepository) ListByOrganizationWithUsers(organizationID uuid.UUID, limit int, offset int) ([]*ActivityDTO, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.event_type, a.level, a.description, a.metadata, a.created_at,
			u.name, u.avatar_url
		FROM activities a
		LEFT JOIN users u ON a.initiator_id = u.id
		WHERE a.organization_id = ?
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?
	`, organizationID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActivityDTOs(rows)
}

func (r *ActivitiesRepository) ListByResourceWithUsers(resourceType string, resourceID uuid.UUID, limit int) ([]*ActivityDTO, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.event_type, a.level, a.description, a.metadata, a.created_at,
			u.name, u.avatar_url
		FROM activities a
		LEFT JOIN users u ON a.initiator_id = u.id
		WHERE a.resource_type = ? AND a.resource_id = ?
		ORDER BY a.created_at DESC
		LIMIT ?
	`, resourceType, resourceID.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActivityDTOs(rows)
}

func (r *ActivitiesRepository) scanActivityDTOs(rows *sql.Rows) ([]*ActivityDTO, error) {
	var result []*ActivityDTO

	for rows.Next() {
		var (
			id          string
			eventType   string
			level       string
			description string
			metadata    string
			createdAt   time.Time
			userName    sql.NullString
			userAvatar  sql.NullString
		)

		if err := rows.Scan(&id, &eventType, &level, &description, &metadata, &createdAt, &userName, &userAvatar); err != nil {
			return nil, err
		}

		var metadataMap map[string]interface{}
		if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
			metadataMap = make(map[string]interface{})
		}

		var initiator *InitiatorDTO
		if userName.Valid {
			initiator = &InitiatorDTO{
				Name:   userName.String,
				Avatar: userAvatar.String,
			}
			// Generate initials from name
			if userName.String != "" {
				parts := strings.Split(strings.TrimSpace(userName.String), " ")
				if len(parts) >= 2 {
					initiator.Initials = strings.ToUpper(string(parts[0][0]) + string(parts[len(parts)-1][0]))
				} else if len(parts[0]) > 0 {
					initiator.Initials = strings.ToUpper(string(parts[0][0]))
				}
			}
		}

		dto := &ActivityDTO{
			ID:          id,
			EventType:   eventType,
			Level:       level,
			Description: description,
			Initiator:   initiator,
			Metadata:    metadataMap,
			CreatedAt:   createdAt,
		}

		result = append(result, dto)
	}

	return result, rows.Err()
}
