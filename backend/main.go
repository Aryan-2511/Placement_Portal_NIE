package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"Github.com/Aryan-2511/Placement_NIE/utils"
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
func withDatabaseAndCORS(db *sql.DB, handler func(http.ResponseWriter, *http.Request, *sql.DB)) http.Handler {
	return utils.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, db)
	}))
}


func main() {
	db := db.InitDB()
	defer db.Close()
	http.Handle("/signup", withDatabaseAndCORS(db,controllers.SignupHandler))												// Route for student signup
	http.Handle("/login", withDatabaseAndCORS(db,controllers.LoginHandler))												// Route for login
	http.Handle("/placed-student/add",withDatabaseAndCORS(db,controllers.AddPlacedStudent))									// Route for adding placed student
	http.Handle("/placed-student/edit",withDatabaseAndCORS(db,controllers.EditPlacedStudent))								// Route for editing details of placed student
	http.Handle("/placed-student/delete",withDatabaseAndCORS(db,controllers.DeletePlacedStudent))							// Route for deleting a placed student
	http.Handle("/get-placed-students",withDatabaseAndCORS(db,controllers.GetPlacedStudents))								// Route for viewing a list of placed students
	http.Handle("/opportunities/add", withDatabaseAndCORS(db,controllers.AddOpportunity))									// Route for applying to an opportunity
	http.Handle("/opportunities/edit", withDatabaseAndCORS(db,controllers.EditOpportunity))									// Route for editing an opportunity
	http.Handle("/opportunities/delete", withDatabaseAndCORS(db,controllers.DeleteOpportunity))								// Route for deleting an opportunity
	http.Handle("/opportunities/update-status", withDatabaseAndCORS(db,controllers.UpdateOpportunityStatusHandler))			// Route for updating the status of an opportunity
	http.Handle("/opportunities/update-completion", withDatabaseAndCORS(db,controllers.UpdateOpportunityCompletionStatus))	// Route for updating the completion status of an opportunity
	http.Handle("/opportunities/details", withDatabaseAndCORS(db,controllers.GetOpportunityDetailsHandler))					// Route for getting details of an opportunity
	http.Handle("/opportunities/by-batch", withDatabaseAndCORS(db,controllers.GetOpportunitiesByBatchHandler))				// Route for getting list of opportunities for a batch
	http.Handle("/admins/add", withDatabaseAndCORS(db,controllers.AddAdmin))												// Route for adding admin
	http.Handle("/admins/edit", withDatabaseAndCORS(db,controllers.EditAdmin))												// Route for editing details of an admin
	http.Handle("/applications/apply", withDatabaseAndCORS(db,controllers.ApplyHandler))                    				// Route for applying to an opportunity
	http.Handle("/applications/student", withDatabaseAndCORS(db,controllers.GetStudentApplicationsHandler)) 				// Route for fetching student applications
	http.Handle("/forgot-password", withDatabaseAndCORS(db,controllers.ForgotPasswordHandler))  							// Route for requesting a password reset link
	http.Handle("/reset-password", withDatabaseAndCORS(db,controllers.ResetPasswordHandler))   								// Route for resetting the password
	http.Handle("/verify-email", withDatabaseAndCORS(db,controllers.VerifyEmailHandler)) 									// Route for verifying email
	http.Handle("/export-student-details",withDatabaseAndCORS(db,controllers.ExportCustomStudentDetailsToCSV)) 			// Route for exporting details of the applicants
	http.Handle("/placement-coordinator/add", withDatabaseAndCORS(db,controllers.AddPlacementCoordinator)) 					// Route for adding a placement coordinator
	http.Handle("/placement-coordinator/edit", withDatabaseAndCORS(db,controllers.EditPlacementCoordinator)) 		// Route for editing details of a placement coordinator
	http.Handle("/placement-coordinator/delete", withDatabaseAndCORS(db,controllers.DeletePlacementCoordinator)) 			// Route for deleting a placement coordinator
	http.Handle("/get-placement-coordinators", withDatabaseAndCORS(db,controllers.GetAllPlacementCoordinators)) 			// Route for viewing a list of placement coordinators

	// http.Handle("/protected", controllers.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "Welcome to the protected route!")
	// })))

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
