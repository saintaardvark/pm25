SHELL=/bin/bash

setup: .venv packages

.venv:
	python3 -mvirtualenv -p python3 .venv

.PHONY: packages
packages:
	source .venv/bin/activate && \
		pip install --upgrade pip && \
		pip install -r requirements.txt
