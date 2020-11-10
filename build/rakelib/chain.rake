desc ""
namespace :chain do
  desc "Chain state operations"
  namespace :state do
    desc "Export chain state"
    task :export, [:file, :node_directory] do |t, args|
      fh = File.open(args[:file], "w")
      if fh.nil?
        puts "unable to open the file #{args[:file]}!"
        exit(1)
      end

      state = `sifnoded export --home #{args[:node_directory]}`
      fh.puts state
      fh.close
    end
  end

  desc "Migrate a chain"
  task :migrate, [:version, :genesis_file, :node_directory] do |t, args|
    system("sifnoded migrate #{args[:version]} #{args[:genesis_file]} --home #{args[:node_directory]}")
  end
end
