package services

import (
	"backend-olimpo/models" // Adjust import path as needed
	"gorm.io/gorm"
)

// ComparisonResult holds the result of comparing study plan with academic history
type ComparisonResult struct {
	SeenSubjects         []models.SubjectResult `json:"seen_subjects"`
	PendingSubjects      []models.SubjectResult `json:"pending_subjects"`
	MissingCreditsByType map[string]int         `json:"missing_credits_by_type"`
}

// CompareStudyPlanWithHistory compares a study plan with academic history
// Returns seen subjects, pending subjects, and missing credits by type
func CompareStudyPlanWithHistory(db *gorm.DB, studyPlan models.StudyPlan, academicHistory models.AcademicHistoryInput) (ComparisonResult, error) {
	var seenSubjects []models.SubjectResult
	var pendingSubjects []models.SubjectResult

	// Load study plan subjects with their equivalences
	err := db.Preload("Subjects").Preload("Subjects.Equivalences").Preload("Subjects.Equivalences.TargetSubject").Find(&studyPlan).Error
	if err != nil {
		return ComparisonResult{}, err
	}

	// Create a map of academic history subjects for quick lookup
	historySubjectsMap := make(map[string]models.SubjectInput)
	for _, subject := range academicHistory.Subjects {
		historySubjectsMap[subject.Code] = subject
	}

	// Process each subject in the study plan
	for _, studyPlanSubject := range studyPlan.Subjects {
		found := false
		var equivalenceInfo *models.EquivalenceResult

		// First, check for direct match
		if _, exists := historySubjectsMap[studyPlanSubject.Code]; exists {
			found = true
		} else {
			// Check for equivalences
			for _, equivalence := range studyPlanSubject.Equivalences {
				if _, exists := historySubjectsMap[equivalence.TargetSubject.Code]; exists {
					found = true
					equivalenceInfo = &models.EquivalenceResult{
						Type:  equivalence.Type,
						Notes: equivalence.Notes,
					}
					break
				}
			}

			// Also check reverse equivalences (where the study plan subject is the target)
			if !found {
				var reverseEquivalences []models.Equivalence
				err := db.Where("target_subject_id = ?", studyPlanSubject.ID).
					Preload("SourceSubject").
					Find(&reverseEquivalences).Error
				if err != nil {
					return ComparisonResult{}, err
				}

				for _, equivalence := range reverseEquivalences {
					if _, exists := historySubjectsMap[equivalence.SourceSubject.Code]; exists {
						found = true
						equivalenceInfo = &models.EquivalenceResult{
							Type:  equivalence.Type,
							Notes: equivalence.Notes,
						}
						break
					}
				}
			}
		}

		// Create subject result
		subjectResult := models.SubjectResult{
			Code:        studyPlanSubject.Code,
			Name:        studyPlanSubject.Name,
			Credits:     studyPlanSubject.Credits,
			Type:        studyPlanSubject.Type,
			Equivalence: equivalenceInfo,
		}

		if found {
			subjectResult.Status = "Seen"
			seenSubjects = append(seenSubjects, subjectResult)
		} else {
			subjectResult.Status = "Pending"
			pendingSubjects = append(pendingSubjects, subjectResult)
		}
	}

	// Calculate credits seen by type
	seenCreditsByType := map[string]int{
		"fund.obligatoria": 0,
		"fund.optativa":    0,
		"dis.obligatoria":  0,
		"dis.optativa":     0,
		"libre":            0,
	}

	for _, subject := range seenSubjects {
		seenCreditsByType[subject.Type] += subject.Credits
	}

	// Calculate missing credits by type (if negative, set to 0)
	missingCreditsByType := map[string]int{
		"fund.obligatoria": max(0, studyPlan.FundObligatoriaCredits-seenCreditsByType["fund.obligatoria"]),
		"fund.optativa":    max(0, studyPlan.FundOptativaCredits-seenCreditsByType["fund.optativa"]),
		"dis.obligatoria":  max(0, studyPlan.DisObligatoriaCredits-seenCreditsByType["dis.obligatoria"]),
		"dis.optativa":     max(0, studyPlan.DisOptativaCredits-seenCreditsByType["dis.optativa"]),
		"libre":            max(0, studyPlan.LibreCredits-seenCreditsByType["libre"]),
	}

	return ComparisonResult{
		SeenSubjects:         seenSubjects,
		PendingSubjects:      pendingSubjects,
		MissingCreditsByType: missingCreditsByType,
	}, nil
}

// Helper function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
