package dto

// Mirrors Dto/AnswerDto.kt. Used by POST /api/interview/evaluate.

type EvaluationRequest struct {
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	Role       string `json:"role"`
	Difficulty string `json:"difficulty"`
}

type EvaluationResponse struct {
	Score              int    `json:"score"`
	Feedback           string `json:"feedback"`
	IdealAnswerSummary string `json:"idealAnswerSummary"`
}
