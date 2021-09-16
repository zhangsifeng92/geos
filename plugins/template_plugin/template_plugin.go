//simple example for eos-go plugins
package template_plugin

import (
	"github.com/urfave/cli"
	. "github.com/zhangsifeng92/geos/plugins/appbase/app"
)

const TemplatePlug = PluginTypeName("TemplatePlugin")

var templatePlugin = App().RegisterPlugin(TemplatePlug, NewTemplatePlugin())

type TemplatePlugin struct {
	AbstractPlugin
	my *TemplatePluginImpl
}

type TemplatePluginImpl struct {
}

func NewTemplatePlugin() *TemplatePlugin {
	plugin := &TemplatePlugin{}
	plugin.my = &TemplatePluginImpl{}
	return plugin
}

func (c *TemplatePlugin) SetProgramOptions(options *[]cli.Flag) {
}

func (c *TemplatePlugin) PluginInitialize(options *cli.Context) {
}

func (c *TemplatePlugin) PluginStartup() {
}

func (c *TemplatePlugin) PluginShutdown() {
}
