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

	http.HandleFunc("/signup", controllers.SignupHandler)
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/add-placed-student",controllers.MarkStudentAsPlaced)
	http.HandleFunc("/get-placed-student",controllers.GetPlacedStudents)
	http.HandleFunc("/opportunities/add", controllers.AddOpportunity)
	http.HandleFunc("/opportunities/edit", controllers.EditOpportunity)
	http.HandleFunc("/opportunities/delete", controllers.DeleteOpportunity)
	http.HandleFunc("/opportunities/update-status", controllers.UpdateOpportunityStatusHandler)
	http.HandleFunc("/opportunities/update-completion", controllers.UpdateOpportunityCompletionStatus)
	http.HandleFunc("/opportunities/details", controllers.GetOpportunityDetailsHandler)
	http.HandleFunc("/opportunities/by-batch", controllers.GetOpportunitiesByBatchHandler)
	http.HandleFunc("/admins/add", controllers.AddAdmin)
	http.HandleFunc("/admins/edit", controllers.EditAdmin)
	http.HandleFunc("/applications/apply", controllers.ApplyHandler)                    // Route for applying to an opportunity
	http.HandleFunc("/applications/student", controllers.GetStudentApplicationsHandler) // Route for fetching student applications
	http.HandleFunc("/forgot-password", controllers.ForgotPasswordHandler)  // Route for requesting a password reset link
	http.HandleFunc("/reset-password", controllers.ResetPasswordHandler)   // Route for resetting the password
	http.HandleFunc("/verify-email", controllers.VerifyEmailHandler) // Route for verifying email


	http.Handle("/protected", controllers.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the protected route!")
	})))

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
