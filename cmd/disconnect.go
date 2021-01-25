package cmd

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/v8platform/rac"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"v8platform/onec-util/rac_addin"
)

type disconnectCmd struct {
	cluster string
	host    string
	name    string

	rac_addin.ServersRemove
}

func (c *disconnectCmd) ClusterSig() (uuid string) {
	return c.cluster
}

func (c *disconnectCmd) Auth() (usr, pwd string) {
	return c.User, c.Pwd
}

func (c *disconnectCmd) run(ctx *cli.Context) error {
	mng, err := rac.NewManager(c.host, rac.ManagerOptions{
		Timeout:         2 * time.Second,
		TryTimeoutCount: 3,
		DetectCluster:   false,
	})

	if err != nil {
		return err
	}

	c.cluster, err = getCluster(mng)

	if err != nil {
		return err
	}

	err = mng.SetDefCluster(c, c)
	if err != nil {
		return err
	}

	switch {
	case len(c.UUID) > 0:
		log.Printf("Delete server by uuid <%s>", c.UUID)
		_, err = rac_addin.Servers(mng, c.ServersRemove)

		if err != nil {
			return err
		}
		return nil
	case len(c.name) > 0:

		log.Printf("Finding server by name <%s>", c.name)

		servers, err := rac_addin.Servers(mng, rac_addin.ServersList{})

		if err != nil {
			return err
		}

		for _, info := range servers.List {

			if strings.EqualFold(strings.Trim(info.Name, "\""), c.name) {
				c.UUID = info.UUID
				break
			}
		}

		if len(c.UUID) == 0 {
			log.Printf("Not found server by name <%s>", c.name)
			return nil
		}

		log.Printf("Delete server by uuid <%s>", c.UUID)

		_, err = rac_addin.Servers(mng, c.ServersRemove)

		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (c *disconnectCmd) Cmd() *cli.Command {

	cmd := &cli.Command{
		Category:    "cluster",
		Name:        "disconnect",
		Usage:       "disconnect work server from 1S.Enterprise cluster",
		Description: `Выполняет отключение рабочего сервера 1С.Предприятие от кластера сервера`,
		ArgsUsage:   "HOST",
		Action:      c.run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Destination: &c.cluster, Name: "cluster-id",
				Value: "", Usage: "cluster uuid for join new work server"},
			&cli.StringFlag{
				Destination: &c.UUID, Name: "server-id", Aliases: []string{"I"},
				Value: "", Usage: "work server uuid", EnvVars: []string{"SERVER_UUID"}},
			&cli.StringFlag{
				Destination: &c.name, Name: "server", Aliases: []string{"N"},
				Value: "", Usage: "work server name", EnvVars: []string{"SERVER_NAME"}},
		},
		Before: func(ctx *cli.Context) error {

			c.host, _ = os.LookupEnv("CLUSTER_HOST")

			if ctx.Args().Len() < 1 && len(c.host) == 0 {
				cli.ShowSubcommandHelp(ctx)
				return errors.New("WRONG USAGE: Requires a HOST argument")
			}

			_, _, err := net.SplitHostPort(ctx.Args().First())
			if err == nil {
				c.host = ctx.Args().First()
				return nil
			}
			cli.ShowSubcommandHelp(ctx)
			return fmt.Errorf("WRONG USAGE: Invalid HOST argument %s", err)
		},
	}
	return cmd
}
