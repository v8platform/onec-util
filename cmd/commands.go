package cmd

import (
	"github.com/urfave/cli/v2"
)

var Commands = []Command{

	&joinCmd{},
	&disconnectCmd{},
	//&clusterCmd{
	//	sub: []Command{
	//		//&disconnectCmd{},
	//		&joinCmd{},
	//	},
	//},
}

type Command interface {
	Cmd() *cli.Command
}
