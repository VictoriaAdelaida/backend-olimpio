package database

import (
	"log"

	"gorm.io/gorm"
	"olimpo-vicedecanatura/models"
)

// RunMigrations ejecuta las migraciones de la base de datos
func RunMigrations(db *gorm.DB) {
	// Auto-migrar los modelos
	err := db.AutoMigrate(
		&models.Career{},
		&models.StudyPlan{},
		&models.Subject{},
		&models.Equivalence{},
	)
	if err != nil {
		log.Fatalf("Error ejecutando migraciones: %v", err)
	}

	// Crear índices adicionales si son necesarios
	// Por ejemplo, para búsquedas frecuentes por código de materia
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_subjects_code ON subjects(code);").Error; err != nil {
		log.Printf("Error creando índice: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_careers_code ON careers(code);").Error; err != nil {
		log.Printf("Error creando índice: %v", err)
	}
}

// SeedInitialData inserta datos iniciales en la base de datos
func SeedInitialData(db *gorm.DB) {
	// Verificar si ya existen datos
	var count int64
	db.Model(&models.Career{}).Count(&count)
	if count > 0 {
		log.Println("La base de datos ya contiene datos iniciales")
		return
	}

	// Crear algunas carreras de ejemplo
	careers := []models.Career{
		{
			Name:        "Ingeniería de Sistemas",
			Code:        "ISIS",
			Description: "Carrera de Ingeniería de Sistemas",
		},
		{
			Name:        "Ingeniería Administrativa",
			Code:        "IADM",
			Description: "Carrera de Ingeniería Administrativa",
		},
		// Agregar más carreras según sea necesario
	}

	// Insertar carreras
	for _, career := range careers {
		if err := db.Create(&career).Error; err != nil {
			log.Printf("Error creando carrera %s: %v", career.Name, err)
		}
	}

	// Nota: Los planes de estudio y materias se pueden cargar desde archivos JSON
	// o mediante una interfaz administrativa
} 