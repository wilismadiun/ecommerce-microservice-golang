package generator

import "github.com/google/uuid"

type GeneratorUUID struct{}

func (h *GeneratorUUID) Generator() string {
	return uuid.New().String()
}
