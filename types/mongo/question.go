package types

type Reply struct {
	Id          string `json:"id"`
	Contents    string `json:"contents"`
	SubmitterId string `json:"submitter_id"`
}

type Answer struct {
	Id          string  `json:"id"`
	Contents    string  `json:"contents"`
	SubmitterId string  `json:"submitter_id"`
	Replies     []Reply `json:"replies"`
}

type Question struct {
	Id         string   `json:"id"`
	SumitterId string   `json:"submitter_id"`
	Title      string   `json:"title"`
	Contents   string   `json:"contents"`
	Answers    []Answer `json:"replies"`
	BestAnswer string   `json:"best_answer_id"`
	Upvoters   []string `json:"upvoters"`
	Downvoters []string `json:"downvoters"`
}

// for questions only
type EditSuggestion struct {
	Id          string `json:"id"`
	QuestionId  string `json:"question_id"`
	SubmitterId string `json:"submitter_id"`
	Contents    string `json:"contents"`
	EditStatus  string `json:"edit_status"` // default,approve,reject
}
