package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Career representa una carrera en la universidad
type Career struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	Code        string    `gorm:"size:20;unique;not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func main() {
	// Configuración de la base de datos
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=olimpo_vicedecanatura port=5432 sslmode=disable"
	}

	// Conectar a la base de datos
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	fmt.Println("🔍 Verificando carrera de Ingeniería de Sistemas...")

	// Buscar la carrera existente con código ISIS
	var existingCareer Career
	if err := db.Where("code = ?", "ISIS").First(&existingCareer).Error; err != nil {
		if err := db.Where("code = ?", "ING-SISTEMAS").First(&existingCareer).Error; err != nil {
			fmt.Println("⚠️  No se encontró carrera de Ingeniería de Sistemas")
			fmt.Println("   Códigos buscados: ISIS, ING-SISTEMAS")
			fmt.Println("   Por favor, crea la carrera manualmente o verifica el código correcto")
			return
		}
	}

	fmt.Printf("✅ Carrera encontrada:\n")
	fmt.Printf("   - ID: %d\n", existingCareer.ID)
	fmt.Printf("   - Nombre: %s\n", existingCareer.Name)
	fmt.Printf("   - Código: %s\n", existingCareer.Code)
	fmt.Printf("   - Descripción: %s\n", existingCareer.Description)

	// Guardar el ID en un archivo temporal para que otros scripts lo usen
	careerInfo := fmt.Sprintf("CAREER_ID=%d\nCAREER_CODE=%s\n", existingCareer.ID, existingCareer.Code)
	if err := os.WriteFile("career_info.txt", []byte(careerInfo), 0644); err != nil {
		log.Printf("⚠️  No se pudo guardar información de la carrera: %v", err)
	} else {
		fmt.Println("💾 Información de la carrera guardada para otros scripts")
	}

	fmt.Printf("\n🎉 ¡Carrera lista para usar!\n")
} 