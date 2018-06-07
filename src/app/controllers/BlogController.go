package controllers

import (
	"strconv"
	// "time"
	"log"

	"net/http"

	"app/models"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

type BlogController struct {
	Db *mgo.Database
}

type ResponseData struct {
	Code    int
	Message string
	Errors  error
}

const COLLECTION = "blog_contents"

// --- List Blog ---
func (ctrl BlogController) IndexBlog(c echo.Context) (err error) {
	modelBlog := []models.Blog{}
	query := bson.M{}
	if c.QueryParam("userId") != "" {
		userIdForm, _ := strconv.Atoi(c.FormValue("userId"))
		query["user_id"] = userIdForm
	}
	if c.QueryParam("slug") != "" {
		query["slug"] = "/.*" + c.QueryParam("slug") + ".*/"
	}
	if c.QueryParam("format") != "" {
		query["format"] = c.QueryParam("format")
	}
	if c.QueryParam("Status") != "" {
		query["format"] = c.QueryParam("status")
	}
	query["deleted_at"] = nil
	getObj := ctrl.Db.C(COLLECTION).Find(query)
	LimitPerPage := 10
	SkipDoc := 0

	if c.QueryParam("limit") != "" {
		LimitPerPage, _ = strconv.Atoi(c.FormValue("limit"))
		getObj = getObj.Limit(LimitPerPage)
	}
	if c.QueryParam("page") != "" {
		page, _ := strconv.Atoi(c.FormValue("page"))
		if page <= 1 {
			page = 1
		}
		page -= 1

		SkipDoc = LimitPerPage * page
		getObj = getObj.Skip(SkipDoc)
	}
	getObj.All(&modelBlog)

	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, modelBlog)
}

// --- Create Blog ---
func (ctrl BlogController) StoreBlog(c echo.Context) (err error) {
	var modelBlog models.Blog
	var blogFormData models.Blog
	if errMap := c.Bind(&blogFormData); err != nil {
		return errMap
	}
	modelBlog.ID = bson.NewObjectId()
	modelBlog.UserId = blogFormData.UserId
	modelBlog.TitleTh = blogFormData.TitleTh
	modelBlog.TitleEn = blogFormData.TitleEn
	modelBlog.ContentTh = blogFormData.ContentTh
	modelBlog.ContentEn = blogFormData.ContentEn
	modelBlog.Slug = blogFormData.Slug
	modelBlog.Status = blogFormData.Status

	formErr := modelBlog.Validate()
	if formErr != nil {
		statusCode := 422
		Res := ResponseData{statusCode, "Validation Failed", formErr}
		return c.JSON(statusCode, Res)
	}
	if blogFormData.Status == "publish" {
		models.SetCurrentTime(&modelBlog, "PublishedAt")
	}
	models.SetCurrentTime(&modelBlog, "CreatedAt", "UpdatedAt")

	ctrl.Db.C(COLLECTION).Insert(&modelBlog)

	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusCreated, &modelBlog)
}

// --- Get Blog by ID ---
func (ctrl BlogController) GetBlog(c echo.Context) (err error) {
	modelBlog := models.Blog{}
	id := c.Param("id")
	getObj := ctrl.Db

	// --- check param id is objectId ?? ---
	if !bson.IsObjectIdHex(id) {
		statusCode := 404
		Res := ResponseData{statusCode, "Blog was not found", nil}
		return c.JSON(statusCode, Res)
	}

	query := getObj.C(COLLECTION).Find(bson.M{"_id": bson.ObjectIdHex(id), "deleted_at": nil})
	count, _ := query.Count()
	if count <= 0 {
		statusCode := 404
		Res := ResponseData{statusCode, "Blog was not found", nil}
		return c.JSON(statusCode, Res)
	}
	//---------------------------------------

	query.One(&modelBlog)
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, modelBlog)
}

// --- Update Blog ---
func (ctrl BlogController) UpdateBlog(c echo.Context) (err error) {
	modelBlog := models.Blog{}
	id := c.Param("id")
	getObj := ctrl.Db

	// --- check param id is objectId ? and check blog by id ---
	if !bson.IsObjectIdHex(id) {
		statusCode := 404
		Res := ResponseData{statusCode, "Blog was not found", nil}
		return c.JSON(statusCode, Res)
	}
	getCl := getObj.C(COLLECTION)
	query := getCl.Find(bson.M{"_id": bson.ObjectIdHex(id), "deleted_at": nil})
	count, _ := query.Count()
	if count <= 0 {
		statusCode := 404
		Res := ResponseData{statusCode, "Blog was not found", nil}
		return c.JSON(statusCode, Res)
	}
	//---------------------------------------

	query.One(&modelBlog)

	// var blogFormData models.Blog
	blogFormData := modelBlog
	if errMap := c.Bind(&blogFormData); err != nil {
		return errMap
	}
	modelBlog.UserId = blogFormData.UserId
	modelBlog.TitleTh = blogFormData.TitleTh
	modelBlog.TitleEn = blogFormData.TitleEn
	modelBlog.ContentTh = blogFormData.ContentTh
	modelBlog.ContentEn = blogFormData.ContentEn
	modelBlog.Slug = blogFormData.Slug
	modelBlog.Status = blogFormData.Status

	formErr := modelBlog.Validate()
	if formErr != nil {
		statusCode := 422
		Res := ResponseData{statusCode, "Validation Failed", formErr}
		return c.JSON(statusCode, Res)
	}
	if blogFormData.Status == "publish" {
		models.SetCurrentTime(&modelBlog, "PublishedAt")
	} else if blogFormData.Status == "draft" {
		modelBlog.PublishedAt = nil
	}
	models.SetCurrentTime(&modelBlog, "UpdatedAt")
	getCl.UpdateId(bson.ObjectIdHex(id), &modelBlog)

	if err != nil {
		log.Println(err)
	}

	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, &modelBlog)
}

// --- Delete Blog ---
func (ctrl BlogController) DeleteBlog(c echo.Context) (err error) {
	modelBlog := models.Blog{}
	id := c.Param("id")
	getObj := ctrl.Db
	typeDelete := "soft"
	if c.QueryParam("typeDel") != "" {
		typeDelete = c.QueryParam("typeDel")
	}
	// --- check param id is objectId ? and check blog by id ---
	if !bson.IsObjectIdHex(id) {
		statusCode := 404
		Res := ResponseData{statusCode, "Blog was not found", nil}
		return c.JSON(statusCode, Res)
	}
	getCl := getObj.C(COLLECTION)
	query := getCl.Find(bson.M{"_id": bson.ObjectIdHex(id), "deleted_at": nil})
	count, _ := query.Count()
	if count <= 0 {
		statusCode := 404
		Res := ResponseData{statusCode, "Blog was not found", nil}
		return c.JSON(statusCode, Res)
	}
	//---------------------------------------

	query.One(&modelBlog)

	if typeDelete == "force" {
		getCl.RemoveId(bson.ObjectIdHex(id))
	} else {
		models.SetCurrentTime(&modelBlog, "DeletedAt")
		getCl.UpdateId(bson.ObjectIdHex(id), &modelBlog)
	}

	if err != nil {
		log.Println(err)
	}
	statusCode := 200
	Res := ResponseData{statusCode, "Delete blog (ID = " + id + ") success", nil}
	return c.JSON(statusCode, Res)
}
