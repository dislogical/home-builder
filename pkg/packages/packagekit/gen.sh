#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

go tool dbus-codegen-go -package packagekit -output $SCRIPT_DIR/packagekit.go \
    <(curl -fsSL https://raw.githubusercontent.com/PackageKit/PackageKit/refs/heads/main/src/org.freedesktop.PackageKit.xml) \
    <(curl -fsSL https://raw.githubusercontent.com/PackageKit/PackageKit/refs/heads/main/src/org.freedesktop.PackageKit.Transaction.xml)
