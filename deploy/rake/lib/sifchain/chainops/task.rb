require 'active_support/all'

module Sifchain
  module Chainops
    class Task
      KLASS_PATH = "Sifchain::Chainops".freeze

      def initialize(opts = {})
        @task = opts.fetch(:task)
        @args = opts.fetch(:args)
      end

      attr_accessor :task, :args

      def build
        klass.new(args: args).generate
      end

      private

      def klass
        "#{KLASS_PATH}::#{tklass}".constantize
      end

      def tklass
        "#{task}".split(":").map { |t| t.capitalize }.join("::")
      end
    end
  end
end
