package gomol

import (
	"path"
	"runtime"
)

type fileRecord struct {
	filename  string
	gomolFile bool
}

var gomolFiles = make(map[string]fileRecord)

var fakeCaller = false
var fakeCallerFile = ""
var fakeCallerLine = 0

func setFakeCallerInfo(file string, line int) {
	if len(file) > 0 {
		fakeCaller = true
		fakeCallerFile = file
		fakeCallerLine = line
	} else {
		fakeCaller = false
		fakeCallerFile = ""
		fakeCallerLine = 0
	}
}

func getCallerInfo() (string, int) {
	if fakeCaller {
		return fakeCallerFile, fakeCallerLine
	}

	file := ""
	line := 0
	/*
	   Start at 3 in the call stack:
	   0 - this function
	   1 - Base.log()
	   2 - Base.<Full Log Function Name>
	   3 - Base.<Short Log Function Name OR external caller>
	*/
	for idx := 3; ; idx++ {
		_, callFile, callLine, _ := runtime.Caller(idx)
		isGomol, filename := isGomolCaller(callFile)

		if isGomol {
			continue
		}

		file = filename
		line = callLine
		break
	}
	return file, line
}
func isGomolCaller(file string) (bool, string) {
	if val, ok := gomolFiles[file]; ok {
		return val.gomolFile, val.filename
	}

	dir := path.Dir(file)
	filename := path.Base(file)

	if len(dir) < 5 {
		gomolFiles[file] = fileRecord{
			filename:  filename,
			gomolFile: false,
		}
		return false, filename
	}

	if dir[len(dir)-5:] == "gomol" {
		if len(filename) < 8 {
			gomolFiles[file] = fileRecord{
				filename:  filename,
				gomolFile: true,
			}
			return true, filename
		}
		if filename[len(filename)-8:] == "_test.go" {
			gomolFiles[file] = fileRecord{
				filename:  filename,
				gomolFile: false,
			}
			return false, filename
		}
	}

	gomolFiles[file] = fileRecord{
		filename:  filename,
		gomolFile: false,
	}
	return false, filename
}
