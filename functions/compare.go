package services

import (
	"backend-olimpo/models" // Adjust import path as needed
	"gorm.io/gorm"
)

// CompareStudyPlanWithHistory compares a study plan with academic history
// Returns seen subjects (present in both or equivalent) and pending subjects
func CompareStudyPlanWithHistory(db *gorm.DB, studyPlan models.StudyPlan, academicHistory models.AcademicHistoryInput) ([]models.SubjectResult, []models.SubjectResult, error) {
	var seenSubjects []models.SubjectResult
	var pendingSubjects []models.SubjectResult

	// Load study plan subjects with their equivalences
	err := db.Preload("Subjects").Preload("Subjects.Equivalences").Preload("Subjects.Equivalences.TargetSubject").Find(&studyPlan).Error
	if err != nil {
		return nil, nil, err
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
					return nil, nil, err
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

	return seenSubjects, pendingSubjects, nil
}
