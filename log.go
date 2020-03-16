package log4go

import (
	"encoding/xml"
	"fmt"
	"net"
	"time"
)

type property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type logevt struct {
	XMLName    string     `xml:"log4j:event"`
	Name       string     `xml:"logger,attr"`
	Level      string     `xml:"level,attr"`
	Timestamp  int64      `xml:"timestamp,attr"`
	Thread     int        `xml:"thread,attr"`
	Message    string     `xml:"log4j:message"`
	Properties []property `xml:"log4j:properties>log4j:datat"`
}

var ch chan *logevt
var loggerName string

const log4JTag = "log4j:event"

type logLevel string

const (
	// Trace low level rapid fire logs
	Trace = "TRACE"
	// Debug Programmer level info about specific components
	Debug = "DEBUG"
	// Info High Level info about specific components
	Info = "INFO"
	// Error problems in the code
	Error = "ERROR"
	// Fatal  oops
	Fatal = "FATAL"
)

// Initialize inits the logging function
func Initialize(programName string) {
	fmt.Println("Init log4go")
	c := make(chan *logevt, 1000)
	ch = c
	loggerName = programName
	go messagePump(ch)
}

func makeLogEvent(level logLevel, message string) *logevt {
	return &logevt{
		log4JTag,
		loggerName,
		string(level),
		makeTimestamp(),
		1,
		message,
		[]property{{"log4japp", "TBD"}, {"log4jmachinename", "Local Machine"}}, //TODO(MWILKINSON):  Determine the best way to get the app name and the Machine name
	}
}

// LogTrace Logs Trace Messages
func LogTrace(s string) {
	var p = makeLogEvent(Trace, s)
	ch <- p
}

// LogDebug Logs Debug Messages
func LogDebug(s string) {
	var p = makeLogEvent(Debug, s)
	ch <- p
}

// LogInfo Logs Info Messages
func LogInfo(s string) {
	var p = makeLogEvent(Info, s)
	ch <- p
}

// LogError Logs Error Messages
func LogError(s string) {
	var p = makeLogEvent(Error, s)
	ch <- p
}

// LogFatal Logs Fatal Messages
func LogFatal(s string) {
	var p = makeLogEvent(Fatal, s)
	p.Level = "Fatal"
	p.Message = s
	ch <- p
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func messagePump(ch chan *logevt) {

	for {
		conn, err := net.Dial("udp4", "127.0.0.1:9998")
		if err != nil {
			fmt.Println("error with the connection", err)
		}
		for {
			p := <-ch
			b, err := xml.Marshal(p)
			if err != nil {
				fmt.Println("An Error was encountered")
			}
			conn.Write(b)
		}
	}
}
