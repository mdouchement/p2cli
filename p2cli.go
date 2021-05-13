package main

import (
	builtinjson "encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/noirbizarre/gonja"
	"github.com/noirbizarre/gonja/exec"
	"github.com/noirbizarre/gonja/nodes"
	"github.com/noirbizarre/gonja/tokens"
	"github.com/spf13/cobra"
)

var (
	version  = "dev"
	revision = "none"
	date     = "unknown"
)

type controller struct {
	prefix   string
	varsfile string

	konf     *koanf.Koanf
	template *exec.Template
}

func main() {
	ctrl := &controller{}

	c := &cobra.Command{
		Use:   "p2cli",
		Short: "Portable Jinja2 (Ansible) template engine",
		Example: strings.Join([]string{
			"p2cli -p 'STANDARDFILE_' example.yml.j2 > example.yml",
			"p2cli -v example.varsfile.yml example.yml.j2 > example.yml",
			"\nSources:\nhttps://github.com/mdouchement/p2cli",
		}, "\n"),
		Args:         cobra.ExactArgs(1),
		Version:      fmt.Sprintf("%s - build %.7s @ %s - %s", version, revision, date, runtime.Version()),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			//
			// Load
			//

			err := ctrl.load()
			if err != nil {
				return err
			}

			//
			// Process template
			//

			ctrl.template, err = gonja.DefaultEnv.GetTemplate(args[0])
			if err != nil {
				return err
			}

			err = ctrl.checkvars()
			if err != nil {
				return err
			}

			b, err := ctrl.template.Execute(ctrl.konf.Raw())
			if err != nil {
				return err
			}

			//
			// Display
			//

			fmt.Println(b)
			return nil
		},
	}
	c.Flags().StringVarP(&ctrl.varsfile, "varsfile", "v", "From Environment", "the vars_file path, format is detected based on the file extension (json|toml|yaml,yml)")
	c.Flags().StringVarP(&ctrl.prefix, "env-prefix", "p", "", "the prefix en environment variables (e.g. 'STANDARDFILE_')")
	c.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Version for p2cli",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(c.Version)
		},
	})

	if err := c.Execute(); err != nil {
		switch {
		case strings.Contains(err.Error(), "unknown shorthand flag"):
			fallthrough
		case strings.Contains(err.Error(), "accepts "):
			c.Println(c.UsageString())
		}
		os.Exit(1)
	}
}

func (ctrl *controller) load() (err error) {
	const delimiter = "."
	const envdelimiter = "__"

	ctrl.konf = koanf.New(delimiter)

	format := ctrl.varsfile
	if format != "From Environment" {
		format = strings.TrimLeft(filepath.Ext(format), ".")
	}

	switch format {
	case "env", "From Environment":
		return ctrl.konf.Load(env.ProviderWithValue(ctrl.prefix, envdelimiter, func(k, v string) (string, interface{}) {
			k = strings.TrimPrefix(k, ctrl.prefix)
			k = strings.ToLower(strings.TrimPrefix(k, ctrl.prefix))

			var value interface{}
			err := builtinjson.Unmarshal([]byte(v), &value)
			if err != nil {
				value = v // Fallback on raw value if v is not a JSON compatible
			}

			return k, value
		}), nil)
	case "json":
		return ctrl.konf.Load(file.Provider(ctrl.varsfile), json.Parser())
	case "toml":
		return ctrl.konf.Load(file.Provider(ctrl.varsfile), toml.Parser())
	case "yml", "yaml":
		return ctrl.konf.Load(file.Provider(ctrl.varsfile), yaml.Parser())
	default:
		return errors.New("unsupported input format")
	}
}

func (ctrl *controller) checkvars() (err error) {
	notfounds := []string{}

	for _, n := range ctrl.template.Parser.Template.Nodes {
		out, ok := n.(*nodes.Output)
		if !ok {
			continue
		}

		tk := out.Expression.Position()
		if tk.Type != tokens.Name {
			continue
		}

		// Use litter.Dump(out) to dig in the data and find the wanted node info

		// Check if there is a default defined
		var asDefault bool
		fexp, ok := out.Expression.(*nodes.FilteredExpression)
		if ok {
			for _, filter := range fexp.Filters {
				if filter.Name == "default" {
					asDefault = true
				}
			}
		}

		if !asDefault && !ctrl.konf.Exists(tk.Val) {
			// no variable and no default
			notfounds = append(notfounds, tk.Val)
		}
	}

	if len(notfounds) > 0 {
		return fmt.Errorf("no value found for key(s) [%s]", strings.Join(notfounds, ", "))
	}
	return nil
}
