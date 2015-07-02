package admin

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/projectatomic/atomic-enterprise/pkg/cmd/admin/node"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/admin/policy"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/admin/project"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/admin/prune"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/admin/registry"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/admin/router"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/cli/cmd"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/experimental/buildchain"
	exipfailover "github.com/projectatomic/atomic-enterprise/pkg/cmd/experimental/ipfailover"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/server/admin"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/templates"
	cmdutil "github.com/projectatomic/atomic-enterprise/pkg/cmd/util"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/util/clientcmd"
	"github.com/projectatomic/atomic-enterprise/pkg/version"
)

const adminLong = `Atomic Enterprise Administrative Commands

Commands for managing an Atomic Enterprise cluster are exposed here. Many administrative
actions involve interaction with the Atomic Enterprise client as well.

Note: This is a beta release of Atomic Enterprise and may change significantly.  See
    https://github.com/projectatomic/atomic-enterprise for the latest information on Atomic Enterprise.`

func NewCommandAdmin(name, fullName string, out io.Writer) *cobra.Command {
	// Main command
	cmds := &cobra.Command{
		Use:   name,
		Short: "Tools for managing an Atomic Enterprise cluster",
		Long:  fmt.Sprintf(adminLong),
		Run:   cmdutil.DefaultSubCommandRun(out),
	}

	f := clientcmd.New(cmds.PersistentFlags())

	cmds.AddCommand(project.NewCmdNewProject(project.NewProjectRecommendedName, fullName+" "+project.NewProjectRecommendedName, f, out))
	cmds.AddCommand(policy.NewCmdPolicy(policy.PolicyRecommendedName, fullName+" "+policy.PolicyRecommendedName, f, out))
	cmds.AddCommand(exipfailover.NewCmdIPFailoverConfig(f, fullName, "ipfailover", out))
	cmds.AddCommand(router.NewCmdRouter(f, fullName, "router", out))
	cmds.AddCommand(registry.NewCmdRegistry(f, fullName, "registry", out))
	cmds.AddCommand(buildchain.NewCmdBuildChain(f, fullName, "build-chain"))
	cmds.AddCommand(node.NewCommandManageNode(f, node.ManageNodeCommandName, fullName+" "+node.ManageNodeCommandName, out))
	cmds.AddCommand(cmd.NewCmdConfig(fullName, "config"))
	cmds.AddCommand(prune.NewCommandPrune(prune.PruneRecommendedName, fullName+" "+prune.PruneRecommendedName, f, out))

	// TODO: these probably belong in a sub command
	cmds.AddCommand(admin.NewCommandCreateKubeConfig(admin.CreateKubeConfigCommandName, fullName+" "+admin.CreateKubeConfigCommandName, out))
	cmds.AddCommand(admin.NewCommandCreateBootstrapPolicyFile(admin.CreateBootstrapPolicyFileCommand, fullName+" "+admin.CreateBootstrapPolicyFileCommand, out))
	cmds.AddCommand(admin.NewCommandCreateBootstrapProjectTemplate(f, admin.CreateBootstrapProjectTemplateCommand, fullName+" "+admin.CreateBootstrapProjectTemplateCommand, out))
	cmds.AddCommand(admin.NewCommandOverwriteBootstrapPolicy(admin.OverwriteBootstrapPolicyCommandName, fullName+" "+admin.OverwriteBootstrapPolicyCommandName, fullName+" "+admin.CreateBootstrapPolicyFileCommand, out))
	cmds.AddCommand(admin.NewCommandNodeConfig(admin.NodeConfigCommandName, fullName+" "+admin.NodeConfigCommandName, out))
	// TODO: these should be rolled up together
	cmds.AddCommand(admin.NewCommandCreateMasterCerts(admin.CreateMasterCertsCommandName, fullName+" "+admin.CreateMasterCertsCommandName, out))
	cmds.AddCommand(admin.NewCommandCreateClient(admin.CreateClientCommandName, fullName+" "+admin.CreateClientCommandName, out))
	cmds.AddCommand(admin.NewCommandCreateKeyPair(admin.CreateKeyPairCommandName, fullName+" "+admin.CreateKeyPairCommandName, out))
	cmds.AddCommand(admin.NewCommandCreateServerCert(admin.CreateServerCertCommandName, fullName+" "+admin.CreateServerCertCommandName, out))
	cmds.AddCommand(admin.NewCommandCreateSignerCert(admin.CreateSignerCertCommandName, fullName+" "+admin.CreateSignerCertCommandName, out))

	// TODO: use groups
	templates.ActsAsRootCommand(cmds)

	if name == fullName {
		cmds.AddCommand(version.NewVersionCommand(fullName))
	}

	cmds.AddCommand(cmd.NewCmdOptions(out))

	return cmds
}
