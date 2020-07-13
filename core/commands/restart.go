package commands

import (
	"os"
	"os/exec"
	"time"

	"github.com/TRON-US/go-btfs-cmds"
	"github.com/cenkalti/backoff/v4"

	path "github.com/TRON-US/go-btfs/core/commands/storage"
)

var daemonStartup = func() *backoff.ExponentialBackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 2 * time.Second
	bo.MaxElapsedTime = 300 * time.Second
	bo.Multiplier = 1
	bo.MaxInterval = 2 * time.Second
	return bo
}()

const (
	postPathModificationName = "post-path-modification"
)

var restartCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Restart the daemon.",
		ShortDescription: `
Shutdown the runnning daemon and start a new daemon process.
And if specified a new btfs path, it will be applied.
`,
	},
	Options: []cmds.Option{
		cmds.BoolOption(postPathModificationName, "p", "post path modification").WithDefault(false),
	}, Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		daemonCmd := exec.Command(path.Excutable, "daemon")
		if req.Options[postPathModificationName].(bool) && path.StorePath != "" && path.OriginPath != "" {
			if err := path.MoveFolder(); err != nil {
				return err
			}

			if err := path.WriteProperties(); err != nil {
				return err
			}

			daemonCmd.Env = os.Environ()
			daemonCmd.Env = append(daemonCmd.Env, "BTFS_PATH="+path.StorePath)
		}
		if err := daemonCmd.Start(); err != nil {
			return err
		}
		os.Exit(0)
		return nil
	},
}
