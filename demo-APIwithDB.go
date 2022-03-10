package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

const coursePath = "courses"
const apiBasePath = "/api"

type Course struct {
	CourseID   int     `json: "courseid"`
	CourseName string  `json: "coursename"`
	Price      float64 `json: "price"`
	ImageURL   string  `json: "imageurl"`
}

func getCourseList() ([]Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := Db.QueryContext(ctx, `SELECT
	courseid,
	coursename,
	price,
	image_url
	FROM courseonline
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer result.Close()
	courses := make([]Course, 0)
	for result.Next() {
		var course Course
		result.Scan(&course.CourseID,
			&course.CourseName,
			&course.Price,
			&course.ImageURL)

		courses = append(courses, course)
	}
	return courses, nil
}

func insertProduct(course Course) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := Db.ExecContext(ctx, `INSERT INTO courseonline (
	courseid,
	coursename,
	price,
	image_url
	) VALUE(?, ?, ?, ?)`,
		&course.CourseID,
		&course.CourseName,
		&course.Price,
		&course.ImageURL)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

func getCourse(courseid int) (*Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := Db.QueryRowContext(ctx, `SELECT
	courseid,
	coursename,
	price,
	image_url
	FROM courseonline
	WHERE courseid = ?`, courseid)

	course := &Course{}
	err := row.Scan(
		&course.CourseID,
		&course.CourseName,
		&course.Price,
		&course.ImageURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return course, nil
}

func removeCourse(CourseID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := Db.ExecContext(ctx, `DELETE
	FROM courseonline
	WHERE courseid = ?`, CourseID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func handlesCourses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		courseList, err := getCourseList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(courseList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		var course Course
		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		CourseID, err := insertProduct(course)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"courseid: %d"}`, CourseID)))
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleCourses(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", coursePath))
	fmt.Println("urlPathSegments:", urlPathSegments)
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	courseID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		course, err := getCourse(courseID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if course == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodDelete:
		err := removeCourse(courseID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, X-CORS")
		handler.ServeHTTP(w, r)
	})
}

func SetupRoutes(apiBasePath string) {
	courseHandler := http.HandlerFunc(handleCourses)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, coursePath), corsMiddleware(courseHandler))
	coursesHandler := http.HandlerFunc(handlesCourses)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, coursePath), corsMiddleware(coursesHandler))
}

func SetupDB() {
	var err error
	Db, err = sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/coursedb")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Db)
	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
}

func main() {
	SetupDB()
	SetupRoutes(apiBasePath)
	log.Fatal(http.ListenAndServe(":80", nil))
}
