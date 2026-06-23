package db

import "testing"

func TestForeignKeysRejectMissingParents(t *testing.T) {
	database := testDB(t)

	profile := StudentProfile{
		UserID:         999,
		School:         "不存在大学",
		EducationLevel: "本科",
		Major:          "计算机科学",
		GraduationYear: 2026,
		SchoolVisible:  true,
	}

	if err := database.Create(&profile).Error; err == nil {
		t.Fatal("expected missing parent user to fail foreign key enforcement")
	}
}
