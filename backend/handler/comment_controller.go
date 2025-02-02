package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/utils"
)

func ReactToPost(w http.ResponseWriter, r *http.Request) {
	var reaction models.Reaction
	err := json.NewDecoder(r.Body).Decode(&reaction)
	if err != nil {
		utils.ErrorHandler(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reaction.UserID == 0 || reaction.PostID == 0 || reaction.Reaction == "" {
		utils.ErrorHandler(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err = repositories.AddReaction(reaction)
	if err != nil {
		log.Printf("Error adding reaction: %v", err)
		utils.ErrorHandler(w, "Failed to add reaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
