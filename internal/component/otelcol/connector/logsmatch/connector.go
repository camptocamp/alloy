package logsmatch

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/connector/logsmatchconnector"
	"github.com/grafana/alloy/internal/component"
	"github.com/grafana/alloy/internal/component/otelcol"
	otelcolCfg "github.com/grafana/alloy/internal/component/otelcol/config"
	"github.com/grafana/alloy/internal/component/otelcol/connector"
	"github.com/grafana/alloy/internal/featuregate"
	"github.com/grafana/alloy/syntax"
	otelcomponent "go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pipeline"
)

var factory = logsmatchconnector.NewFactory()

type Arguments struct {
	Match        MatchArguments                   `alloy:"match,block,optional"`
	Metric       MetricArguments                  `alloy:"metric,block,optional"`
	Output       *otelcol.ConsumerArguments       `alloy:"output,block"`
	DebugMetrics otelcolCfg.DebugMetricsArguments `alloy:"debug_metrics,block,optional"`
}

type MatchArguments struct {
	Attribute string `alloy:"attribute,attr,optional"`
	Regexp    string `alloy:"regexp,attr,optional"`
}

type MetricArguments struct {
	Name        string `alloy:"name,attr,optional"`
	Description string `alloy:"description,attr,optional"`
	KeepTimestamp bool `alloy:"keep_timestamp,attr,optional"`
}

var (
	_ syntax.Defaulter    = (*Arguments)(nil)
	_ connector.Arguments = (*Arguments)(nil)
)

func init() {
	component.Register(component.Registration{
		Name:      "otelcol.connector.logsmatch",
		Stability: featuregate.StabilityGenerallyAvailable,
		Args:      Arguments{},
		Exports:   otelcol.ConsumerExports{},
		Build: func(opts component.Options, args component.Arguments) (component.Component, error) {
			return connector.New(opts, factory, args.(Arguments))
		},
	})
}

func (args *Arguments) SetToDefault() {
	config := factory.CreateDefaultConfig().(*logsmatchconnector.Config)

	args.Match.Attribute = config.Match.Attribute
	args.Match.Regexp = config.Match.Regexp
	args.Metric.Name = config.Metric.Name
	args.Metric.Description = config.Metric.Description
	args.Metric.KeepTimestamp = config.Metric.KeepTimestamp
	args.DebugMetrics.SetToDefault()
}

func (args Arguments) Convert() (otelcomponent.Config, error) {
	config := &logsmatchconnector.Config{
		Match: logsmatchconnector.MatchConfig{
			Attribute: args.Match.Attribute,
			Regexp:    args.Match.Regexp,
		},
		Metric: logsmatchconnector.MetricConfig{
			Name:        args.Metric.Name,
			Description: args.Metric.Description,
			KeepTimestamp: args.Metric.KeepTimestamp,
		},
	}

	return config, nil
}

func (args Arguments) Extensions() map[otelcomponent.ID]otelcomponent.Component {
	return nil
}

func (args Arguments) Exporters() map[pipeline.Signal]map[otelcomponent.ID]otelcomponent.Component {
	return nil
}

func (args Arguments) NextConsumers() *otelcol.ConsumerArguments {
	return args.Output
}

func (args Arguments) ConnectorType() int {
	return connector.ConnectorLogsToMetrics
}

func (args Arguments) DebugMetricsConfig() otelcolCfg.DebugMetricsArguments {
	return args.DebugMetrics
}
