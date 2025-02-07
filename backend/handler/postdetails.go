package handler

import (
	"fmt"
	"log"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

func PostDetails(posts []models.Post) ([]models.Post, error) {
	for i, post := range posts {
		comments, err1 := repositories.GetComments(util.DB, post.ID)
		if err1 != nil {
			log.Println("Failed to get comments:", err1)
			return nil, fmt.Errorf("failed to get comments: %v", err1)
		}

		// Getting comment reactions
		for j, comment := range comments {
			commentLikes, errLikes := repositories.GetReactions(util.DB, comment.ID, "Like")
			if errLikes != nil {
				log.Println("Failed to get likes", errLikes)
				return nil, fmt.Errorf("failed to get likes: %v", errLikes)
			}

			commentDislikes, errDislikes := repositories.GetReactions(util.DB, comment.ID, "Dislike")
			if errDislikes != nil {
				log.Println("Failed to get dislikes", errDislikes)
				return nil, fmt.Errorf("failed to get dislikes: %v", errDislikes)
			}

			comments[j].Likes = len(commentLikes)
			comments[j].Dislikes = len(commentDislikes)
		}

		categories, err3 := repositories.GetCategories(util.DB, post.ID)
		if err3 != nil {
			log.Println("Failed to get categories", err3)
			return nil, fmt.Errorf("failed to get categories: %v", err3)
		}
		likes, err4 := repositories.GetReactions(util.DB, post.ID, "Like")
		if err4 != nil {
			log.Println("Failed to get likes", err4)
			return nil, fmt.Errorf("failed to get likes: %v", err4)
		}
		dislikes, err := repositories.GetReactions(util.DB, post.ID, "Dislike")
		if err != nil {
			log.Printf("Failed to get dislikes: %v", err)
			return nil, fmt.Errorf("failed to get dislikes: %v", err)
		}

		posts[i].Comments = comments
		posts[i].CommentCount = len(comments)
		posts[i].Categories = categories
		posts[i].Likes = len(likes)
		posts[i].Dislikes = len(dislikes)
	}
	return posts, nil
}
