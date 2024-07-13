package server

import (
	"fmt"
	"log"
	"math/rand"
	frequency "team00/generated"
	"team00/types"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	MessageGenerator types.MessageGenerator

	frequency.UnimplementedFrequencyRandomizerServer
}

func (s *Server) SpawnFrequencies(_ *empty.Empty, stream frequency.FrequencyRandomizer_SpawnFrequenciesServer) error {
	errorsCount := 0

	for {
		message, err := s.MessageGenerator.Generate()
		if err != nil {
			if err := checkErrors(err, &errorsCount); err != nil {
				return err
			}
			continue
		}

		if err := stream.Send(message); err != nil {
			if err := checkErrors(err, &errorsCount); err != nil {
				return err
			}
			continue
		}

		fmt.Println(message.String())
		time.Sleep(SendCooldown)
	}
}

func checkErrors(err error, errorsCount *int) error {
	if err != nil {
		log.Println(err)

		*errorsCount++
		if *errorsCount >= MaxErrorsCount {
			return fmt.Errorf("reached max (%d) retries count for generating message", MaxErrorsCount)
		}
	}

	return nil
}

type messageGenerator struct {
	Mean float64
	Std  float64
}

func NewMessageGenerator() *messageGenerator {
	return &messageGenerator{
		Mean: MinMean + rand.Float64()*(MaxMean-MinMean),
		Std:  MinStd + rand.Float64()*(MaxStd-MinStd),
	}
}

const (
	MinMean float64 = -10
	MaxMean float64 = 10

	MinStd float64 = 0.3
	MaxStd float64 = 1.5

	MaxErrorsCount int = 5

	SendCooldown time.Duration = time.Millisecond * 100
)

func (mg *messageGenerator) Generate() (*frequency.Message, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	return mg.generateWithSessionId(id), nil
}

func (mg *messageGenerator) generateWithSessionId(id uuid.UUID) *frequency.Message {
	return &frequency.Message{
		SessionId:        id.String(),
		Frequency:        rand.NormFloat64()*mg.Std + mg.Mean,
		CurrentTimestamp: timestamppb.Now(),
	}
}
