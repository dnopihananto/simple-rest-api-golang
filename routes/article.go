package routes

import (
	"errors"
	"strconv"
	"time"

	"github.com/dnopihananto/gin-full-api/config"
	"github.com/dnopihananto/gin-full-api/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func GetHome(c *gin.Context) {

	items := []models.Article{}
	config.DB.Find(&items)

	c.JSON(200, gin.H{
		"status":  "berhasil",
		"message": "berhasil akses home",
		"data":    items,
	})
}

func GetArticle(c *gin.Context) {
	slug := c.Param("slug")

	var item models.Article

	dbRresult := config.DB.First(&item, "slug = ?", slug)
	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "berhasil",
		"message": slug,
		"data":    item,
	})
}

func PostArticle(c *gin.Context) {

	var oldItem models.Article
	slug := slug.Make(c.PostForm("title"))

	dbRresult := config.DB.First(&oldItem, "slug = ?", slug)
	if !errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	item := models.Article{
		Title:  c.PostForm("title"),
		Desc:   c.PostForm("desc"),
		Tag:    c.PostForm("tag"),
		Slug:   slug,
		UserID: uint(c.MustGet("jwt_user_id").(float64)),
	}

	config.DB.Create(&item)

	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   item,
	})
}

func GetArticleByTag(c *gin.Context) {
	tag := c.Param("tag")
	items := []models.Article{}

	config.DB.Where("tag LIKE ?", "%"+tag+"%").Find(&items)

	c.JSON(200, gin.H{"data": items})
}

func UpdateArticle(c *gin.Context) {
	id := c.Param("id")

	var item models.Article

	dbRresult := config.DB.First(&item, "id = ?", id)
	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found",
		})
		c.Abort()
		return
	}

	if uint(c.MustGet("jwt_user_id").(float64)) != item.UserID {
		c.JSON(403, gin.H{
			"status":  "error",
			"message": "this data is forbidden",
		})
		c.Abort()
		return
	}

	config.DB.Model(&item).Where("id = ?", id).Updates(models.Article{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Tag:   c.PostForm("tag"),
	})

	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   item,
	})
}

func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	config.DB.Where("id = ?", id).Delete(&article)

	c.JSON(200, gin.H{
		"status": "berhasil delete",
		"data":   article,
	})
}
