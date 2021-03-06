package thehive

import (
	"context"
	"fmt"
	"net/http"
)

const (
	// caseMain is used to list and create cases
	caseMain       = APIRoute + "/case"    // GET, POST
	caseSearch     = caseMain + "/_search" // POST
	caseBulkUpdate = caseMain + "/_bulk"   // PATCH
	caseStats      = caseMain + "/_stats"  // POST

	// caseRoute is used in Sprintf parameter interpolation
	caseRoute = caseMain + "/%s"         // GET, PATCH, DELETE
	caseLinks = caseRoute + "/links"     // GET
	caseMerge = caseRoute + "/_merge/%s" // POST
)

// Entity is a common object representing struct
type Entity struct {
	ID   string `json:"_id"`
	Type string `json:"_type"`
}

// CustomField is a custom field in the case
type CustomField struct {
	String string `json:"string,omitempty"`
	Order  int    `json:"int"`
}

// Case represents TheHive Case
type Case struct {
	Entity
	ArtifactCount        int                    `json:"artifactCount"`
	CaseID               int                    `json:"caseId"`
	CreatedAt            int64                  `json:"createdAt"`
	CreatedBy            string                 `json:"createdBy"`
	CustomFields         map[string]CustomField `json:"customFields"`
	Description          string                 `json:"description"`
	EndDate              int64                  `json:"endDate,omitempty"`
	Flag                 bool                   `json:"flag"`
	ID                   string                 `json:"id"`
	IOCCount             int                    `json:"iocCount"`
	ImpactStatus         string                 `json:"impactStatus,omitempty"`
	Metrics              map[string]interface{} `json:"metrics"`
	Owner                string                 `json:"owner"`
	ResolutionStatus     string                 `json:"resolutionStatus,omitempty"`
	Severity             int                    `json:"severity"`
	SimilarArtifactCount int                    `json:"similarArtifactCount,omitempty"`
	SimilarIocCount      int                    `json:"similarIocCount,omitempty"`
	StartDate            int64                  `json:"startDate"`
	Status               string                 `json:"status"`
	Summary              string                 `json:"summary,omitempty"`
	TLP                  int                    `json:"tlp"`
	Tags                 []string               `json:"tags"`
	Title                string                 `json:"title"`
	UpdatedAt            int64                  `json:"updatedAt"`
	UpdatedBy            string                 `json:"updatedBy"`
}

// SendableCase represents a case to import in TheHive
type SendableCase struct {
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Severity     int                    `json:"severity,omitempty"`
	TLP          int                    `json:"tlp,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	Tasks        []SendableTask         `json:"tasks,omitempty"`
	CustomFields map[string]interface{} `json:"customFields,omitempty"`
}

// CasesService is an interface for managing cases
type CasesService interface {
	Get(context.Context, string) (*Case, *http.Response, error)
	List(context.Context) ([]Case, *http.Response, error)
}

// CasesServiceOp handles cases methods from TheHive API
type CasesServiceOp struct {
	client *Client
}

// Get a case from TheHive
func (c *CasesServiceOp) Get(ctx context.Context, id string) (*Case, *http.Response, error) {
	req, err := c.client.NewRequest("GET", fmt.Sprintf(caseRoute, id), nil)
	if err != nil {
		return nil, nil, err
	}

	hc := &Case{}
	resp, err := c.client.Do(ctx, req, hc)
	if err != nil {
		return nil, resp, err
	}

	return hc, resp, nil
}

// List cases from TheHive with pagination
func (c *CasesServiceOp) List(ctx context.Context) ([]Case, *http.Response, error) {
	var cases []Case
	var resp *http.Response
	start := 0

	for {
		pagedCases := fmt.Sprintf("%s?range=%d-%d", caseMain, start, c.client.PageSize)
		req, err := c.client.NewRequest("GET", pagedCases, nil)
		if err != nil {
			return nil, nil, err
		}

		var hcs []Case
		resp, err := c.client.Do(ctx, req, &hcs)
		if err != nil {
			return nil, resp, err
		}

		if len(hcs) < c.client.PageSize {
			cases = append(cases, hcs...)
			break
		}

		cases = append(cases, hcs...)
		start = start + c.client.PageSize
	}

	return cases, resp, nil
}
