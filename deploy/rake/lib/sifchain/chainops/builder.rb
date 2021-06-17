module Sifchain
  module Chainops
    class Builder
      def initialize(opts = {})
        @args = opts.fetch(:args)
      end

      attr_accessor :args

      def build!(arg_map)
        arg_map.map { |k, v| "#{v} #{args[k]}" }.join(" ")
      end
    end
  end
end
