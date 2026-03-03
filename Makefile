# Simple Makefile for OVMS setup on Windows

.PHONY: help install export run clean

help:
	@echo Available commands:
	@echo   make install  - Create venv and install dependencies
	@echo   make export   - Run model export script
	@echo   make run      - Start OVMS server
	@echo   make clean    - Remove venv

install:
	uv python install 3.12 --install-dir python
	for /d %%d in (python\cpython-*) do uv venv export --python %%d\python.exe --relocatable
	uv pip install --python export\Scripts\python.exe -r export-model\requirements.txt

export:
	export\Scripts\python.exe export-model\export_model.py $(ARGS)

run:
	ovms\setupvars.ps1 && ovms\ovms.exe --rest_port 8000 --config_path config.json

clean:
	rmdir /s /q export
	rmdir /s /q python
