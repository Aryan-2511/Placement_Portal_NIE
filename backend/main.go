package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/controllers"
	"Github.com/Aryan-2511/Placement_NIE/db"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)

//	func main() {
//	    password := "Admin123@"
//	    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	    fmt.Println("Hash:", string(hash))
//	    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//	    if err != nil {
//	        fmt.Println("Password does not match!")
//	    } else {
//	        fmt.Println("Password matches!")
//	    }
//	}
func withDatabaseAndCORS(db *sql.DB, handler func(http.ResponseWriter, *http.Request, *sql.DB)) http.Handler {
	return utils.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, db)
	}))
}
func withDatabaseAndAuth(db *sql.DB, handler func(http.ResponseWriter, *http.Request, *sql.DB)) http.Handler {
	return withDatabaseAndCORS(db, func(w http.ResponseWriter, r *http.Request, db *sql.DB) {
		controllers.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, db)
		})).ServeHTTP(w, r)
	})
}

func main() {

	db := db.InitDB()
	defer db.Close()

	http.Handle("/student/details", withDatabaseAndAuth(db, controllers.GetStudentDetailsHandler))
	http.Handle("/signup", withDatabaseAndCORS(db, controllers.SignupHandler))                                       // Route for student signup
	http.Handle("/login", withDatabaseAndCORS(db, controllers.LoginHandler))                                         // Route for login
	http.Handle("/logout", withDatabaseAndCORS(db, controllers.LogoutHandler))                                       // Route for logout
	http.Handle("/student/edit", withDatabaseAndCORS(db, controllers.EditStudentDetailsHandler))                     // Route for editing student details
	http.Handle("/placed-student/add", withDatabaseAndAuth(db, controllers.AddPlacedStudent))                        // Route for adding placed student
	http.Handle("/placed-student/edit", withDatabaseAndAuth(db, controllers.EditPlacedStudent))                      // Route for editing details of placed student
	http.Handle("/placed-student/delete", withDatabaseAndAuth(db, controllers.DeletePlacedStudent))                  // Route for deleting a placed student
	http.Handle("/get-placed-students", withDatabaseAndAuth(db, controllers.GetPlacedStudents))                      // Route for viewing a list of placed students
	http.Handle("/opportunities/add", withDatabaseAndAuth(db, controllers.AddOpportunity))                           // Route for applying to an opportunity
	http.Handle("/opportunities/edit", withDatabaseAndAuth(db, controllers.EditOpportunity))                         // Route for editing an opportunity
	http.Handle("/opportunities/delete", withDatabaseAndAuth(db, controllers.DeleteOpportunity))                     // Route for deleting an opportunity
	http.Handle("/opportunities/update-status", withDatabaseAndAuth(db, controllers.UpdateOpportunityStatusHandler)) // Route for updating the status of an opportunity
	// http.Handle("/opportunities/update-completion", withDatabaseAndAuth(db,controllers.UpdateOpportunityCompletionStatus))	// Route for updating the completion status of an opportunity
	http.Handle("/opportunities/details", withDatabaseAndAuth(db, controllers.GetOpportunityDetailsHandler))            // Route for getting details of an opportunity
	http.Handle("/opportunities/by-batch", withDatabaseAndAuth(db, controllers.GetOpportunitiesByBatchHandler))         // Route for getting list of opportunities for a batch
	http.Handle("/admins/add", withDatabaseAndAuth(db, controllers.AddAdmin))                                           // Route for adding admin
	http.Handle("/admins/edit", withDatabaseAndAuth(db, controllers.EditAdmin))                                         // Route for editing details of an admin
	http.Handle("/applications/apply", withDatabaseAndAuth(db, controllers.ApplyHandler))                               // Route for applying to an opportunity
	http.Handle("/applications/student", withDatabaseAndAuth(db, controllers.GetStudentApplicationsHandler))            // Route for fetching student applications
	http.Handle("/forgot-password", withDatabaseAndCORS(db, controllers.ForgotPasswordHandler))                         // Route for requesting a password reset link
	http.Handle("/reset-password", withDatabaseAndCORS(db, controllers.ResetPasswordHandler))                           // Route for resetting the password
	http.Handle("/verify-email", withDatabaseAndCORS(db, controllers.VerifyEmailHandler))                               // Route for verifying email
	http.Handle("/export-student-details", withDatabaseAndAuth(db, controllers.ExportCustomStudentDetailsToCSV))        // Route for exporting details of the applicants
	http.Handle("/placement-coordinator/add", withDatabaseAndAuth(db, controllers.AddPlacementCoordinator))             // Route for adding a placement coordinator
	http.Handle("/placement-coordinator/edit", withDatabaseAndAuth(db, controllers.EditPlacementCoordinator))           // Route for editing details of a placement coordinator
	http.Handle("/placement-coordinator/delete", withDatabaseAndAuth(db, controllers.DeletePlacementCoordinator))       // Route for deleting a placement coordinator
	http.Handle("/get-placement-coordinators", withDatabaseAndAuth(db, controllers.GetAllPlacementCoordinators))        // Route for viewing a list of placement coordinators
	http.Handle("/schedule/add", withDatabaseAndAuth(db, controllers.AddEvent))                                         // Route for adding an event to the schedule
	http.Handle("/schedule/edit", withDatabaseAndAuth(db, controllers.EditEvent))                                       // Route for editing an event
	http.Handle("/schedule/delete", withDatabaseAndAuth(db, controllers.DeleteEvent))                                   // Route for deleting an event
	http.Handle("/schedule/all", withDatabaseAndAuth(db, controllers.GetAllEvents))                                     // Route for getting all the events
	http.Handle("/schedule/student", withDatabaseAndAuth(db, controllers.GetStudentEvents))                             // Route for getting events of a student
	http.Handle("/student-dash/active-opportunities", withDatabaseAndAuth(db, controllers.GetActiveOpportunitiesCount)) // Route for getting events of a student
	http.Handle("/student-dash/recent-opportunities", withDatabaseAndAuth(db, controllers.GetRecentOpportunities))      // Route for getting events of a student
	http.Handle("/student-dash/total-placed-students", withDatabaseAndAuth(db, controllers.GetPlacedStudentsCount))     // Route for getting events of a student
	http.Handle("/student-dash/total-applications", withDatabaseAndAuth(db, controllers.GetTotalApplicationsByStudent)) // Route for getting events of a student
	http.Handle("/students/all-by-batch", withDatabaseAndAuth(db, controllers.FilterByBatch))                           // Route for getting events of a student
	http.Handle("/students/all-by-branch", withDatabaseAndAuth(db, controllers.FilterByBranch))                   // Route for getting events of a student

	// http.Handle("/protected", controllers.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "Welcome to the protected route!")
	// })))

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
