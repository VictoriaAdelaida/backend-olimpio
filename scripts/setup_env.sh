#!/bin/bash

# Script para configurar las variables de entorno desde el archivo .env

# Obtener el directorio del proyecto
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Cargar variables desde .env
if [ -f "$PROJECT_DIR/.env" ]; then
    echo "📁 Cargando variables de entorno desde $PROJECT_DIR/.env"
    
    # Leer el archivo .env y exportar las variables
    while IFS='=' read -r key value; do
        # Ignorar líneas vacías y comentarios
        if [[ ! -z "$key" && ! "$key" =~ ^# ]]; then
            export "$key=$value"
            echo "   ✅ $key configurado"
        fi
    done < "$PROJECT_DIR/.env"
    
    # Construir DATABASE_URL para Supabase
    export DATABASE_URL="host=$DB_HOST user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME port=5432 sslmode=require"
    
    echo "🔗 DATABASE_URL configurado para Supabase"
    echo "   Host: $DB_HOST"
    echo "   Usuario: $DB_USER"
    echo "   Base de datos: $DB_NAME"
    
else
    echo "❌ Archivo .env no encontrado en $PROJECT_DIR"
    exit 1
fi

echo ""
echo "🚀 Variables de entorno configuradas. Puedes ejecutar los scripts ahora." 