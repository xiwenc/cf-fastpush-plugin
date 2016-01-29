#!/bin/bash

set -e

(cf uninstall-plugin "FastPushPlugin" || true) && go build -o cf-fastpush-plugin main.go && cf install-plugin cf-fastpush-plugin
