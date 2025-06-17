#!/bin/bash

# Script completo para poblar la base de datos con datos de Ingeniería de Sistemas

echo "🚀 Iniciando proceso completo de población de datos para Ingeniería de Sistemas..."
echo ""

# 1. Configurar variables de entorno
echo "📁 Paso 1: Configurando variables de entorno..."
source ./setup_env.sh

if [ $? -ne 0 ]; then
    echo "❌ Error configurando variables de entorno"
    exit 1
fi

echo ""

# 2. Verificar carrera existente
echo "📚 Paso 2: Verificando carrera de Ingeniería de Sistemas..."
go run create_career.go

if [ $? -ne 0 ]; then
    echo "❌ Error verificando la carrera"
    exit 1
fi

echo ""

# 3. Poblar el pensum
echo "📖 Paso 3: Poblando pensum de Ingeniería de Sistemas..."
go run seed_ing_sistemas.go

if [ $? -ne 0 ]; then
    echo "❌ Error poblando el pensum"
    exit 1
fi

echo ""

# 4. Limpiar archivos temporales
echo "🧹 Paso 4: Limpiando archivos temporales..."
if [ -f "career_info.txt" ]; then
    rm career_info.txt
    echo "✅ Archivo temporal eliminado"
fi

echo ""
echo "🎉 ¡Proceso completado exitosamente!"
echo "✅ La base de datos ha sido poblada con:"
echo "   - Verificación de carrera ISIS (ID: 1)"
echo "   - Plan de estudio 2023-1"
echo "   - Todas las materias del pensum (65 materias)"
echo "   - Equivalencias entre materias (12 equivalencias)"
echo ""
echo "📊 Resumen de datos:"
echo "   - Total créditos: 215"
echo "   - Fundamentación Obligatoria: 28 créditos"
echo "   - Fundamentación Optativa: 58 créditos"
echo "   - Disciplinar Obligatoria: 54 créditos"
echo "   - Disciplinar Optativa: 75 créditos" 