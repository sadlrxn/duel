package log

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/TwiN/go-color"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	logrus_cloudwatchlogs "github.com/kdar/logrus-cloudwatchlogs"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init() {
	Refresh()
	c := cron.New(cron.WithLocation(time.UTC))
	c.AddFunc("0 0 * * *", Refresh)
	c.Start()
}

func Refresh() {
	config := config.Get()

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(config.AWSRegion),
		Credentials: credentials.NewStaticCredentials(config.AWSAccessID, config.AWSSecretKey, ""),
	}))

	streamName, err := getLogStream(sess, config.CloudWatchLogGroup)
	if err != nil {
		// LogMessage("log refresher", "failed to get log stream", "error", logrus.Fields{"error": err.Error()})
		return
	}

	lgr := logrus.New()

	hook, _ := logrus_cloudwatchlogs.NewHook(config.CloudWatchLogGroup, streamName, sess)

	lgr.Hooks.Add(hook)
	lgr.Out = io.Discard
	lgr.Formatter = &logrus.TextFormatter{FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05"}
	Logger = lgr
}

func getLogStream(sess *session.Session, logGroupName string) (string, error) {
	cwl := cloudwatchlogs.New(sess)
	name := time.Now().String()[:10]

	var descending = true
	resp, err := cwl.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{LogGroupName: &logGroupName, Descending: &descending})
	if err != nil {
		return name, err
	}

	for _, logStream := range resp.LogStreams {
		if *logStream.LogStreamName == name {
			return name, nil
		}
	}

	_, err = cwl.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  &logGroupName,
		LogStreamName: &name,
	})

	return name, err
}

func LogMessage(caller string, message string, level string, fields logrus.Fields) {
	jsonString, _ := json.Marshal(fields)
	switch level {
	case "error":
		fmt.Println(color.Colorize(color.Red, "-- "+caller+" -> "+message+" -- "+string(jsonString)))
		// Logger.WithFields(fields).Errorln(caller + " -> " + message + " -- " + string(jsonString))
	case "info":
		fmt.Println(color.Colorize(color.Cyan, "-- "+caller+" -> "+message+" -- "+string(jsonString)))
		// Logger.WithFields(fields).Infoln(caller + " -> " + message + " -- " + string(jsonString))
	case "success":
		fmt.Println(color.Colorize(color.Green, "-- "+caller+" -> "+message+" -- "+string(jsonString)))
		// Logger.WithFields(fields).Println(caller + " -> " + message + " -- " + string(jsonString))
	default:
		fmt.Println(color.Colorize(color.Purple, "-- "+caller+" -> "+message+" : "+level))
		break
	}
}
