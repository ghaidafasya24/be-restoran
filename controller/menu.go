package controller

import (
	"be/config" // Sesuaikan dengan nama package project Anda
	"be/model"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertMenu(c *fiber.Ctx) error {
	// Bind data menu dari request body
	var menu model.Menu
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	// Validasi: Periksa jika MenuName kosong
	if menu.MenuName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Menu name is required and cannot be empty",
		})
	}

	// Validasi: Periksa jika Price kosong atau nol
	if menu.Price == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Price is required and cannot be zero",
		})
	}

	// Validasi: Periksa jika Description kosong
	if menu.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Description is required and cannot be empty",
		})
	}

	// Validasi: Periksa jika Stock kosong atau nol
	if menu.Stock == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Stock is required and cannot be zero",
		})
	}

	// Validasi: Periksa jika MenuCategories kosong
	if menu.MenuCategories == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Menu categories is required and cannot be empty",
		})
	}

	// // Proses upload gambar
	// file, err := c.FormFile("Image")
	// if err != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
	// 		"status":  http.StatusBadRequest,
	// 		"message": "Image file is required: " + err.Error(),
	// 	})
	// }
	// imageURL, err := UploadImageToGitHub(file, menu.MenuName)
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":  http.StatusInternalServerError,
	// 		"message": err.Error(),
	// 	})
	// }

	// menu.Image = imageURL // Tambahkan ID unik dan waktu pembuatan
	menu.ID = primitive.NewObjectID()
	menu.CreatedAt = time.Now()

	// Connect ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil koleksi menu
	menusCollection := config.Ulbimongoconn.Collection("menu")

	// Masukkan data menu ke MongoDB
	insertedID, err := menusCollection.InsertOne(ctx, menu)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert menu",
		})
	}
	// Response sukses
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Product data saved successfully.",
		"inserted_id": insertedID,
		// "image_url":   imageURL,
	})
}

// func UploadImageToGitHub(file *multipart.FileHeader, productName string) (string, error) {
// 	githubToken := os.Getenv("GH_ACCESS_TOKEN")
// 	repoOwner := "ghaidafasya24"
// 	repoName := "images-restoran"
// 	filePath := fmt.Sprintf("menu/%d_%s.jpg", time.Now().Unix(), productName)

// 	fileContent, err := file.Open()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to open image file: %w", err)
// 	}
// 	defer fileContent.Close()

// 	imageData, err := ioutil.ReadAll(fileContent)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read image file: %w", err)
// 	}

// 	encodedImage := base64.StdEncoding.EncodeToString(imageData)
// 	payload := map[string]string{
// 		"message": fmt.Sprintf("Add image for product %s", productName),
// 		"content": encodedImage,
// 	}
// 	payloadBytes, _ := json.Marshal(payload)

// 	req, _ := http.NewRequest("PUT", fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, filePath), bytes.NewReader(payloadBytes))
// 	req.Header.Set("Authorization", "Bearer "+githubToken)
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to upload image to GitHub: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusCreated {
// 		body, _ := ioutil.ReadAll(resp.Body)
// 		return "", fmt.Errorf("GitHub API error: %s", body)
// 	}

// 	var result map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return "", fmt.Errorf("failed to parse GitHub API response: %w", err)
// 	}

// 	content, ok := result["content"].(map[string]interface{})
// 	if !ok || content["download_url"] == nil {
// 		return "", fmt.Errorf("GitHub API response missing download_url")
// 	}

// 	return content["download_url"].(string), nil
// }

// GetAllMenu function untuk mengambil semua menu
func GetAllMenu(c *fiber.Ctx) error {
	// Connect ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil semua menu dari collection menus
	menusCollection := config.Ulbimongoconn.Collection("menu")

	// Mencari semua data menu
	cursor, err := menusCollection.Find(ctx, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No menus found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch menus",
		})
	}
	defer cursor.Close(ctx)

	// Menyimpan hasil menu ke dalam slice
	var menus []model.Menu
	for cursor.Next(ctx) {
		var menu model.Menu
		if err := cursor.Decode(&menu); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode menu",
			})
		}
		menus = append(menus, menu)
	}

	// Cek jika ada error pada cursor
	if err := cursor.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cursor error",
		})
	}

	// Response sukses dengan semua menu
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"menus": menus,
	})
}

// GetMenuByID function untuk mengambil menu berdasarkan ID
func GetMenuByID(c *fiber.Ctx) error {
	// Ambil parameter ID dari URL
	menuID := c.Params("id")

	// Parse ID menjadi ObjectID
	objectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid menu ID",
		})
	}

	// Connect ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil menu berdasarkan ID dari collection menus
	menusCollection := config.Ulbimongoconn.Collection("menu")
	var menu model.Menu
	err = menusCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&menu)
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Menu not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch menu",
		})
	}

	// Response sukses dengan menu yang ditemukan
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"menu": menu,
	})
}

// UpdateMenu function untuk memperbarui menu berdasarkan ID
func UpdateMenu(c *fiber.Ctx) error {
	// Ambil parameter ID dari URL
	menuID := c.Params("id")

	// Parse ID menjadi ObjectID
	objectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid menu ID",
		})
	}

	// Bind data menu baru dari request body
	var menuData model.Menu
	if err := c.BodyParser(&menuData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	// Connect ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil menu collection
	menusCollection := config.Ulbimongoconn.Collection("menu")

	// Update menu berdasarkan ID
	update := bson.M{
		"$set": bson.M{
			"menu_name":       menuData.MenuName,
			"price":           menuData.Price,
			"description":     menuData.Description,
			"stock":           menuData.Stock,
			"menu_categories": menuData.MenuCategories,
			"created_at":      time.Now(),
		},
	}

	// Melakukan update
	result, err := menusCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update menu",
		})
	}

	// Cek apakah menu ditemukan untuk di-update
	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Menu not found",
		})
	}

	// Response sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Menu updated successfully",
	})
}

// DeleteMenu function untuk menghapus menu berdasarkan ID
func DeleteMenu(c *fiber.Ctx) error {
	// Ambil parameter ID dari URL
	menuID := c.Params("id")

	// Parse ID menjadi ObjectID
	objectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid menu ID",
		})
	}

	// Connect ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil menu collection
	menusCollection := config.Ulbimongoconn.Collection("menu")

	// Hapus menu berdasarkan ID
	result, err := menusCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete menu",
		})
	}

	// Cek apakah menu ditemukan untuk dihapus
	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error deleting data for id %s", menuID),
		})
	}

	// Response sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Product data with id %s deleted successfully", menuID),
	})
}
