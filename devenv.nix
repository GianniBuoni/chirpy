{pkgs, ...}: {
  # packages
  packages = with pkgs; [
    git

    # language support/tools
    gopls
    goose
    sqlc

    # extras
    jq
  ];

  # languages
  languages.go.enable = true;

  # services
  # services.postgres.enable = true;

  # https://devenv.sh/scripts/
  enterShell = ''
    git --version
    export PATH="$HOME/go/bin:$PATH"
    tmux
  '';

  # https://devenv.sh/tests/
  enterTest = ''
    go test ./...
  '';

  # https://devenv.sh/git-hooks/
  # git-hooks.hooks.shellcheck.enable = true;

  # See full reference at https://devenv.sh/reference/options/
}
