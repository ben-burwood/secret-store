set shell := ["powershell.exe", "-c"]

default:
    just --list

format:
    gofmt -w .
    Set-Location frontend; vp check --fix
