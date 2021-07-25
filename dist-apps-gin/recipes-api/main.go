package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

func main() {

	recipes = make([]Recipe, 0)

	router := gin.Default()

	router.POST("/recipes", newRecipeHandler)
	router.GET("/recipes", listRecipesHandler)
	router.PUT("/recipes/:id", updateRecipeHandler)
	router.DELETE("/recipes/:id", deleteRecipeHandler)
	router.GET("/recipes/search", searchRecipeHandler)

	log.Print("Starting up ...")
	_ = router.Run()
}

// --------
// Handlers

func newRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

func listRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func updateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idx := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			idx = i
		}
	}
	if idx == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	recipes[idx] = recipe
	c.JSON(http.StatusOK, recipe)
}

func deleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	idx := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			idx = i
		}
	}
	if idx == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	recipes = append(recipes[:idx], recipes[idx+1:]...)
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
}

func searchRecipeHandler(c *gin.Context) {
	t := c.Query("tag")
	res := make([]Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, tt := range recipes[i].Tags {
			if strings.EqualFold(tt, t) {
				found = true
			}
		}
		if found {
			res = append(res, recipes[i])
		}
	}
	c.JSON(http.StatusOK, res)
}
