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
	Properties []property `xml:"log4j:properties>log4j:data"`
}

var ch chan *logevt

// Initialize inits the logging function
func Initialize() {
	fmt.Println("Init log4go")
	c := make(chan *logevt, 1000)
	ch = c

	go messagePump(ch)
}

// LogTrace Logs Trace Messages
func LogTrace(s string) {
	var p = makeLogEvt()
	p.Level = "Trace"
	p.Message = s
	ch <- p
}

// LogDebug Logs Debug Messages
func LogDebug(s string) {
	var p = makeLogEvt()
	p.Level = "DEBUG"
	p.Message = s
	ch <- p
}

// LogInfo Logs Info Messages
func LogInfo(s string) {
	var p = makeLogEvt()
	p.Level = "INFO"
	p.Message = s
	ch <- p
}

// LogError Logs Error Messages
func LogError(s string) {
	var p = makeLogEvt()
	p.Level = "ERROR"
	p.Message = s
	ch <- p
}

// LogFatal Logs Fatal Messages
func LogFatal(s string) {
	var p = makeLogEvt()
	p.Level = "Fatal"
	p.Message = s
	ch <- p
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func makeLogEvt() *logevt {
	//NOTE(MWILKINSON): This method is hampered by the lack of reflection in this implementation. The goal here is to make this generic enough for any appliation to have relevant meta data
	return &logevt{
		"log4j:event",
		"GetToWork",
		"Trace",
		makeTimestamp(),
		1,
		"",
		[]property{{"log4japp", "TBD"}, {"log4jmachinename", "Local Machine"}}, //TODO(MWILKINSON):  Determine the best way to get the app name and the Machine name
	}
}

func messagePump(ch chan *logevt) {

	for {
		conn, err := net.Dial("udp4", "127.0.0.1:9998")
		if err != nil {
			fmt.Println("error with the connection", err)
		}
		for {
			p := <-ch
			//conn.Write([]byte(fmt.Sprintf("A sample Message %d", i)))
			b, err := xml.Marshal(p)
			if err != nil {
				fmt.Println("An Error was encountered")
			}
			conn.Write(b)

		}
	}
}
