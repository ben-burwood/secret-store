set shell := ["powershell.exe", "-c"]

default:
    just --list

format:
    Set-Location backend; ruff format
    Set-Location frontend; vp check

backend *args:
    Set-Location backend; uv run python manage.py {{ args }}

frontend *args:
    Set-Location frontend; npm run {{ args }}
