package orm

type Release struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	Prerelease bool   `gorm:"not_null;default:false" json:"prerelease"`
	Name       string `gorm:"not_null" json:"name"`
	Body       string `gorm:"not_null" json:"body"`
}
