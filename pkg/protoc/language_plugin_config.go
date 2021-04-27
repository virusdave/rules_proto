package protoc

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/bazelbuild/bazel-gazelle/label"
)

// LanguagePluginConfig associates metadata with a plugin implementation.
type LanguagePluginConfig struct {
	// Name of the config, for the sake of configuration
	Name string
	// Label is the bazel label of the PluginInfo provider
	Label label.Label
	// Tool is the bazel label for the binary tool
	Tool label.Label
	// Options is a set of option strings.
	Options map[string]bool
	// Implementation is the Plugin implementation registered.
	Implementation Plugin
	// Enabled flag
	Enabled bool
}

func newLanguagePluginConfig(name string) *LanguagePluginConfig {
	return &LanguagePluginConfig{
		Name:    name,
		Options: make(map[string]bool),
		Enabled: true,
	}
}

// GetOptions returns the sorted list of options
func (c *LanguagePluginConfig) GetOptions() []string {
	opts := make([]string, 0)
	for opt, want := range c.Options {
		if !want {
			continue
		}
		opts = append(opts, opt)
	}
	sort.Strings(opts)
	return opts
}

func (c *LanguagePluginConfig) clone() *LanguagePluginConfig {
	clone := &LanguagePluginConfig{
		Label:          c.Label,
		Name:           c.Name,
		Tool:           c.Tool,
		Implementation: c.Implementation,
		Enabled:        c.Enabled,
	}
	for k, v := range c.Options {
		clone.Options[k] = v
	}
	return clone
}

// parseDirective parses the directive string or returns error.
func (c *LanguagePluginConfig) parseDirective(cfg *PackageConfig, d, param, value string) error {
	intent := parseIntent(param)
	switch intent.Value {
	case "enabled", "enable":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("enabled %s: %w", value, err)
		}
		c.Enabled = enabled
	case "label":
		l, err := label.Parse(value)
		if err != nil {
			return fmt.Errorf("label %q: %w", value, err)
		}
		c.Label = l
	case "tool":
		l, err := label.Parse(value)
		if err != nil {
			return fmt.Errorf("tool %q: %w", value, err)
		}
		c.Tool = l
	case "option":
		if intent.Negative {
			delete(c.Options, value)
		} else {
			c.Options[value] = true
		}
	default:
		return fmt.Errorf("invalid directive %q: unknown parameter %q", d, value)
	}

	return nil
}