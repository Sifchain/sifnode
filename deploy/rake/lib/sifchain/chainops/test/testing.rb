module Sifchain
  module Chainops
    module Test
      class Testing < ::Sifchain::Chainops::Builder
        ARGS = { arg1: "--first-arg", arg2: "--second-arg", arg3: "--third-arg" }.freeze

        def initialize(opts = {})
          super
        end

        def generate
          build! ARGS
        end
      end
    end
  end
end
