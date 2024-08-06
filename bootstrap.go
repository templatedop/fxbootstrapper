package fxbootstrapper

import (
	"context"	
	"github.com/templatedop/fxdb"	
	"testing"
	config "github.com/templatedop/fxconfig"
	logger "github.com/templatedop/fxlogger"
	fxgin "github.com/templatedop/fxrouter"	
	"github.com/templatedop/fxvalidator"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// func RegisterServices() fx.Option {
// 	return fx.Provide(
// 		// health check probes
// 		fxhealthchecker.AsHealthCheckerProbe(repo.NewDBProbe),
// 	)
// }

var FxValidatorModule = fx.Module("validator",
	fx.Provide(fxvalidator.NewValidService),
	//fx.Invoke(CustomValidation)
)

type Bootstrapper struct {
	defaultOptions []fx.Option
}

var Ds = NewBootstrapper().WithOptions()

func BootstrapServer(ctx context.Context) *fx.App {
	return Ds.BoostrapApp(
		fxdb.FxDBModule,
		fxgin.RouterModule,
	)
}

func NewBootstrapper() *Bootstrapper {
	return &Bootstrapper{
		defaultOptions: []fx.Option{
			config.FxConfigModule,
			logger.FxLoggerModule,
			FxValidatorModule,
		},
	}
}

func (b *Bootstrapper) WithOptions(option ...fx.Option) *Bootstrapper {
	b.defaultOptions = append(b.defaultOptions, option...)

	return b
}

func (b *Bootstrapper) BoostrapApp(bootstrapOptions ...fx.Option) *fx.App {
	return fx.New(
		fx.WithLogger(logger.FxEventLogger),
		fx.Options(b.defaultOptions...),
		fx.Options(bootstrapOptions...),
	)
}

func (b *Bootstrapper) BoostrapAndRunApp(bootstrapOptions ...fx.Option) {
	b.BoostrapApp(bootstrapOptions...).Run()
}

func (b *Bootstrapper) BoostrapTestApp(t testing.TB, bootstrapOptions ...fx.Option) *fxtest.App {

	t.Setenv("APP_ENV", "test")

	return fxtest.New(
		t,
		fx.NopLogger,
		fx.Options(b.defaultOptions...),
		fx.Options(bootstrapOptions...),
	)
}

func (b *Bootstrapper) BoostrapAndRunTestApp(t testing.TB, bootstrapOptions ...fx.Option) {
	b.BoostrapTestApp(t, bootstrapOptions...).RequireStart().RequireStop()
}
