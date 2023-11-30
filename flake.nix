{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs =
    { self
    , nixpkgs
    }:
    let
      pkgs = import nixpkgs {
        system = "x86_64-linux";
      };
    in
    {
      devShells.x86_64-linux.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          delve
        ];
        # Required for delve to run correctly
        # https://github.com/go-delve/delve/issues/3085
        hardeningDisable = [ "all" ];
        shellHook = ''
          alias go-build="go build -ldflags \"-w -X main.version=robsversion\" -o ./scionnats"
        '';
      };
    };
}
