package models

import (
	"time"
)

// Career representa una carrera en la universidad
type Career struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	Code        string    `gorm:"size:20;unique;not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	StudyPlans  []StudyPlan `gorm:"foreignKey:CareerID"`
}

// StudyPlan representa un plan de estudio de una carrera
type StudyPlan struct {
	ID          uint      `gorm:"primaryKey"`
	CareerID    uint      `gorm:"not null"`
	Version     string    `gorm:"size:20;not null"` // Ejemplo: "2023-1"
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Subjects    []Subject `gorm:"many2many:study_plan_subjects;"`
	Career      Career    `gorm:"foreignKey:CareerID"`
	// Creditos divididos por tipo de materia
	FundObligatoriaCredits  int       `gorm:"not null;default:0"` // Créditos requeridos de fundamentación obligatoria
	FundOptativaCredits     int       `gorm:"not null;default:0"` // Créditos requeridos de fundamentación optativa
	DisObligatoriaCredits   int       `gorm:"not null;default:0"` // Créditos requeridos de disciplinar obligatoria
	DisOptativaCredits      int       `gorm:"not null;default:0"` // Créditos requeridos de disciplinar optativa
	LibreCredits            int       `gorm:"not null;default:0"` // Créditos requeridos de libre elección
	TotalCredits            int       `gorm:"not null;default:0"` // Total de créditos del plan
}

// Subject representa una materia del plan de estudio
type Subject struct {
	ID          uint      `gorm:"primaryKey"`
	Code        string    `gorm:"size:20;unique;not null"` // Código de la materia
	Name        string    `gorm:"size:100;not null"`
	Credits     int       `gorm:"not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Type        string    `gorm:"size:20;not null;check:type IN ('fund.obligatoria','fund.optativa','dis.obligatoria','dis.optativa','libre')"` // Tipo de materia
	// Relaciones
	Prerequisites []Subject `gorm:"many2many:subject_prerequisites;"`
	Equivalences  []Equivalence `gorm:"foreignKey:SourceSubjectID"`
	StudyPlans    []StudyPlan   `gorm:"many2many:study_plan_subjects;"`
}

// Equivalence representa una equivalencia entre materias de diferentes planes
type Equivalence struct {
	ID              uint      `gorm:"primaryKey"`
	SourceSubjectID uint      `gorm:"not null"` // Materia origen
	TargetSubjectID uint      `gorm:"not null"` // Materia destino
	Type            string    `gorm:"size:20;not null"` // Tipo de equivalencia (total, parcial, etc)
	Notes           string    `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	// Relaciones
	SourceSubject Subject `gorm:"foreignKey:SourceSubjectID"`
	TargetSubject Subject `gorm:"foreignKey:TargetSubjectID"`
}

// AcademicHistoryInput representa la entrada de historia académica para procesar
// Este es un DTO (Data Transfer Object) y no se almacena en la base de datos
type AcademicHistoryInput struct {
	CareerCode    string   `json:"career_code" binding:"required"`
	Subjects      []SubjectInput `json:"subjects" binding:"required"`
}

// SubjectInput representa una materia en la historia académica de entrada
type SubjectInput struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Credits     int     `json:"credits" binding:"required"`
	Grade       float64 `json:"grade" binding:"required"`
	Status      string  `json:"status" binding:"required"` // Aprobada, Reprobada, En curso, etc.
	Semester    string  `json:"semester" binding:"required"` // Semestre en que se cursó
	Type        string  `json:"type" binding:"required"` // Tipo de materia
	
}

// ComparisonResult representa el resultado de la comparación de planes
// Este es un DTO y no se almacena en la base de datos
type ComparisonResult struct {
	EquivalentSubjects []SubjectResult `json:"equivalent_subjects"`
	MissingSubjects    []SubjectResult `json:"missing_subjects"`
	TotalCredits       int             `json:"total_credits"`
	MissingCredits     int             `json:"missing_credits"`
}

// SubjectResult representa una materia en el resultado de la comparación
type SubjectResult struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Credits     int     `json:"credits"`
	Status      string  `json:"status"` // Equivalente, Falta, etc.
	Equivalence *EquivalenceResult `json:"equivalence,omitempty"`
}

// EquivalenceResult representa una equivalencia en el resultado
type EquivalenceResult struct {
	Type  string `json:"type"`
	Notes string `json:"notes"`
} 
