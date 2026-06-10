package main

import (
	"context"
	"fmt"
	"log"

	"github.com/godbus/dbus/v5"
	"github.com/urfave/cli/v3"

	"github.com/dislogical/home-builder/pkg/packages/packagekit"
)

var test *cli.Command = &cli.Command{
	Name: "test",
	Action: func(ctx context.Context, c *cli.Command) error {
		conn, err := dbus.SystemBusPrivate()
		if err != nil {
			return fmt.Errorf("connecting to system d-bus: %w", err)
		}
		defer conn.Close() //nolint:errcheck

		err = conn.Auth(nil)
		if err != nil {
			return fmt.Errorf("auth'ing on the system d-bus: %w", err)
		}

		err = conn.Hello()
		if err != nil {
			return fmt.Errorf("hello'ing on the system d-bus: %w", err)
		}

		pko := packagekit.NewOrg_Freedesktop_PackageKit(
			conn.Object(packagekit.InterfaceOrg_Freedesktop_PackageKit, "/org/freedesktop/PackageKit"),
		)

		backend, err := pko.GetBackendName(context.TODO())
		if err != nil {
			return fmt.Errorf("getting packagekit backend: %w", err)
		}

		log.Println(backend)

		tkPath, err := pko.CreateTransaction(context.TODO())
		if err != nil {
			return fmt.Errorf("creating transaction: %w", err)
		}

		log.Println(tkPath)

		tk := packagekit.NewOrg_Freedesktop_PackageKit_Transaction(
			conn.Object(packagekit.InterfaceOrg_Freedesktop_PackageKit_Transaction, tkPath),
		)
		err = tk.WhatProvides(context.TODO(), 0, []string{"vim"})
		if err != nil {
			_ = tk.Cancel(context.TODO())

			return fmt.Errorf("resolving vim: %w", err)
		}

		// packagekit.AddMatchSignal(conn, )

		return nil
	},
}
