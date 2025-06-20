package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"social-nework/pkg/models"

	"github.com/gorilla/mux"
)

type FollowHandler struct {
	FollowModel *models.FollowModel
}

func (h *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	followedID := vars["userID"]
	if followedID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.FollowModel.Follow(ctx, userID, followedID)
	if err != nil {
		if err.Error() == "cannot follow yourself" || err.Error() == "follow already exists" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Failed to follow user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Followed successfully"})
}

func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	followedID := vars["userID"]
	if followedID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.FollowModel.Unfollow(ctx, userID, followedID)
	if err != nil {
		if err.Error() == "follow does not exist" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Failed to unfollow user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Unfollowed successfully"})
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	fmt.Println(userID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	fmt.Println(ctx)
	defer cancel()

	followers, err := h.FollowModel.GetFollowers(ctx, userID)
	fmt.Println(followers)
	if err != nil {
		log.Printf("Failed to get followers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(followers)
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	following, err := h.FollowModel.GetFollowing(ctx, userID)
	fmt.Println(following)
	if err != nil {
		log.Printf("Failed to get following: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(following)
}

func (h *FollowHandler) AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUserID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	requesterID := vars["requesterID"]
	if requesterID == "" {
		http.Error(w, "Requester ID required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: Replace with actual model call once implemented
	// err := h.FollowModel.AcceptFollowRequest(ctx, currentUserID, requesterID)
	// if err != nil {
	// 	// Handle specific errors like "request not found", "not authorized"
	// 	log.Printf("Failed to accept follow request: %v", err)
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	_ = ctx // Placeholder to avoid unused variable error until model call is added

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Follow request accepted successfully"})
}

func (h *FollowHandler) DeclineFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUserID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	requesterID := vars["requesterID"]
	if requesterID == "" {
		http.Error(w, "Requester ID required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: Replace with actual model call once implemented
	// err := h.FollowModel.DeclineFollowRequest(ctx, currentUserID, requesterID)
	// if err != nil {
	// 	// Handle specific errors
	// 	log.Printf("Failed to decline follow request: %v", err)
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }
	_ = ctx // Placeholder to avoid unused variable error

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Follow request declined successfully"})
}

func (h *FollowHandler) GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: Replace with actual model call once implemented
	// pendingRequests, err := h.FollowModel.GetPendingFollowRequests(ctx, userID)
	// if err != nil {
	// 	log.Printf("Failed to get pending follow requests: %v", err)
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	_ = ctx // Placeholder to avoid unused variable error
	var pendingRequests []models.User // Placeholder
	pendingRequests = []models.User{} // Initialize to avoid nil response if model call isn't ready

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pendingRequests)
}
