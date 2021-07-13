package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var numTargetNodesPoolScale int

var kubernetesNodePoolScaleCmd = &cobra.Command{
	Use:     "scale",
	Short:   "Scale a node pool in a Kubernetes cluster",
	Example: "civo kubernetes node-pool scale CLUSTER_NAME NODEPOOL_ID [flags]",
	Args:    cobra.MinimumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return getAllKubernetesList(), cobra.ShellCompDirectiveNoFileComp
		}
		return getKubernetesList(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubernetesFindCluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		nodePool := []civogo.KubernetesClusterPoolConfig{}
		for _, v := range kubernetesFindCluster.Pools {
			nodePool = append(nodePool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
		}

		nodePool = utility.UpdateNodePool(nodePool, args[1], numTargetNodesPoolScale)

		configKubernetes := &civogo.KubernetesClusterConfig{
			Pools: nodePool,
		}

		kubernetesCluster, err := client.UpdateKubernetesCluster(kubernetesFindCluster.ID, configKubernetes)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": kubernetesCluster.ID, "name": kubernetesCluster.Name, "pool_id": args[1]})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The pool (%s) was scaled to (%s) in the cluster (%s)\n", utility.Green(args[1]), utility.Green(strconv.Itoa(numTargetNodesPoolScale)), utility.Green(kubernetesCluster.Name))
		}
	},
}
