package shared

import (
	"fmt"
	"strings"
)

type PluginNotFoundError struct {
	Name string
}

func (e PluginNotFoundError) Error() string {
	return "Plugin {{.Name}} does not exist."
}

func (e PluginNotFoundError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"Name": e.Name,
	})
}

type NoPluginRepositoriesError struct{}

func (e NoPluginRepositoriesError) Error() string {
	return "No plugin repositories registered to search for plugin updates."
}

func (e NoPluginRepositoriesError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error())
}

// GettingPluginRepositoryError is returned when there's an error
// accessing the plugin repository
type GettingPluginRepositoryError struct {
	Name    string
	Message string
}

func (e GettingPluginRepositoryError) Error() string {
	return "Could not get plugin repository '{{.RepositoryName}}': {{.ErrorMessage}}"
}

func (e GettingPluginRepositoryError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{"RepositoryName": e.Name, "ErrorMessage": e.Message})
}

// RepositoryNameTakenError is returned when adding a plugin repository
// fails due to a repository already existing with the same name
type RepositoryNameTakenError struct {
	Name string
}

func (e RepositoryNameTakenError) Error() string {
	return "Plugin repo named '{{.RepositoryName}}' already exists, please use another name."
}

func (e RepositoryNameTakenError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{"RepositoryName": e.Name})
}

// RepositoryURLTakenError is returned when adding a plugin repository
// fails due to a repository already existing with the same URL
type RepositoryURLTakenError struct {
	Name string
	URL  string
}

func (e RepositoryURLTakenError) Error() string {
	return "{{.RepositoryURL}} ({{.RepositoryName}}) already exists."
}

func (e RepositoryURLTakenError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"RepositoryName": e.Name,
		"RepositoryURL":  e.URL,
	})
}

type AddPluginRepositoryError struct {
	Name    string
	URL     string
	Message string
}

func (e AddPluginRepositoryError) Error() string {
	return "Could not add repository '{{.RepositoryName}}' from {{.RepositoryURL}}: {{.Message}}"
}

func (e AddPluginRepositoryError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"RepositoryName": e.Name,
		"RepositoryURL":  e.URL,
		"Message":        e.Message,
	})
}

// FileNotFoundError is returned when a local plugin binary is not found during
// installation.
type FileNotFoundError struct {
	Path string
}

func (e FileNotFoundError) Error() string {
	return "File not found locally, make sure the file exists at given path {{.FilePath}}"
}

func (e FileNotFoundError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"FilePath": e.Path,
	})
}

// PluginInstallationCancelled is returned when the user chooses no during
// plugin installation confirmation.
type PluginInstallationCancelled struct {
}

func (e PluginInstallationCancelled) Error() string {
	return "Plugin installation cancelled"
}

func (e PluginInstallationCancelled) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error())
}

// PluginInvalidError is returned with a plugin is invalid because it is
// missing a name or has 0 commands.
type PluginInvalidError struct {
	Path              string
	WrappedErrMessage string
}

func (e PluginInvalidError) Error() string {
	return e.WrappedErrMessage
}

func (e PluginInvalidError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"Path": e.Path,
	})
}

// PluginCommandConflictError is returned when a plugin command name conflicts
// with a native or existing plugin command name.
type PluginCommandsConflictError struct {
	PluginName     string
	PluginVersion  string
	CommandNames   []string
	CommandAliases []string
}

func (e PluginCommandsConflictError) Error() string {
	switch {
	case len(e.CommandNames) > 0 && len(e.CommandAliases) > 0:
		return "Plugin {{.PluginName}} v{{.PluginVersion}} could not be installed as it contains commands with names and aliases that are already used: {{.CommandNamesAndAliases}}."
	case len(e.CommandNames) > 0:
		return "Plugin {{.PluginName}} v{{.PluginVersion}} could not be installed as it contains commands with names that are already used: {{.CommandNames}}."
	case len(e.CommandAliases) > 0:
		return "Plugin {{.PluginName}} v{{.PluginVersion}} could not be installed as it contains commands with aliases that are already used: {{.CommandAliases}}."
	default:
		return "Plugin {{.PluginName}} v{{.PluginVersion}} could not be installed as it contains commands with names or aliases that are already used."
	}
}

func (e PluginCommandsConflictError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"PluginName":             e.PluginName,
		"PluginVersion":          e.PluginVersion,
		"CommandNames":           strings.Join(e.CommandNames, ", "),
		"CommandAliases":         strings.Join(e.CommandAliases, ", "),
		"CommandNamesAndAliases": strings.Join(append(e.CommandNames, e.CommandAliases...), ", "),
	})
}

// PluginAlreadyInstalledError is returned when the plugin has the same name as
// an installed plugin.
type PluginAlreadyInstalledError struct {
	Name    string
	Version string
	Path    string
}

func (e PluginAlreadyInstalledError) Error() string {
	return "Plugin {{.Name}} {{.Version}} could not be installed. A plugin with that name is already installed.\nTIP: Use '{{.Command}}' to force a reinstall."
}

func (e PluginAlreadyInstalledError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"Name":    e.Name,
		"Version": e.Version,
		"Command": fmt.Sprintf("cf install-plugin %s -f", e.Path),
	})
}
