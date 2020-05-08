package cmd

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var (
	shortDescription                                string
	Name, description, defaultUsername, cloudConfig string
)

var templateUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"change", "modify"},
	Short:   "Update a template",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		template, err := client.GetTemplateByCode(args[0])
		if err != nil {
			fmt.Printf("Unable to find the template: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		configTemplateUpdate := &civogo.Template{
			ImageID:  template.ImageID,
			VolumeID: template.VolumeID,
		}

		if Name != "" {
			configTemplateUpdate.Name = Name
		} else {
			configTemplateUpdate.Name = template.Name
		}

		if shortDescription != "" {
			configTemplateUpdate.ShortDescription = shortDescription
		} else {
			configTemplateUpdate.ShortDescription = template.ShortDescription
		}

		if description != "" {
			configTemplateUpdate.Description = description
		} else {
			configTemplateUpdate.Description = template.Description
		}

		if defaultUsername != "" {
			configTemplateUpdate.DefaultUsername = defaultUsername
		} else {
			configTemplateUpdate.DefaultUsername = template.DefaultUsername
		}

		if cloudConfig != "" {
			// reading the file
			data, err := ioutil.ReadFile(cloudConfig)
			if err != nil {
				fmt.Printf("Unable to read the cloud config file: %s\n", aurora.Red(err))
				os.Exit(1)
			}

			configTemplateUpdate.CloudConfig = string(data)
		} else {
			configTemplateUpdate.CloudConfig = template.CloudConfig
		}

		templateUpdate, err := client.UpdateTemplate(template.ID, configTemplateUpdate)
		if err != nil {
			fmt.Printf("Unable to update the template: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": templateUpdate.ID, "Name": templateUpdate.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Updated template with name %s with ID %s\n", aurora.Green(templateUpdate.Name), aurora.Green(templateUpdate.ID))
		}
	},
}