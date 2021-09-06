package main

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaginationResult struct {
	Page    int         `json:"page"`
	PerPage int         `json:"perPage"`
	Count   int64       `json:"count"`
	Items   interface{} `json:"items"`
}

func GetList(inputStruct interface{}, c *gin.Context, db *gorm.DB, preloads []string) (result PaginationResult, err error) {
	var (
		resultSliceT = reflect.SliceOf(reflect.TypeOf(inputStruct))
		resultSlice  = reflect.New(resultSliceT).Interface()
		ctx, done    = context.WithTimeout(context.Background(), time.Second*30)
		tx           = db.WithContext(ctx)
	)
	result = SetPagination(c)

	tx = SetPreload(db, preloads)
	tx = tx.
		Scopes(SetFilters(c)).
		Scopes(SetSort(c))
	err = tx.
		Find(resultSlice).
		Count(&result.Count).
		Error
	if err != nil {
		return
	}

	err = tx.
		Scopes(Paginate(c)).
		Find(resultSlice).
		Error
	if err != nil {
		return
	}

	result.Items = resultSlice
	done()
	return result, err
}

func SetSort(c *gin.Context) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		var (
			sort  = c.Query("sort")
			field = c.Query("sort_field")
		)
		if sort != "" && field != "" {
			return tx.Order(field + " " + sort)
		}
		return tx
	}
}

func SetFilters(c *gin.Context) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		var (
			desiredValues []interface{}
			logic         = fmt.Sprintf(` %s `, c.DefaultQuery("filter_logic", "and"))
			ftype         = c.DefaultQuery("filter_type", "=")
			values        = strings.Split(strings.Replace(c.DefaultQuery("filter_values", ""), " ", "", -1), ",")
			fields        = strings.Split(strings.Replace(c.DefaultQuery("filter_fields", ""), " ", "", -1), ",")
		)
		if fields[0] == "" || values[0] == "" {
			return tx
		}
		if len(fields) != len(values) {
			return tx
		}
		if ftype == "like" {
			for index, filter := range values {
				filter = fmt.Sprintf(`%%%s%%`, filter)
				values[index] = filter
			}
		}
		for index := range fields {
			fields[index] = fmt.Sprintf(`%s::text %s ?`, fields[index], ftype)
			desiredValues = append(desiredValues, values[index])
		}
		arguments := strings.Join(fields, logic)
		return tx.Where(arguments, desiredValues...)
	}
}

func SetPreload(db *gorm.DB, preloads []string) (tx *gorm.DB) {
	if len(preloads) == 0 {
		return db
	}

	for _, val := range preloads {
		tx = db.Preload(val)
	}
	return tx
}

func SetPagination(c *gin.Context) (pagination PaginationResult) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page == 0 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))
	if perPage > 100 || perPage <= 0 {
		perPage = 100
	}
	pagination = PaginationResult{
		Page:    page,
		PerPage: perPage,
	}
	return
}
func Paginate(c *gin.Context) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		p := SetPagination(c)
		offset := (p.Page - 1) * p.PerPage
		return tx.Offset(offset).Limit(p.PerPage)
	}
}
