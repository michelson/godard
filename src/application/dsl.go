/*module Bluepill
  def self.application(app_name, options = {}, &block)
    app_proxy = AppProxy.new(app_name, options)
    if block.arity == 0
      app_proxy.instance_eval &block
    else
      app_proxy.instance_exec(app_proxy, &block)
    end
    app_proxy.app.load
  end
end
*/
package application

import (
 cfg "godard_config"
)

func InitApplication(app_name string , options *cfg.GodardConfig ){
  app_proxy := NewAppProxy(app_name, options)
  app_proxy.App.Load()
}

