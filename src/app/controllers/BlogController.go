package controllers

import (
	// "fmt"
	// "time"

	"log"

	"app/models"
	"net/http"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type BlogController struct {
	Db *gorm.DB
}

type ResponseData struct {
	Code    int
	Message string
	Errors  error
}

const createdFormat = "2006-01-02 15:04:05"

func (ctrl BlogController) TestTimeIndex(c echo.Context) (err error) {
	// // get time current
	// loc,_ := time.LoadLocation("Asia/Bangkok")
	// timenow := time.Now().In(loc).Format(createdFormat)

	// public,_ := time.Parse(createdFormat,"2018-05-17 13:53:50")
	// publics := public.Unix()
	// now := time.Unix(publics,0).Format(createdFormat)

	// local := time.Now().Local()
	// un := now.Unix()

	// // get time from date data
	// now , _ := time.Parse(createdFormat,"2018-05-17 13:53:50")
	// un := now.Unix()
	// tt := time.Unix(un, 0).Format(createdFormat)
	DBHost, _ := os.LookupEnv("DB_HOST")
	log.Println(DBHost)
	return c.JSON(200, DBHost)

	// cf := config.Get().LocalTimeZone
	// log.Println(cf)
	// return c.JSON(200, cf)
}

// --- List Blog ---
func (ctrl BlogController) IndexBlog(c echo.Context) (err error) {
	var modelBlog []models.Blog
	getObj := ctrl.Db
	if c.QueryParam("UserId") != "" {
		getObj = getObj.Where("user_id = ?", c.QueryParam("UserId"))
	}
	if c.QueryParam("Slug") != "" {
		getObj = getObj.Where("slug = ?", c.QueryParam("Slug"))
	}
	if c.QueryParam("Format") != "" {
		getObj = getObj.Where("format = ?", c.QueryParam("Format"))
	}
	if c.QueryParam("Status") != "" {
		getObj = getObj.Where("status = ?", c.QueryParam("Status"))
	}
	getObj.Find(&modelBlog)

	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, modelBlog)
}

// --- Create Blog ---
func (ctrl BlogController) StoreBlog(c echo.Context) (err error) {
	var modelBlog models.Blog
	modelBlog.UserId, _ = strconv.Atoi(c.FormValue("UserId"))
	modelBlog.TitleTh = c.FormValue("TitleTh")
	modelBlog.TitleEn = c.FormValue("TitleEn")
	modelBlog.ContentTh = c.FormValue("ContentTh")
	modelBlog.ContentEn = c.FormValue("ContentEn")
	modelBlog.Slug = c.FormValue("Slug")
	modelBlog.Status = c.FormValue("Status")

	formErr := modelBlog.Validate()
	if formErr != nil {
		statusCode := 422
		Res := ResponseData{statusCode, "Validation Failed", formErr}
		return c.JSON(statusCode, Res)
	}
	if c.FormValue("Status") == "publish" {
		models.SetCurrentTime(&modelBlog, "PublishedAt")
	}
	models.SetCurrentTime(&modelBlog, "CreatedAt", "UpdatedAt")
	ctrl.Db.Create(&modelBlog)

	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusCreated, &modelBlog)
}

// --- Get Blog by ID ---
func (ctrl BlogController) GetBlog(c echo.Context) (err error) {
	var modelBlog models.Blog
	var count int
	id := c.Param("id")

	ctrl.Db.Where("id = ?", id).Find(&modelBlog).Count(&count)
	if count > 0 {
		ctrl.Db.First(&modelBlog, id)
	} else {
		statusCode := 500
		Res := ResponseData{statusCode, "Blog was not found", nil}
		return c.JSON(statusCode, Res)
	}

	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, modelBlog)
}

// --- Update Blog ---
func (ctrl BlogController) UpdateBlog(c echo.Context) (err error) {
	var modelBlog models.Blog
	var count int
	id := c.Param("id")

	ctrl.Db.Where("id = ?", id).Find(&modelBlog).Count(&count)
	if count > 0 {
		// ctrl.Db.First(&modelBlog, id)
		if c.FormValue("UserId") != "" {
			modelBlog.UserId, _ = strconv.Atoi(c.FormValue("UserId"))
		}
		if c.FormValue("TitleTh") != "" {
			modelBlog.TitleTh = c.FormValue("TitleTh")
		}
		if c.FormValue("TitleEn") != "" {
			modelBlog.TitleEn = c.FormValue("TitleEn")
		}
		if c.FormValue("ContentTh") != "" {
			modelBlog.ContentTh = c.FormValue("ContentTh")
		}
		if c.FormValue("ContentEn") != "" {
			modelBlog.ContentEn = c.FormValue("ContentEn")
		}
		if c.FormValue("Slug") != "" {
			modelBlog.Slug = c.FormValue("Slug")
		}
		if c.FormValue("Status") != "" {
			modelBlog.Status = c.FormValue("Status")
			if c.FormValue("Status") == "publish" {
				models.SetCurrentTime(&modelBlog, "PublishedAt")
			} else if c.FormValue("Status") == "draft" {
				modelBlog.PublishedAt = "0000-00-00 00:00:00"
			}
		}

		formErr := modelBlog.Validate()
		if formErr != nil {
			statusCode := 422
			Res := ResponseData{statusCode, "Validation Failed", formErr}
			return c.JSON(statusCode, Res)
		}
		// models.SetCurrentTime(&modelBlog, "UpdatedAt")
		ctrl.Db.Model(&modelBlog).Where("id = ?", id).Updates(modelBlog)
	} else {
		statusCode := 500
		Res := ResponseData{statusCode, "Blog was not found", nil}
		return c.JSON(statusCode, Res)
	}

	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, &modelBlog)
}

// --- Delete Blog ---
func (ctrl BlogController) DeleteBlog(c echo.Context) (err error) {
	var modelBlog models.Blog
	var count int
	id := c.Param("id")

	data := ctrl.Db.Model(&modelBlog).Where("id = ?", id)
	data.Find(&modelBlog).Count(&count)
	if count <= 0 {
		statusCode := 500
		Res := ResponseData{statusCode, "Blog was not found", nil}
		return c.JSON(statusCode, Res)
	}
	models.SetCurrentTime(&modelBlog, "DeletedAt")
	// data.Delete(&modelBlog)
	data.Update("DeletedAt", modelBlog.DeletedAt)

	if err != nil {
		log.Println(err)
	}

	statusCode := 200
	Res := ResponseData{statusCode, "Delete blog (ID = " + id + ") success", nil}
	return c.JSON(statusCode, Res)
}
