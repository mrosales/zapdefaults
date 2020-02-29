zapdefaults
===========

Easy configuration for ⚡[Zap](https://github.com/uber-go/zap)!

Motivation
----------

I use Zap in just about every Go project I set up. I almost always make a couple of
modifications to the zap defaults for `ProductionConfig` and `DevelopmentConfig`, and
after several times of copying them into new projects, I've just set up this project to
capture the setup I always use.

Getting Started
---------------


```go
package main

import (
	"github.com/mrosales/zapdefaults"
	"go.uber.org/zap"
)

func main() {
	// Use the defaults
	logger := zapdefaults.MustLogger()

	logger.Info("Hello, World!")

	logger = zapdefaults.MustLogger(
		zapdefaults.Dynamic,
		zapdefaults.EnvironmentOverrides("ZAP"),
	)
	logger.Info("this is the same the last one!")

	logger = zapdefaults.MustLogger(zapdefaults.Development)
	logger.Info("explicitly use the development config")

	logger = zapdefaults.MustLogger(zapdefaults.Production)
	logger.Info("explicitly use the production config")

	logger = zapdefaults.MustLogger(
		zapdefaults.Dynamic,
		zapdefaults.OptionFunc(func(cfg *zap.Config) error {
			cfg.DisableCaller = true
			return nil
		}),
	)
	logger.Info("use a custom option to tweak a setting")
}
```
