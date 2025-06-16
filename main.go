package main

import (
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"olimpo-vicedecanatura/config"
	"olimpo-vicedecanatura/database"
)

// TipologiaAsignatura representa los tipos permitidos de asignaturas
type TipologiaAsignatura string

const (
	TipologiaDisciplinarOptativa   TipologiaAsignatura = "DISCIPLINAR OPTATIVA"
	TipologiaFundamentalObligatoria TipologiaAsignatura = "FUND. OBLIGATORIA"
	TipologiaFundamentalOptativa    TipologiaAsignatura = "FUND. OPTATIVA"
	TipologiaDisciplinarObligatoria TipologiaAsignatura = "DISCIPLINAR OBLIGATORIA"
	TipologiaLibreEleccion         TipologiaAsignatura = "LIBRE ELECCIÓN"
	TipologiaTrabajoGrado          TipologiaAsignatura = "TRABAJO DE GRADO"
)

// ValidarTipologia verifica si una tipología es válida
func ValidarTipologia(tipo string) bool {
	switch TipologiaAsignatura(tipo) {
	case TipologiaDisciplinarOptativa,
		 TipologiaFundamentalObligatoria,
		 TipologiaFundamentalOptativa,
		 TipologiaDisciplinarObligatoria,
		 TipologiaLibreEleccion,
		 TipologiaTrabajoGrado:
		return true
	default:
		return false
	}
}

type HistoriaAcademicaRequest struct {
	Historia string `json:"historia" binding:"required"`
}

type Asignatura struct {
	Nombre      string            `json:"nombre"`
	Codigo      string            `json:"codigo"`
	Creditos    int               `json:"creditos"`
	Tipo        TipologiaAsignatura `json:"tipo"`
	Periodo     string            `json:"periodo"`
	Calificacion float64           `json:"calificacion"`
	Estado      string            `json:"estado"`
}

type ResumenCreditos struct {
	Tipologia  TipologiaAsignatura `json:"tipologia"`
	Exigidos   int                 `json:"exigidos"`
	Aprobados  int                 `json:"aprobados"`
	Pendientes int                 `json:"pendientes"`
	Inscritos  int                 `json:"inscritos"`
	Cursados   int                 `json:"cursados"`
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

	// Endpoint para recibir la historia académica
	r.POST("/api/historia-academica", func(c *gin.Context) {
		var req HistoriaAcademicaRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'historia' es requerido"})
			return
		}

		// Respuesta de ejemplo (mock)
		resp := HistoriaAcademicaResponse{
			PlanEstudios:     "INGENIERÍA DE SISTEMAS E INFORMÁTICA",
			Facultad:         "FACULTAD DE MINAS",
			PAPA:             4.3,
			Promedio:         4.3,
			Asignaturas: []Asignatura{
				{Nombre: "Desarrollo móvil", Codigo: "3011171", Creditos: 3, Tipo: TipologiaDisciplinarOptativa, Periodo: "2024-2S", Calificacion: 4.7, Estado: "APROBADA"},
				{Nombre: "Desarrollo web I", Codigo: "3011019", Creditos: 3, Tipo: TipologiaDisciplinarOptativa, Periodo: "2024-2S", Calificacion: 4.8, Estado: "APROBADA"},
			},
			ResumenCreditos: []ResumenCreditos{
				{Tipologia: TipologiaDisciplinarOptativa, Exigidos: 22, Aprobados: 9, Pendientes: 13, Inscritos: 9, Cursados: 9},
				{Tipologia: TipologiaFundamentalObligatoria, Exigidos: 27, Aprobados: 27, Pendientes: 0, Inscritos: 0, Cursados: 27},
			},
			PorcentajeAvance: 76.9,
		}
		c.JSON(http.StatusOK, resp)
	})

	// Ejecuta el servidor en el puerto 8080 (por defecto)
	log.Println("🚀 Servidor iniciado en http://localhost:8080")
	if err := r.Run(); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}
