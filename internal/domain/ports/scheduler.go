package ports

import (
	"context"
	"time"
)

type ScheduleState string

const (
	ScheduleStateEnabled  ScheduleState = "ENABLED"
	ScheduleStateDisabled ScheduleState = "DISABLED"
)

type ScheduleResult struct {
	ScheduleArn  string
	ScheduleName string
	State        ScheduleState
}

type ScheduleInfo struct {
	Arn                string
	Name               string
	Description        string
	State              ScheduleState
	ScheduleExpression string
	TargetArn          string
	CreatedAt          time.Time
	LastModifiedAt     time.Time
	NextExecutionTime  *time.Time
}

type CreateScheduleParams struct {
	Name               string
	Description        string
	ScheduleExpression string
	ScheduledTime      *time.Time  
	TargetArn          string
	RoleArn           string
	Input              interface{} 
	TimeZone          string
	StartDate         *time.Time
	EndDate           *time.Time
}

type RetryPolicy struct {
	MaxRetries        int32
	MaxEventAge       time.Duration
}


type Scheduler interface {
	Get(ctx context.Context, name string) (*ScheduleInfo, error)
	Create(ctx context.Context, params CreateScheduleParams) (*ScheduleResult, error)
	Delete(ctx context.Context, name string) error
}