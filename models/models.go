package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email string `gorm:"not null;uniqueIndex" json:"email"`
	Password string `gorm:"not null" json:"-"`
}

type SignUpOrLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRating struct {
	gorm.Model
	UserID uint `gorm:"not null;uniqueIndex"`
	CurrentRating uint `gorm:"not null"`
}

/*
Type1 <= 1200
Type2 1201-1500
Type3 1501-1900
Type4 >= 1901
*/

type ProblemStat struct {
	gorm.Model
	UserID uint `gorm:"not null;uniqueIndex"`
	GraphCount uint `gorm:"default:0"`
	TreeCount uint `gorm:"default:0"`
	GreedyCount uint `gorm:"default:0"`
	DynamicProgramming uint `gorm:"default:0"`
	BinarySearch uint `gorm:"default:0"`
	NumberTheory uint `gorm:"default:0"`
	TotalSolved uint `gorm:"default:0"`
	Type1 uint `gorm:"default:0"`
	Type2 uint `gorm:"default:0"`
	Type3 uint `gorm:"default:0"`
	Type4 uint `gorm:"default:0"`
}

type PostProblemRequest struct {
	ProblemRating uint `json:"problemrating"`
	GraphCount bool `json:"graph"`
	TreeCount bool `json:"tree"`
	GreedyCount bool `json:"greedy"`
	DynamicProgramming bool `json:"dp"`
	BinarySearch bool `json:"binarysearch"`
	NumberTheory bool `json:"numbertheory"`
}

type UpdateRatingRequest struct {
	CurrentRating uint `json:"currentrating"`
}