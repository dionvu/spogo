package main

import (
	"github.com/dionvu/spogo/cli"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/session"
)

func main() {
	c, err := config.New()
	errors.Catch(err)
	errors.Catch(c.Load())

	session, err := session.New(c)
	errors.Catch(err)

	player, err := player.New(c)
	errors.Catch(err)

	cli := cli.New(session, player)

	err = cli.SetUp(c)
	errors.Catch(err)

	cli.Run()
}
