package cmd

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/v8platform/rac"
	"net"
	"os"
	"time"
	"v8platform/onec-util/rac_addin"
)

type joinCmd struct {
	cluster string
	host    string

	usr, pwd string

	rac_addin.ServersInsert
}

func (c *joinCmd) ClusterSig() (uuid string) {
	return c.cluster
}
func (c *joinCmd) Auth() (usr, pwd string) {
	return c.usr, c.pwd
}

func (c *joinCmd) run(context *cli.Context) error {

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

	//servers, err := rac_addin.Servers(mng, rac_addin.ServersList{})
	//
	//if err != nil {
	//	return err
	//}
	//pp.Println(servers.List)

	_, err = rac_addin.Servers(mng, c.ServersInsert)

	if err != nil {
		return err
	}
	//pp.Println(servers)

	return nil
}

func getCluster(mng *rac.Manager) (string, error) {
	clusters, err := mng.Clusters(rac.ClustersList{})
	if err != nil {
		return "", err
	}

	if len(clusters.List) > 1 {
		return "", errors.New("too many clusters in server. Need set cluster uuid")
	}

	if len(clusters.List) == 0 {
		return "", errors.New("cluster not found in server")
	}

	return clusters.List[0].UUID, nil
}

func (c *joinCmd) Cmd() *cli.Command {

	cmd := &cli.Command{
		Category:    "cluster",
		Name:        "join",
		Usage:       "join work server to 1S.Enterprise cluster",
		Description: `выполняет присоединение рабочего сервера 1С.Предприятие к кластеру сервера `,
		ArgsUsage:   "HOST",
		Action:      c.run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Destination: &c.cluster, Name: "cluster-id",
				Value: "", Usage: "cluster uuid for join new work server"},
			&cli.StringFlag{
				Destination: &c.Name, Name: "name", Aliases: []string{"N"},
				Value: "", Usage: "work server name", EnvVars: []string{"SERVER_NAME"}, Required: true},
			&cli.StringFlag{
				Destination: &c.AgentHost, Name: "host",
				Value: "", Usage: "work server agent host", EnvVars: []string{"SERVER_AGENT_HOST"}, Required: true},
			&cli.IntFlag{
				Destination: &c.AgentPort, Name: "port",
				Value: 1540, Usage: "work server agent port", EnvVars: []string{"SERVER_AGENT_HOST"}, DefaultText: "1540"},
			&cli.StringFlag{
				Destination: &c.PortRange, Name: "port-range",
				Value: "1560:1591", Usage: "work server port range", EnvVars: []string{"SERVER_PORT_RANGE"}, DefaultText: "1560:1591"},
			&cli.StringFlag{
				Destination: &c.Using, Name: "using",
				Value: "normal", Usage: "variant of using work server (main, normal)", EnvVars: []string{"SERVER_USING"}, DefaultText: "normal"},
			&cli.StringFlag{
				Destination: &c.DedicateManagers, Name: "dedicate-managers",
				Value: "none", Usage: "вариант размещения менеджеров сервисов (all, none)", EnvVars: []string{"SERVER_DEDICATE_MANAGERS"}, DefaultText: "none"},
			&cli.IntFlag{
				Destination: &c.ClusterPort, Name: "cluster-port",
				Value: 1541, Usage: "номер порта главного менеджера кластера", EnvVars: []string{"SERVER_CLUSTER_PORT"}, DefaultText: "1541"},
			&cli.Int64Flag{
				Destination: &c.MemoryLimit, Name: "memory-limit",
				Value: 0, Usage: "предел использования памяти рабочими процессами (kilobytes)", EnvVars: []string{"SERVER_MEMORY_LIMIT"}},
			&cli.Int64Flag{
				Destination: &c.ConnectionsLimit, Name: "connections-limit",
				Value: 128, Usage: "максимальное количество соединения на рабочий процесс", EnvVars: []string{"SERVER_CONNECTIONS_LIMIT"}, DefaultText: "128"},
			&cli.Int64Flag{
				Destination: &c.CriticalTotalMemory, Name: "total-memory",
				Value: 0, Usage: "максимальный объем памяти процессов рабочего сервера (bytes)", EnvVars: []string{"SERVER_TOTAL_MEMORY"}},
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
	//cmd.Args = func(cmd *cobra.Command, args []string) error {
	//	if len(args) < 1 {
	//		return errors.New("requires a host argument")
	//	}
	//
	//	_, _, err := net.SplitHostPort(args[0])
	//	if err == nil {
	//		return nil
	//	}
	//	return fmt.Errorf("invalid host %s", err)
	//}

	return cmd
}
