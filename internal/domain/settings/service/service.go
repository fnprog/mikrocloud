package service

import (
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mikrocloud/mikrocloud/internal/domain/settings"
	"github.com/mikrocloud/mikrocloud/internal/domain/settings/repository"
)

type SettingsService struct {
	repo *repository.SettingsRepository
}

func NewSettingsService(repo *repository.SettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) GetGeneralSettings() (*settings.GeneralSettings, error) {
	return s.repo.GetGeneralSettings()
}

func (s *SettingsService) SaveGeneralSettings(generalSettings *settings.GeneralSettings) error {
	return s.repo.SaveGeneralSettings(generalSettings)
}

func (s *SettingsService) GetAdvancedSettings() (*settings.AdvancedSettings, error) {
	return s.repo.GetAdvancedSettings()
}

func (s *SettingsService) SaveAdvancedSettings(advancedSettings *settings.AdvancedSettings) error {
	return s.repo.SaveAdvancedSettings(advancedSettings)
}

func (s *SettingsService) GetUpdateSettings() (*settings.UpdateSettings, error) {
	return s.repo.GetUpdateSettings()
}

func (s *SettingsService) SaveUpdateSettings(updateSettings *settings.UpdateSettings) error {
	return s.repo.SaveUpdateSettings(updateSettings)
}

func (s *SettingsService) GetInstanceInfo() (*settings.InstanceInfo, error) {
	generalSettings, err := s.repo.GetGeneralSettings()
	if err != nil {
		return nil, err
	}

	return &settings.InstanceInfo{
		FQDN: generalSettings.Domain,
		IPv4: generalSettings.IPv4,
		IPv6: generalSettings.IPv6,
	}, nil
}

func (s *SettingsService) DetectIPAddresses() (*settings.DetectedIPs, error) {
	detected := &settings.DetectedIPs{}

	ipv4, err := detectPublicIPv4()
	if err == nil {
		detected.IPv4 = ipv4
	}

	ipv6, err := detectPublicIPv6()
	if err == nil {
		detected.IPv6 = ipv6
	}

	return detected, nil
}

func detectPublicIPv4() (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

func detectPublicIPv6() (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api64.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(body))
	if net.ParseIP(ip).To4() != nil {
		return "", nil
	}

	return ip, nil
}
