package types

type UserLoginResponse struct {
	Token string `json:"token"`
}

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Rating   int64  `json:"rating"`
}

type UserData struct {
	Id     string `json:"id"`
	Email  string `json:"email"`
	Rating int64  `json:"rating"`
}

type Reply struct {
	Id             string `json:"id"`
	Contents       string `json:"contents"`
	SubmitterId    string `json:"submitter_id"`
	SubmitterEmail string `json:"submitter_email"`
}

type Answer struct {
	Id             string  `json:"id"`
	Contents       string  `json:"contents"`
	SubmitterId    string  `json:"submitter_id"`
	SubmitterEmail string  `json:"submitter_email"`
	Replies        []Reply `json:"replies"`
}

type Question struct {
	Id             string   `json:"id"`
	SubmitterId    string   `json:"submitter_id"`
	SubmitterEmail string   `json:"submitter_email"`
	Title          string   `json:"title"`
	Contents       string   `json:"contents"`
	Answers        []Answer `json:"answers"`
	BestAnswer     string   `json:"best_answer_id"`
	Upvoters       []string `json:"upvoters"`
	Downvoters     []string `json:"downvoters"`
}

// for questions only
type EditSuggestion struct {
	Id             string `json:"id"`
	QuestionId     string `json:"question_id"`
	ApproverId     string `json:"approver_id"`
	SubmitterId    string `json:"submitter_id"`
	SubmitterEmail string `json:"submitter_email"`
	Contents       string `json:"contents"`
	EditStatus     string `json:"edit_status"` // default,approve,reject
}

type FavoriteCommentRequest struct {
	QuestionId string `json:"question_id"`
	AnswerId   string `json:"answer_id"`
}

type UpdateQuestionVotesRequest struct {
	Type       string `json:"type"` // upvote, downvote
	QuestionId string `json:"question_id"`
}

type EditContentRequest struct {
	QuestionId string `json:"question_id"`
	Title      string `json:"title"`
	Content    string `json:"contents"`
}

type AddAnswerRequest struct {
	QuestionId string `json:"question_id"`
	Contents   string `json:"contents"`
}

type AddReplyRequest struct {
	QuestionId string `json:"question_id"`
	AnswerId   string `json:"answer_id"`
	Contents   string `json:"contents"`
}
