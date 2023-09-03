package lib

type ItemJson struct {
	Category    string
	Subcategory string
	Screenshot  string
	URI         string
	Name        string
	Info        string
	Size        string
	Site        string
	Downloads   string
	Supreme     string
	Sourcecode  string
	Shareware   string
	Noinstall   string
}

func GetUniqueCategories(items []ItemJsonStruct) []string {
	uniqueCategories := make(map[string]bool)
	var result []string

	for _, item := range items {
		if _, exists := uniqueCategories[item.Category]; !exists {
			uniqueCategories[item.Category] = true
			result = append(result, item.Category)
		}
	}

	return result
}

func FilterItemsByCategory(items []ItemJsonStruct, targetCategory string) []ItemJsonStruct {
	var filteredItems []ItemJsonStruct

	for _, item := range items {
		if item.Category == targetCategory {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func FilterItemsAndUniqueSubcategories(items []ItemJsonStruct, targetCategory string) []string {
	var uniqueSubcategories []string
	seenSubcategories := make(map[string]bool)

	for _, item := range items {
		if item.Category == targetCategory {
			if _, exists := seenSubcategories[item.Subcategory]; !exists {
				seenSubcategories[item.Subcategory] = true
				uniqueSubcategories = append(uniqueSubcategories, item.Subcategory)
			}
		}
	}

	return uniqueSubcategories
}

func FilterItemsByCategoryAndSubcategory(items []ItemJsonStruct, targetCategory, targetSubcategory string) []ItemJsonStruct {
	var matchingSubcategories []ItemJsonStruct

	for _, item := range items {
		if item.Category == targetCategory && item.Subcategory == targetSubcategory {
			matchingSubcategories = append(matchingSubcategories, item)
		}
	}

	return matchingSubcategories
}

func FilterItemsByCategoryAndSubcategoryFiltered(items []ItemJsonStruct, targetCategory, targetSubcategory string,
	nofilter bool, supreme bool, shareware bool, noinstall bool, source bool) []ItemJsonStruct {
	var matchingSubcategories []ItemJsonStruct

	for _, item := range items {
		if item.Category == targetCategory && item.Subcategory == targetSubcategory && (item.Name != "") {
			if nofilter {
				matchingSubcategories = append(matchingSubcategories, item)
			} else if shareware && item.Shareware == "true" || supreme && item.Supreme == "🌱" || noinstall && item.Noinstall == "true" || source && item.Sourcecode == "true" {
				matchingSubcategories = append(matchingSubcategories, item)
			}
		}
	}

	return matchingSubcategories
}

//func GetUniqueSubcategory(items []ItemJson) []string {
//	uniqueCategories := make(map[string]bool)
//	var result []string
//
//	for _, item := range items {
//		if _, exists := uniqueCategories[item.Subcategory]; !exists {
//			uniqueCategories[item.Subcategory] = true
//			result = append(result, item.Subcategory)
//		}
//	}
//
//	return result
//}
