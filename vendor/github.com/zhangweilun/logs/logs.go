package logs

import (
	"time"
	"sync"
	"os"
	"log"
	"fmt"
	"runtime"
	"strconv"
	"math/rand"
)

/**
*
* @author willian
* @created 2017-01-22 19:02
* @email 18702515157@163.com
**/

const (
	ALL int = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)
const (
	date_format = "2006-01-02"
	_        = iota
	KB int64 = 1 << (iota * 10)
	MB
	GB
	TB
)
var log_level int = 1
var max_file_size int64 = MB*20
var consoleAppender = true
var log_obj *file
type file struct {
	dir      string
	filename string
	suffix int
	date    *time.Time
	mu       *sync.RWMutex
	logfile  *os.File
	lg       *log.Logger
}

func Init(dir string)  {
	mkdirlog(dir)
	now := time.Now().Format(date_format)
	time_now := time.Now()
	log_obj= &file{dir:dir, filename:now+".log", date:&time_now,mu:new(sync.RWMutex)}
	log_obj.mu.Lock()
	defer log_obj.mu.Unlock()
	log_obj.logfile, _ = os.OpenFile(dir+"/"+log_obj.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	log_obj.lg = log.New(log_obj.logfile, "[--]", log.Ldate|log.Ltime|log.Lshortfile)
	check_size := time.NewTicker(4 * time.Hour)
	go func() {
		for true {
			select {
			case <-check_size.C:
				if fileSize(log_obj.logfile.Name()) >= max_file_size {
					log_obj.changeFile()
				}
			}
		}
	}()
	day_change := time.NewTicker(24 * time.Hour)
	go func() {
		for true {
			select {
			case <-day_change.C :
				now := time.Now().Format(date_format)
				time_now :=time.Now()
				log_obj= &file{dir:dir, filename:now+".log", date:&time_now,mu:new(sync.RWMutex)}
				log_obj.mu.Lock()
				defer log_obj.mu.Unlock()
				log_obj.logfile, _ = os.OpenFile(dir+"/"+log_obj.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
				log_obj.lg = log.New(log_obj.logfile, "[--]", log.Ldate|log.Ltime|log.Lshortfile)
			}
		}
	}()
}

func Info(v ...interface{}) {
	if log_obj != nil {
		log_obj.mu.RLock()
		defer log_obj.mu.RUnlock()
		if log_level <= INFO {
			log_obj.lg.Output(2, fmt.Sprintln("info", v))
			console("info", v)
		}
	}


}
func Debug(format string,v ...interface{}) {
	if log_obj != nil {
		log_obj.mu.RLock()
		defer log_obj.mu.RUnlock()
		if log_level <= DEBUG {
			log_obj.lg.Printf(format,v)
			//log_obj.lg.Output(2, fmt.Sprintln("debug", v))
			console("debug", v)
		}
	}

}

func Warn(format string,v ...interface{}) {
	if log_obj != nil {
		log_obj.mu.RLock()
		defer log_obj.mu.RUnlock()
		if log_level <= WARN {

			log_obj.lg.Printf(format,v)
			//log_obj.lg.Output(2, fmt.Sprintln("warn", v))

			console("warn", v)
		}
	}


}
func Error(format string,v ...interface{}) {
	if log_obj != nil {
		log_obj.mu.RLock()
		defer log_obj.mu.RUnlock()
		if log_level <= ERROR {

			log_obj.lg.Printf(format,v)
			//log_obj.lg.Output(2, fmt.Sprintln("error", v))

			console("error", v)
		}
	}

}
func Fatal(format string,v ...interface{}) {
	if log_obj != nil {
		log_obj.mu.RLock()
		defer log_obj.mu.RUnlock()
		if log_level <= FATAL {

			log_obj.lg.Printf(format,v)
			//log_obj.lg.Output(2, fmt.Sprintln("fatal", v))

			console("fatal", v)
		}
	}

}


//console output


func console(s ...interface{}) {
	if consoleAppender {
		_, file, line, _ := runtime.Caller(2)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		log.Println(file, strconv.Itoa(line), s)
	}
}

//create a dir that every child of the dir will have same permission mode.
func mkdirlog(dir string) (e error) {
	_, er := os.Stat(dir)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			if os.IsPermission(err) {
				log.Fatalln("create dir error:", err.Error())
			}
		}
	}
	return nil
}


func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}

func (f *file) changeFile(){

	now := time.Now().Format(date_format)
	log_obj.mu.Lock()
	defer log_obj.mu.Unlock()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand_number := r.Intn(10)
	f.filename = now + strconv.Itoa(rand_number)
	log_obj.logfile, _ = os.OpenFile(f.dir+"/"+f.filename+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	log_obj.lg = log.New(f.logfile, "[--]", log.Ldate|log.Ltime|log.Lshortfile)
}