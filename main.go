package main

import (
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"olimpo-vicedecanatura/config"
	"olimpo-vicedecanatura/database"
	"olimpo-vicedecanatura/models"
	"olimpo-vicedecanatura/services"
)

type HistoriaAcademicaRequest struct {
	Historia string `json:"historia" binding:"required"`
}

type Asignatura struct {
	Nombre      string  `json:"nombre"`
	Codigo      string  `json:"codigo"`
	Creditos    int     `json:"creditos"`
	Tipo        string  `json:"tipo"`
	Periodo     string  `json:"periodo"`
	Calificacion float64 `json:"calificacion"`
	Estado      string  `json:"estado"`
}

type ResumenCreditos struct {
	Tipologia  string `json:"tipologia"`
	Exigidos   int    `json:"exigidos"`
	Aprobados  int    `json:"aprobados"`
	Pendientes int    `json:"pendientes"`
	Inscritos  int    `json:"inscritos"`
	Cursados   int    `json:"cursados"`
}

type HistoriaAcademicaResponse struct {
	PlanEstudios      string            `json:"plan_estudios"`
	Facultad          string            `json:"facultad"`
	PAPA              float64           `json:"papa"`
	Promedio          float64           `json:"promedio"`
	Asignaturas       []Asignatura      `json:"asignaturas"`
	ResumenCreditos   []ResumenCreditos `json:"resumen_creditos"`
	PorcentajeAvance  float64           `json:"porcentaje_avance"`
}

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env (puede que no exista o ya estén las variables en el entorno)")
	}

	// Inicializar la base de datos
	config.InitDB()

	// Verificar la conexión
	sqlDB, err := config.DB.DB()
	if err != nil {
		log.Fatalf("Error obteniendo la conexión SQL: %v", err)
	}
	defer sqlDB.Close()

	// Probar la conexión
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	log.Println("✅ Conexión a la base de datos establecida exitosamente")

	// Ejecutar migraciones
	database.RunMigrations(config.DB)
	log.Println("✅ Migraciones ejecutadas exitosamente")

	// Insertar datos iniciales (opcional)
	database.SeedInitialData(config.DB)
	log.Println("✅ Datos iniciales cargados (si era necesario)")

	// Creamos un router con middlewares por defecto (logger, recovery)
	r := gin.Default()

	// Definimos la ruta raíz con un manejador de tipo GET
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API de Olimpo Vicedecanatura",
			"status":  "online",
			"db":      "connected",
		})
	})

	// NUEVO ENDPOINT 
	r.POST("/api/compare", func(c *gin.Context) {
		var req struct {
			StudyPlanID     uint                        `json:"study_plan_id" binding:"required"`
			AcademicHistory models.AcademicHistoryInput `json:"academic_history" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos: " + err.Error()})
			return
		}
		
		// Obtener el plan de estudio
		var studyPlan models.StudyPlan
		if err := config.DB.First(&studyPlan, req.StudyPlanID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plan de estudio no encontrado"})
			return
		}
		
		// USAR LA FUNCIÓN DE COMPARACIÓN REAL
		result, err := services.CompareStudyPlanWithHistory(config.DB, studyPlan, req.AcademicHistory)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// Retornar el resultado real
		c.JSON(http.StatusOK, result)
	})

	// Ejecuta el servidor en el puerto 8080 (por defecto)
	log.Println("🚀 Servidor iniciado en http://localhost:8080")
	if err := r.Run(); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}
