package logruster

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
//MODFORZIPPING = WHAT WOULD YOU LIKE MOD TO BE ? , IN WHAT INTERVAL YOUR LOGS SHOULD BE ZIPPED ? - E.G : IF YOU GIVE 10, YOUR LOGS WILL BE ZIPPED WHEN THEY REACH 10 AND 10K
//LOGFILEPATH = WHAT IS THE PATH IN THE ROOT, WHERE LOGS SHOULD BE LOCATED
func New(timeToReLog int, modForZipping int, logFilePath string) Logruster {

	createLogDirectories(logFilePath)
	//WE WILL INIT THE FILE FIRST
	file := initTheFile(logFilePath)
	//THEN WE WILL GET A LOGRUS INSTANCE
	log := initNewLogrusInstance(file)
	//WE WILL TRIGGER OUR RELOGGER IN A BRAND NEW GOROUTINE
	//IT WILL BE TRIGGERED IN THE GIVEN INTERVAL IN AN INFINITE LOOP
	go relogger(timeToReLog, modForZipping, logFilePath, log, file)
	return Logruster{Log: log}
}

func relogger(timeToReLog int, modForZipping int, logFilePath string, logger *logrus.Logger, oldFile *os.File) {
	defer checkRecover()
	//WE WANT TO CLOSE OLDFILES THAT WE HAVE NOTHING TO DO ANYTHING ANYMORE
	old := oldFile
	for {
		checkForZip(modForZipping, logFilePath)
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
	file, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	checkError(err)
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

func checkForZip(modForZipping int, path string) {
	if numberOfLogs%modForZipping == 0 {

		zipFile, err := os.Create(fmt.Sprintf("%s_archive/%d_logs.zip", path, numberOfLogs))
		checkError(err)
		zipWriter := zip.NewWriter(zipFile)
		defer zipFile.Close()
		defer zipWriter.Close()
		files, _ := ioutil.ReadDir(fmt.Sprintf("./%s", path))

		for i, f := range files {
			fmt.Println(i)
			filePath := fmt.Sprintf("./%s/%s", path, f.Name())
			logFile, _ := os.Open(filePath)
			w, err := zipWriter.Create(f.Name())
			checkError(err)
			io.Copy(w, logFile)
			logFile.Close()
			err5 := os.Remove(filePath)
			checkError(err5)
			fmt.Println(filePath)
		}

	}
}

func createLogDirectories(path string) {
	os.Mkdir(fmt.Sprintf("%s/", path), 0755)
	os.Mkdir(fmt.Sprintf("%s_archive/", path), 0755)
}

//IF THERE IS AN ERROR
//WE WILL LOG IT USING A BRAND NEW LOGRUS CLIENT
//SINCE THE OLD CLEINT MIGHT STILL BE ERRONEUS
func checkError(err error) {
	if err != nil {
		l := logrus.New()
		l.Error(err, " an unexpected error occured while logging")
	}
}

func checkRecover() {
	if err := recover(); err != nil {
		checkError(errors.New("AN UNEXPECTED ERROR OCCURED, CHECK YOUR CONFIGURATIONS"))
	}
}
