package condition 

 /* MB := 1024 ** 2
  FORMAT_STR := "%d%s"
  MB_LABEL := "MB"
  KB_LABEL := "KB"

*/

type FileTime struct {
  Below int
  filename string
}


/*

      def initialize(options = {})
        @below = options[:below]
        @filename = options[:filename]
      end

      def run(pid, include_children)
        if File.exists?(@filename)
          Time.now()-File::mtime(@filename)
        else
          nil
        end
      rescue
        $!
      end

      def check(value)
        return false if value.nil?
        return value < @below
      end

*/
