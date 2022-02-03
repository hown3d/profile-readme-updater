package github

// collectIssue adds an issue to the map if the id of the issue isn't present already
func (c *Client) collectIssue(key int64, val IssueWithRepository) {
	_, exists := c.collectedEvents.Issues[key]
	if exists {
		return
	}
	c.collectedEvents.Issues[key] = val
}

// collectIssue adds a pr to the map if the id of the pr isn't present already
func (c *Client) collectPullRequest(key int64, val PullRequestWithRepository) {
	_, exists := c.collectedEvents.PullRequests[key]
	if exists {
		return
	}
	c.collectedEvents.PullRequests[key] = val
}
