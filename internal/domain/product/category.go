package product

type Category string

const (
    CategoryFood     Category = "FOOD"
    CategoryDrinks   Category = "DRINKS"
    CategoryCleaning Category = "CLEANING"
    CategoryHygiene  Category = "HYGIENE"
    CategoryOthers   Category = "OTHERS"
)

func (c Category) isValid() bool {
    switch c {
    case CategoryFood,
        CategoryDrinks,
        CategoryCleaning,
        CategoryHygiene,
        CategoryOthers:
        return true
    default:
        return false
    }
}
