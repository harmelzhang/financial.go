package category

import (
	"bytes"
	"financial/config"
	"financial/models"
	"financial/utils/http"
	"financial/utils/tools"
	"github.com/antchfx/htmlquery"
	"log"
	"strings"
)

var categorys = make([]models.Category, 0)

func init() {
	log.Println("初始化证监会行业分类信息")
	root, err := htmlquery.Parse(bytes.NewReader(http.Get(config.FetchCategoryUrl)))
	if err != nil {
		log.Fatalf("解析HTML出错 : %s", err)
	}

	nodes := htmlquery.Find(root, "//*[@id=\"f0-f7\"]/ul/li")
	var categoryOrder uint8 = 1
	for _, node := range nodes {
		categoryId := strings.Split(tools.FindAttrVal(node, "qquery"), ":")[1]
		categoryName := tools.FindAttrVal(htmlquery.Find(node, "./a[1]")[0], "title")
		category := models.Category{
			Id:           categoryId,
			Name:         categoryName,
			DisplayOrder: categoryOrder,
		}
		categorys = append(categorys, category)
		categoryOrder++

		var subCategoryOrder uint8 = 1
		for _, subNode := range htmlquery.Find(node, "./ul/li") {
			subCategoryId := tools.FindAttrVal(subNode, "qid")
			subCategoryName := tools.FindAttrVal(htmlquery.Find(subNode, "./a[1]")[0], "title")
			subCategory := models.Category{
				Id:           subCategoryId,
				Name:         subCategoryName,
				DisplayOrder: subCategoryOrder,
				ParentId:     categoryId,
			}
			categorys = append(categorys, subCategory)
			subCategoryOrder++
		}
	}
}

// GetCategorys 获取所有分类信息
func GetCategorys() []models.Category {
	return categorys
}
