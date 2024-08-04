FROM mcr.microsoft.com/devcontainers/go:1-1.22-bullseye

RUN apt-get update && apt-get upgrade -y
RUN curl -fsSL https://bun.sh/install | bash
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN chsh -s $(which zsh) $USER
