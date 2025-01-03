package main

import (
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

func main() {
	db.InitDB()

	http.HandleFunc("/signup", controllers.SignupHandler)												// Route for student signup
	http.HandleFunc("/login", controllers.LoginHandler)													// Route for login
	http.HandleFunc("/placed-student/add",controllers.AddPlacedStudent)									// Route for adding placed student
	http.HandleFunc("/placed-student/edit",controllers.EditPlacedStudent)								// Route for editing details of placed student
	http.HandleFunc("/placed-student/delete",controllers.DeletePlacedStudent)							// Route for deleting a placed student
	http.HandleFunc("/get-placed-students",controllers.GetPlacedStudents)								// Route for viewing a list of placed students
	http.HandleFunc("/opportunities/add", controllers.AddOpportunity)									// Route for applying to an opportunity
	http.HandleFunc("/opportunities/edit", controllers.EditOpportunity)									// Route for editing an opportunity
	http.HandleFunc("/opportunities/delete", controllers.DeleteOpportunity)								// Route for deleting an opportunity
	http.HandleFunc("/opportunities/update-status", controllers.UpdateOpportunityStatusHandler)			// Route for updating the status of an opportunity
	http.HandleFunc("/opportunities/update-completion", controllers.UpdateOpportunityCompletionStatus)	// Route for updating the completion status of an opportunity
	http.HandleFunc("/opportunities/details", controllers.GetOpportunityDetailsHandler)					// Route for getting details of an opportunity
	http.HandleFunc("/opportunities/by-batch", controllers.GetOpportunitiesByBatchHandler)				// Route for getting list of opportunities for a batch
	http.HandleFunc("/admins/add", controllers.AddAdmin)												// Route for adding admin
	http.HandleFunc("/admins/edit", controllers.EditAdmin)												// Route for editing details of an admin
	http.HandleFunc("/applications/apply", controllers.ApplyHandler)                    				// Route for applying to an opportunity
	http.HandleFunc("/applications/student", controllers.GetStudentApplicationsHandler) 				// Route for fetching student applications
	http.HandleFunc("/forgot-password", controllers.ForgotPasswordHandler)  							// Route for requesting a password reset link
	http.HandleFunc("/reset-password", controllers.ResetPasswordHandler)   								// Route for resetting the password
	http.HandleFunc("/verify-email", controllers.VerifyEmailHandler) 									// Route for verifying email
	http.HandleFunc("/export-student-details", controllers.ExportCustomStudentDetailsToCSV) 			// Route for exporting details of the applicants
	http.HandleFunc("/placement-coordinator/add", controllers.AddPlacementCoordinator) 					// Route for adding a placement coordinator
	http.HandleFunc("/placement-coordinator/edit", controllers.EditPlacementCoordinator) 		// Route for editing details of a placement coordinator
	http.HandleFunc("/placement-coordinator/delete", controllers.DeletePlacementCoordinator) 			// Route for deleting a placement coordinator
	http.HandleFunc("/get-placement-coordinators", controllers.GetAllPlacementCoordinators) 			// Route for viewing a list of placement coordinators

	http.Handle("/protected", controllers.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the protected route!")
	})))

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
