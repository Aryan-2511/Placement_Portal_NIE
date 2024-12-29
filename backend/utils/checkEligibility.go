package utils

import (
	"fmt"
	"log"
	"time"
	"strings"
	"Github.com/Aryan-2511/Placement_NIE/db"
	"Github.com/Aryan-2511/Placement_NIE/models"
)
type EligibilityCriterion struct {
	Criterion string `json:"criterion"`
	Passed    bool   `json:"passed"`
}
type EligibilityEvaluation struct {
	IsEligible bool                   `json:"is_eligible"`
	Criteria   []EligibilityCriterion `json:"criteria"`
}
func CheckEligibility(student models.User, opportunity models.Opportunity) []EligibilityCriterion{
	var criteria []EligibilityCriterion


	var isPlaced bool
	query := `SELECT COUNT(*) FROM placed_students WHERE usn = $1`
	err := db.DB.QueryRow(query, student.USN).Scan(&isPlaced)
	if err != nil {
		fmt.Printf("Error checking placement status: %v", err)
		criteria = append(criteria, EligibilityCriterion{
			Criterion: "Placement Status",
			Passed:    false,
		})
	} else if isPlaced {
		criteria = append(criteria, EligibilityCriterion{
			Criterion: "Placement Status",
			Passed:    false,
		})
	} else {
		criteria = append(criteria, EligibilityCriterion{
			Criterion: "Placement Status",
			Passed:    true,
		})
	}
	//CGPA 
	if student.CurrentCGPA >= opportunity.CGPA {
		criteria = append(criteria, EligibilityCriterion{Criterion: "CGPA Criteria", Passed: true})
	} else {
		criteria = append(criteria, EligibilityCriterion{
			Criterion: "CGPA Criteria",
			Passed:    false,
		})
	}
	//10th_percentage
	if float64(student.Class_10_Percentage) >= opportunity.Class_10_Percentage_Criteria {
		criteria = append(criteria, EligibilityCriterion{Criterion: "10th percentage Criteria", Passed: true})
	} else {
		criteria = append(criteria, EligibilityCriterion{
			Criterion: "10th percentage Criteria",
			Passed:    false,
		})
	}

	//12th_percentage
	if float64(student.Class_12_Percentage) >= opportunity.Class_12_Percentage_Criteria {
		criteria = append(criteria, EligibilityCriterion{Criterion: "12th percentage Criteria", Passed: true})
	} else {
		criteria = append(criteria, EligibilityCriterion{
			Criterion: "12th percentage Criteria",
			Passed:    false,
		})
	}
	//branch
	branchAllowed := false
	for _,allowedBranch := range opportunity.AllowedBranches{
		if models.Branch(student.Branch) == allowedBranch{
			branchAllowed = true
			break
	}}
	if branchAllowed{
		criteria = append(criteria, EligibilityCriterion{Criterion: "Branch Criteria", Passed: true})
	}else{
		criteria = append(criteria, EligibilityCriterion{
			Criterion: "Branch Criteria",
			Passed:    false,
		})
	}

	//batch
	log.Print(student.Batch)
	log.Print(opportunity.Batch)
	if fmt.Sprintf("%d", student.Batch) == strings.TrimSpace(opportunity.Batch) {
		criteria = append(criteria, EligibilityCriterion{Criterion: "Batch Criteria", Passed: true})
	} else {
		criteria = append(criteria, EligibilityCriterion{
			Criterion: "Batch Criteria",
			Passed:    false,
		})
	}
	if student.Backlogs <= opportunity.Backlog{
		criteria = append(criteria, EligibilityCriterion{Criterion: "Backlog Criteria", Passed: true})
	}else{
		criteria = append(criteria, EligibilityCriterion{
			Criterion: "Backlog Criteria",
			Passed: false,
		})
	}
	//gender
	genderAllowed := false
	for _,allowedGender := range opportunity.AllowedGenders{
		if student.Gender == allowedGender{
			genderAllowed = true
			break
		}
	}
	if genderAllowed {
		criteria = append(criteria, EligibilityCriterion{Criterion: "Gender Criteria", Passed: true})
	} else {
		criteria = append(criteria, EligibilityCriterion{
			Criterion: "Gender Criteria",
			Passed:    false,
		})
	}
	//deadline
	if time.Now().Before(opportunity.RegistrationDate) {
		criteria = append(criteria, EligibilityCriterion{Criterion: "Registration Criteria", Passed: true})
	} else {
		criteria = append(criteria, EligibilityCriterion{
			Criterion: "Registration Criteria",
			Passed:    false,
		})
	}
	return criteria
}