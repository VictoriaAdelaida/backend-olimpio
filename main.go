package main

import (
	"log"
	"net/http"
	"strconv"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"olimpo-vicedecanatura/config"
	"olimpo-vicedecanatura/database"
	"olimpo-vicedecanatura/models"
	"olimpo-vicedecanatura/functions"
	"strings"
	"errors"
	"regexp"
)


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

	// Configurar CORS y middlewares
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Ruta raíz
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API de Olimpo Vicedecanatura",
			"status":  "online",
			"db":      "connected",
			"endpoints": []string{
				"GET /api/careers - Obtener todas las carreras",
				"GET /api/careers/:code/study-plans - Obtener planes de estudio de una carrera",
				"GET /api/study-plans/:id - Obtener detalles de un plan de estudio",
				"POST /api/compare - Comparar historia académica con plan de estudio",
				"POST /api/compare-by-career - Comparar por código de carrera",
				"POST /api/api-compare - Comparar historia académica en texto plano",
			},
		})
	})


	// API Routes
	api := r.Group("/api")
	{
		// Obtener todas las carreras disponibles
		api.GET("/careers", getCareers)
		
		// Obtener planes de estudio de una carrera específica
		api.GET("/careers/:code/study-plans", getStudyPlansByCareer)
		
		// Obtener detalles de un plan de estudio específico
		api.GET("/study-plans/:id", getStudyPlanDetails)
		
		// Comparar historia académica con plan de estudio
		api.POST("/compare", compareAcademicHistory)
		
		// Endpoint adicional para comparar por código de carrera (más simple)
		api.POST("/compare-by-career", compareByCareerCode)
		
		// Nuevo endpoint para comparar historia académica en texto plano
		api.POST("/api-compare", compareAcademicHistoryFromText)
	}


	// Ejecutar servidor
	log.Println("🚀 Servidor iniciado en http://localhost:8080")
	if err := r.Run(); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}

// getCareers obtiene todas las carreras disponibles
func getCareers(c *gin.Context) {
	var careers []models.Career
	if err := config.DB.Find(&careers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo carreras"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"careers": careers,
	})
}

// getStudyPlansByCareer obtiene los planes de estudio de una carrera específica
func getStudyPlansByCareer(c *gin.Context) {
	careerCode := c.Param("code")
	
	var studyPlans []models.StudyPlan
	if err := config.DB.Preload("Career").
		Joins("JOIN careers ON careers.id = study_plans.career_id").
		Where("careers.code = ?", careerCode).
		Find(&studyPlans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo planes de estudio"})
		return
	}
	
	if len(studyPlans) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron planes de estudio para esta carrera"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"study_plans": studyPlans,
	})
}

// getStudyPlanDetails obtiene los detalles completos de un plan de estudio
func getStudyPlanDetails(c *gin.Context) {
	studyPlanID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de plan de estudio inválido"})
		return
	}
	
	var studyPlan models.StudyPlan
	if err := config.DB.Preload("Career").Preload("Subjects").
		First(&studyPlan, uint(studyPlanID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan de estudio no encontrado"})
		return
	}
	
	// Calcular estadísticas del plan
	subjectsByType := make(map[string][]models.Subject)
	creditsByType := make(map[string]int)
	
	for _, subject := range studyPlan.Subjects {
		subjectsByType[string(subject.Type)] = append(subjectsByType[string(subject.Type)], subject)
		creditsByType[string(subject.Type)] += subject.Credits
	}
	
	c.JSON(http.StatusOK, gin.H{
		"study_plan":        studyPlan,
		"subjects_by_type":  subjectsByType,
		"credits_by_type":   creditsByType,
		"total_subjects":    len(studyPlan.Subjects),
	})
}

// CompareRequest estructura para la solicitud de comparación
type CompareRequest struct {
	StudyPlanID     uint                        `json:"study_plan_id" binding:"required"`
	AcademicHistory models.AcademicHistoryInput `json:"academic_history" binding:"required"`
}

// compareAcademicHistory compara una historia académica con un plan de estudio específico
func compareAcademicHistory(c *gin.Context) {
	var req CompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos: " + err.Error()})
		return
	}
	
	// Realizar la comparación usando la función que creamos
	result, err := functions.CompareAcademicHistoryWithStudyPlan(config.DB, req.AcademicHistory, req.StudyPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Obtener información adicional del plan de estudio para el contexto
	var studyPlan models.StudyPlan
	config.DB.Preload("Career").First(&studyPlan, req.StudyPlanID)
	
	c.JSON(http.StatusOK, gin.H{
		"comparison_result": result,
		"study_plan_info": gin.H{
			"id":      studyPlan.ID,
			"version": studyPlan.Version,
			"career":  studyPlan.Career.Name,
		},
		"summary": gin.H{
			"total_subjects_in_plan":     len(result.EquivalentSubjects) + len(result.MissingSubjects),
			"approved_subjects":          len(result.EquivalentSubjects),
			"missing_subjects":           len(result.MissingSubjects),
			"completion_percentage":      calculateCompletionPercentage(result.CreditsSummary),
		},
	})
}

// compareByCareerCode compara usando el código de carrera (más simple)
func compareByCareerCode(c *gin.Context) {
	var academicHistory models.AcademicHistoryInput
	if err := c.ShouldBindJSON(&academicHistory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos: " + err.Error()})
		return
	}
	
	// Realizar la comparación usando el código de carrera
	result, err := functions.CompareAcademicHistoryByCareerCode(config.DB, academicHistory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Obtener información del plan de estudio usado
	studyPlan, _ := functions.GetStudyPlanByCareerCode(config.DB, academicHistory.CareerCode)
	
	c.JSON(http.StatusOK, gin.H{
		"comparison_result": result,
		"study_plan_info": gin.H{
			"id":      studyPlan.ID,
			"version": studyPlan.Version,
			"career":  studyPlan.Career.Name,
		},
		"summary": gin.H{
			"total_subjects_in_plan":     len(result.EquivalentSubjects) + len(result.MissingSubjects),
			"approved_subjects":          len(result.EquivalentSubjects),
			"missing_subjects":           len(result.MissingSubjects),
			"completion_percentage":      calculateCompletionPercentage(result.CreditsSummary),
		},
	})
}

// calculateCompletionPercentage calcula el porcentaje de completitud basado en créditos
func calculateCompletionPercentage(summary models.CreditsSummary) float64 {
	if summary.Total.Required == 0 {
		return 0.0
	}
	return (float64(summary.Total.Completed) / float64(summary.Total.Required)) * 100.0
}

// APICompareRequest estructura para la solicitud de comparación desde texto
type APICompareRequest struct {
	AcademicHistoryText string `json:"academic_history_text" binding:"required"`
	TargetCareerCode    string `json:"target_career_code" binding:"required"`
}

// ParsedSubject representa una materia extraída del texto de historia académica
type ParsedSubject struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Credits     int     `json:"credits"`
	Type        string  `json:"type"`
	Grade       float64 `json:"grade"`
	Status      string  `json:"status"`
	Semester    string  `json:"semester"`
}

// parseAcademicHistoryText extrae las materias de la historia académica en texto
func parseAcademicHistoryText(text string) ([]ParsedSubject, error) {
	var subjects []ParsedSubject
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "(") && strings.Contains(line, ")") {
			subject, err := parseSubjectLineUltraTolerant(line)
			if err == nil {
				subjects = append(subjects, subject)
			}
		}
	}
	return subjects, nil
}

// Versión ultra tolerante del parser de materias
func parseSubjectLineUltraTolerant(line string) (ParsedSubject, error) {
	codeStart := strings.Index(line, "(")
	codeEnd := strings.Index(line, ")")
	if codeStart == -1 || codeEnd == -1 || codeEnd <= codeStart {
		return ParsedSubject{}, errors.New("código no encontrado")
	}
	code := line[codeStart+1 : codeEnd]
	name := strings.TrimSpace(line[:codeStart])
	remaining := strings.TrimSpace(line[codeEnd+1:])
	parts := strings.Fields(remaining)
	if len(parts) < 2 {
		return ParsedSubject{}, errors.New("información insuficiente")
	}
	reNum := regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	if !reNum.MatchString(parts[0]) {
		return ParsedSubject{}, errors.New("créditos no es un número válido: " + parts[0])
	}
	creditsFloat, _ := strconv.ParseFloat(parts[0], 64)
	credits := int(creditsFloat)
	subjectType := determineSubjectType(parts[1:])
	var grade float64
	for _, part := range parts {
		if g, err := strconv.ParseFloat(part, 64); err == nil && g >= 0.0 && g <= 5.0 {
			grade = g
			break
		}
	}
	status := "APROBADA"
	if strings.Contains(strings.ToUpper(line), "REPROBADA") {
		status = "REPROBADA"
	} else if strings.Contains(strings.ToUpper(line), "EN CURSO") {
		status = "EN CURSO"
	}
	var semester string
	for _, part := range parts {
		if strings.Contains(part, "-") && len(part) >= 6 {
			semester = part
			break
		}
	}
	return ParsedSubject{
		Code:     code,
		Name:     name,
		Credits:  credits,
		Type:     subjectType,
		Grade:    grade,
		Status:   status,
		Semester: semester,
	}, nil
}

// determineSubjectType determina el tipo de materia basándose en palabras clave
func determineSubjectType(parts []string) string {
	text := strings.Join(parts, " ")
	text = strings.ToUpper(text)
	
	if strings.Contains(text, "FUND. OBLIGATORIA") || strings.Contains(text, "FUNDAMENTACIÓN OBLIGATORIA") {
		return "FUND. OBLIGATORIA"
	}
	if strings.Contains(text, "FUND. OPTATIVA") || strings.Contains(text, "FUNDAMENTACIÓN OPTATIVA") {
		return "FUND. OPTATIVA"
	}
	if strings.Contains(text, "DISCIPLINAR OBLIGATORIA") {
		return "DISCIPLINAR OBLIGATORIA"
	}
	if strings.Contains(text, "DISCIPLINAR OPTATIVA") {
		return "DISCIPLINAR OPTATIVA"
	}
	if strings.Contains(text, "LIBRE ELECCIÓN") {
		return "LIBRE ELECCIÓN"
	}
	if strings.Contains(text, "TRABAJO DE GRADO") {
		return "TRABAJO DE GRADO"
	}
	if strings.Contains(text, "NIVELACIÓN") {
		return "NIVELACIÓN"
	}
	
	return "LIBRE ELECCIÓN" // Por defecto
}

// Limpieza y normalización del texto de historia académica
func preprocessAcademicHistoryText(raw string) string {
	// 1. Reemplazar saltos de línea de Windows por Unix
	cleaned := strings.ReplaceAll(raw, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", "\n")
	// 2. Reemplazar múltiples saltos de línea por uno solo
	cleaned = regexp.MustCompile(`\n+`).ReplaceAllString(cleaned, "\n")
	// 3. Quitar espacios en blanco al inicio y final de cada línea
	lines := strings.Split(cleaned, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	cleaned = strings.Join(lines, "\n")
	// 4. Quitar espacios en blanco al inicio y final del texto
	cleaned = strings.TrimSpace(cleaned)
	return cleaned
}

// compareAcademicHistoryFromText compara historia académica en texto con el pensum
func compareAcademicHistoryFromText(c *gin.Context) {
	var academicHistoryText, targetCareerCode string

	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		var req APICompareRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos: " + err.Error()})
			return
		}
		academicHistoryText = req.AcademicHistoryText
		targetCareerCode = req.TargetCareerCode
	} else if strings.HasPrefix(contentType, "multipart/form-data") || strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		// Leer desde form-data o x-www-form-urlencoded
		academicHistoryText = c.PostForm("academic_history_text")
		targetCareerCode = c.PostForm("target_career_code")
		if academicHistoryText == "" || targetCareerCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan campos en el formulario: academic_history_text y target_career_code son requeridos"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type no soportado. Usa application/json o form-data."})
		return
	}

	// Limpieza y normalización del texto
	cleanedText := preprocessAcademicHistoryText(academicHistoryText)

	// Parsear la historia académica del texto limpio
	parsedSubjects, err := parseAcademicHistoryText(cleanedText)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parseando historia académica: " + err.Error()})
		return
	}

	// Convertir a formato de entrada de la API
	var subjects []models.SubjectInput
	for _, ps := range parsedSubjects {
		subject := models.SubjectInput{
			Code:     ps.Code,
			Name:     ps.Name,
			Credits:  ps.Credits,
			Type:     models.TipologiaAsignatura(ps.Type),
			Grade:    ps.Grade,
			Status:   ps.Status,
			Semester: ps.Semester,
		}
		subjects = append(subjects, subject)
	}

	academicHistory := models.AcademicHistoryInput{
		CareerCode: targetCareerCode,
		Subjects:   subjects,
	}

	// Realizar la comparación
	result, err := functions.CompareAcademicHistoryByCareerCode(config.DB, academicHistory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Obtener información del plan de estudio usado
	studyPlan, _ := functions.GetStudyPlanByCareerCode(config.DB, targetCareerCode)

	c.JSON(http.StatusOK, gin.H{
		"parsed_subjects": parsedSubjects,
		"comparison_result": result,
		"study_plan_info": gin.H{
			"id":      studyPlan.ID,
			"version": studyPlan.Version,
			"career":  studyPlan.Career.Name,
		},
		"summary": gin.H{
			"total_subjects_parsed":     len(parsedSubjects),
			"total_subjects_in_plan":    len(result.EquivalentSubjects) + len(result.MissingSubjects),
			"approved_subjects":         len(result.EquivalentSubjects),
			"missing_subjects":          len(result.MissingSubjects),
			"completion_percentage":     calculateCompletionPercentage(result.CreditsSummary),
		},
	})
}
