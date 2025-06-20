package api

import (
	"database/sql"
	"net/http"

	"social-nework/pkg/handlers"
	"social-nework/pkg/middleware"
	"social-nework/pkg/models"

	"github.com/gorilla/mux"
)

// NewRouter creates and configures a new router.
// It takes a *sql.DB connection to initialize handlers.
func NewRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Initialize handlers
	// We need UserHandler for auth middleware and potentially other user-related routes
	userHandler := &handlers.UserHandler{DB: db} // Assuming UserHandler uses DB directly or via a UserModel

	followHandler := &handlers.FollowHandler{
		FollowModel: &models.FollowModel{DB: db},
	}
	// postHandler := &handlers.PostHandler{PostModel: &models.PostModel{DB: db}} // Example for other handlers

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(userHandler) // Pass userHandler or DB as needed by your AuthMiddleware

	// API v1 Subrouter
	// All routes will be versioned under /api/v1
	apiV1Router := router.PathPrefix("/api/v1").Subrouter()

	// --- Authentication Routes ---
	// These usually don't require auth middleware themselves
	apiV1Router.HandleFunc("/register", userHandler.RegisterUser).Methods(http.MethodPost)
	apiV1Router.HandleFunc("/login", userHandler.LoginUser).Methods(http.MethodPost)
	// apiV1Router.HandleFunc("/logout", authMiddleware.RequireAuth(userHandler.LogoutUser)).Methods(http.MethodPost) // Example logout

	// --- Follow Handler Routes ---
	// Routes for managing follow requests (actioned by the authenticated user)
	// {requesterID} is the ID of the user who sent the follow request.
	// The authenticated user (currentUserID) is the one accepting/declining.
	acceptDeclineRouter := apiV1Router.PathPrefix("/users/{requesterID}").Subrouter()
	acceptDeclineRouter.Use(authMiddleware.RequireAuth)
	acceptDeclineRouter.HandleFunc("/accept_follow", followHandler.AcceptFollowRequest).Methods(http.MethodPost)
	acceptDeclineRouter.HandleFunc("/decline_follow", followHandler.DeclineFollowRequest).Methods(http.MethodPost, http.MethodDelete)

	// Route for the authenticated user to get their pending follow requests
	pendingRequestsRouter := apiV1Router.PathPrefix("/follow_requests").Subrouter()
	pendingRequestsRouter.Use(authMiddleware.RequireAuth)
	pendingRequestsRouter.HandleFunc("/pending", followHandler.GetPendingFollowRequests).Methods(http.MethodGet)

	// Routes for an authenticated user to follow/unfollow another user ({userID})
	// And to get follower/following lists for any {userID}
	userActionsRouter := apiV1Router.PathPrefix("/users/{userID}").Subrouter()
	userActionsRouter.Use(authMiddleware.RequireAuth)
	userActionsRouter.HandleFunc("/follow", followHandler.Follow).Methods(http.MethodPost)
	userActionsRouter.HandleFunc("/unfollow", followHandler.Unfollow).Methods(http.MethodDelete) // Using DELETE for unfollow
	userActionsRouter.HandleFunc("/followers", followHandler.GetFollowers).Methods(http.MethodGet)
	userActionsRouter.HandleFunc("/following", followHandler.GetFollowing).Methods(http.MethodGet)

	// --- User Profile Routes --- (Example, assuming userHandler has these)
	// profileRouter := apiV1Router.PathPrefix("/profile").Subrouter()
	// profileRouter.Use(authMiddleware.RequireAuth)
	// profileRouter.HandleFunc("", userHandler.GetMyProfile).Methods(http.MethodGet)
	// profileRouter.HandleFunc("/update", userHandler.UpdateMyProfile).Methods(http.MethodPut)

	// --- Post Routes --- (Example, assuming postHandler exists)
	// postRouter := apiV1Router.PathPrefix("/posts").Subrouter()
	// postRouter.Use(authMiddleware.RequireAuth)
	// postRouter.HandleFunc("", postHandler.CreatePost).Methods(http.MethodPost)
	// postRouter.HandleFunc("", postHandler.GetGlobalFeed).Methods(http.MethodGet) // Or GetUserFeed
	// postRouter.HandleFunc("/{postID}", postHandler.GetPost).Methods(http.MethodGet)
	// postRouter.HandleFunc("/{postID}", postHandler.UpdatePost).Methods(http.MethodPut)
	// postRouter.HandleFunc("/{postID}", postHandler.DeletePost).Methods(http.MethodDelete)
	// postRouter.HandleFunc("/{postID}/comments", postHandler.AddComment).Methods(http.MethodPost)
	// postRouter.HandleFunc("/{postID}/comments", postHandler.GetComments).Methods(http.MethodGet)


	// Add more routes for other handlers (posts, comments, groups, etc.)

	return router
}
