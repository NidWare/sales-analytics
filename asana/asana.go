package asana

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sales-count/queryBuilder"
	"time"
)

type Task struct {
	Gid          string        `json:"gid"`
	Assignee     Assignee      `json:"assignee"`
	CompletedAt  string        `json:"completed_at"`
	CustomFields []CustomField `json:"custom_fields"`
}

type Assignee struct {
	Gid          string `json:"gid"`
	ResourceType string `json:"resource_type"`
}

type CustomField struct {
	Gid          string  `json:"gid"`
	DisplayValue *string `json:"display_value"`
}

type Response struct {
	Data []Task `json:"data"`
}

func CalculateSumByManagerID(managerID string, startDate, endDate time.Time, projectIDs []string, asanaToken string) (float64, error) {
	url := buildAsanaQueryURL(managerID, startDate, endDate, projectIDs)

	response, err := makeAsanaRequest(url, asanaToken)
	if err != nil {
		return 0, err
	}

	sum, err := calculateSumFromResponse(response)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func buildAsanaQueryURL(managerID string, startDate, endDate time.Time, projectIDs []string) string {
	QB := queryBuilder.NewAsanaTaskSearchBuilder("1206405818803094")
	QB.AddField("assignee").AddField("completed_at").AddField("custom_fields").AddField("custom_fields.display_value")
	QB.SetPretty(true)
	QB.SetResourceSubtype("default_task")
	QB.SetCompletedBefore(endDate)
	QB.SetCompletedAfter(startDate)
	QB.SetCompleted(true)
	QB.SetSortBy("modified_at")
	QB.SetSortAscending(false)
	QB.AddAssigneeID(managerID)

	for _, projectID := range projectIDs {
		QB.AddProjectID(projectID)
	}

	return QB.Build()
}

func makeAsanaRequest(url, asanaToken string) (Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Response{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", asanaToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func calculateSumFromResponse(response Response) (float64, error) {
	var sum float64
	for _, task := range response.Data {
		for _, field := range task.CustomFields {
			if field.DisplayValue != nil {
				var value float64
				_, err := fmt.Sscanf(*field.DisplayValue, "%f", &value)
				if err == nil {
					sum += value
				}
			}
		}
	}
	return sum, nil
}
