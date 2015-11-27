#!/bin/bash

set -e

(cf uninstall-plugin "FastPushPlugin" || true) && go build -o fastpush-plugin main.go && cf install-plugin fastpush-plugin
