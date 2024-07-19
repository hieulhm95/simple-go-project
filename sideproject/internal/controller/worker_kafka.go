package controller

import (
	"context"
)

type KafkaWorker struct {
	ctx context.Context
}

func (w *KafkaWorker) Run() error {
	return nil
}
