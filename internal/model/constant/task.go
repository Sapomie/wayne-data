package constant

//task name
const (
	CodeInput             = "codeInput"             //10
	CodeOutput            = "codeOutput"            //8
	CodeInfoAndArrange    = "codeInfoAndArrange"    //6
	Sleep                 = "sleeping"              //0
	Vacant                = "vacant"                //0
	VideoGame             = "videoGame"             //2
	SocialEntertain       = "socialEntertain"       //5
	OtherAndCommute       = "other"                 //4
	Routine               = "routine"               //3
	WorkOther             = "workOther"             //6
	GraphInput            = "graphInput"            //10
	GraphOutput           = "graphOutput"           //8
	GraphInfoAndArrange   = "graphInfoAndArrange"   //6
	EnglishInput          = "englishInput"          //10
	EnglishOutput         = "englishOutput"         //8
	EnglishInfoAndArrange = "englishInfoAndArrange" //6
	BallGame              = "ballGame"              //6
	Running               = "running"               //10
	Riding                = "riding"                //10
	Anaerobic             = "anaerobic"             //8
	RecordPlan            = "recordPlan"            //8
	Housework             = "housework"             //6
	VirtueArrange         = "virtueArrange"         //6
	Nonfiction            = "nonfiction"            //10
	Literal               = "literal"               //6
	AnimationAndEpisode   = "animationAndEpisode"   //5
	Movie                 = "movie"                 //5
	WatchGame             = "watchGame"             //4
	WatchOther            = "watchOther"            //4
	OtherTechInput        = "otherTechInput"        //8
	OtherTechOutput       = "otherTechOutput"       //6
)

//main task
var MainTasks = []string{
	//TaskName
	CodeInput,
	CodeOutput,
	CodeInfoAndArrange,
	GraphInput,
	GraphOutput,
	EnglishInput,
	EnglishOutput,
	Running,
	Anaerobic,
	Nonfiction,
	Literal,
	Movie,

	WatchGame,
	AnimationAndEpisode,
}

//daily task

var DailyTasks = []string{
	EnglishInput,
	Running,
	Anaerobic,
	Nonfiction,

	//CodeInput,
	//CodeOutput,
	//Movie,
}

func IsDailyTask(taskName string) bool {
	for _, dailyTask := range DailyTasks {
		if dailyTask == taskName {
			return true
		}
	}
	return false
}
