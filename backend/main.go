package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"Github.com/Aryan-2511/Placement_NIE/controllers"
	"Github.com/Aryan-2511/Placement_NIE/db"
	"Github.com/Aryan-2511/Placement_NIE/utils"
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
func withDatabaseAndCORSAndKEY(db *sql.DB, secretKey string, handler func(http.ResponseWriter, *http.Request, *sql.DB, string)) http.Handler {
	return utils.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, db, secretKey)
	}))
}



func main() {
	err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
	jwtKey := os.Getenv("jwt_secret_key")
	db := db.InitDB()
	defer db.Close()
	http.Handle("/signup", withDatabaseAndCORS(db,controllers.SignupHandler))												// Route for student signup
	http.Handle("/login", withDatabaseAndCORS(db,controllers.LoginHandler))													// Route for login
	http.Handle("/student/details", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.GetStudentDetailsHandler))													// Route for fetching student details
	http.Handle("/student/edit", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.EditStudentDetailsHandler))													// Route for editing student details
	http.Handle("/placed-student/add",withDatabaseAndCORSAndKEY(db,jwtKey,controllers.AddPlacedStudent))									// Route for adding placed student
	http.Handle("/placed-student/edit",withDatabaseAndCORSAndKEY(db,jwtKey,controllers.EditPlacedStudent))								// Route for editing details of placed student
	http.Handle("/placed-student/delete",withDatabaseAndCORSAndKEY(db,jwtKey,controllers.DeletePlacedStudent))							// Route for deleting a placed student
	http.Handle("/get-placed-students",withDatabaseAndCORSAndKEY(db,jwtKey,controllers.GetPlacedStudents))								// Route for viewing a list of placed students
	http.Handle("/opportunities/add", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.AddOpportunity))									// Route for applying to an opportunity
	http.Handle("/opportunities/edit", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.EditOpportunity))									// Route for editing an opportunity
	http.Handle("/opportunities/delete", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.DeleteOpportunity))								// Route for deleting an opportunity
	http.Handle("/opportunities/update-status", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.UpdateOpportunityStatusHandler))			// Route for updating the status of an opportunity
	http.Handle("/opportunities/update-completion", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.UpdateOpportunityCompletionStatus))	// Route for updating the completion status of an opportunity
	http.Handle("/opportunities/details", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.GetOpportunityDetailsHandler))					// Route for getting details of an opportunity
	http.Handle("/opportunities/by-batch", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.GetOpportunitiesByBatchHandler))				// Route for getting list of opportunities for a batch
	http.Handle("/admins/add", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.AddAdmin))												// Route for adding admin
	http.Handle("/admins/edit", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.EditAdmin))												// Route for editing details of an admin
	http.Handle("/applications/apply", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.ApplyHandler))                    				// Route for applying to an opportunity
	http.Handle("/applications/student", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.GetStudentApplicationsHandler)) 				// Route for fetching student applications
	http.Handle("/forgot-password", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.ForgotPasswordHandler))  							// Route for requesting a password reset link
	http.Handle("/reset-password", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.ResetPasswordHandler))   								// Route for resetting the password
	http.Handle("/verify-email", withDatabaseAndCORS(db,controllers.VerifyEmailHandler)) 									// Route for verifying email
	http.Handle("/export-student-details",withDatabaseAndCORSAndKEY(db,jwtKey,controllers.ExportCustomStudentDetailsToCSV)) 			// Route for exporting details of the applicants
	http.Handle("/placement-coordinator/add", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.AddPlacementCoordinator)) 					// Route for adding a placement coordinator
	http.Handle("/placement-coordinator/edit", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.EditPlacementCoordinator)) 		// Route for editing details of a placement coordinator
	http.Handle("/placement-coordinator/delete", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.DeletePlacementCoordinator)) 			// Route for deleting a placement coordinator
	http.Handle("/get-placement-coordinators", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.GetAllPlacementCoordinators)) 			// Route for viewing a list of placement coordinators
	http.Handle("/schedule/add", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.AddEvent)) 												// Route for adding an event to the schedule
	http.Handle("/schedule/edit", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.EditEvent)) 												// Route for editing an event
	http.Handle("/schedule/delete", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.DeleteEvent)) 												// Route for deleting an event 
	http.Handle("/schedule/all", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.GetAllEvents)) 												// Route for getting all the events
	http.Handle("/schedule/student", withDatabaseAndCORSAndKEY(db,jwtKey,controllers.GetStudentEvents)) 												// Route for getting events of a student


	// http.Handle("/protected", controllers.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "Welcome to the protected route!")
	// })))

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
