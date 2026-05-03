set shell := ["powershell.exe", "-c"]

default:
    just --list

format:
    Set-Location backend; ruff format
    Set-Location frontend; vp check --fix
