#
{ pkgs ? import <nixpkgs> { } }:

let

in pkgs.mkShell {
  # nativeBuildInputs is usually what you want -- tools you need to run
  nativeBuildInputs = [
    pkgs.buildPackages.ncurses
    pkgs.buildPackages.erlang
    pkgs.buildPackages.elixir
    pkgs.buildPackages.elixir_ls
    pkgs.go
    pkgs.gnumake
    pkgs.sqlite-interactive
    pkgs.gopls
    pkgs.inotify-tools
    pkgs.go-task
    pkgs.gotools
  ];
  shellHook = ''
    echo "Starting nix-shell for the project (yeah)..."
  '';

}
