package wlog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	mName          string
	mDirectory     string
	mFileInterval  int64
	mDuration      int64
	mWriteInterval int64
	mRunning       bool = true
	mNext          int64
	mFileName      string
	mSDF           string
	SDF_MINUTE     string = "20060102150405"
	SDF_HOUR       string = "200601021504"
	SDF_MILLIS     string = "15:04:05"
	mWriter               = bufio.NewWriter(nil)
	mMessageList   []string
	mTopSize       int
	mLevel         int
	mFile          *os.File
)

func Log(sName string, sDirectory string, iFileInterval int) {
	mName = "NOMANE"
	if len(sName) > 0 && sName != "" {
		mName = sName
	}

	sDefaultDirectory := "./log/"
	if len(sDirectory) <= 0 && sDirectory == "" {
		mDirectory = sDefaultDirectory
	} else {
		mDirectory = sDirectory
	}

	switch iFileInterval {
	case 1:
		mDuration = 60000
		mSDF = SDF_MINUTE
		break
	case 2:
		mDuration = 3600000
		mSDF = SDF_HOUR
		break
	default:
		break
	}

	mWriteInterval = 200
	mFileInterval = 200
	mNext = 0
	mFileName = ""
}

func isRunning() bool {
	return mRunning
}

func createDirectory() bool {
	result := false

	if _, err := os.Stat(mDirectory); os.IsNotExist(err) {
		os.Mkdir(mDirectory, os.ModePerm)
	} else {
		result = true
	}

	return result
}

func createFile(iTime int64) bool {
	result := false
	curr := iTime - (iTime % mDuration)

	exists := true
	newName := createFileName(curr)

	directory := mDirectory + "/" + newName
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if mFile != nil {
			mFile.Close()
		}

		file, err := os.Create(directory)
		if err != nil {
			exists = false
		}

		mFile = file
		//defer file.Close()

		if exists {
			mWriter = bufio.NewWriter(file)
			mWriter.Flush()

			mFileName = newName

			result = true
		}
	} else {
		if mFile == nil {
			file, err := os.OpenFile(directory, os.O_APPEND, 0644)

			if err != nil {
				exists = false
			}

			mFile = file

			mWriter = bufio.NewWriter(file)
			mWriter.Flush()

			result = true
		}
	}

	mNext = curr + mDuration

	return result
}

func createFileName(iTime int64) string {
	var fileName strings.Builder
	t := time.Unix(0, iTime*int64(time.Millisecond))

	fileName.WriteString(t.Format(mSDF))
	if mName != "" {
		fileName.WriteString("_")
		fileName.WriteString(mName)
	}

	fileName.WriteString(".")
	fileName.WriteString("log")

	return fileName.String()
}

func GetFileThread() {
	running := true

	for running {
		time.Sleep(time.Duration(mFileInterval))
		running = isRunning()
		T0 := CurrentTimeMillis()
		if T0 > mNext && createDirectory() {
			createFile(T0)
		}

	}
}

func GetWriteThread() {
	running := true

	for running {
		time.Sleep(time.Duration(mWriteInterval))

		running = isRunning()
		var cloned []string = mMessageList
		mMessageList = []string{}

		for _, el := range cloned {
			mWriter.WriteString(el)
		}

		mWriter.Flush()
	}
}

func setRunning(bRunning bool) {
	mRunning = bRunning
}

func SetLevel(iLevel int) {
	switch iLevel {
	case 1:
	case 2:
	case 3:
	case 4:
		mLevel = iLevel
		break
	default:
		mLevel = 1
		break

	}
}

func D(sTag string, sMeta string, sMessage string) {
	if mLevel <= 1 {
		write("DBG", sTag, sMeta, sMessage)
	}
}

func I(sTag string, sMeta string, sMessage string) {
	if mLevel <= 2 {
		write("INF", sTag, sMeta, sMessage)
	}
}

func W(sTag string, sMeta string, sMessage string) {
	if mLevel <= 3 {
		write("WRN", sTag, sMeta, sMessage)
	}
}

func E(sTag string, sMeta string, sMessage string) {
	if mLevel <= 4 {
		write("ERR", sTag, sMeta, sMessage)
	}
}

func write(sLevel string, sTag string, sMeta string, sMessage string) {
	T0 := CurrentTimeMillis()
	t := time.Unix(0, T0*int64(time.Millisecond))

	var sb strings.Builder

	sb.WriteString(sLevel)
	sb.WriteString(" ")
	sb.WriteString(t.Format(SDF_MILLIS))
	sb.WriteString("   | ")
	sb.WriteString(sTag)
	sb.WriteString("|")
	sb.WriteString(sMeta)
	sb.WriteString("|   ")

	sb.WriteString(sMessage)
	sb.WriteString("\n")

	if isRunning() {
		mMessageList = append(mMessageList, sb.String())
	} else {
		fmt.Println(sb.String())
	}
}

func Start() {
	if isRunning() {
		if !createDirectory() {
			log.Fatal("Can't create directory")
		}

		if !createFile(CurrentTimeMillis()) {
			mNext = 0
			log.Fatal("Can't create file")
		}

		setRunning(true)
		go GetFileThread()
		go GetWriteThread()

		fmt.Println("logging start")
	}
}
