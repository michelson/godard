package trigger

import(
  process "process"
)

type Trigger struct {
  Process string
  Logger string
  //mutex 
  Name string
  ScheduledEvents []string

}

func NewTrigger(process process.Process, options map[string]string) *Trigger {
  c := &Trigger{}
  c.Name = options["name"]
  c.Process = process
  c.Logger = options["logger"]
  c.ScheduledEvents = make([]string, 0)
  return c
}

func (c*Trigger) Reset(){
  //self.cancel_all_events

}

func (c*Trigger) Notify(transition string){
  //raise "Implement in subclass"
}

func (c*Trigger) Dispatch(){
  //self.process.dispatch!(event, self.class.name.split("::").last)

}

func (c*Trigger) ScheduleEvent(){
  /*

      # TODO: maybe wrap this in a ScheduledEvent class with methods like cancel
      thread = Thread.new(self) do |trigger|
        begin
          sleep delay.to_f
          trigger.dispatch!(event)
          trigger.mutex.synchronize do
            trigger.scheduled_events.delete_if { |_, thread| thread == Thread.current }
          end
        rescue StandardError => e
          trigger.logger.err(e)
          trigger.logger.err(e.backtrace.join("\n"))
        end
      end

      self.scheduled_events.push([event, thread])

  */
}

func (c*Trigger) CancellAllEvents(){
  /*
     self.logger.info "Canceling all scheduled events"
      self.mutex.synchronize do
        self.scheduled_events.each {|_, thread| thread.kill}
      end
  */
}