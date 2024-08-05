package models

type Student struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

type Test struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Question struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	TestID   string `json:"test_id"`
}

type Enrollment struct {
	StudentID string `json:"student_id"`
	TestID    string `json:"test_id"`
}

type TestScore struct {
	TestID    string `json:"test_id"`
	StudentID string `json:"student_id"`
	Ok        int32  `json:"ok"`
	Ko        int32  `json:"ko"`
	Total     int32  `json:"total"`
	Score     int32  `json:"score"`
}

type Answer struct {
	StudentId  string `json:"student_id"`
	TestId     string `json:"test_id"`
	QuestionId string `json:"question_id"`
	Answer     string `json:"answer"`
	Correct    bool   `json:"correct"`
}
