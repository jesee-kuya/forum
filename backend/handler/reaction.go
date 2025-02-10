package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	session := struct {
		UserId int
	}{
		1,
	}
	if r.Method != http.MethodPost {
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		util.ErrorHandler(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	reactionType := r.FormValue("reaction")
	postID, _ := strconv.Atoi(r.FormValue("post_id"))

	fmt.Println("Reaction: ", reactionType)
	fmt.Println("Post ID: ", postID)

	check, reaction := repositories.CheckReactions(util.DB, session.UserId, postID)

	if !check {
		_, err := repositories.InsertRecord(util.DB, "tblReactions", []string{"user_id", "post_id", "reaction"}, session.UserId, postID, reactionType)
		if err != nil {
			log.Println("Failed to insert record:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if reactionType == reaction {
		err := repositories.UpdateReactionStatus(util.DB, session.UserId, postID)
		if err != nil {
			log.Println("Failed to update reaction status:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	} else {
		err := repositories.UpdateReaction(util.DB, reactionType, session.UserId, postID)
		if err != nil {
			log.Println("Failed to update reaction:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	r.Method = http.MethodGet
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
