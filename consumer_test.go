package kafka

import (
	"errors"
	"testing"
)

func Test_consumer_process(t *testing.T) {
	t.Run("When_Processing_Is_Successful", func(t *testing.T) {
		// Given
		c := consumer{
			base: &base{metric: &ConsumerMetric{}},
			consumeFn: func(Message) error {
				return nil
			},
		}

		// When
		c.process(Message{})

		// Then
		if c.metric.TotalProcessedMessagesCounter != 1 {
			t.Fatalf("Total Processed Message Counter must equal to 3")
		}
		if c.metric.TotalUnprocessedMessagesCounter != 0 {
			t.Fatalf("Total Unprocessed Message Counter must equal to 0")
		}
	})
	t.Run("When_Re-processing_Is_Successful", func(t *testing.T) {
		// Given
		gotOnlyOneTimeException := true
		c := consumer{
			base: &base{metric: &ConsumerMetric{}, logger: NewZapLogger(LogLevelDebug)},
			consumeFn: func(Message) error {
				if gotOnlyOneTimeException {
					gotOnlyOneTimeException = false
					return errors.New("simulate only one time exception")
				}
				return nil
			},
		}

		// When
		c.process(Message{})

		// Then
		if c.metric.TotalProcessedMessagesCounter != 1 {
			t.Fatalf("Total Processed Message Counter must equal to 3")
		}
		if c.metric.TotalUnprocessedMessagesCounter != 0 {
			t.Fatalf("Total Unprocessed Message Counter must equal to 0")
		}
	})
	t.Run("When_Re-processing_Is_Failed_And_Retry_Disabled", func(t *testing.T) {
		// Given
		c := consumer{
			base: &base{metric: &ConsumerMetric{}, logger: NewZapLogger(LogLevelDebug)},
			consumeFn: func(Message) error {
				return errors.New("error case")
			},
		}

		// When
		c.process(Message{})

		// Then
		if c.metric.TotalProcessedMessagesCounter != 0 {
			t.Fatalf("Total Processed Message Counter must equal to 0")
		}
		if c.metric.TotalUnprocessedMessagesCounter != 1 {
			t.Fatalf("Total Unprocessed Message Counter must equal to 1")
		}
	})
	t.Run("When_Re-processing_Is_Failed_And_Retry_Enabled", func(t *testing.T) {
		// Given
		mc := mockCronsumer{}
		c := consumer{
			base: &base{metric: &ConsumerMetric{}, logger: NewZapLogger(LogLevelDebug), retryEnabled: true, cronsumer: &mc},
			consumeFn: func(Message) error {
				return errors.New("error case")
			},
		}

		// When
		c.process(Message{})

		// Then
		if c.metric.TotalProcessedMessagesCounter != 0 {
			t.Fatalf("Total Processed Message Counter must equal to 0")
		}
		if c.metric.TotalUnprocessedMessagesCounter != 1 {
			t.Fatalf("Total Unprocessed Message Counter must equal to 1")
		}
	})
	t.Run("When_Re-processing_Is_Failed_And_Retry_Failed", func(t *testing.T) {
		// Given
		mc := mockCronsumer{wantErr: true}
		c := consumer{
			base: &base{metric: &ConsumerMetric{}, logger: NewZapLogger(LogLevelDebug), retryEnabled: true, cronsumer: &mc},
			consumeFn: func(Message) error {
				return errors.New("error case")
			},
		}

		// When
		c.process(Message{})

		// Then
		if c.metric.TotalProcessedMessagesCounter != 0 {
			t.Fatalf("Total Processed Message Counter must equal to 0")
		}
		if c.metric.TotalUnprocessedMessagesCounter != 1 {
			t.Fatalf("Total Unprocessed Message Counter must equal to 1")
		}
	})
}
