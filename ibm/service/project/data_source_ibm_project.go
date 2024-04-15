// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package project

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/project-go-sdk/projectv1"
)

func DataSourceIbmProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmProjectRead,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique project ID.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An IBM Cloud resource name that uniquely identifies a resource.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
			},
			"cumulative_needs_attention_view": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The cumulative list of needs attention items for a project. If the view is successfully retrieved, an empty or nonempty array is returned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The event name.",
						},
						"event_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique ID for this individual event.",
						},
						"config_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique ID for the configuration.",
						},
						"config_version": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version number of the configuration.",
						},
					},
				},
			},
			"cumulative_needs_attention_view_error": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A value of `true` indicates that the fetch of the needs attention items failed. This property only exists if there was an error when you retrieved the cumulative needs attention view.",
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IBM Cloud location where a resource is deployed.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group ID where the project's data and tools are created.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project status value.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name where the project's data and tools are created.",
			},
			"event_notifications_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN of the Event Notifications instance if one is connected to this project.",
			},
			"configs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The project configurations. These configurations are only included in the response of creating a project if a configuration array is specified in the request payload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"approved_version": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A summary of a project configuration version.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "A summary of the definition in a project configuration version.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"environment_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The ID of the project environment.",
												},
												"locator_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the Schematics workspace has one.> There are 3 scenarios:> 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.> 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the existing schematics workspace.> 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.> For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform template](/docs/schematics?topic=schematics-sch-create-wks).",
												},
											},
										},
									},
									"state": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The state of the configuration.",
									},
									"version": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The version number of the configuration.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A URL.",
									},
								},
							},
						},
						"deployed_version": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A summary of a project configuration version.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "A summary of the definition in a project configuration version.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"environment_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The ID of the project environment.",
												},
												"locator_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the Schematics workspace has one.> There are 3 scenarios:> 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.> 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the existing schematics workspace.> 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.> For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform template](/docs/schematics?topic=schematics-sch-create-wks).",
												},
											},
										},
									},
									"state": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The state of the configuration.",
									},
									"version": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The version number of the configuration.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A URL.",
									},
								},
							},
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.",
						},
						"version": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version of the configuration.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the configuration.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
						},
						"modified_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL.",
						},
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The description of a project configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A project configuration description.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The configuration name. It's unique within the account across projects and regions.",
									},
									"locator_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the Schematics workspace has one.> There are 3 scenarios:> 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.> 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the existing schematics workspace.> 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.> For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform template](/docs/schematics?topic=schematics-sch-create-wks).",
									},
								},
							},
						},
						"project": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The project that is referenced by this resource.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique ID.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A URL.",
									},
									"definition": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The definition of the project reference.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the project.",
												},
											},
										},
									},
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An IBM Cloud resource name that uniquely identifies a resource.",
									},
								},
							},
						},
						"deployment_model": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configuration type.",
						},
					},
				},
			},
			"environments": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The project environment. These environments are only included in the response if project environments were created on the project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The environment ID as a friendly name.",
						},
						"project": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The project that is referenced by this resource.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique ID.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A URL.",
									},
									"definition": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The definition of the project reference.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the project.",
												},
											},
										},
									},
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An IBM Cloud resource name that uniquely identifies a resource.",
									},
								},
							},
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL.",
						},
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The environment definition that is used in the project collection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the environment.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the environment. It's unique within the account across projects and regions.",
									},
								},
							},
						},
					},
				},
			},
			"definition": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The definition of the project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the project.  It's unique within the account across regions.",
						},
						"destroy_on_delete": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The policy that indicates whether the resources are destroyed or not when a project is deleted.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A brief explanation of the project's use in the configuration of a deployable architecture. You can create a project without providing a description.",
						},
						"monitoring_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "A boolean flag to enable automatic drift detection. Use this field to run a daily check to compare your configurations to your deployed resources to detect any difference.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmProjectRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getProjectOptions := &projectv1.GetProjectOptions{}

	getProjectOptions.SetID(d.Get("project_id").(string))

	project, _, err := projectClient.GetProjectWithContext(context, getProjectOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetProjectWithContext failed: %s", err.Error()), "(Data) ibm_project", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s", *getProjectOptions.ID))

	if err = d.Set("crn", project.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(project.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	cumulativeNeedsAttentionView := []map[string]interface{}{}
	if project.CumulativeNeedsAttentionView != nil {
		for _, modelItem := range project.CumulativeNeedsAttentionView {
			modelMap, err := dataSourceIbmProjectCumulativeNeedsAttentionToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project", "read")
				return tfErr.GetDiag()
			}
			cumulativeNeedsAttentionView = append(cumulativeNeedsAttentionView, modelMap)
		}
	}
	if err = d.Set("cumulative_needs_attention_view", cumulativeNeedsAttentionView); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting cumulative_needs_attention_view: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("cumulative_needs_attention_view_error", project.CumulativeNeedsAttentionViewError); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting cumulative_needs_attention_view_error: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("location", project.Location); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting location: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_group_id", project.ResourceGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting resource_group_id: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state", project.State); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("href", project.Href); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_group", project.ResourceGroup); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting resource_group: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("event_notifications_crn", project.EventNotificationsCrn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting event_notifications_crn: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	configs := []map[string]interface{}{}
	if project.Configs != nil {
		for _, modelItem := range project.Configs {
			modelMap, err := dataSourceIbmProjectProjectConfigSummaryToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project", "read")
				return tfErr.GetDiag()
			}
			configs = append(configs, modelMap)
		}
	}
	if err = d.Set("configs", configs); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting configs: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	environments := []map[string]interface{}{}
	if project.Environments != nil {
		for _, modelItem := range project.Environments {
			modelMap, err := dataSourceIbmProjectProjectEnvironmentSummaryToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project", "read")
				return tfErr.GetDiag()
			}
			environments = append(environments, modelMap)
		}
	}
	if err = d.Set("environments", environments); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting environments: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	definition := []map[string]interface{}{}
	if project.Definition != nil {
		modelMap, err := dataSourceIbmProjectProjectDefinitionPropertiesToMap(project.Definition)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project", "read")
			return tfErr.GetDiag()
		}
		definition = append(definition, modelMap)
	}
	if err = d.Set("definition", definition); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting definition: %s", err), "(Data) ibm_project", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIbmProjectCumulativeNeedsAttentionToMap(model *projectv1.CumulativeNeedsAttention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Event != nil {
		modelMap["event"] = model.Event
	}
	if model.EventID != nil {
		modelMap["event_id"] = model.EventID
	}
	if model.ConfigID != nil {
		modelMap["config_id"] = model.ConfigID
	}
	if model.ConfigVersion != nil {
		modelMap["config_version"] = flex.IntValue(model.ConfigVersion)
	}
	return modelMap, nil
}

func dataSourceIbmProjectProjectConfigSummaryToMap(model *projectv1.ProjectConfigSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ApprovedVersion != nil {
		approvedVersionMap, err := dataSourceIbmProjectProjectConfigVersionSummaryToMap(model.ApprovedVersion)
		if err != nil {
			return modelMap, err
		}
		modelMap["approved_version"] = []map[string]interface{}{approvedVersionMap}
	}
	if model.DeployedVersion != nil {
		deployedVersionMap, err := dataSourceIbmProjectProjectConfigVersionSummaryToMap(model.DeployedVersion)
		if err != nil {
			return modelMap, err
		}
		modelMap["deployed_version"] = []map[string]interface{}{deployedVersionMap}
	}
	modelMap["id"] = model.ID
	modelMap["version"] = flex.IntValue(model.Version)
	modelMap["state"] = model.State
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["modified_at"] = model.ModifiedAt.String()
	modelMap["href"] = model.Href
	definitionMap, err := dataSourceIbmProjectProjectConfigSummaryDefinitionToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	projectMap, err := dataSourceIbmProjectProjectReferenceToMap(model.Project)
	if err != nil {
		return modelMap, err
	}
	modelMap["project"] = []map[string]interface{}{projectMap}
	if model.DeploymentModel != nil {
		modelMap["deployment_model"] = model.DeploymentModel
	}
	return modelMap, nil
}

func dataSourceIbmProjectProjectConfigVersionSummaryToMap(model *projectv1.ProjectConfigVersionSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	definitionMap, err := dataSourceIbmProjectProjectConfigVersionDefinitionSummaryToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["state"] = model.State
	modelMap["version"] = flex.IntValue(model.Version)
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIbmProjectProjectConfigVersionDefinitionSummaryToMap(model *projectv1.ProjectConfigVersionDefinitionSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnvironmentID != nil {
		modelMap["environment_id"] = model.EnvironmentID
	}
	if model.LocatorID != nil {
		modelMap["locator_id"] = model.LocatorID
	}
	return modelMap, nil
}

func dataSourceIbmProjectProjectConfigSummaryDefinitionToMap(model *projectv1.ProjectConfigSummaryDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["description"] = model.Description
	modelMap["name"] = model.Name
	if model.LocatorID != nil {
		modelMap["locator_id"] = model.LocatorID
	}
	return modelMap, nil
}

func dataSourceIbmProjectProjectReferenceToMap(model *projectv1.ProjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["href"] = model.Href
	definitionMap, err := dataSourceIbmProjectProjectDefinitionReferenceToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["crn"] = model.Crn
	return modelMap, nil
}

func dataSourceIbmProjectProjectDefinitionReferenceToMap(model *projectv1.ProjectDefinitionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIbmProjectProjectEnvironmentSummaryToMap(model *projectv1.ProjectEnvironmentSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	projectMap, err := dataSourceIbmProjectProjectReferenceToMap(model.Project)
	if err != nil {
		return modelMap, err
	}
	modelMap["project"] = []map[string]interface{}{projectMap}
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["href"] = model.Href
	definitionMap, err := dataSourceIbmProjectProjectEnvironmentSummaryDefinitionToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	return modelMap, nil
}

func dataSourceIbmProjectProjectEnvironmentSummaryDefinitionToMap(model *projectv1.ProjectEnvironmentSummaryDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["description"] = model.Description
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIbmProjectProjectDefinitionPropertiesToMap(model *projectv1.ProjectDefinitionProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["destroy_on_delete"] = model.DestroyOnDelete
	modelMap["description"] = model.Description
	if model.MonitoringEnabled != nil {
		modelMap["monitoring_enabled"] = model.MonitoringEnabled
	}
	return modelMap, nil
}
