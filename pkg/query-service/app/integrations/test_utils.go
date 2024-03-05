package integrations

import (
	"context"
	"os"
	"slices"
	"testing"

	"github.com/jmoiron/sqlx"
	"go.signoz.io/signoz/pkg/query-service/app/dashboards"
	"go.signoz.io/signoz/pkg/query-service/app/logparsingpipeline"
	"go.signoz.io/signoz/pkg/query-service/model"
	v3 "go.signoz.io/signoz/pkg/query-service/model/v3"
	"go.signoz.io/signoz/pkg/query-service/rules"
)

func NewTestSqliteDB(t *testing.T) (
	db *sqlx.DB, dbFilePath string,
) {
	testDBFile, err := os.CreateTemp("", "test-signoz-db-*")
	if err != nil {
		t.Fatalf("could not create temp file for test db: %v", err)
	}
	testDBFilePath := testDBFile.Name()
	t.Cleanup(func() { os.Remove(testDBFilePath) })
	testDBFile.Close()

	testDB, err := sqlx.Open("sqlite3", testDBFilePath)
	if err != nil {
		t.Fatalf("could not open test db sqlite file: %v", err)
	}

	return testDB, testDBFilePath
}

func NewTestIntegrationsManager(t *testing.T) *Manager {
	testDB, _ := NewTestSqliteDB(t)

	installedIntegrationsRepo, err := NewInstalledIntegrationsSqliteRepo(testDB)
	if err != nil {
		t.Fatalf("could not init sqlite DB for installed integrations: %v", err)
	}

	return &Manager{
		availableIntegrationsRepo: &TestAvailableIntegrationsRepo{},
		installedIntegrationsRepo: installedIntegrationsRepo,
	}
}

type TestAvailableIntegrationsRepo struct{}

func (t *TestAvailableIntegrationsRepo) list(
	ctx context.Context,
) ([]IntegrationDetails, *model.ApiError) {
	return []IntegrationDetails{
		{
			IntegrationSummary: IntegrationSummary{
				Id:          "test-integration-1",
				Title:       "Test Integration 1",
				Description: "A test integration",
				Author: IntegrationAuthor{
					Name:     "signoz",
					Email:    "integrations@signoz.io",
					HomePage: "https://signoz.io",
				},
				Icon: `data:image/svg+xml;utf8,<svg ... > ... </svg>`,
			},
			Categories: []string{"testcat1", "testcat2"},
			Overview:   "test integration overview",
			Configuration: []IntegrationConfigStep{
				{
					Title:        "Step 1",
					Instructions: "Set source attrib on your signals",
				},
			},
			DataCollected: DataCollectedForIntegration{
				Logs:    []CollectedLogAttribute{},
				Metrics: []CollectedMetric{},
			},
			Assets: IntegrationAssets{
				Logs: LogsAssets{
					Pipelines: []logparsingpipeline.PostablePipeline{
						{
							Name:    "pipeline1",
							Alias:   "pipeline1",
							Enabled: true,
							Filter: &v3.FilterSet{
								Operator: "AND",
								Items: []v3.FilterItem{
									{
										Key: v3.AttributeKey{
											Key:      "source",
											DataType: v3.AttributeKeyDataTypeString,
											Type:     v3.AttributeKeyTypeTag,
										},
										Operator: "=",
										Value:    "nginx",
									},
								},
							},
							Config: []logparsingpipeline.PipelineOperator{
								{
									OrderId: 1,
									ID:      "add",
									Type:    "add",
									Field:   "attributes.test",
									Value:   "val",
									Enabled: true,
									Name:    "test add",
								},
							},
						},
					},
				},
				Dashboards: []dashboards.Dashboard{},
				Alerts:     []rules.PostableRule{},
			},
			ConnectionTests: &IntegrationConnectionTests{
				Logs: &v3.FilterSet{
					Operator: "AND",
					Items: []v3.FilterItem{
						{
							Key: v3.AttributeKey{
								Key:      "source",
								DataType: v3.AttributeKeyDataTypeString,
								Type:     v3.AttributeKeyTypeTag,
							},
							Operator: "=",
							Value:    "nginx",
						},
					},
				},
			},
		}, {
			IntegrationSummary: IntegrationSummary{
				Id:          "test-integration-2",
				Title:       "Test Integration 2",
				Description: "Another test integration",
				Author: IntegrationAuthor{
					Name:     "signoz",
					Email:    "integrations@signoz.io",
					HomePage: "https://signoz.io",
				},
				Icon: `data:image/svg+xml;utf8,<svg ... > ... </svg>`,
			},
			Categories: []string{"testcat1", "testcat2"},
			Overview:   "test integration overview",
			Configuration: []IntegrationConfigStep{
				{
					Title:        "Step 1",
					Instructions: "Set source attrib on your signals",
				},
			},
			DataCollected: DataCollectedForIntegration{
				Logs:    []CollectedLogAttribute{},
				Metrics: []CollectedMetric{},
			},
			Assets: IntegrationAssets{
				Logs: LogsAssets{
					Pipelines: []logparsingpipeline.PostablePipeline{
						{
							Name:    "pipeline2",
							Alias:   "pipeline2",
							Enabled: true,
							Filter: &v3.FilterSet{
								Operator: "AND",
								Items: []v3.FilterItem{
									{
										Key: v3.AttributeKey{
											Key:      "source",
											DataType: v3.AttributeKeyDataTypeString,
											Type:     v3.AttributeKeyTypeTag,
										},
										Operator: "=",
										Value:    "redis",
									},
								},
							},
							Config: []logparsingpipeline.PipelineOperator{
								{
									OrderId: 1,
									ID:      "add",
									Type:    "add",
									Field:   "attributes.test",
									Value:   "val",
									Enabled: true,
									Name:    "test add",
								},
							},
						},
					},
				},
				Dashboards: []dashboards.Dashboard{},
				Alerts:     []rules.PostableRule{},
			},
			ConnectionTests: &IntegrationConnectionTests{
				Logs: &v3.FilterSet{
					Operator: "AND",
					Items: []v3.FilterItem{
						{
							Key: v3.AttributeKey{
								Key:      "source",
								DataType: v3.AttributeKeyDataTypeString,
								Type:     v3.AttributeKeyTypeTag,
							},
							Operator: "=",
							Value:    "nginx",
						},
					},
				},
			},
		},
	}, nil
}

func (t *TestAvailableIntegrationsRepo) get(
	ctx context.Context, ids []string,
) (map[string]IntegrationDetails, *model.ApiError) {
	availableIntegrations, apiErr := t.list(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	result := map[string]IntegrationDetails{}

	for _, ai := range availableIntegrations {
		if slices.Contains(ids, ai.Id) {
			result[ai.Id] = ai
		}
	}

	return result, nil
}
