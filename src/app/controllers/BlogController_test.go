package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"app/db"
	"app/models"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// var (
// 	mockModelBlog = &models.Blog{
// 		TitleEn:   "Test Title EN",
// 		TitleTh:   "Test Title TH",
// 		ContentEn: "Test Content EN",
// 		ContentTh: "Test Content TH",
// 		Slug:      "Test Slug",
// 		Status:    "draft",
// 		UserId:    555,
// 	}
// )

var blogId string

func connectDb() db.EnvDB {
	var cfEnv db.EnvDB
	errOp := godotenv.Load("../.env")
	if errOp != nil {
		log.Println(errOp)
	}
	cfEnv.Host = os.Getenv("DB_HOST")
	cfEnv.Dbname = os.Getenv("DB_DATABASE")
	cfEnv.Port = os.Getenv("DB_PORT")
	return cfEnv
}

func getValueBlogValidate(validate string) *models.Blog {
	mockModelBlog := &models.Blog{
		TitleEn:   "Test Title EN",
		TitleTh:   "Test Title TH",
		ContentEn: "Test Content EN",
		ContentTh: "Test Content TH",
		Slug:      "Test Slug",
		Status:    "draft",
		UserId:    555,
	}

	switch validate {
	case "TitleEn":
		mockModelBlog.TitleEn = ""
	case "TitleTh":
		mockModelBlog.TitleTh = ""
	case "ContentEn":
		mockModelBlog.ContentEn = ""
	case "ContentTh":
		mockModelBlog.ContentTh = ""
	case "Slug":
		mockModelBlog.Slug = ""
	case "Status":
		mockModelBlog.Status = ""
	case "UserId":
		mockModelBlog.UserId = 0
	}
	return mockModelBlog
}

func createBlogForTest(validate string) (error, *httptest.ResponseRecorder) {
	cfEnv := connectDb()
	e := echo.New()

	mockModelBlog := getValueBlogValidate(validate)
	m, _ := json.Marshal(mockModelBlog)
	req := httptest.NewRequest(echo.POST, "/blog", strings.NewReader(string(m)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	b := BlogController{cfEnv.Connect()}

	return b.StoreBlog(c), rec
}

func updateBlogForTest(validate string, id string) (error, *httptest.ResponseRecorder) {
	cfEnv := connectDb()
	e := echo.New()

	mockModelBlog := getValueBlogValidate(validate)
	m, _ := json.Marshal(mockModelBlog)
	req := httptest.NewRequest(echo.PATCH, "/", strings.NewReader(string(m)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/blog/:id")
	c.SetParamNames("id")
	c.SetParamValues(id)
	b := BlogController{cfEnv.Connect()}

	return b.UpdateBlog(c), rec
}

// ----- test create with validate ------
func TestCreateBlogValidateTitleEn(t *testing.T) {
	err, rec := createBlogForTest("TitleEn")
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestCreateBlogValidateTitleTh(t *testing.T) {
	err, rec := createBlogForTest("TitleTh")
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestCreateBlogValidateContentEn(t *testing.T) {
	err, rec := createBlogForTest("ContentEn")
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestCreateBlogValidateContentTh(t *testing.T) {
	err, rec := createBlogForTest("ContentTh")
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestCreateBlogValidateSlug(t *testing.T) {
	err, rec := createBlogForTest("Slug")
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestCreateBlogValidateStatus(t *testing.T) {
	err, rec := createBlogForTest("Status")
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestCreateBlogValidateUserId(t *testing.T) {
	err, rec := createBlogForTest("UserId")
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

// -------------------------------

// ---- test create blog success ----
func TestCreateBlog(t *testing.T) {
	err, rec := createBlogForTest("")
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		if rec.Code == 201 {
			decoder := json.NewDecoder(rec.Body)
			str := make(map[string]interface{})
			decoder.Decode(&str)
			blogId = str["id"].(string)
		}
	}
}

// ---- test list all blog ----
func TestIndexBlog(t *testing.T) {
	cfEnv := connectDb()
	e := echo.New()

	req := httptest.NewRequest(echo.GET, "/blog?page=1&limit=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	b := BlogController{cfEnv.Connect()}

	// Assertions
	if assert.NoError(t, b.IndexBlog(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body)
	}
}

// ---- test get blog by id ----
func TestGetViewBlog(t *testing.T) {
	cfEnv := connectDb()
	e := echo.New()

	req := httptest.NewRequest(echo.GET, "/blog/:id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// c.SetPath("/blog/:id")
	c.SetParamNames("id")
	c.SetParamValues(blogId)
	b := BlogController{cfEnv.Connect()}

	// Assertions
	if assert.NoError(t, b.GetBlog(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body)
	}
}

// ----- test update with validate ------
func TestUpdateBlogValidateTitleEn(t *testing.T) {
	if blogId != "" {
		err, rec := updateBlogForTest("TitleEn", blogId)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		}
	}
}

func TestUpdateBlogValidateTitleTh(t *testing.T) {
	if blogId != "" {
		err, rec := updateBlogForTest("TitleTh", blogId)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		}
	}
}

func TestUpdateBlogValidateContentEn(t *testing.T) {
	if blogId != "" {
		err, rec := updateBlogForTest("ContentEn", blogId)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		}
	}
}

func TestUpdateBlogValidateContentTh(t *testing.T) {
	if blogId != "" {
		err, rec := updateBlogForTest("ContentTh", blogId)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		}
	}
}

func TestUpdateBlogValidateSlug(t *testing.T) {
	if blogId != "" {
		err, rec := updateBlogForTest("Slug", blogId)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		}
	}
}

func TestUpdateBlogValidateStatus(t *testing.T) {
	if blogId != "" {
		err, rec := updateBlogForTest("Status", blogId)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		}
	}
}

func TestUpdateBlogValidateUserId(t *testing.T) {
	if blogId != "" {
		err, rec := updateBlogForTest("UserId", blogId)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		}
	}
}

// -------------------------------

// ---- test update blog success ----
func TestUpdateBlog(t *testing.T) {
	if blogId != "" {
		err, rec := updateBlogForTest("", blogId)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	}
}

// ---- test soft delete blog success ----
func TestDeleteBlog(t *testing.T) {
	if blogId != "" {
		cfEnv := connectDb()
		e := echo.New()

		req := httptest.NewRequest(echo.DELETE, "/blog/:id?typeDel=force", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(blogId)
		b := BlogController{cfEnv.Connect()}

		// Assertions
		if assert.NoError(t, b.DeleteBlog(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.NotEmpty(t, rec.Body)
		}
	}
}
