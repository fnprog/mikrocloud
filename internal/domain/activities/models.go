package activities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	id             ActivityID
	eventType      EventType
	level          ActivityLevel
	description    string
	initiatorID    *uuid.UUID
	resourceType   *string
	resourceID     *uuid.UUID
	resourceName   *string
	metadata       string
	organizationID uuid.UUID
	createdAt      time.Time
}

type ActivityID struct {
	value uuid.UUID
}

func NewActivityID() ActivityID {
	return ActivityID{value: uuid.Must(uuid.NewV7())}
}

func ActivityIDFromUUID(id uuid.UUID) ActivityID {
	return ActivityID{value: id}
}

func ActivityIDFromString(id string) (ActivityID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ActivityID{}, fmt.Errorf("invalid activity ID: %w", err)
	}
	return ActivityID{value: parsedID}, nil
}

func (id ActivityID) String() string {
	return id.value.String()
}

func (id ActivityID) UUID() uuid.UUID {
	return id.value
}

type EventType string

const (
	EventTypeProjectCreated     EventType = "project.created"
	EventTypeProjectUpdated     EventType = "project.updated"
	EventTypeProjectDeleted     EventType = "project.deleted"
	EventTypeAppCreated         EventType = "app.created"
	EventTypeAppUpdated         EventType = "app.updated"
	EventTypeAppDeleted         EventType = "app.deleted"
	EventTypeAppDeployed        EventType = "app.deployed"
	EventTypeAppStarted         EventType = "app.started"
	EventTypeAppStopped         EventType = "app.stopped"
	EventTypeAppRestarted       EventType = "app.restarted"
	EventTypeDatabaseCreated    EventType = "database.created"
	EventTypeDatabaseUpdated    EventType = "database.updated"
	EventTypeDatabaseDeleted    EventType = "database.deleted"
	EventTypeDatabaseStarted    EventType = "database.started"
	EventTypeDatabaseStopped    EventType = "database.stopped"
	EventTypeDatabaseRestarted  EventType = "database.restarted"
	EventTypeEnvironmentCreated EventType = "environment.created"
	EventTypeEnvironmentUpdated EventType = "environment.updated"
	EventTypeEnvironmentDeleted EventType = "environment.deleted"
	EventTypeDiskCreated        EventType = "disk.created"
	EventTypeDiskResized        EventType = "disk.resized"
	EventTypeDiskDeleted        EventType = "disk.deleted"
	EventTypeDiskAttached       EventType = "disk.attached"
	EventTypeDiskDetached       EventType = "disk.detached"
	EventTypeProxyCreated       EventType = "proxy.created"
	EventTypeProxyUpdated       EventType = "proxy.updated"
	EventTypeProxyDeleted       EventType = "proxy.deleted"
	EventTypeUserRegistered     EventType = "user.registered"
	EventTypeUserLogin          EventType = "user.login"
	EventTypeSettingsUpdated    EventType = "settings.updated"
	EventTypeBackupCreated      EventType = "backup.created"
	EventTypeBackupRestored     EventType = "backup.restored"
	EventTypeSystemStarted      EventType = "system.started"
	EventTypeSystemStopped      EventType = "system.stopped"
)

type ActivityLevel string

const (
	ActivityLevelInfo    ActivityLevel = "info"
	ActivityLevelError   ActivityLevel = "error"
	ActivityLevelWarn    ActivityLevel = "warn"
	ActivityLevelSuccess ActivityLevel = "success"
)

func NewActivity(
	eventType EventType,
	level ActivityLevel,
	description string,
	initiatorID *uuid.UUID,
	resourceType *string,
	resourceID *uuid.UUID,
	resourceName *string,
	metadata string,
	organizationID uuid.UUID,
) *Activity {
	return &Activity{
		id:             NewActivityID(),
		eventType:      eventType,
		level:          level,
		description:    description,
		initiatorID:    initiatorID,
		resourceType:   resourceType,
		resourceID:     resourceID,
		resourceName:   resourceName,
		metadata:       metadata,
		organizationID: organizationID,
		createdAt:      time.Now(),
	}
}

func (a *Activity) ID() ActivityID {
	return a.id
}

func (a *Activity) EventType() EventType {
	return a.eventType
}

func (a *Activity) Level() ActivityLevel {
	return a.level
}

func (a *Activity) Description() string {
	return a.description
}

func (a *Activity) InitiatorID() *uuid.UUID {
	return a.initiatorID
}

func (a *Activity) ResourceType() *string {
	return a.resourceType
}

func (a *Activity) ResourceID() *uuid.UUID {
	return a.resourceID
}

func (a *Activity) ResourceName() *string {
	return a.resourceName
}

func (a *Activity) Metadata() string {
	return a.metadata
}

func (a *Activity) OrganizationID() uuid.UUID {
	return a.organizationID
}

func (a *Activity) CreatedAt() time.Time {
	return a.createdAt
}

func ReconstructActivity(
	id ActivityID,
	eventType EventType,
	level ActivityLevel,
	description string,
	initiatorID *uuid.UUID,
	resourceType *string,
	resourceID *uuid.UUID,
	resourceName *string,
	metadata string,
	organizationID uuid.UUID,
	createdAt time.Time,
) *Activity {
	return &Activity{
		id:             id,
		eventType:      eventType,
		level:          level,
		description:    description,
		initiatorID:    initiatorID,
		resourceType:   resourceType,
		resourceID:     resourceID,
		resourceName:   resourceName,
		metadata:       metadata,
		organizationID: organizationID,
		createdAt:      createdAt,
	}
}
