package entity

import "errors"

var ErrAnswersCannotBeEmpty = errors.New("Answers param cannot be empty")

var ErrNothingFound = errors.New("Nothing found with this filters")
