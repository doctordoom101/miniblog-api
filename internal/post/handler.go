package post

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreatePostHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	post, err := CreatePost(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Post created successfully",
		"data":    post,
	})
}

func GetAllPostsHandler(c *fiber.Ctx) error {
	posts, err := GetAllPosts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Posts retrieved successfully",
		"data":    posts,
	})
}

func GetPostByIDHandler(c *fiber.Ctx) error {
	postID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	post, err := GetPostByID(uint(postID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Post retrieved successfully",
		"data":    post,
	})
}

func GetPostsByUserHandler(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("user_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	posts, err := GetPostsByUserID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User posts retrieved successfully",
		"data":    posts,
	})
}

func UpdatePostHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	postID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	var req UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	post, err := UpdatePost(uint(postID), userID, req)
	if err != nil {
		if err.Error() == "unauthorized to update this post" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Post updated successfully",
		"data":    post,
	})
}

func DeletePostHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	postID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	err = DeletePost(uint(postID), userID)
	if err != nil {
		if err.Error() == "unauthorized to delete this post" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}
