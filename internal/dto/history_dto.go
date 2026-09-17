package dto

// Mirrors Dto/HistoryDto.kt. NOTE: unused by any handler in the original
// backend (see internal/models/session.go); kept for 1:1 source parity.

type SessionSummary struct {
	ID             string  `json:"id"`
	Role           string  `json:"role"`
	Difficulty     string  `json:"difficulty"`
	AverageScore   float64 `json:"averageScore"`
	TotalQuestions int     `json:"totalQuestions"`
	Completed      bool    `json:"completed"`
	CreatedAt      string  `json:"createdAt"`
}

type AnswerDetail struct {
	ID            string  `json:"id"`
	QuestionText  string  `json:"questionText"`
	UserAnswer    string  `json:"userAnswer"`
	Rating        float64 `json:"rating"`
	Feedback      string  `json:"feedback"`
	IdealAnswer   string  `json:"idealAnswer"`
	AccuracyScore int     `json:"accuracyScore"`
	DepthScore    int     `json:"depthScore"`
	ClarityScore  int     `json:"clarityScore"`
}
