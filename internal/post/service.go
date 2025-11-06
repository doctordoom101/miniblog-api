package post

import (
	"errors"
	"go-miniblog/config"

	"gorm.io/gorm"
)

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func CreatePost(userID uint, req CreatePostRequest) (*Post, error) {
	post := Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func GetAllPosts() ([]Post, error) {
	var posts []Post
	if err := config.DB.Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostByID(postID uint) (*Post, error) {
	var post Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func GetPostsByUserID(userID uint) ([]Post, error) {
	var posts []Post
	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func UpdatePost(postID, userID uint, req UpdatePostRequest) (*Post, error) {
	post, err := GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	if post.UserID != userID {
		return nil, errors.New("unauthorized to update this post")
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := config.DB.Save(&post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func DeletePost(postID, userID uint) error {
	post, err := GetPostByID(postID)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return errors.New("unauthorized to delete this post")
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		return err
	}

	return nil
}
