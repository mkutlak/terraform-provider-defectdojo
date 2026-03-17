package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type findingDataSource struct {
	terraformDatasource
}

var _ datasource.DataSource = &findingDataSource{}
var _ datasource.DataSourceWithConfigure = &findingDataSource{}

func NewFindingDataSource() datasource.DataSource {
	return &findingDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: findingDataProvider{},
		},
	}
}

func (d findingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_finding"
}

func (d findingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Finding data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Identifier",
			},
			"title":              schema.StringAttribute{Computed: true, MarkdownDescription: "A short description of the flaw"},
			"severity":           schema.StringAttribute{Computed: true, MarkdownDescription: "Severity level"},
			"description":        schema.StringAttribute{Computed: true, MarkdownDescription: "Longer description"},
			"numerical_severity": schema.StringAttribute{Computed: true, MarkdownDescription: "Numerical severity"},
			"test":               schema.Int64Attribute{Computed: true, MarkdownDescription: "ID of the Test"},
			"found_by": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "IDs of test types that found this finding",
				ElementType:         types.Int64Type,
			},
			"active":                     schema.BoolAttribute{Computed: true, MarkdownDescription: "Is the flaw active"},
			"verified":                   schema.BoolAttribute{Computed: true, MarkdownDescription: "Has flaw been verified"},
			"duplicate":                  schema.BoolAttribute{Computed: true, MarkdownDescription: "Is a duplicate"},
			"false_p":                    schema.BoolAttribute{Computed: true, MarkdownDescription: "Is false positive"},
			"out_of_scope":               schema.BoolAttribute{Computed: true, MarkdownDescription: "Is out of scope"},
			"risk_accepted":              schema.BoolAttribute{Computed: true, MarkdownDescription: "Risk accepted"},
			"is_mitigated":               schema.BoolAttribute{Computed: true, MarkdownDescription: "Is mitigated"},
			"dynamic_finding":            schema.BoolAttribute{Computed: true, MarkdownDescription: "DAST finding"},
			"static_finding":             schema.BoolAttribute{Computed: true, MarkdownDescription: "SAST finding"},
			"under_review":               schema.BoolAttribute{Computed: true, MarkdownDescription: "Under review"},
			"under_defect_review":        schema.BoolAttribute{Computed: true, MarkdownDescription: "Under defect review"},
			"push_to_jira":               schema.BoolAttribute{Computed: true, MarkdownDescription: "Push to Jira"},
			"known_exploited":            schema.BoolAttribute{Computed: true, MarkdownDescription: "Known exploited"},
			"ransomware_used":            schema.BoolAttribute{Computed: true, MarkdownDescription: "Ransomware used"},
			"fix_available":              schema.BoolAttribute{Computed: true, MarkdownDescription: "Fix available"},
			"date":                       schema.StringAttribute{Computed: true, MarkdownDescription: "Discovery date"},
			"mitigated":                  schema.StringAttribute{Computed: true, MarkdownDescription: "Mitigation datetime"},
			"last_reviewed":              schema.StringAttribute{Computed: true, MarkdownDescription: "Last reviewed datetime"},
			"planned_remediation_date":   schema.StringAttribute{Computed: true, MarkdownDescription: "Planned remediation date"},
			"publish_date":               schema.StringAttribute{Computed: true, MarkdownDescription: "Publish date"},
			"kev_date":                   schema.StringAttribute{Computed: true, MarkdownDescription: "KEV catalog date"},
			"cvssv3":                     schema.StringAttribute{Computed: true, MarkdownDescription: "CVSSv3 vector"},
			"cvssv3_score":               schema.Float64Attribute{Computed: true, MarkdownDescription: "CVSSv3 score"},
			"cvssv4":                     schema.StringAttribute{Computed: true, MarkdownDescription: "CVSSv4 vector"},
			"cvssv4_score":               schema.Float64Attribute{Computed: true, MarkdownDescription: "CVSSv4 score"},
			"epss_score":                 schema.Float64Attribute{Computed: true, MarkdownDescription: "EPSS score"},
			"epss_percentile":            schema.Float64Attribute{Computed: true, MarkdownDescription: "EPSS percentile"},
			"cwe":                        schema.Int64Attribute{Computed: true, MarkdownDescription: "CWE number"},
			"reporter":                   schema.Int64Attribute{Computed: true, MarkdownDescription: "Reporter ID"},
			"mitigated_by":               schema.Int64Attribute{Computed: true, MarkdownDescription: "Mitigated by user ID"},
			"review_requested_by":        schema.Int64Attribute{Computed: true, MarkdownDescription: "Review requested by user ID"},
			"defect_review_requested_by": schema.Int64Attribute{Computed: true, MarkdownDescription: "Defect review requested by user ID"},
			"sonarqube_issue":            schema.Int64Attribute{Computed: true, MarkdownDescription: "SonarQube issue ID"},
			"line":                       schema.Int64Attribute{Computed: true, MarkdownDescription: "Source line number"},
			"nb_occurences":              schema.Int64Attribute{Computed: true, MarkdownDescription: "Number of occurrences"},
			"sast_source_line":           schema.Int64Attribute{Computed: true, MarkdownDescription: "SAST source line"},
			"scanner_confidence":         schema.Int64Attribute{Computed: true, MarkdownDescription: "Scanner confidence"},
			"file_path":                  schema.StringAttribute{Computed: true, MarkdownDescription: "File path"},
			"component_name":             schema.StringAttribute{Computed: true, MarkdownDescription: "Component name"},
			"component_version":          schema.StringAttribute{Computed: true, MarkdownDescription: "Component version"},
			"fix_version":                schema.StringAttribute{Computed: true, MarkdownDescription: "Fix version"},
			"planned_remediation_version": schema.StringAttribute{Computed: true, MarkdownDescription: "Planned remediation version"},
			"unique_id_from_tool":        schema.StringAttribute{Computed: true, MarkdownDescription: "Unique ID from tool"},
			"vuln_id_from_tool":          schema.StringAttribute{Computed: true, MarkdownDescription: "Vulnerability ID from tool"},
			"service":                    schema.StringAttribute{Computed: true, MarkdownDescription: "Service name"},
			"param":                      schema.StringAttribute{Computed: true, MarkdownDescription: "Parameter"},
			"payload":                    schema.StringAttribute{Computed: true, MarkdownDescription: "Payload"},
			"impact":                     schema.StringAttribute{Computed: true, MarkdownDescription: "Impact description"},
			"mitigation":                 schema.StringAttribute{Computed: true, MarkdownDescription: "Mitigation description"},
			"references":                 schema.StringAttribute{Computed: true, MarkdownDescription: "External references"},
			"steps_to_reproduce":         schema.StringAttribute{Computed: true, MarkdownDescription: "Steps to reproduce"},
			"severity_justification":     schema.StringAttribute{Computed: true, MarkdownDescription: "Severity justification"},
			"effort_for_fixing":          schema.StringAttribute{Computed: true, MarkdownDescription: "Effort for fixing"},
			"hash_code":                  schema.StringAttribute{Computed: true, MarkdownDescription: "Hash code"},
			"sast_source_object":         schema.StringAttribute{Computed: true, MarkdownDescription: "SAST source object"},
			"sast_sink_object":           schema.StringAttribute{Computed: true, MarkdownDescription: "SAST sink object"},
			"sast_source_file_path":      schema.StringAttribute{Computed: true, MarkdownDescription: "SAST source file path"},
			"tags": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "Tags",
				ElementType:         types.StringType,
			},
			"endpoints": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "Endpoint IDs",
				ElementType:         types.Int64Type,
			},
			"reviewers": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "Reviewer user IDs",
				ElementType:         types.Int64Type,
			},
		},
	}
}
