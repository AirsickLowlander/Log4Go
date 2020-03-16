# Log4Go

## Usage example

    package main

    import (
        "fmt"
        "log4go"
        "time"
    )

    func main() {

        fmt.Println("Launching the server...")
        log4go.Initialize("Message Bus Example")

        for {
            //message, _ := bufio.NewReader(conn).ReadString('\n')
            log4go.LogTrace("A Trace Message")
            time.Sleep(10 * time.Millisecond)

            log4go.LogDebug("A Debug Message")
            time.Sleep(10 * time.Millisecond)

            log4go.LogInfo("A Info Message")
            time.Sleep(10 * time.Millisecond)

            log4go.LogError("A Error Message")
            time.Sleep(10 * time.Millisecond)
            log4go.LogFatal("A Fatal Mesage")

            time.Sleep(500 * time.Millisecond)
        }
    }

