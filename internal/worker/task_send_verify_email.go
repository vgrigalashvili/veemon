package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/vgrigalashvili/veemon/internal/helper"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type (
	// Define parameters for creating a verify email.
	CreateVerifyEmailParams struct {
		Email      string `json:"email"`
		SecretCode string `json:"secret_code"`
	}

	// Define the payload structure for sending verification emails.
	PayloadSendVerifyEmail struct {
		Email string `json:"email"`
	}
)

// DistributeTaskSendVerifyEmail enqueues a task to send a verification email.
func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	// Marshal the payload to JSON.
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	// Create a new Asynq task.
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	// Log the successful enqueue of the task.
	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")
	return nil
}

// ProcessTaskSendVerifyEmail processes a task to send a verification email.
func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	// Unmarshal the payload from the task.
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		// Skip retrying if the payload is invalid.
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// Create a verification email entry in the database.
	verifyEmail, err := processor.createVerifyEmail(ctx, CreateVerifyEmailParams{
		Email:      payload.Email,
		SecretCode: helper.RandomString(32), // Generate a random secret code.
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	// Prepare email content.
	subject := "Welcome to Veemon"
	// TODO: Use an environment variable for the frontend URL.
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>`, verifyUrl)
	to := payload.Email

	// Send the verification email.
	err = processor.mailer.SendEmail(ctx, []string{to}, subject, content)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	// Log the successful processing of the task.
	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", payload.Email).
		Msg("processed task")
	return nil
}

// createVerifyEmail is a placeholder for the actual implementation of creating a verification email entry.
func (processor *RedisTaskProcessor) createVerifyEmail(ctx context.Context, params CreateVerifyEmailParams) (*VerifyEmail, error) {
	// Implementation should interact with the database to create the email verification record.
	// Replace the below code with your database logic.
	return &VerifyEmail{
		ID:         1,                 // Example ID
		SecretCode: params.SecretCode, // Generated secret code
	}, nil
}

// VerifyEmail represents the verification email record (simplified example).
type VerifyEmail struct {
	ID         int    `json:"id"`
	SecretCode string `json:"secret_code"`
}
