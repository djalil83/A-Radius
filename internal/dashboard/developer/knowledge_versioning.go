package developer

import (
	"errors"
	"time"
)

type KnowledgeLifecycle string

const (
	KnowledgeDiscovered     KnowledgeLifecycle = "DISCOVERED"
	KnowledgeAnalyzing      KnowledgeLifecycle = "ANALYZING"
	KnowledgeValidated      KnowledgeLifecycle = "VALIDATED"
	KnowledgeReviewRequired KnowledgeLifecycle = "REVIEW_REQUIRED"
	KnowledgeApproved       KnowledgeLifecycle = "APPROVED"
	KnowledgeStaged         KnowledgeLifecycle = "STAGED"
	KnowledgeActive         KnowledgeLifecycle = "ACTIVE"
	KnowledgeSuperseded     KnowledgeLifecycle = "SUPERSEDED"
	KnowledgeArchived       KnowledgeLifecycle = "ARCHIVED"
)

type KnowledgeEnvironment string

const (
	KnowledgeEnvironmentIsolated   KnowledgeEnvironment = "ai-knowledge-db"
	KnowledgeEnvironmentCache      KnowledgeEnvironment = "analysis-cache"
	KnowledgeEnvironmentStaging    KnowledgeEnvironment = "staging"
	KnowledgeEnvironmentProduction KnowledgeEnvironment = "production"
)

var ErrInvalidKnowledgeTransition = errors.New("invalid knowledge lifecycle transition")
var ErrAIActiveTransitionDenied = errors.New("AI cannot transition knowledge to ACTIVE")

// KnowledgeVersionRecord is isolated from application release state.
type KnowledgeVersionRecord struct {
	KnowledgeKey       string               `json:"knowledge_key"`
	Version            string               `json:"version"`
	Source             string               `json:"source"`
	DiscoveredAt       time.Time            `json:"discovered_at"`
	Confidence         int                  `json:"confidence"`
	RelevanceToARadius string               `json:"relevance_to_a_radius"`
	AffectedModules    []string             `json:"affected_modules"`
	AnalysisResult     string               `json:"analysis_result"`
	Recommendation     string               `json:"recommendation"`
	Status             KnowledgeLifecycle   `json:"status"`
	Environment        KnowledgeEnvironment `json:"environment"`
	Findings           int                  `json:"findings"`
	ProductionChanged  bool                 `json:"production_changed"`
}

func (s KnowledgeLifecycle) Valid() bool {
	switch s {
	case KnowledgeDiscovered, KnowledgeAnalyzing, KnowledgeValidated, KnowledgeReviewRequired, KnowledgeApproved, KnowledgeStaged, KnowledgeActive, KnowledgeSuperseded, KnowledgeArchived:
		return true
	default:
		return false
	}
}

// CanTransition defines the only valid lifecycle edges. ACTIVE requires human-controlled gates.
func CanTransition(from, to KnowledgeLifecycle, actor string) error {
	allowed := map[KnowledgeLifecycle][]KnowledgeLifecycle{
		KnowledgeDiscovered:     {KnowledgeAnalyzing},
		KnowledgeAnalyzing:      {KnowledgeValidated, KnowledgeReviewRequired},
		KnowledgeValidated:      {KnowledgeReviewRequired},
		KnowledgeReviewRequired: {KnowledgeApproved, KnowledgeArchived},
		KnowledgeApproved:       {KnowledgeStaged},
		KnowledgeStaged:         {KnowledgeActive},
		KnowledgeActive:         {KnowledgeSuperseded},
		KnowledgeSuperseded:     {KnowledgeArchived},
	}
	if !from.Valid() || !to.Valid() {
		return ErrInvalidKnowledgeTransition
	}
	for _, candidate := range allowed[from] {
		if candidate == to {
			if to == KnowledgeActive && actor == "ai" {
				return ErrAIActiveTransitionDenied
			}
			return nil
		}
	}
	return ErrInvalidKnowledgeTransition
}

// KnowledgePromotionPolicy prevents learning state from becoming application state.
type KnowledgePromotionPolicy struct {
	KnowledgeDBIsolated       bool     `json:"knowledge_db_isolated"`
	AnalysisCacheReadOnly     bool     `json:"analysis_cache_read_only"`
	AutoProductionPromotion   bool     `json:"auto_production_promotion"`
	RequiredDeveloperApproval bool     `json:"required_developer_approval"`
	RequiredStages            []string `json:"required_stages"`
}

var DefaultKnowledgePromotionPolicy = KnowledgePromotionPolicy{
	KnowledgeDBIsolated:       true,
	AnalysisCacheReadOnly:     true,
	AutoProductionPromotion:   false,
	RequiredDeveloperApproval: true,
	RequiredStages:            []string{"DISCOVERED", "ANALYZING", "VALIDATED", "REVIEW_REQUIRED", "APPROVED", "STAGED", "ACTIVE"},
}

var FeaturedKnowledgeVersion = KnowledgeVersionRecord{
	KnowledgeKey:       "a-radius.api-security",
	Version:            "SK-2.4.7",
	Source:             "Security Intelligence Feed",
	DiscoveredAt:       time.Date(2026, time.August, 17, 14, 20, 0, 0, time.FixedZone("WITA", 8*60*60)),
	Confidence:         91,
	RelevanceToARadius: "HIGH",
	AffectedModules:    []string{"API Langganan", "API Administrator", "API Technician"},
	AnalysisResult:     "Pola authorization middleware tidak konsisten dengan RBAC baseline.",
	Recommendation:     "Pastikan setiap endpoint sensitif memiliki authorization policy berdasarkan role dan permission.",
	Status:             KnowledgeActive,
	Environment:        KnowledgeEnvironmentIsolated,
	Findings:           23,
	ProductionChanged:  false,
}

var FeaturedNewIntelligence = KnowledgeVersionRecord{
	KnowledgeKey:       "INT-2026-00821",
	Version:            "SK-2.4.8",
	Source:             "Security Intelligence Feed",
	DiscoveredAt:       time.Date(2026, time.August, 17, 14, 2, 0, 0, time.FixedZone("WITA", 8*60*60)),
	Confidence:         91,
	RelevanceToARadius: "HIGH",
	AffectedModules:    []string{"/api/langganan", "/api/pelanggan", "/api/admin", "/api/technician"},
	AnalysisResult:     "Ditemukan pola vulnerability baru yang relevan dengan modul API.",
	Recommendation:     "Pastikan setiap endpoint sensitif memiliki authorization policy berdasarkan role dan permission.",
	Status:             KnowledgeReviewRequired,
	Environment:        KnowledgeEnvironmentIsolated,
	Findings:           0,
	ProductionChanged:  false,
}
