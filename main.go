package main

import (
	"fmt"
	"time"

	"github.com/canergulay/gopractice/logruster"
	logrus "github.com/sirupsen/logrus"
)

func main() {
	logruster := logruster.New(int(time.Second*10), "logs")

	for i := 0; i < 10000; i++ {
		time.Sleep(time.Second * 2)
		info := fmt.Sprintf("Logrus has logged for %dth time ! ", i)
		logruster.Log.WithFields(logrus.Fields{"location": "main", "reason": "testing the log"}).Info(info)
	}

}

//
