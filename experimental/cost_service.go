package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain"
	"github.com/mikrocloud/mikrocloud/pkg/cloudprovider"
)

// CostService handles cost management and monitoring
type CostService struct {
	repos        domain.RepositoryManager
	cloudFactory *cloudprovider.CloudProviderFactory
}

// NewCostService creates a new cost service
func NewCostService(
	repos domain.RepositoryManager,
	cloudFactory *cloudprovider.CloudProviderFactory,
) *CostService {
	return &CostService{
		repos:        repos,
		cloudFactory: cloudFactory,
	}
}

// CostBreakdown represents cost breakdown by service
type CostBreakdown struct {
	Service string  `json:"service"`
	Cost    float64 `json:"cost"`
	Usage   string  `json:"usage"`
	Trend   string  `json:"trend"` // "increasing", "decreasing", "stable"
}

// ProjectCosts represents comprehensive project cost information
type ProjectCosts struct {
	ProjectID       uuid.UUID            `json:"project_id"`
	ProjectName     string               `json:"project_name"`
	CloudProvider   domain.CloudProvider `json:"cloud_provider"`
	TotalCost       float64              `json:"total_cost"`
	Currency        string               `json:"currency"`
	Period          string               `json:"period"`
	Breakdown       []CostBreakdown      `json:"breakdown"`
	Budget          *BudgetInfo          `json:"budget,omitempty"`
	Projection      *CostProjection      `json:"projection,omitempty"`
	Recommendations []CostRecommendation `json:"recommendations,omitempty"`
	Alerts          []CostAlert          `json:"alerts,omitempty"`
}

// BudgetInfo represents budget information
type BudgetInfo struct {
	MonthlyLimit      float64 `json:"monthly_limit"`
	CurrentSpend      float64 `json:"current_spend"`
	AlertThreshold    float64 `json:"alert_threshold"`
	DaysRemaining     int     `json:"days_remaining"`
	BudgetUtilization float64 `json:"budget_utilization"` // Percentage
}

// CostProjection represents cost projections
type CostProjection struct {
	ProjectedMonthly float64 `json:"projected_monthly"`
	ProjectedDaily   float64 `json:"projected_daily"`
	Confidence       float64 `json:"confidence"` // 0-1 scale
	BasedOnDays      int     `json:"based_on_days"`
}

// CostRecommendation represents cost optimization recommendations
type CostRecommendation struct {
	Type             string  `json:"type"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	PotentialSavings float64 `json:"potential_savings"`
	Priority         string  `json:"priority"` // "high", "medium", "low"
	Actionable       bool    `json:"actionable"`
}

// CostAlert represents cost alerts
type CostAlert struct {
	ID           uuid.UUID `json:"id"`
	Type         string    `json:"type"`
	Severity     string    `json:"severity"`
	Message      string    `json:"message"`
	Amount       float64   `json:"amount"`
	Threshold    float64   `json:"threshold"`
	CreatedAt    time.Time `json:"created_at"`
	Acknowledged bool      `json:"acknowledged"`
}

// SetCostBudgetRequest represents a request to set cost budget
type SetCostBudgetRequest struct {
	ProjectID      uuid.UUID `json:"project_id" validate:"required"`
	MonthlyLimit   float64   `json:"monthly_limit" validate:"required,gt=0"`
	AlertThreshold float64   `json:"alert_threshold" validate:"required,gte=50,lte=100"`
	Currency       string    `json:"currency" validate:"required"`
}

// GetProjectCosts retrieves comprehensive cost information for a project
func (s *CostService) GetProjectCosts(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, timeRange cloudprovider.TimeRange) (*ProjectCosts, error) {
	// Verify project ownership
	project, err := s.repos.Project().GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Get cloud provider
	provider, err := s.getCloudProvider(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud provider: %w", err)
	}

	// Get costs from cloud provider
	costData, err := provider.GetCosts(ctx, projectID.String(), timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get costs from cloud provider: %w", err)
	}

	// Get budget information
	budget, err := s.getBudgetInfo(ctx, projectID, project.CloudProvider.Provider)
	if err != nil {
		// Budget might not exist, that's okay
		budget = nil
	}

	// Calculate projection
	projection := s.calculateCostProjection(costData, timeRange)

	// Generate recommendations
	recommendations := s.generateCostRecommendations(ctx, project, costData)

	// Get active alerts
	alerts, err := s.getActiveAlerts(ctx, projectID)
	if err != nil {
		// Continue without alerts if there's an error
		alerts = []CostAlert{}
	}

	// Convert breakdown
	breakdown := make([]CostBreakdown, len(costData.Breakdown))
	for i, item := range costData.Breakdown {
		breakdown[i] = CostBreakdown{
			Service: item.Service,
			Cost:    item.Cost,
			Usage:   item.Usage,
			Trend:   "stable", // TODO: Calculate actual trend
		}
	}

	return &ProjectCosts{
		ProjectID:       projectID,
		ProjectName:     project.Name,
		CloudProvider:   project.CloudProvider.Provider,
		TotalCost:       costData.TotalCost,
		Currency:        costData.Currency,
		Period:          fmt.Sprintf("%s to %s", timeRange.Start.Format("2006-01-02"), timeRange.End.Format("2006-01-02")),
		Breakdown:       breakdown,
		Budget:          budget,
		Projection:      projection,
		Recommendations: recommendations,
		Alerts:          alerts,
	}, nil
}

// SetCostBudget sets or updates a cost budget for a project
func (s *CostService) SetCostBudget(ctx context.Context, userID uuid.UUID, req *SetCostBudgetRequest) (*BudgetInfo, error) {
	// Verify project ownership
	project, err := s.repos.Project().GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Check if budget already exists
	existingBudget, err := s.repos.CostBudget().GetByProjectAndProvider(ctx, req.ProjectID, project.CloudProvider.Provider)
	if err == nil {
		// Update existing budget
		existingBudget.MonthlyLimit = req.MonthlyLimit
		existingBudget.AlertThreshold = req.AlertThreshold
		if err := s.repos.CostBudget().Update(ctx, existingBudget); err != nil {
			return nil, fmt.Errorf("failed to update budget: %w", err)
		}
	} else {
		// Create new budget
		budget := &domain.CostBudget{
			ProjectID:      req.ProjectID,
			CloudProvider:  project.CloudProvider.Provider,
			MonthlyLimit:   req.MonthlyLimit,
			AlertThreshold: req.AlertThreshold,
			CurrentSpend:   0,
			IsActive:       true,
		}

		if err := s.repos.CostBudget().Create(ctx, budget); err != nil {
			return nil, fmt.Errorf("failed to create budget: %w", err)
		}
	}

	// Set up cloud provider budget
	provider, err := s.getCloudProvider(ctx, project)
	if err == nil {
		cloudBudget := cloudprovider.CostBudget{
			Amount:    req.MonthlyLimit,
			Currency:  req.Currency,
			Threshold: req.AlertThreshold,
		}
		provider.SetCostBudget(ctx, req.ProjectID.String(), cloudBudget)
	}

	// Return budget info
	return s.getBudgetInfo(ctx, req.ProjectID, project.CloudProvider.Provider)
}

// GetCostAlerts retrieves active cost alerts for a project
func (s *CostService) GetCostAlerts(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]CostAlert, error) {
	// Verify project ownership
	project, err := s.repos.Project().GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return s.getActiveAlerts(ctx, projectID)
}

// CompareCostsAcrossClouds compares costs for the same workload across different cloud providers
func (s *CostService) CompareCostsAcrossClouds(ctx context.Context, userID uuid.UUID, workloadConfig map[string]interface{}) (*MultiCloudCostComparison, error) {
	comparison := &MultiCloudCostComparison{
		Workload:    workloadConfig,
		Comparisons: make([]CloudCostEstimate, 0),
		GeneratedAt: time.Now(),
	}

	// Get estimates from each cloud provider
	for _, provider := range []domain.CloudProvider{domain.CloudProviderAWS, domain.CloudProviderAzure, domain.CloudProviderGCP} {
		estimate, err := s.estimateCostForProvider(ctx, provider, workloadConfig)
		if err != nil {
			// Log error but continue with other providers
			continue
		}
		comparison.Comparisons = append(comparison.Comparisons, *estimate)
	}

	// Find the most cost-effective option
	if len(comparison.Comparisons) > 0 {
		minCost := comparison.Comparisons[0].EstimatedMonthlyCost
		minIndex := 0
		for i, comp := range comparison.Comparisons {
			if comp.EstimatedMonthlyCost < minCost {
				minCost = comp.EstimatedMonthlyCost
				minIndex = i
			}
		}
		comparison.Recommendation = &comparison.Comparisons[minIndex]
	}

	return comparison, nil
}

// Helper methods

func (s *CostService) getCloudProvider(ctx context.Context, project *domain.Project) (cloudprovider.CloudProvider, error) {
	cloudConfig, err := s.repos.CloudProvider().GetByID(ctx, project.CloudProviderID)
	if err != nil {
		return nil, err
	}

	credentials := map[string]string{
		"region": cloudConfig.Region,
		// TODO: Decrypt actual credentials
	}

	return s.cloudFactory.Create(cloudConfig.Provider, credentials)
}

func (s *CostService) getBudgetInfo(ctx context.Context, projectID uuid.UUID, provider domain.CloudProvider) (*BudgetInfo, error) {
	budget, err := s.repos.CostBudget().GetByProjectAndProvider(ctx, projectID, provider)
	if err != nil {
		return nil, err
	}

	// Calculate days remaining in current month
	now := time.Now()
	endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
	daysRemaining := int(endOfMonth.Sub(now).Hours() / 24)

	// Calculate budget utilization
	utilization := 0.0
	if budget.MonthlyLimit > 0 {
		utilization = (budget.CurrentSpend / budget.MonthlyLimit) * 100
	}

	return &BudgetInfo{
		MonthlyLimit:      budget.MonthlyLimit,
		CurrentSpend:      budget.CurrentSpend,
		AlertThreshold:    budget.AlertThreshold,
		DaysRemaining:     daysRemaining,
		BudgetUtilization: utilization,
	}, nil
}

func (s *CostService) calculateCostProjection(costData *cloudprovider.CostData, timeRange cloudprovider.TimeRange) *CostProjection {
	days := int(timeRange.End.Sub(timeRange.Start).Hours() / 24)
	if days == 0 {
		days = 1
	}

	dailyAverage := costData.TotalCost / float64(days)
	projectedMonthly := dailyAverage * 30

	// Simple confidence calculation based on data period
	confidence := 0.5
	if days >= 7 {
		confidence = 0.7
	}
	if days >= 30 {
		confidence = 0.9
	}

	return &CostProjection{
		ProjectedMonthly: projectedMonthly,
		ProjectedDaily:   dailyAverage,
		Confidence:       confidence,
		BasedOnDays:      days,
	}
}

func (s *CostService) generateCostRecommendations(ctx context.Context, project *domain.Project, costData *cloudprovider.CostData) []CostRecommendation {
	recommendations := []CostRecommendation{}

	// Analyze cost breakdown for optimization opportunities
	for _, item := range costData.Breakdown {
		switch item.Service {
		case "EC2", "ECS", "Container Apps", "Cloud Run":
			if item.Cost > 50 { // Arbitrary threshold
				recommendations = append(recommendations, CostRecommendation{
					Type:             "compute_optimization",
					Title:            "Optimize Compute Resources",
					Description:      fmt.Sprintf("Consider right-sizing or using spot instances for %s", item.Service),
					PotentialSavings: item.Cost * 0.3, // Estimated 30% savings
					Priority:         "medium",
					Actionable:       true,
				})
			}
		case "RDS", "Azure SQL", "Cloud SQL":
			if item.Cost > 30 {
				recommendations = append(recommendations, CostRecommendation{
					Type:             "database_optimization",
					Title:            "Database Cost Optimization",
					Description:      "Consider using reserved instances or adjusting database size",
					PotentialSavings: item.Cost * 0.25,
					Priority:         "high",
					Actionable:       true,
				})
			}
		}
	}

	// Add general recommendations
	if costData.TotalCost > 100 {
		recommendations = append(recommendations, CostRecommendation{
			Type:             "monitoring",
			Title:            "Set up Cost Monitoring",
			Description:      "Enable detailed cost monitoring and alerts to prevent cost overruns",
			PotentialSavings: 0,
			Priority:         "high",
			Actionable:       true,
		})
	}

	return recommendations
}

func (s *CostService) getActiveAlerts(ctx context.Context, projectID uuid.UUID) ([]CostAlert, error) {
	// This would typically query a separate alerts table or get alerts from the cloud provider
	// For now, return empty slice
	return []CostAlert{}, nil
}

func (s *CostService) estimateCostForProvider(ctx context.Context, provider domain.CloudProvider, workloadConfig map[string]interface{}) (*CloudCostEstimate, error) {
	// This would implement cost estimation logic for each provider
	// For now, return a placeholder estimate
	baseEstimate := 100.0 // Base cost

	// Adjust based on provider
	switch provider {
	case domain.CloudProviderAWS:
		baseEstimate *= 1.0 // AWS baseline
	case domain.CloudProviderAzure:
		baseEstimate *= 0.9 // Azure typically 10% cheaper
	case domain.CloudProviderGCP:
		baseEstimate *= 0.85 // GCP typically 15% cheaper
	}

	return &CloudCostEstimate{
		Provider:             provider,
		EstimatedMonthlyCost: baseEstimate,
		Currency:             "USD",
		Confidence:           0.7,
		Breakdown: []CostBreakdown{
			{Service: "Compute", Cost: baseEstimate * 0.6, Usage: "Standard"},
			{Service: "Storage", Cost: baseEstimate * 0.2, Usage: "Standard"},
			{Service: "Network", Cost: baseEstimate * 0.2, Usage: "Standard"},
		},
	}, nil
}

// Additional types for multi-cloud cost comparison
type MultiCloudCostComparison struct {
	Workload       map[string]interface{} `json:"workload"`
	Comparisons    []CloudCostEstimate    `json:"comparisons"`
	Recommendation *CloudCostEstimate     `json:"recommendation,omitempty"`
	GeneratedAt    time.Time              `json:"generated_at"`
}

type CloudCostEstimate struct {
	Provider             domain.CloudProvider `json:"provider"`
	EstimatedMonthlyCost float64              `json:"estimated_monthly_cost"`
	Currency             string               `json:"currency"`
	Confidence           float64              `json:"confidence"`
	Breakdown            []CostBreakdown      `json:"breakdown"`
}
