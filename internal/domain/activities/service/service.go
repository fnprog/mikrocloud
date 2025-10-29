package service

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/activities"
	"github.com/mikrocloud/mikrocloud/internal/domain/activities/repository"
)

type ActivitiesService struct {
	repo *repository.ActivitiesRepository
}

func NewActivitiesService(repo *repository.ActivitiesRepository) *ActivitiesService {
	return &ActivitiesService{repo: repo}
}

func (s *ActivitiesService) LogActivity(
	eventType activities.EventType,
	description string,
	initiatorID *uuid.UUID,
	resourceType *string,
	resourceID *uuid.UUID,
	resourceName *string,
	metadata map[string]any,
	organizationID uuid.UUID,
) error {
	level := s.determineLevel(eventType)

	metadataJSON := "{}"
	if metadata != nil {
		jsonBytes, err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataJSON = string(jsonBytes)
	}

	activity := activities.NewActivity(
		eventType,
		level,
		description,
		initiatorID,
		resourceType,
		resourceID,
		resourceName,
		metadataJSON,
		organizationID,
	)

	return s.repo.Create(activity)
}

func (s *ActivitiesService) determineLevel(eventType activities.EventType) activities.ActivityLevel {
	switch eventType {
	case activities.EventTypeAppDeleted,
		activities.EventTypeDatabaseDeleted,
		activities.EventTypeEnvironmentDeleted,
		activities.EventTypeDiskDeleted,
		activities.EventTypeProxyDeleted,
		activities.EventTypeUserLogin:
		return activities.ActivityLevelWarn

	case activities.EventTypeSystemStarted,
		activities.EventTypeSystemStopped:
		return activities.ActivityLevelInfo

	default:
		return activities.ActivityLevelSuccess
	}
}

func (s *ActivitiesService) GetRecentActivities(organizationID uuid.UUID, limit int, offset int) ([]*repository.ActivityDTO, error) {
	return s.repo.ListByOrganizationWithUsers(organizationID, limit, offset)
}

func (s *ActivitiesService) GetResourceActivities(resourceType string, resourceID uuid.UUID, limit int) ([]*repository.ActivityDTO, error) {
	return s.repo.ListByResourceWithUsers(resourceType, resourceID, limit)
}

func (s *ActivitiesService) CleanupOldActivities(days int) error {
	return s.repo.DeleteOlderThan(days)
}
