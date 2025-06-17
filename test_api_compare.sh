#!/bin/bash

# Script para probar la API de comparación de historia académica

echo "🧪 Probando API de comparación de historia académica..."
echo ""

# Historia académica del estudiante (Ingeniería de Minas y Metalurgia)
ACADEMIC_HISTORY='Portal de Servicios Académicos

joroblesr

Datos personales

Información académica

Mi historia académica
Mis Calificaciones
Mi horario
Mis planes
Mis tutores

Proceso de inscripción

Buscador de cursos

Catálogo prog. curriculares

Información Financiera

Trámites y solicitudes

Evaluación docente

Historia Académica

Plan de estudios

INGENIERÍA DE MINAS Y METALURGIA
Facultad: FACULTAD DE MINASHist. Acad.: 202ESTADO BLOQUEADOCausas de bloqueo: B - 41 Readmisión Res.235 de 2009 Vic. Académica
Resumen
4.1 (Acumulado)Pregrado - Promedio académico2021-2S
4.1 (Acumulado)Pregrado - P.A.P.A2021-2S
Asignaturas

Asignaturas
Créditos
Tipo
Periodo
Calificación
/api-compare
Fundamentos de programación (3010435)
3
FUND. OBLIGATORIA
2021-2S Ordinaria
4.6
APROBADA
ÁLGEBRA LINEAL (1000003-M)
4
FUND. OBLIGATORIA
2021-1S Ordinaria
4.0
APROBADA
CÁLCULO INTEGRAL (1000005-M)
4
FUND. OBLIGATORIA
2021-1S Ordinaria
3.5
APROBADA
FÍSICA MECÁNICA (1000019-M)
4
FUND. OBLIGATORIA
2021-1S Ordinaria
3.7
APROBADA
Cátedra estudiantil: universidad, participación y sociedad (3010348)
3
LIBRE ELECCIÓN
2021-1S Ordinaria
4.5
APROBADA
CÁLCULO DIFERENCIAL (1000004-M)
4
FUND. OBLIGATORIA
2020-2S Ordinaria
4.9
APROBADA
GEOMETRÍA VECTORIAL Y ANALÍTICA (1000008-M)
4
FUND. OBLIGATORIA
2020-2S Ordinaria
3.6
APROBADA
CIENCIA DE LOS MATERIALES (3007309)
3
FUND. OPTATIVA
2020-2S Ordinaria
3.5
APROBADA
Introducción a la Ingeniería de Minas y Metalurgia (3007476)
1
DISCIPLINAR OBLIGATORIA
2020-2S Ordinaria
4.6
APROBADA
ECOLOGÍA GENERAL (3007022)
3
DISCIPLINAR OPTATIVA
2020-1S Ordinaria
4.4
APROBADA
Química general (3006829)
3
FUND. OBLIGATORIA
2020-1S Ordinaria
4.3
APROBADA
LECTO-ESCRITURA (1000002-M)
4
NIVELACIÓN
2020-1S Ordinaria
4.1
APROBADA
MATEMÁTICAS BÁSICAS (1000001-M)
4
NIVELACIÓN
2020-1S Ordinaria
4.1
APROBADA
Cátedra nacional de inducción y preparación para la vida universitaria (1000089-O)
2
LIBRE ELECCIÓN
2020-1S Ordinaria
APROBADA
INGLÉS I (1000044-M)
3
NIVELACIÓN
2020-1S Validacion por suficiencia
APROBADA
INGLÉS II (1000045-M)
3
NIVELACIÓN
2020-1S Validacion por suficiencia
APROBADA

Resumen de créditos

Tipologías
Exigidos
Aprobados
Pendientes
Inscritos
Cursados

DISCIPLINAR OPTATIVA
21
3
18
0
3
FUND. OBLIGATORIA
29
26
3
0
26
FUND. OPTATIVA
16
3
13
0
3
DISCIPLINAR OBLIGATORIA
72
1
71
0
1
LIBRE ELECCIÓN
36
5
31
0
5
TRABAJO DE GRADO
6
0
6
0
0
TOTAL
180
38
142
0
38
NIVELACIÓN
20
14
6
0
14
TOTAL ESTUDIANTE
200
52
148
0
52

Total Créditos Excedentes0

Total de Créditos Cancelados en los Periodos Cursado0

Porcentaje de Avance21,1%

Cupo de créditos

Créditos adicionales80Cupo de créditos228Créditos disponibles80Créditos de estudio doble titulación80

Universidad Nacional de Colombia--Dirección Nacional de Información Académica
Portal de Servicios Académicos (V. 4.3.21) | Todos los derechos reservados'

# Escapar caracteres especiales para JSON
ACADEMIC_HISTORY_ESCAPED=$(echo "$ACADEMIC_HISTORY" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')

echo "📤 Enviando petición a la API..."
echo ""

# Crear el JSON para la petición usando jq para manejar correctamente el escape
JSON_DATA=$(jq -n \
  --arg history "$ACADEMIC_HISTORY" \
  --arg career "ISIS" \
  '{
    "academic_history_text": $history,
    "target_career_code": $career
  }')

# Realizar la petición POST
curl -X POST http://localhost:8080/api/api-compare \
  -H "Content-Type: application/json" \
  -d "$JSON_DATA" \
  | jq '.'

echo ""
echo "✅ Prueba completada" 