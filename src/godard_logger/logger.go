package godard_logger
//http://technosophos.com/2013/09/14/using-gos-built-logger-log-syslog.html
//http://stackoverflow.com/questions/380172/reading-syslog-output-on-a-mac
import (
  "log"
  "log/syslog"
  //"github.com/hashicorp/go-syslog"
  "reflect"
  "os"
)

//LOG_METHODS = [:emerg, :alert, :crit, :err, :warning, :notice, :info, :debug]

type GodardLogger struct {
  options map[string]interface{}
  Logger   *log.Logger
  Prefix   string
  Stdout   bool
  Prefixes map[string]*GodardLogger
}

var logger *log.Logger

func NewGodardLogger(options map[string]interface{}) *GodardLogger{

    c := &GodardLogger{}
    c.options  = options
    if _,ok := c.options["logger"]; ok {
      c.Logger = options["logger"].(*log.Logger)
    }else{
      c.Logger = c.CreateLogger()
    }
    if _,ok := c.options["prefix"]; ok {
      c.Prefix   = options["prefix"].(string)
    }
    if _,ok := c.options["stdout"]; ok {
      c.Stdout   = options["stdout"].(bool)
    }
    c.Prefixes = make(map[string]*GodardLogger , 0)
    return c
}

func (c*GodardLogger) PrefixWith(prefix string) *GodardLogger{
  //@prefixes[prefix] ||= self.class.new(:logger => self, :prefix => prefix)
  if _, ok := c.Prefixes[prefix]; ok {
    return c.Prefixes[prefix]
  }else{
    opts := make(map[string]interface{}, 0)
    opts["logger"] = c.Logger
    opts["prefix"] = c.Prefix
    return NewGodardLogger(opts)
  }
}


func (c*GodardLogger) Reopen(){
  if reflect.TypeOf(c.Logger) == reflect.TypeOf(c) {
    //c.Logger.Reopen()
  }else{
    c.Logger = c.CreateLogger()
  }
}

func (c*GodardLogger) CreateLogger() *log.Logger {

  if len(c.options["log_file"].(string)) > 0 {
    log.Println("LOGGING TO:" , c.options["log_file"].(string))
    LoggerAdapter(c.options["log_file"].(string))
    //LoggerAdapter.new(@options[:log_file])
  }else{

    log.Println("LOGGING TO: syslog myprog")
    // Configure logger to write to the syslog. You could do this in init(), too.
    logwriter, e := syslog.New(syslog.LOG_LOCAL6, "myprog")
    if e == nil {
      log.SetOutput(logwriter)
    }
    /*err := syslog.Close()
    if err != nil {
      log.Println("error ", err)
    }*/
    //Syslog.close if Syslog.opened? # need to explictly close it before reopening it
    //Syslog.open(@options[:identity] || 'godardd', Syslog::LOG_PID, Syslog::LOG_LOCAL6)
  }
  return logger
}

func LoggerAdapter(log_file string) {

  f, _ := os.Open(log_file)
  logger = log.New(os.Stderr, "xxx: ", log.Ldate | log.Ltime | log.Lshortfile)
  log.SetOutput(f)
}

/*
    class LoggerAdapter < ::Logger
      LOGGER_EQUIVALENTS =
        {:debug => :debug, :err => :error, :warning => :warn, :info => :info, :emerg => :fatal, :alert => :warn, :crit => :fatal, :notice => :info}

      LOG_METHODS.each do |method|
        next if method == LOGGER_EQUIVALENTS[method]
        alias_method method, LOGGER_EQUIVALENTS[method]
      end
    end
*/
