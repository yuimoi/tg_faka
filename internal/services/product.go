package services

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"tg_go_faka/internal/models"
	"tg_go_faka/internal/utils/db"
	"tg_go_faka/internal/utils/functions"
	_type "tg_go_faka/internal/utils/type"
)

func GetValidProductByID(id int64) (*models.Product, error) {
	var product *models.Product
	result := db.DB.Model(models.Product{}).Where("id=? and status=?", id, _type.ProductStatusValid).Find(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}
func GetProductItemByID(productItemID int64) (*models.ProductItem, error) {
	var productItem *models.ProductItem
	result := db.DB.Model(models.ProductItem{}).Where("id=?", productItemID).Find(&productItem)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}
	return productItem, nil
}
func GetProductsByPage(page int) ([]*models.Product, _type.PaginationQueryDataStruct, error) {
	var pagination _type.PaginationQueryDataStruct

	var products []*models.Product

	limit := 10

	query := db.DB.Model(models.Product{}).Where("status=?", _type.ProductStatusValid).Order("create_time desc")
	functions.ApplyPaginationQueryData(query, _type.PaginationQueryDataStruct{
		Page:  page,
		Limit: limit,
	})

	// 计算总记录数
	var totalRecords int64
	db.DB.Model(models.Product{}).Where("status=?", _type.ProductStatusValid).Count(&totalRecords)
	// 计算总页数
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))

	result := query.Find(&products)
	if result.Error != nil {
		return nil, pagination, result.Error
	}

	var items []*models.Product
	for _, product := range products {
		// 获取每个产品的product_item数量, TODO:合并为一个sql
		productItemsCount, _ := GetProductItemValidCounts(product.ID)
		product.InStockCount = productItemsCount
		items = append(items, product)
	}

	pagination = _type.PaginationQueryDataStruct{
		Page:      page,
		Limit:     limit,
		TotalPage: totalPages,
	}

	return items, pagination, nil
}

func GetProductItemValidCounts(productID int64) (int64, error) {
	var productItemsCount int64
	result := db.DB.Model(&models.ProductItem{}).Where("product_id = ? and status=?", productID, _type.ProductItemStatusValid).Count(&productItemsCount)
	if result.Error != nil {
		return 0, result.Error
	}
	return productItemsCount, nil
}

func GetAllProductItems() ([]*models.ProductItem, error) {
	var productItems []*models.ProductItem
	result := db.DB.Model(models.ProductItem{}).Find(&productItems)
	if result.Error != nil {
		return nil, result.Error
	}

	return productItems, nil
}
func GetAllProducts() ([]*models.Product, error) {
	var products []*models.Product
	result := db.DB.Model(models.Product{}).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
func DeleteAllProductItems() error {
	result := db.DB.Where("1=1").Delete(&models.ProductItem{})
	return result.Error
}
func DeleteAllProducts() error {
	result := db.DB.Where("1=1").Delete(&models.Product{})
	return result.Error
}

func GenerateExcelFromItems(items []*models.ProductItem, specificColumns []string, secret *string) []byte {
	var maps []map[string]interface{}
	for _, item := range items {
		mapData := functions.StructToMap(item)
		maps = append(maps, mapData)
	}

	// 创建一个新的Excel文件
	var f *excelize.File
	if secret != nil {
		f = excelize.NewFile(excelize.Options{Password: *secret})

	} else {
		f = excelize.NewFile()
	}

	sheetName := "Sheet1"
	// 创建一个工作表
	index, _ := f.NewSheet(sheetName)
	// 设置表头
	for colIndex, colName := range specificColumns {
		colLetter, _ := excelize.ColumnNumberToName(colIndex + 1)
		cell := colLetter + "1"
		_ = f.SetCellValue(sheetName, cell, colName)
	}
	// 填充数据
	for rowIndex, mapData := range maps {
		row := rowIndex + 2
		for colIndex, colName := range specificColumns {
			value, exists := mapData[colName]
			if exists {
				colLetter, _ := excelize.ColumnNumberToName(colIndex + 1)
				cell := colLetter + strconv.Itoa(row)
				_ = f.SetCellValue(sheetName, cell, value)

			}
		}
	}
	f.SetActiveSheet(index)

	// 写入到Buffer而不是文件
	buf, _ := f.WriteToBuffer()

	return buf.Bytes()
}
