// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	oci_database "github.com/oracle/oci-go-sdk/database"
)

func DatabaseDatabasesDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readDatabaseDatabases,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"compartment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_home_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(DatabaseDatabaseDataSource()),
			},
		},
	}
}

func readDatabaseDatabases(d *schema.ResourceData, m interface{}) error {
	sync := &DatabaseDatabasesDataSourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).databaseClient

	return ReadResource(sync)
}

type DatabaseDatabasesDataSourceCrud struct {
	D      *schema.ResourceData
	Client *oci_database.DatabaseClient
	Res    *oci_database.ListDatabasesResponse
}

func (s *DatabaseDatabasesDataSourceCrud) VoidState() {
	s.D.SetId("")
}

func (s *DatabaseDatabasesDataSourceCrud) Get() error {
	request := oci_database.ListDatabasesRequest{}

	if compartmentId, ok := s.D.GetOkExists("compartment_id"); ok {
		tmp := compartmentId.(string)
		request.CompartmentId = &tmp
	}

	if dbHomeId, ok := s.D.GetOkExists("db_home_id"); ok {
		tmp := dbHomeId.(string)
		request.DbHomeId = &tmp
	}

	if dbName, ok := s.D.GetOkExists("db_name"); ok {
		tmp := dbName.(string)
		request.DbName = &tmp
	}

	if state, ok := s.D.GetOkExists("state"); ok {
		request.LifecycleState = oci_database.DatabaseSummaryLifecycleStateEnum(state.(string))
	}

	request.RequestMetadata.RetryPolicy = getRetryPolicy(false, "database")

	response, err := s.Client.ListDatabases(context.Background(), request)
	if err != nil {
		return err
	}

	s.Res = &response
	request.Page = s.Res.OpcNextPage

	for request.Page != nil {
		listResponse, err := s.Client.ListDatabases(context.Background(), request)
		if err != nil {
			return err
		}

		s.Res.Items = append(s.Res.Items, listResponse.Items...)
		request.Page = listResponse.OpcNextPage
	}

	return nil
}

func (s *DatabaseDatabasesDataSourceCrud) SetData() error {
	if s.Res == nil {
		return nil
	}

	s.D.SetId(GenerateDataSourceID())
	resources := []map[string]interface{}{}

	for _, r := range s.Res.Items {
		database := map[string]interface{}{
			"compartment_id": *r.CompartmentId,
			"db_home_id":     *r.DbHomeId,
		}

		if r.CharacterSet != nil {
			database["character_set"] = *r.CharacterSet
		}

		if r.ConnectionStrings != nil {
			database["connection_strings"] = []interface{}{DatabaseConnectionStringsToMap(r.ConnectionStrings)}
		} else {
			database["connection_strings"] = nil
		}

		if r.DbBackupConfig != nil {
			database["db_backup_config"] = []interface{}{DbBackupConfigToMap(r.DbBackupConfig)}
		} else {
			database["db_backup_config"] = nil
		}

		if r.DbName != nil {
			database["db_name"] = *r.DbName
		}

		if r.DbUniqueName != nil {
			database["db_unique_name"] = *r.DbUniqueName
		}

		if r.DbWorkload != nil {
			database["db_workload"] = *r.DbWorkload
		}

		if r.DefinedTags != nil {
			database["defined_tags"] = definedTagsToMap(r.DefinedTags)
		}

		database["freeform_tags"] = r.FreeformTags

		if r.Id != nil {
			database["id"] = *r.Id
		}

		if r.LifecycleDetails != nil {
			database["lifecycle_details"] = *r.LifecycleDetails
		}

		if r.NcharacterSet != nil {
			database["ncharacter_set"] = *r.NcharacterSet
		}

		if r.PdbName != nil {
			database["pdb_name"] = *r.PdbName
		}

		database["state"] = r.LifecycleState

		if r.TimeCreated != nil {
			database["time_created"] = r.TimeCreated.String()
		}

		resources = append(resources, database)
	}

	if f, fOk := s.D.GetOkExists("filter"); fOk {
		resources = ApplyFilters(f.(*schema.Set), resources, DatabaseDatabasesDataSource().Schema["databases"].Elem.(*schema.Resource).Schema)
	}

	if err := s.D.Set("databases", resources); err != nil {
		return err
	}

	return nil
}
