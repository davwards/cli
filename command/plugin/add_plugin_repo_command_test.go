package plugin_test

import (
	"code.cloudfoundry.org/cli/actor/pluginaction"
	"code.cloudfoundry.org/cli/command/commandfakes"
	. "code.cloudfoundry.org/cli/command/plugin"
	"code.cloudfoundry.org/cli/command/plugin/pluginfakes"
	"code.cloudfoundry.org/cli/command/plugin/shared"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("add-plugin-repo command", func() {
	var (
		cmd        AddPluginRepoCommand
		testUI     *ui.UI
		fakeConfig *commandfakes.FakeConfig
		executeErr error
		fakeActor  *pluginfakes.FakeAddPluginRepoActor
	)

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
		fakeConfig = new(commandfakes.FakeConfig)
		fakeActor = new(pluginfakes.FakeAddPluginRepoActor)
		cmd = AddPluginRepoCommand{UI: testUI, Config: fakeConfig, Actor: fakeActor}

		fakeConfig.BinaryNameReturns("faceman")
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	Context("when the provided repo name already exists", func() {
		BeforeEach(func() {
			cmd.RequiredArgs.PluginRepoName = "some-repo"
			cmd.RequiredArgs.PluginRepoURL = "some-repo-URL"
			fakeActor.AddPluginRepositoryReturns(pluginaction.RepositoryNameTakenError{Name: "some-repo"})
		})

		It("errors with the RepositoryNameTakenError", func() {
			Expect(executeErr).To(MatchError(shared.RepositoryNameTakenError{Name: "some-repo"}))
		})
	})

	Context("when the provided repo URL already exists", func() {
		BeforeEach(func() {
			cmd.RequiredArgs.PluginRepoName = "some-repo"
			cmd.RequiredArgs.PluginRepoURL = "some-repo-URL"
			fakeActor.AddPluginRepositoryReturns(pluginaction.RepositoryURLTakenError{Name: "existing-repo", URL: "some-URL"})
		})

		It("errors with the RepositoryNameTakenError", func() {
			Expect(executeErr).To(MatchError(shared.RepositoryURLTakenError{Name: "existing-repo", URL: "some-URL"}))
		})
	})

	Context("when the actor returns a pluginaction.AddPluginRepoError", func() {
		BeforeEach(func() {
			cmd.RequiredArgs.PluginRepoName = "some-repo"
			cmd.RequiredArgs.PluginRepoURL = "some-URL"
			fakeActor.AddPluginRepositoryReturns(pluginaction.AddPluginRepositoryError{Name: "some-repo", URL: "some-URL", Message: "404"})
		})

		It("handles the error", func() {
			Expect(executeErr).To(MatchError(shared.AddPluginRepositoryError{Name: "some-repo", URL: "some-URL", Message: "404"}))
		})
	})

	Context("when no errors are encountered", func() {
		BeforeEach(func() {
			cmd.RequiredArgs.PluginRepoName = "some-repo"
			cmd.RequiredArgs.PluginRepoURL = "https://some-repo-URL"
			fakeActor.AddPluginRepositoryReturns(nil)
		})

		It("displays OK and message", func() {
			Expect(executeErr).ToNot(HaveOccurred())

			repoName, repoURL := fakeActor.AddPluginRepositoryArgsForCall(0)
			Expect(repoName).To(Equal("some-repo"))
			Expect(repoURL).To(Equal("https://some-repo-URL"))
			Expect(testUI.Out).To(Say("https://some-repo-URL added as 'some-repo'"))
		})
	})
})
