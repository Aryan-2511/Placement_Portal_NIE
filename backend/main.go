package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/controllers"
	"Github.com/Aryan-2511/Placement_NIE/db"
)

// func main() {
//     password := "Admin123@"
//     hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//     fmt.Println("Hash:", string(hash))
//     err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//     if err != nil {
//         fmt.Println("Password does not match!")
//     } else {
//         fmt.Println("Password matches!")
//     }
// }
func withDatabase(db *sql.DB, handler func(http.ResponseWriter, *http.Request, *sql.DB)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, db)
	}
}

func main() {
	db := db.InitDB()
	defer db.Close()
	http.HandleFunc("/signup", withDatabase(db,controllers.SignupHandler))												// Route for student signup
	http.HandleFunc("/login", withDatabase(db,controllers.LoginHandler))												// Route for login
	http.HandleFunc("/placed-student/add",withDatabase(db,controllers.AddPlacedStudent))									// Route for adding placed student
	http.HandleFunc("/placed-student/edit",withDatabase(db,controllers.EditPlacedStudent))								// Route for editing details of placed student
	http.HandleFunc("/placed-student/delete",withDatabase(db,controllers.DeletePlacedStudent))							// Route for deleting a placed student
	http.HandleFunc("/get-placed-students",withDatabase(db,controllers.GetPlacedStudents))								// Route for viewing a list of placed students
	http.HandleFunc("/opportunities/add", withDatabase(db,controllers.AddOpportunity))									// Route for applying to an opportunity
	http.HandleFunc("/opportunities/edit", withDatabase(db,controllers.EditOpportunity))									// Route for editing an opportunity
	http.HandleFunc("/opportunities/delete", withDatabase(db,controllers.DeleteOpportunity))								// Route for deleting an opportunity
	http.HandleFunc("/opportunities/update-status", withDatabase(db,controllers.UpdateOpportunityStatusHandler))			// Route for updating the status of an opportunity
	http.HandleFunc("/opportunities/update-completion", withDatabase(db,controllers.UpdateOpportunityCompletionStatus))	// Route for updating the completion status of an opportunity
	http.HandleFunc("/opportunities/details", withDatabase(db,controllers.GetOpportunityDetailsHandler))					// Route for getting details of an opportunity
	http.HandleFunc("/opportunities/by-batch", withDatabase(db,controllers.GetOpportunitiesByBatchHandler))				// Route for getting list of opportunities for a batch
	http.HandleFunc("/admins/add", withDatabase(db,controllers.AddAdmin))												// Route for adding admin
	http.HandleFunc("/admins/edit", withDatabase(db,controllers.EditAdmin))												// Route for editing details of an admin
	http.HandleFunc("/applications/apply", withDatabase(db,controllers.ApplyHandler))                    				// Route for applying to an opportunity
	http.HandleFunc("/applications/student", withDatabase(db,controllers.GetStudentApplicationsHandler)) 				// Route for fetching student applications
	http.HandleFunc("/forgot-password", withDatabase(db,controllers.ForgotPasswordHandler))  							// Route for requesting a password reset link
	http.HandleFunc("/reset-password", withDatabase(db,controllers.ResetPasswordHandler))   								// Route for resetting the password
	http.HandleFunc("/verify-email", withDatabase(db,controllers.VerifyEmailHandler)) 									// Route for verifying email
	http.HandleFunc("/export-student-details",withDatabase(db,controllers.ExportCustomStudentDetailsToCSV)) 			// Route for exporting details of the applicants
	http.HandleFunc("/placement-coordinator/add", withDatabase(db,controllers.AddPlacementCoordinator)) 					// Route for adding a placement coordinator
	http.HandleFunc("/placement-coordinator/edit", withDatabase(db,controllers.EditPlacementCoordinator)) 		// Route for editing details of a placement coordinator
	http.HandleFunc("/placement-coordinator/delete", withDatabase(db,controllers.DeletePlacementCoordinator)) 			// Route for deleting a placement coordinator
	http.HandleFunc("/get-placement-coordinators", withDatabase(db,controllers.GetAllPlacementCoordinators)) 			// Route for viewing a list of placement coordinators

	http.Handle("/protected", controllers.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the protected route!")
	})))

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
