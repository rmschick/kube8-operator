package types

// Status is used by dispatchers to tell the runner the status of a batch.
// 'success' and 'failure' denote whether dispatching the batch was successful.
// 'ignored' denotes that received batch was not applicable to that part of the pipeline,
// e.g. an Azure publisher receiving data tagged with the Jira destination will ignore it.
type Status string

const (
	Success Status = "success"
	Ignored Status = "ignored"
	Failure Status = "failure"
)
