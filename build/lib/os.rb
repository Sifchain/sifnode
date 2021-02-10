require 'rbconfig'

def detect_os
  @os ||= (
  host_os = RbConfig::CONFIG['host_os']
  case host_os
  when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
    :windows
  when /darwin|mac os/
    :macosx
  when /linux/
    :linux
  when /solaris|bsd/
    :unix
  else
    raise Error::WebDriverError, "unknown os: #{host_os.inspect}"
  end
  )
end
