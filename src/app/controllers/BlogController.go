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
	// // tt := time.Unix(un, 0).Format(createdFormat)
	// DBHost, _ := os.LookupEnv("DB_HOST")
	// log.Println(DBHost)
	// errOp := godotenv.Load(".env")
	// if errOp != nil {
	// 	log.Println(errOp)
	// }
	// // return c.JSON(200, DBHost)
	// return c.JSON(200, os.Getenv("FOO"))
	// cf := config.Get().LocalTimeZone
	// log.Println(cf)
	return c.JSON(200, "")
}

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
	getObj := ctrl.Db.C(models.COLLECTION).Find(query)
	LimitPerPage := 10
	SkipDoc := 0

	if c.QueryParam("perPage") != "" {
		LimitPerPage, _ = strconv.Atoi(c.FormValue("perPage"))
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

	ctrl.Db.C(models.COLLECTION).Insert(&modelBlog)

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

	query := getObj.C(models.COLLECTION).Find(bson.M{"_id": bson.ObjectIdHex(id), "deleted_at": nil})
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
	getCl := getObj.C(models.COLLECTION)
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

	// var modelBlog models.Blog
	// var count int
	// id := c.Param("id")

	// ctrl.Db.Where("id = ?", id).Find(&modelBlog).Count(&count)
	// if count > 0 {
	// 	// ctrl.Db.First(&modelBlog, id)
	// 	if c.FormValue("UserId") != "" {
	// 		modelBlog.UserId, _ = strconv.Atoi(c.FormValue("UserId"))
	// 	}
	// 	if c.FormValue("TitleTh") != "" {
	// 		modelBlog.TitleTh = c.FormValue("TitleTh")
	// 	}
	// 	if c.FormValue("TitleEn") != "" {
	// 		modelBlog.TitleEn = c.FormValue("TitleEn")
	// 	}
	// 	if c.FormValue("ContentTh") != "" {
	// 		modelBlog.ContentTh = c.FormValue("ContentTh")
	// 	}
	// 	if c.FormValue("ContentEn") != "" {
	// 		modelBlog.ContentEn = c.FormValue("ContentEn")
	// 	}
	// 	if c.FormValue("Slug") != "" {
	// 		modelBlog.Slug = c.FormValue("Slug")
	// 	}
	// 	if c.FormValue("Status") != "" {
	// 		modelBlog.Status = c.FormValue("Status")
	// 		if c.FormValue("Status") == "publish" {
	// 			models.SetCurrentTime(&modelBlog, "PublishedAt")
	// 		} else if c.FormValue("Status") == "draft" {
	// 			modelBlog.PublishedAt = "0000-00-00 00:00:00"
	// 		}
	// 	}

	// 	formErr := modelBlog.Validate()
	// 	if formErr != nil {
	// 		statusCode := 422
	// 		Res := ResponseData{statusCode, "Validation Failed", formErr}
	// 		return c.JSON(statusCode, Res)
	// 	}
	// 	// models.SetCurrentTime(&modelBlog, "UpdatedAt")
	// 	ctrl.Db.Model(&modelBlog).Where("id = ?", id).Updates(modelBlog)
	// } else {
	// 	statusCode := 500
	// 	Res := ResponseData{statusCode, "Blog was not found", nil}
	// 	return c.JSON(statusCode, Res)
	// }

	// if err != nil {
	// 	log.Println(err)
	// }
	// return c.JSON(http.StatusOK, &modelBlog)
}

// --- Delete Blog ---
func (ctrl BlogController) DeleteBlog(c echo.Context) (err error) {
	modelBlog := models.Blog{}
	id := c.Param("id")
	getObj := ctrl.Db

	// --- check param id is objectId ? and check blog by id ---
	if !bson.IsObjectIdHex(id) {
		statusCode := 404
		Res := ResponseData{statusCode, "Blog was not found", nil}
		return c.JSON(statusCode, Res)
	}
	getCl := getObj.C(models.COLLECTION)
	query := getCl.Find(bson.M{"_id": bson.ObjectIdHex(id), "deleted_at": nil})
	count, _ := query.Count()
	if count <= 0 {
		statusCode := 404
		Res := ResponseData{statusCode, "Blog was not found", nil}
		return c.JSON(statusCode, Res)
	}
	//---------------------------------------

	query.One(&modelBlog)
	models.SetCurrentTime(&modelBlog, "DeletedAt")
	getCl.UpdateId(bson.ObjectIdHex(id), &modelBlog)

	if err != nil {
		log.Println(err)
	}
	statusCode := 200
	Res := ResponseData{statusCode, "Delete blog (ID = " + id + ") success", nil}
	return c.JSON(statusCode, Res)
}
