package main

import "log"

var allCategories []Category
var categoryByID map[int]Category
var childCategories map[int][]int

func loadCategories() {
	err := dbx.Select(&allCategories, "SELECT * FROM `categories`")
	if err != nil {
		log.Fatal(err)
	}

	categoryByID = make(map[int]Category)
	for _, c := range allCategories {
		categoryByID[c.ID] = c
	}

	childCategories = make(map[int][]int)

	for _, c := range categoryByID {
		if c.ParentID != 0 {
			c.ParentCategoryName = categoryByID[c.ParentID].CategoryName
		}

		categoryByID[c.ID] = c

		var children []int

		for id, child := range categoryByID {
			if child.ParentID == c.ID {
				children = append(children, id)
			}
		}

		if len(children) > 0 {
			childCategories[c.ID] = children
		}
	}
}

func getChildCategories(parentID int) []int {
	return childCategories[parentID]
}
