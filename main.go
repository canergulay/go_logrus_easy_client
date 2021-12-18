package main

import (
	"fmt"
	"time"

	"github.com/canergulay/gopractice/logruster"
	logrus "github.com/sirupsen/logrus"
)

func main() {
	//LETS GIVE A TRY TO IT

	//WE WILL INITIALISE OUR LOGRUSTER WITH AN INTERVAL 10 SECONDS
	//WHICH MEANS IT WILL RECREATE A BRAND NEW LOG FILE IN A INTERVAL OF 10 SECONDS
	//IT WILL ALSO SAVE ALL OF THOSE STUFF IN /LOGS FOLDER AS WE GIVE IT SO.
	logruster := logruster.New(int(time.Second*10), "logs")

	//LETS LOG SOME DUMMY DATA AND SEE WHAT HAPPENS
	for i := 0; i < 10000; i++ {
		time.Sleep(time.Second * 2)
		info := fmt.Sprintf("Logrus has logged for %dth time ! ", i)
		logruster.Log.WithFields(logrus.Fields{"location": "main", "reason": "testing the log"}).Info(info)
	}

}
