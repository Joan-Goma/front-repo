package errorController

import(
	"io"
	"os"
	"log"
	"bytes"
  "github.com/fatih/color"
)
 
 
var (

    WarningLogger *log.Logger
    InfoLogger    *log.Logger
		DebugLogger    *log.Logger
    ErrorLogger   *log.Logger
    ErrorColor = color.New(color.Bold, color.FgRed).SprintFunc()
    InfoColor = color.New(color.Bold, color.FgWhite).SprintFunc()
    DebugColor = color.New(color.Bold, color.FgGreen).SprintFunc()
    WarningColor = color.New(color.Bold, color.FgYellow).SprintFunc()

)
func InitLog(debugEnabled bool){
  
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
  
	InfoLogger = log.New(wrt, InfoColor("\nINFO: "), log.Ldate|log.Ltime|log.Lshortfile)
  InfoLogger.SetOutput(wrt)
	
	WarningLogger = log.New(wrt, WarningColor("\nWARNING: "), log.Ldate|log.Ltime|log.Lshortfile)
  WarningLogger.SetOutput(wrt)
  
	DebugLogger = log.New(wrt, DebugColor("\nDEBUG: "), log.Ldate|log.Ltime|log.Lshortfile)
  DebugLogger.SetOutput(wrt)
	if !debugEnabled {
		var buff bytes.Buffer
		DebugLogger.SetOutput(&buff)
	}
	
	ErrorLogger = log.New(wrt, ErrorColor("\nERROR: "), log.Ldate|log.Ltime|log.Lshortfile)
  ErrorLogger.SetOutput(wrt)
}