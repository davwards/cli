package shared

import "code.cloudfoundry.org/cli/actor/pluginaction"

func HandleError(err error) error {
	switch e := err.(type) {
	case pluginaction.PluginNotFoundError:
		return PluginNotFoundError{Name: e.Name}
	case pluginaction.GettingPluginRepositoryError:
		return GettingPluginRepositoryError{Name: e.Name, Message: e.Message}
	case pluginaction.RepositoryNameTakenError:
		return RepositoryNameTakenError{Name: e.Name}
	case pluginaction.RepositoryURLTakenError:
		return RepositoryURLTakenError{Name: e.Name, URL: e.URL}
	case pluginaction.AddPluginRepositoryError:
		return AddPluginRepositoryError{Name: e.Name, URL: e.URL, Message: e.Message}
	case pluginaction.PluginInvalidError:
		return PluginInvalidError{Path: e.Path, WrappedErrMessage: e.Error()}
	case pluginaction.PluginCommandsConflictError:
		return PluginCommandsConflictError{
			PluginName:     e.PluginName,
			PluginVersion:  e.PluginVersion,
			CommandNames:   e.CommandNames,
			CommandAliases: e.CommandAliases,
		}
	}
	return err
}
