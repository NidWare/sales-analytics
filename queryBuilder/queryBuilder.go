package queryBuilder

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const baseURL = "https://app.asana.com/api/1.0"

type AsanaTaskSearchBuilder struct {
	baseURL         string
	workspaceID     string
	fields          []string
	pretty          bool
	resourceSubtype string
	assigneeIDs     []string
	projectIDs      []string
	completedBefore time.Time
	completedAfter  time.Time
	completed       bool
	sortBy          string
	sortAscending   bool
}

func NewAsanaTaskSearchBuilder(workspaceID string) *AsanaTaskSearchBuilder {
	return &AsanaTaskSearchBuilder{
		baseURL:     baseURL,
		workspaceID: workspaceID,
	}
}

func (b *AsanaTaskSearchBuilder) AddField(field string) *AsanaTaskSearchBuilder {
	b.fields = append(b.fields, field)
	return b
}

func (b *AsanaTaskSearchBuilder) SetPretty(pretty bool) *AsanaTaskSearchBuilder {
	b.pretty = pretty
	return b
}

func (b *AsanaTaskSearchBuilder) SetResourceSubtype(subtype string) *AsanaTaskSearchBuilder {
	b.resourceSubtype = subtype
	return b
}

func (b *AsanaTaskSearchBuilder) AddAssigneeID(assigneeID string) *AsanaTaskSearchBuilder {
	b.assigneeIDs = append(b.assigneeIDs, assigneeID)
	return b
}

func (b *AsanaTaskSearchBuilder) AddProjectID(projectID string) *AsanaTaskSearchBuilder {
	b.projectIDs = append(b.projectIDs, projectID)
	return b
}

func (b *AsanaTaskSearchBuilder) SetCompletedBefore(completedBefore time.Time) *AsanaTaskSearchBuilder {
	b.completedBefore = completedBefore
	return b
}

func (b *AsanaTaskSearchBuilder) SetCompletedAfter(completedAfter time.Time) *AsanaTaskSearchBuilder {
	b.completedAfter = completedAfter
	return b
}

func (b *AsanaTaskSearchBuilder) SetCompleted(completed bool) *AsanaTaskSearchBuilder {
	b.completed = completed
	return b
}

func (b *AsanaTaskSearchBuilder) SetSortBy(sortBy string) *AsanaTaskSearchBuilder {
	b.sortBy = sortBy
	return b
}

func (b *AsanaTaskSearchBuilder) SetSortAscending(sortAscending bool) *AsanaTaskSearchBuilder {
	b.sortAscending = sortAscending
	return b
}

func (b *AsanaTaskSearchBuilder) Build() string {
	u, _ := url.Parse(fmt.Sprintf("%s/workspaces/%s/tasks/search", b.baseURL, b.workspaceID))
	q := u.Query()

	if len(b.fields) > 0 {
		q.Set("opt_fields", strings.Join(b.fields, ","))
	}

	q.Set("opt_pretty", strconv.FormatBool(b.pretty))
	q.Set("resource_subtype", b.resourceSubtype)

	if len(b.assigneeIDs) > 0 {
		q.Set("assignee.any", strings.Join(b.assigneeIDs, ","))
	}

	if len(b.projectIDs) > 0 {
		q.Set("projects.any", strings.Join(b.projectIDs, ","))
	}

	if !b.completedBefore.IsZero() {
		q.Set("completed_since", b.completedBefore.Format(time.RFC3339))
	}

	q.Set("sort_by", b.sortBy)
	q.Set("sort_ascending", strconv.FormatBool(b.sortAscending))

	u.RawQuery = q.Encode()
	return u.String()
}
