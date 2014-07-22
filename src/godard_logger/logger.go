package godard_logger
//http://technosophos.com/2013/09/14/using-gos-built-logger-log-syslog.html
import (
  "log"
  "gsyslog"
  "reflect"
)

//LOG_METHODS = [:emerg, :alert, :crit, :err, :warning, :notice, :info, :debug]


type GodardLogger struct {
  options map[string]{}interface
  Logger string
  Prefix string
  Prefixes map[string]string
}

func NewGodardLogger(options map[string]{}interface) *GodardLogger{

    c := &GodardLogger{}
    c.options = options
    c.Logger  = options["logger"] || c.createLogger()
    c.Prefix  = options["prefix"]
    //@prefix   = options[:prefix]
    //@stdout   = options[:stdout]
    c.Prefixes = make(map[string]string , 0)
}

func (c*GodardLogger) PrefixWith(prefix string){
  //@prefixes[prefix] ||= self.class.new(:logger => self, :prefix => prefix)
}


func (c*GodardLogger) Reopen(){
  if reflect.TypeOf(c.Logger) == reflect.TypeOf(c) {
    c.Logger.Reopen()
  }else{
    c.Logger = c.CreateLogger()
  }
}

func (c*GodardLogger) CreateLogger(){
  if _,ok := c.options["log_file"]; ok {
    //LoggerAdapter.new(@options[:log_file])
  }else{

    // Configure logger to write to the syslog. You could do this in init(), too.
    logwriter, e := syslog.New(syslog.LOG_LOCAL6, "myprog")
    if e == nil {
        log.SetOutput(logwriter)
    }
    //Syslog.close if Syslog.opened? # need to explictly close it before reopening it
    //Syslog.open(@options[:identity] || 'godardd', Syslog::LOG_PID, Syslog::LOG_LOCAL6)
  }
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
