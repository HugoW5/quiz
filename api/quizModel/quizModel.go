package quizModel

type Question struct {
	Id      string
	Text    string
	Choiecs [3]string
}

type Answers struct {
	Id     string
	Choice int
}

type ScoreResponse struct {
	Score      int `json:"score"`
	Total      int `json:"total"`
	Percentile int `json:"percentile"`
}

// Map the question id to the correct answer index
var CorrectAnswers = map[string]int{
	"1": 2,
	"2": 1,
	"3": 0,
	"4": 2,
	"5": 0,
	"6": 2,
}

var Questions = []Question{
	{
		Id:      "1",
		Text:    "Hur långt är det mellan Jönköping och Varberg med bil? (snabbaste vägen)",
		Choiecs: [3]string{"110km", "92km", "169km"},
	},
	{
		Id:      "2",
		Text:    "Hur många elever går i SUT24?",
		Choiecs: [3]string{"19st", "24st", "32st"},
	},
	{
		Id:      "3",
		Text:    "Vilken är Sveriges största insjö?",
		Choiecs: [3]string{"Vänern", "Vättern", "Mälaren"},
	},
	{
		Id:      "4",
		Text:    "Hur lång är Vättern?",
		Choiecs: [3]string{"102km", "231km", "135km"},
	},
	{
		Id:      "5",
		Text:    "Hur lång är tunnelen under Varberg?",
		Choiecs: [3]string{"3km", "6km", "4.8km"},
	},
	{
		Id:      "6",
		Text:    "Vilket år byggdes Societén - Varberg?",
		Choiecs: [3]string{"1902", "1807", "1886"},
	},
}
