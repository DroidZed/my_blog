#!/bin/zsh

echo "Installing tools in env... "

go install github.com/air-verse/air@latest
go install github.com/a-h/templ/cmd/templ@latest

echo "Done !"
