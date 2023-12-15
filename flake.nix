{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
    # utils.url = "github:numtide/flake-utils";
    # gomod2nix = {
    #   url = "github:tweag/gomod2nix";
    #   inputs.nixpkgs.follows = "nixpkgs";
    #   inputs.utils.follows = "utils";
    # };
  };

  outputs = { self, nixpkgs /*, utils, gomod2nix */ }:
    let
      # Generate a user-friendly version number.
      version = builtins.substring 0 8 self.lastModifiedDate;

      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

    in
   {
     packages = forAllSystems (system:
       let
         pkgs = import nixpkgs {
           inherit system;
           # overlays = [ gomod2nix.overlays.default ];
         };
       in
         {
           default = pkgs.buildGoModule {
             pname = "econ-epub";
             inherit version;
             src = ./.;
             vendorHash = null;  # this should probably be something else...
           };
         });

     apps = forAllSystems(system: {
       default = {
         type = "app";
         program = "${self.packages.${system}.default}/bin/econ-epub";
       };
     });

     devShells = forAllSystems (system:
       let pkgs = nixpkgsFor.${system};
       in {
         default = pkgs.mkShell {
           buildInputs = with pkgs; [
             go
             gopls
             gotools
             go-tools
             # gomod2nix.packages.${system}.default
           ];
         };
       });
   };
}
