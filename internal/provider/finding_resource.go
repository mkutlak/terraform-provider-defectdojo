package provider

import (
	"context"
	"fmt"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r findingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Finding",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "A short description of the flaw",
				Required:            true,
			},
			"severity": schema.StringAttribute{
				MarkdownDescription: "The severity level of this flaw (Critical, High, Medium, Low, Info)",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Longer more descriptive information about the flaw",
				Required:            true,
			},
			"numerical_severity": schema.StringAttribute{
				MarkdownDescription: "The numerical representation of the severity (S0, S1, S2, S3, S4)",
				Required:            true,
			},
			"test": schema.Int64Attribute{
				MarkdownDescription: "ID of the Test this Finding belongs to",
				Required:            true,
			},
			"found_by": schema.SetAttribute{
				MarkdownDescription: "IDs of test types that found this finding",
				Required:            true,
				ElementType:         types.Int64Type,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Denotes if this flaw is active or not",
				Optional:            true,
			},
			"verified": schema.BoolAttribute{
				MarkdownDescription: "Denotes if this flaw has been manually verified by the tester",
				Optional:            true,
			},
			"duplicate": schema.BoolAttribute{
				MarkdownDescription: "Denotes if this flaw is a duplicate of other flaws reported",
				Optional:            true,
			},
			"false_p": schema.BoolAttribute{
				MarkdownDescription: "Denotes if this flaw has been deemed a false positive by the tester",
				Optional:            true,
			},
			"out_of_scope": schema.BoolAttribute{
				MarkdownDescription: "Denotes if this flaw falls outside the scope of the test and/or engagement",
				Optional:            true,
			},
			"risk_accepted": schema.BoolAttribute{
				MarkdownDescription: "Denotes if this finding has been marked as an accepted risk",
				Optional:            true,
			},
			"is_mitigated": schema.BoolAttribute{
				MarkdownDescription: "Denotes if this flaw has been fixed",
				Optional:            true,
			},
			"dynamic_finding": schema.BoolAttribute{
				MarkdownDescription: "Flaw has been detected from a DAST tool",
				Optional:            true,
			},
			"static_finding": schema.BoolAttribute{
				MarkdownDescription: "Flaw has been detected from a SAST tool",
				Optional:            true,
			},
			"under_review": schema.BoolAttribute{
				MarkdownDescription: "Denotes if this flaw is currently being reviewed",
				Optional:            true,
			},
			"under_defect_review": schema.BoolAttribute{
				MarkdownDescription: "Denotes if this finding is under defect review",
				Optional:            true,
			},
			"push_to_jira": schema.BoolAttribute{
				MarkdownDescription: "Whether to push this finding to Jira",
				Optional:            true,
			},
			"known_exploited": schema.BoolAttribute{
				MarkdownDescription: "Whether this vulnerability is known to have been exploited in the wild",
				Optional:            true,
			},
			"ransomware_used": schema.BoolAttribute{
				MarkdownDescription: "Whether this vulnerability is known to have been leveraged as part of a ransomware campaign",
				Optional:            true,
			},
			"fix_available": schema.BoolAttribute{
				MarkdownDescription: "Denotes if there is a fix available for this flaw",
				Optional:            true,
			},
			"date": schema.StringAttribute{
				MarkdownDescription: "The date the flaw was discovered (format: 2006-01-02)",
				Optional:            true,
			},
			"mitigated": schema.StringAttribute{
				MarkdownDescription: "Datetime when this finding was mitigated (RFC3339 format)",
				Optional:            true,
			},
			"last_reviewed": schema.StringAttribute{
				MarkdownDescription: "Datetime the flaw was last reviewed (RFC3339 format)",
				Optional:            true,
			},
			"planned_remediation_date": schema.StringAttribute{
				MarkdownDescription: "The date the flaw is expected to be remediated (format: 2006-01-02)",
				Optional:            true,
			},
			"publish_date": schema.StringAttribute{
				MarkdownDescription: "Date when this vulnerability was made publicly available (format: 2006-01-02)",
				Optional:            true,
			},
			"kev_date": schema.StringAttribute{
				MarkdownDescription: "The date the vulnerability was added to the KEV catalog (format: 2006-01-02)",
				Optional:            true,
			},
			"cvssv3": schema.StringAttribute{
				MarkdownDescription: "CVSSv3 vector string",
				Optional:            true,
			},
			"cvssv3_score": schema.Float64Attribute{
				MarkdownDescription: "Numerical CVSSv3 score (0-10)",
				Optional:            true,
			},
			"cvssv4": schema.StringAttribute{
				MarkdownDescription: "CVSSv4 vector string",
				Optional:            true,
			},
			"cvssv4_score": schema.Float64Attribute{
				MarkdownDescription: "Numerical CVSSv4 score (0-10)",
				Optional:            true,
			},
			"epss_score": schema.Float64Attribute{
				MarkdownDescription: "EPSS score for the CVE",
				Optional:            true,
			},
			"epss_percentile": schema.Float64Attribute{
				MarkdownDescription: "EPSS percentile for the CVE",
				Optional:            true,
			},
			"cwe": schema.Int64Attribute{
				MarkdownDescription: "The CWE number associated with this flaw",
				Optional:            true,
			},
			"reporter": schema.Int64Attribute{
				MarkdownDescription: "ID of the reporter",
				Optional:            true,
			},
			"mitigated_by": schema.Int64Attribute{
				MarkdownDescription: "ID of the user who mitigated this finding",
				Optional:            true,
			},
			"review_requested_by": schema.Int64Attribute{
				MarkdownDescription: "ID of user who requested a review for this finding",
				Optional:            true,
			},
			"defect_review_requested_by": schema.Int64Attribute{
				MarkdownDescription: "ID of user who requested a defect review",
				Optional:            true,
			},
			"sonarqube_issue": schema.Int64Attribute{
				MarkdownDescription: "The SonarQube issue ID associated with this finding",
				Optional:            true,
			},
			"line": schema.Int64Attribute{
				MarkdownDescription: "Source line number of the attack vector",
				Optional:            true,
			},
			"nb_occurences": schema.Int64Attribute{
				MarkdownDescription: "Number of occurrences in the source tool",
				Optional:            true,
			},
			"sast_source_line": schema.Int64Attribute{
				MarkdownDescription: "Source line number of the attack vector (SAST)",
				Optional:            true,
			},
			"scanner_confidence": schema.Int64Attribute{
				MarkdownDescription: "Confidence level of vulnerability from scanner",
				Optional:            true,
			},
			"file_path": schema.StringAttribute{
				MarkdownDescription: "Identified file(s) containing the flaw",
				Optional:            true,
			},
			"component_name": schema.StringAttribute{
				MarkdownDescription: "Name of the affected component",
				Optional:            true,
			},
			"component_version": schema.StringAttribute{
				MarkdownDescription: "Version of the affected component",
				Optional:            true,
			},
			"fix_version": schema.StringAttribute{
				MarkdownDescription: "Version of the affected component in which the flaw is fixed",
				Optional:            true,
			},
			"planned_remediation_version": schema.StringAttribute{
				MarkdownDescription: "The target version when the vulnerability should be fixed",
				Optional:            true,
			},
			"unique_id_from_tool": schema.StringAttribute{
				MarkdownDescription: "Vulnerability technical id from the source tool",
				Optional:            true,
			},
			"vuln_id_from_tool": schema.StringAttribute{
				MarkdownDescription: "Non-unique technical id from the source tool",
				Optional:            true,
			},
			"service": schema.StringAttribute{
				MarkdownDescription: "A service is a self-contained piece of functionality within a Product",
				Optional:            true,
			},
			"param": schema.StringAttribute{
				MarkdownDescription: "Parameter used to trigger the issue (DAST)",
				Optional:            true,
			},
			"payload": schema.StringAttribute{
				MarkdownDescription: "Payload used to attack the service / application",
				Optional:            true,
			},
			"impact": schema.StringAttribute{
				MarkdownDescription: "Text describing the impact this flaw has on systems",
				Optional:            true,
			},
			"mitigation": schema.StringAttribute{
				MarkdownDescription: "Text describing how to best fix the flaw",
				Optional:            true,
			},
			"references": schema.StringAttribute{
				MarkdownDescription: "The external documentation available for this flaw",
				Optional:            true,
			},
			"steps_to_reproduce": schema.StringAttribute{
				MarkdownDescription: "Text describing the steps to reproduce the flaw",
				Optional:            true,
			},
			"severity_justification": schema.StringAttribute{
				MarkdownDescription: "Text describing why a certain severity was associated with this flaw",
				Optional:            true,
			},
			"effort_for_fixing": schema.StringAttribute{
				MarkdownDescription: "Effort for fixing / remediating the vulnerability (Low, Medium, High)",
				Optional:            true,
			},
			"hash_code": schema.StringAttribute{
				MarkdownDescription: "A hash over a configurable set of fields used for deduplication",
				Optional:            true,
			},
			"sast_source_object": schema.StringAttribute{
				MarkdownDescription: "Source object (variable, function...) of the attack vector",
				Optional:            true,
			},
			"sast_sink_object": schema.StringAttribute{
				MarkdownDescription: "Sink object (variable, function...) of the attack vector",
				Optional:            true,
			},
			"sast_source_file_path": schema.StringAttribute{
				MarkdownDescription: "Source file path of the attack vector",
				Optional:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags for this Finding",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"endpoints": schema.SetAttribute{
				MarkdownDescription: "IDs of endpoints susceptible to this flaw",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"reviewers": schema.SetAttribute{
				MarkdownDescription: "IDs of users who reviewed the flaw",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

// findingResourceData uses ddField tags for all fields except FoundBy ([]int, not *[]int)
// and Test (*int mapped to Int64). Both are handled manually in the API calls.
type findingResourceData struct {
	Id                        types.String  `tfsdk:"id" ddField:"Id"`
	Title                     types.String  `tfsdk:"title" ddField:"Title"`
	Severity                  types.String  `tfsdk:"severity" ddField:"Severity"`
	Description               types.String  `tfsdk:"description" ddField:"Description"`
	NumericalSeverity         types.String  `tfsdk:"numerical_severity" ddField:"NumericalSeverity"`
	// Test and FoundBy are handled manually (Test is *int, FoundBy is []int not *[]int)
	Test                      types.Int64   `tfsdk:"test"`
	FoundBy                   types.Set     `tfsdk:"found_by"`
	Active                    types.Bool    `tfsdk:"active" ddField:"Active"`
	Verified                  types.Bool    `tfsdk:"verified" ddField:"Verified"`
	Duplicate                 types.Bool    `tfsdk:"duplicate" ddField:"Duplicate"`
	FalseP                    types.Bool    `tfsdk:"false_p" ddField:"FalseP"`
	OutOfScope                types.Bool    `tfsdk:"out_of_scope" ddField:"OutOfScope"`
	RiskAccepted              types.Bool    `tfsdk:"risk_accepted" ddField:"RiskAccepted"`
	IsMitigated               types.Bool    `tfsdk:"is_mitigated" ddField:"IsMitigated"`
	DynamicFinding            types.Bool    `tfsdk:"dynamic_finding" ddField:"DynamicFinding"`
	StaticFinding             types.Bool    `tfsdk:"static_finding" ddField:"StaticFinding"`
	UnderReview               types.Bool    `tfsdk:"under_review" ddField:"UnderReview"`
	UnderDefectReview         types.Bool    `tfsdk:"under_defect_review" ddField:"UnderDefectReview"`
	PushToJira                types.Bool    `tfsdk:"push_to_jira" ddField:"PushToJira"`
	KnownExploited            types.Bool    `tfsdk:"known_exploited" ddField:"KnownExploited"`
	RansomwareUsed            types.Bool    `tfsdk:"ransomware_used" ddField:"RansomwareUsed"`
	FixAvailable              types.Bool    `tfsdk:"fix_available" ddField:"FixAvailable"`
	Date                      types.String  `tfsdk:"date" ddField:"Date"`
	Mitigated                 types.String  `tfsdk:"mitigated" ddField:"Mitigated"`
	LastReviewed              types.String  `tfsdk:"last_reviewed" ddField:"LastReviewed"`
	PlannedRemediationDate    types.String  `tfsdk:"planned_remediation_date" ddField:"PlannedRemediationDate"`
	PublishDate               types.String  `tfsdk:"publish_date" ddField:"PublishDate"`
	KevDate                   types.String  `tfsdk:"kev_date" ddField:"KevDate"`
	Cvssv3                    types.String  `tfsdk:"cvssv3" ddField:"Cvssv3"`
	Cvssv3Score               types.Float64 `tfsdk:"cvssv3_score" ddField:"Cvssv3Score"`
	Cvssv4                    types.String  `tfsdk:"cvssv4" ddField:"Cvssv4"`
	Cvssv4Score               types.Float64 `tfsdk:"cvssv4_score" ddField:"Cvssv4Score"`
	EpssScore                 types.Float64 `tfsdk:"epss_score" ddField:"EpssScore"`
	EpssPercentile            types.Float64 `tfsdk:"epss_percentile" ddField:"EpssPercentile"`
	Cwe                       types.Int64   `tfsdk:"cwe" ddField:"Cwe"`
	Reporter                  types.Int64   `tfsdk:"reporter" ddField:"Reporter"`
	MitigatedBy               types.Int64   `tfsdk:"mitigated_by" ddField:"MitigatedBy"`
	ReviewRequestedBy         types.Int64   `tfsdk:"review_requested_by" ddField:"ReviewRequestedBy"`
	DefectReviewRequestedBy   types.Int64   `tfsdk:"defect_review_requested_by" ddField:"DefectReviewRequestedBy"`
	SonarqubeIssue            types.Int64   `tfsdk:"sonarqube_issue" ddField:"SonarqubeIssue"`
	Line                      types.Int64   `tfsdk:"line" ddField:"Line"`
	NbOccurences              types.Int64   `tfsdk:"nb_occurences" ddField:"NbOccurences"`
	SastSourceLine            types.Int64   `tfsdk:"sast_source_line" ddField:"SastSourceLine"`
	ScannerConfidence         types.Int64   `tfsdk:"scanner_confidence" ddField:"ScannerConfidence"`
	FilePath                  types.String  `tfsdk:"file_path" ddField:"FilePath"`
	ComponentName             types.String  `tfsdk:"component_name" ddField:"ComponentName"`
	ComponentVersion          types.String  `tfsdk:"component_version" ddField:"ComponentVersion"`
	FixVersion                types.String  `tfsdk:"fix_version" ddField:"FixVersion"`
	PlannedRemediationVersion types.String  `tfsdk:"planned_remediation_version" ddField:"PlannedRemediationVersion"`
	UniqueIdFromTool          types.String  `tfsdk:"unique_id_from_tool" ddField:"UniqueIdFromTool"`
	VulnIdFromTool            types.String  `tfsdk:"vuln_id_from_tool" ddField:"VulnIdFromTool"`
	Service                   types.String  `tfsdk:"service" ddField:"Service"`
	Param                     types.String  `tfsdk:"param" ddField:"Param"`
	Payload                   types.String  `tfsdk:"payload" ddField:"Payload"`
	Impact                    types.String  `tfsdk:"impact" ddField:"Impact"`
	Mitigation                types.String  `tfsdk:"mitigation" ddField:"Mitigation"`
	References                types.String  `tfsdk:"references" ddField:"References"`
	StepsToReproduce          types.String  `tfsdk:"steps_to_reproduce" ddField:"StepsToReproduce"`
	SeverityJustification     types.String  `tfsdk:"severity_justification" ddField:"SeverityJustification"`
	EffortForFixing           types.String  `tfsdk:"effort_for_fixing" ddField:"EffortForFixing"`
	HashCode                  types.String  `tfsdk:"hash_code" ddField:"HashCode"`
	SastSourceObject          types.String  `tfsdk:"sast_source_object" ddField:"SastSourceObject"`
	SastSinkObject            types.String  `tfsdk:"sast_sink_object" ddField:"SastSinkObject"`
	SastSourceFilePath        types.String  `tfsdk:"sast_source_file_path" ddField:"SastSourceFilePath"`
	Tags                      types.Set     `tfsdk:"tags" ddField:"Tags"`
	Endpoints                 types.Set     `tfsdk:"endpoints" ddField:"Endpoints"`
	Reviewers                 types.Set     `tfsdk:"reviewers" ddField:"Reviewers"`
}

type findingDefectdojoResource struct {
	dd.Finding
}

// findingPopulateManualFields copies the manually-managed Test and FoundBy fields
// from Terraform state into the Finding struct.
func findingPopulateManualFields(ctx context.Context, diags *diag.Diagnostics, data *findingResourceData, f *dd.Finding) {
	// Test: types.Int64 -> *int
	if !data.Test.IsNull() && !data.Test.IsUnknown() {
		v := int(data.Test.ValueInt64())
		f.Test = &v
	}
	// FoundBy: types.Set -> []int
	if !data.FoundBy.IsNull() && !data.FoundBy.IsUnknown() {
		int64s := []int64{}
		d := data.FoundBy.ElementsAs(ctx, &int64s, false)
		diags.Append(d...)
		ints := make([]int, len(int64s))
		for i, v := range int64s {
			ints[i] = int(v)
		}
		f.FoundBy = ints
	}
}

// findingReadManualFields copies Test and FoundBy back from Finding to Terraform state.
func findingReadManualFields(ctx context.Context, diags *diag.Diagnostics, data *findingResourceData, f *dd.Finding) {
	// Test: *int -> types.Int64
	if f.Test != nil {
		data.Test = types.Int64Value(int64(*f.Test))
	} else {
		data.Test = types.Int64Null()
	}
	// FoundBy: []int -> types.Set
	elems := []attr.Value{}
	for _, v := range f.FoundBy {
		elems = append(elems, types.Int64Value(int64(v)))
	}
	set, d := types.SetValue(types.Int64Type, elems)
	diags.Append(d...)
	data.FoundBy = set
}

// findingToCreateRequest converts a Finding to a FindingCreateRequest (used for POST).
// FindingCreateRequest requires Test int, Active bool, Verified bool as non-pointer required fields.
func findingToCreateRequest(f dd.Finding) dd.FindingCreateRequest {
	req := dd.FindingCreateRequest{
		Title:                     f.Title,
		Severity:                  f.Severity,
		Description:               f.Description,
		NumericalSeverity:         f.NumericalSeverity,
		FoundBy:                   f.FoundBy,
		Duplicate:                 f.Duplicate,
		FalseP:                    f.FalseP,
		OutOfScope:                f.OutOfScope,
		RiskAccepted:              f.RiskAccepted,
		IsMitigated:               f.IsMitigated,
		DynamicFinding:            f.DynamicFinding,
		StaticFinding:             f.StaticFinding,
		UnderReview:               f.UnderReview,
		UnderDefectReview:         f.UnderDefectReview,
		PushToJira:                f.PushToJira,
		KnownExploited:            f.KnownExploited,
		RansomwareUsed:            f.RansomwareUsed,
		FixAvailable:              f.FixAvailable,
		Date:                      f.Date,
		Mitigated:                 f.Mitigated,
		PlannedRemediationDate:    f.PlannedRemediationDate,
		PublishDate:               f.PublishDate,
		KevDate:                   f.KevDate,
		Cvssv3:                    f.Cvssv3,
		Cvssv3Score:               f.Cvssv3Score,
		Cvssv4:                    f.Cvssv4,
		Cvssv4Score:               f.Cvssv4Score,
		EpssScore:                 f.EpssScore,
		EpssPercentile:            f.EpssPercentile,
		Cwe:                       f.Cwe,
		Reporter:                  f.Reporter,
		MitigatedBy:               f.MitigatedBy,
		ReviewRequestedBy:         f.ReviewRequestedBy,
		DefectReviewRequestedBy:   f.DefectReviewRequestedBy,
		SonarqubeIssue:            f.SonarqubeIssue,
		Line:                      f.Line,
		NbOccurences:              f.NbOccurences,
		SastSourceLine:            f.SastSourceLine,
		FilePath:                  f.FilePath,
		ComponentName:             f.ComponentName,
		ComponentVersion:          f.ComponentVersion,
		FixVersion:                f.FixVersion,
		PlannedRemediationVersion: f.PlannedRemediationVersion,
		UniqueIdFromTool:          f.UniqueIdFromTool,
		VulnIdFromTool:            f.VulnIdFromTool,
		Service:                   f.Service,
		Impact:                    f.Impact,
		Mitigation:                f.Mitigation,
		References:                f.References,
		StepsToReproduce:          f.StepsToReproduce,
		SeverityJustification:     f.SeverityJustification,
		EffortForFixing:           f.EffortForFixing,
		SastSourceObject:          f.SastSourceObject,
		SastSinkObject:            f.SastSinkObject,
		SastSourceFilePath:        f.SastSourceFilePath,
		Tags:                      f.Tags,
		Reviewers:                 f.Reviewers,
	}
	// Test, Active, Verified are required non-pointer fields in FindingCreateRequest
	if f.Test != nil {
		req.Test = *f.Test
	}
	if f.Active != nil {
		req.Active = *f.Active
	}
	if f.Verified != nil {
		req.Verified = *f.Verified
	}
	return req
}

// findingToRequest converts a Finding to a FindingRequest (used for PUT updates).
// FindingRequest does not have Test, ScannerConfidence, Param, or Payload fields.
func findingToRequest(f dd.Finding) dd.FindingRequest {
	return dd.FindingRequest{
		Title:                     f.Title,
		Severity:                  f.Severity,
		Description:               f.Description,
		NumericalSeverity:         f.NumericalSeverity,
		FoundBy:                   f.FoundBy,
		Active:                    f.Active,
		Verified:                  f.Verified,
		Duplicate:                 f.Duplicate,
		FalseP:                    f.FalseP,
		OutOfScope:                f.OutOfScope,
		RiskAccepted:              f.RiskAccepted,
		IsMitigated:               f.IsMitigated,
		DynamicFinding:            f.DynamicFinding,
		StaticFinding:             f.StaticFinding,
		UnderReview:               f.UnderReview,
		UnderDefectReview:         f.UnderDefectReview,
		PushToJira:                f.PushToJira,
		KnownExploited:            f.KnownExploited,
		RansomwareUsed:            f.RansomwareUsed,
		FixAvailable:              f.FixAvailable,
		Date:                      f.Date,
		Mitigated:                 f.Mitigated,
		PlannedRemediationDate:    f.PlannedRemediationDate,
		PublishDate:               f.PublishDate,
		KevDate:                   f.KevDate,
		Cvssv3:                    f.Cvssv3,
		Cvssv3Score:               f.Cvssv3Score,
		Cvssv4:                    f.Cvssv4,
		Cvssv4Score:               f.Cvssv4Score,
		EpssScore:                 f.EpssScore,
		EpssPercentile:            f.EpssPercentile,
		Cwe:                       f.Cwe,
		Reporter:                  f.Reporter,
		MitigatedBy:               f.MitigatedBy,
		ReviewRequestedBy:         f.ReviewRequestedBy,
		DefectReviewRequestedBy:   f.DefectReviewRequestedBy,
		SonarqubeIssue:            f.SonarqubeIssue,
		Line:                      f.Line,
		NbOccurences:              f.NbOccurences,
		SastSourceLine:            f.SastSourceLine,
		FilePath:                  f.FilePath,
		ComponentName:             f.ComponentName,
		ComponentVersion:          f.ComponentVersion,
		FixVersion:                f.FixVersion,
		PlannedRemediationVersion: f.PlannedRemediationVersion,
		UniqueIdFromTool:          f.UniqueIdFromTool,
		VulnIdFromTool:            f.VulnIdFromTool,
		Service:                   f.Service,
		Impact:                    f.Impact,
		Mitigation:                f.Mitigation,
		References:                f.References,
		StepsToReproduce:          f.StepsToReproduce,
		SeverityJustification:     f.SeverityJustification,
		EffortForFixing:           f.EffortForFixing,
		SastSourceObject:          f.SastSourceObject,
		SastSinkObject:            f.SastSinkObject,
		SastSourceFilePath:        f.SastSourceFilePath,
		Tags:                      f.Tags,
		Reviewers:                 f.Reviewers,
	}
}

func (ddr *findingDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "findingDefectdojoResource createApiCall")
	reqBody := findingToCreateRequest(ddr.Finding)
	apiResp, err := client.FindingsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		// FindingsCreate returns FindingCreate, copy fields back into Finding
		fc := apiResp.JSON201
		ddr.Finding.Id = fc.Id
		ddr.Finding.Title = fc.Title
		ddr.Finding.Severity = fc.Severity
		ddr.Finding.Description = fc.Description
		ddr.Finding.NumericalSeverity = fc.NumericalSeverity
		ddr.Finding.Test = &fc.Test
		ddr.Finding.FoundBy = fc.FoundBy
		ddr.Finding.Active = &fc.Active
		ddr.Finding.Verified = &fc.Verified
		ddr.Finding.Duplicate = fc.Duplicate
		ddr.Finding.FalseP = fc.FalseP
		ddr.Finding.OutOfScope = fc.OutOfScope
		ddr.Finding.RiskAccepted = fc.RiskAccepted
		ddr.Finding.IsMitigated = fc.IsMitigated
		ddr.Finding.DynamicFinding = fc.DynamicFinding
		ddr.Finding.StaticFinding = fc.StaticFinding
		ddr.Finding.UnderReview = fc.UnderReview
		ddr.Finding.UnderDefectReview = fc.UnderDefectReview
		ddr.Finding.PushToJira = fc.PushToJira
		ddr.Finding.KnownExploited = fc.KnownExploited
		ddr.Finding.RansomwareUsed = fc.RansomwareUsed
		ddr.Finding.FixAvailable = fc.FixAvailable
		ddr.Finding.Date = fc.Date
		ddr.Finding.Mitigated = fc.Mitigated
		ddr.Finding.PlannedRemediationDate = fc.PlannedRemediationDate
		ddr.Finding.PublishDate = fc.PublishDate
		ddr.Finding.KevDate = fc.KevDate
		ddr.Finding.Cvssv3 = fc.Cvssv3
		ddr.Finding.Cvssv3Score = fc.Cvssv3Score
		ddr.Finding.Cvssv4 = fc.Cvssv4
		ddr.Finding.Cvssv4Score = fc.Cvssv4Score
		ddr.Finding.EpssScore = fc.EpssScore
		ddr.Finding.EpssPercentile = fc.EpssPercentile
		ddr.Finding.Cwe = fc.Cwe
		ddr.Finding.Reporter = fc.Reporter
		ddr.Finding.MitigatedBy = fc.MitigatedBy
		ddr.Finding.ReviewRequestedBy = fc.ReviewRequestedBy
		ddr.Finding.DefectReviewRequestedBy = fc.DefectReviewRequestedBy
		ddr.Finding.SonarqubeIssue = fc.SonarqubeIssue
		ddr.Finding.Line = fc.Line
		ddr.Finding.NbOccurences = fc.NbOccurences
		ddr.Finding.SastSourceLine = fc.SastSourceLine
		ddr.Finding.ScannerConfidence = fc.ScannerConfidence
		ddr.Finding.FilePath = fc.FilePath
		ddr.Finding.ComponentName = fc.ComponentName
		ddr.Finding.ComponentVersion = fc.ComponentVersion
		ddr.Finding.FixVersion = fc.FixVersion
		ddr.Finding.PlannedRemediationVersion = fc.PlannedRemediationVersion
		ddr.Finding.UniqueIdFromTool = fc.UniqueIdFromTool
		ddr.Finding.VulnIdFromTool = fc.VulnIdFromTool
		ddr.Finding.Service = fc.Service
		ddr.Finding.Param = fc.Param
		ddr.Finding.Payload = fc.Payload
		ddr.Finding.Impact = fc.Impact
		ddr.Finding.Mitigation = fc.Mitigation
		ddr.Finding.References = fc.References
		ddr.Finding.StepsToReproduce = fc.StepsToReproduce
		ddr.Finding.SeverityJustification = fc.SeverityJustification
		ddr.Finding.EffortForFixing = fc.EffortForFixing
		ddr.Finding.SastSourceObject = fc.SastSourceObject
		ddr.Finding.SastSinkObject = fc.SastSinkObject
		ddr.Finding.SastSourceFilePath = fc.SastSourceFilePath
		ddr.Finding.Tags = fc.Tags
		ddr.Finding.Reviewers = fc.Reviewers
		ddr.Finding.HashCode = fc.HashCode
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *findingDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "findingDefectdojoResource readApiCall")
	apiResp, err := client.FindingsRetrieveWithResponse(ctx, idNumber, &dd.FindingsRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Finding = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *findingDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "findingDefectdojoResource updateApiCall")
	reqBody := findingToRequest(ddr.Finding)
	apiResp, err := client.FindingsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Finding = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *findingDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "findingDefectdojoResource deleteApiCall")
	apiResp, err := client.FindingsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, nil
}

// findingTerraformResource wraps terraformResource and overrides Create/Read/Update
// to also handle the manually-managed Test and FoundBy fields.
type findingTerraformResource struct {
	terraformResource
}

type findingResource struct {
	findingTerraformResource
}

var _ resource.Resource = &findingResource{}
var _ resource.ResourceWithImportState = &findingResource{}
var _ resource.ResourceWithConfigure = &findingResource{}

func NewFindingResource() resource.Resource {
	return &findingResource{
		findingTerraformResource: findingTerraformResource{
			terraformResource: terraformResource{
				dataProvider: findingDataProvider{},
			},
		},
	}
}

func (r findingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_finding"
}

func (r findingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data, diags := r.getData(ctx, req.Config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured HTTP Client", "Expected configured HTTP client.")
		return
	}
	findingData := data.(*findingResourceData)
	ddResource := data.defectdojoResource()
	populateDefectdojoResource(ctx, &diags, data, &ddResource)
	findingDD := ddResource.(*findingDefectdojoResource)
	findingPopulateManualFields(ctx, &diags, findingData, &findingDD.Finding)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	statusCode, body, err := ddResource.createApiCall(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Resource", fmt.Sprintf("%s", err))
		return
	}
	if statusCode == 201 {
		populateResourceData(ctx, &diags, &data, ddResource)
		findingReadManualFields(ctx, &diags, findingData, &findingDD.Finding)
	} else {
		resp.Diagnostics.AddError("API Error Creating Resource",
			fmt.Sprintf("Unexpected response code from API: %d\n\nbody:\n\n%s", statusCode, string(body)))
		return
	}
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, findingData)
	resp.Diagnostics.Append(diags...)
}

func (r findingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data, diags := r.getData(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	findingData := data.(*findingResourceData)
	if findingData.Id.IsNull() {
		resp.Diagnostics.AddError("Could not Retrieve Resource", "The Id field was null")
		return
	}
	idNumber := 0
	fmt.Sscanf(findingData.Id.ValueString(), "%d", &idNumber)
	ddResource := data.defectdojoResource()
	statusCode, body, err := ddResource.readApiCall(ctx, r.client, idNumber)
	if err != nil {
		resp.Diagnostics.AddError("Error Retrieving Resource", fmt.Sprintf("%s", err))
		return
	}
	findingDD := ddResource.(*findingDefectdojoResource)
	if statusCode == 200 {
		populateResourceData(ctx, &diags, &data, ddResource)
		findingReadManualFields(ctx, &diags, findingData, &findingDD.Finding)
	} else if statusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	} else {
		resp.Diagnostics.AddError("API Error Retrieving Resource",
			fmt.Sprintf("Unexpected response code from API: %d\n\nbody:\n\n%s", statusCode, string(body)))
		return
	}
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, findingData)
	resp.Diagnostics.Append(diags...)
}

func (r findingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data, diags := r.getData(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured HTTP Client", "Expected configured HTTP client.")
		return
	}
	findingData := data.(*findingResourceData)
	if findingData.Id.IsNull() {
		resp.Diagnostics.AddError("Could not Update Resource", "The Id field was null")
		return
	}
	idNumber := 0
	fmt.Sscanf(findingData.Id.ValueString(), "%d", &idNumber)
	ddResource := data.defectdojoResource()
	populateDefectdojoResource(ctx, &diags, data, &ddResource)
	findingDD := ddResource.(*findingDefectdojoResource)
	findingPopulateManualFields(ctx, &diags, findingData, &findingDD.Finding)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	statusCode, body, err := ddResource.updateApiCall(ctx, r.client, idNumber)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Resource", fmt.Sprintf("%s", err))
		return
	}
	if statusCode == 200 {
		populateResourceData(ctx, &diags, &data, ddResource)
		findingReadManualFields(ctx, &diags, findingData, &findingDD.Finding)
	} else {
		resp.Diagnostics.AddError("API Error Updating Resource",
			fmt.Sprintf("Unexpected response code from API: %d\n\nbody:\n\n%s", statusCode, string(body)))
		return
	}
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, findingData)
	resp.Diagnostics.Append(diags...)
}

func (r findingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.terraformResource.Delete(ctx, req, resp)
}

func (r findingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.terraformResource.ImportState(ctx, req, resp)
}

func (r findingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.terraformResource.Configure(ctx, req, resp)
}

type findingDataProvider struct{}

func (r findingDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data findingResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *findingResourceData) id() types.String {
	return d.Id
}

func (d *findingResourceData) defectdojoResource() defectdojoResource {
	return &findingDefectdojoResource{
		Finding: dd.Finding{},
	}
}
