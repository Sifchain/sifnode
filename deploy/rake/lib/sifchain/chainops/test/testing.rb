module Sifchain
  module Chainops
    module Test
      class Testing < ::Sifchain::Chainops::Builder
        ARGS = { arg1: "--first-arg", arg2: "--second-arg", arg3: "--third-arg" }.freeze

        def initialize(opts = {})
          super
        end

        def generate
          "#{exec} #{build! ARGS}"
        end

        private

        def exec
          ::Sifchain::Chainops::Cli::KUBECTL
        end
      end
    end
  end
end
