package logruster

import (
	"fmt"
	"os"
	"time"

	logrus "github.com/sirupsen/logrus"
)

// THIS VARIABLE WILL COUNT HOW MANY LOGGING FILE HAS BEEN CREATED
var numberOfLogs = 1

//OUR LOGRUSTER STRUCT
type Logruster struct {
	Log *logrus.Logger
}

//NEW FUNCTION WILL TAKE A TIME AND A LOGFILEPATH PARAMETER
//IN ESSENCE ;
//TIMETORELOG = WHAT IS THE DESIRED INTERVAL FOR CREATING A NEW LOG FILE AND CONTINUE LOGGING IN IT
//LOGFILEPATH = WHAT IS THE PATH IN THE ROOT, WHERE LOGS SHOULD BE LOCATED
func New(timeToReLog int, logFilePath string) Logruster {
	//WE WILL INIT THE FILE FIRST
	file := initTheFile(logFilePath)
	//THEN WE WILL GET A LOGRUS INSTANCE
	log := initNewLogrusInstance(file)
	//WE WILL TRIGGER OUR RELOGGER IN A BRAND NEW GOROUTINE
	//IT WILL BE TRIGGERED IN THE GIVEN INTERVAL IN AN INFINITE LOOP
	go relogger(timeToReLog, logFilePath, log, file)
	return Logruster{Log: log}
}

func relogger(timeToReLog int, logFilePath string, logger *logrus.Logger, oldFile *os.File) {
	//WE WANT TO CLOSE OLDFILES THAT WE HAVE NOTHING TO DO ANYTHING ANYMORE
	old := oldFile
	for {
		//THIS IS THE DESIRED INTERVAL DELAY
		time.Sleep(time.Duration(timeToReLog))
		//THIS VARIABLE REPRESENTS THE NUMBER OF THE LOG FILE
		numberOfLogs += 1
		//LETS CREATE A NEW FILE WITH THE INCREMENTED NUMBEROFLOGS VALUE
		file := initTheFile(logFilePath)
		//LETS CLOSE THE OLD FILE
		old.Close()
		//LETS SET OLDFILE, TO NEW FILE, SO THAT WE ARE ABLE TO CLOSE IN THE NEXT ITERATION
		old = file
		//LASTLY, LETS SET NEWLY CREATED FILE TO LOGGER
		logger.SetOutput(file)

	}
}

//INITIALISES A NEW FILE DEPENDING ON THE GIVEN PATH AND NUMBEROFLOGS & CURRENT TIME
func initTheFile(logFilePath string) *os.File {
	fn := fmt.Sprintf("%s/log_%d_%d.log", logFilePath, numberOfLogs, time.Now().Unix())
	file, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		panic("Error occured oppening file to log , check your log file !")
	}
	return file
}

//SIMPLY CREATES LOGRUS INSTANCE
//ANY FUNCTIONALITY RELATED TO LOGRUS INSTANCE CAN BE GIVEN IN THIS LAYER.
func initNewLogrusInstance(file *os.File) *logrus.Logger {
	log := logrus.New()
	log.SetOutput(file)
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}
