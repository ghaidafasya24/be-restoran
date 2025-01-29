package controller

import (
	"be/config" // Sesuaikan dengan nama package project Anda
	"be/model"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertMenu function untuk menambahkan menu dengan token
func InsertMenu(c *fiber.Ctx) error {
	// Bind data menu dari request body
	var menu model.Menu
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	// Tambahkan ID unik dan waktu pembuatan
	menu.ID = primitive.NewObjectID()
	menu.CreatedAt = time.Now()

	// Connect ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil koleksi menu
	menusCollection := config.Ulbimongoconn.Collection("menu")

	// Masukkan data menu ke MongoDB
	_, err := menusCollection.InsertOne(ctx, menu)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert menu",
		})
	}

	// Response sukses
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Menu inserted successfully",
		"menu":    menu,
	})
}

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
			"message": "Menu not found",
		})
	}

	// Response sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
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
			"message": "Menu not found",
		})
	}

	// Response sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Menu deleted successfully",
	})
}
