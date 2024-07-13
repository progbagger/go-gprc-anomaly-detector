package server

import (
	"fmt"
	"log"
	"math/rand"
	frequency "team00/generated"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	frequency.UnimplementedFrequencyRandomizerServer
}

func (s *Server) SpawnFrequencies(_ *empty.Empty, stream frequency.FrequencyRandomizer_SpawnFrequenciesServer) error {
	errorsCount := 0
	sessionId, err := uuid.NewUUID()
	if err != nil {
		log.Fatalln(err)
	}
	messageGenerator := NewMessageGenerator(sessionId)
	fmt.Printf("generated mean=%f; std=%f for session \"%s\"\n", messageGenerator.Mean, messageGenerator.Std, sessionId.String())

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("session with id \"%s\" is closed", messageGenerator.SessionId.String())
			return nil
		default:
			message, err := messageGenerator.Generate()
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
		}
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
	SessionId uuid.UUID

	Mean float64
	Std  float64
}

func NewMessageGenerator(sessionId uuid.UUID) *messageGenerator {
	return &messageGenerator{
		SessionId: sessionId,
		Mean:      MinMean + rand.Float64()*(MaxMean-MinMean),
		Std:       MinStd + rand.Float64()*(MaxStd-MinStd),
	}
}

const (
	MinMean float64 = -10
	MaxMean float64 = 10

	MinStd float64 = 0.3
	MaxStd float64 = 1.5

	MaxErrorsCount int = 5
)

func (mg *messageGenerator) Generate() (*frequency.Message, error) {
	return mg.generateWithSessionId(mg.SessionId), nil
}

func (mg *messageGenerator) generateWithSessionId(id uuid.UUID) *frequency.Message {
	return &frequency.Message{
		SessionId:        id.String(),
		Frequency:        rand.NormFloat64()*mg.Std + mg.Mean,
		CurrentTimestamp: timestamppb.Now(),
	}
}
