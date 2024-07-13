package types

import frequency "team00/generated"

type MessageGenerator interface {
	Generate() (*frequency.Message, error)
}
