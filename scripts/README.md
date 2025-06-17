# Scripts de Población de Datos

Este directorio contiene scripts para poblar la base de datos con datos del pensum de Ingeniería de Sistemas y sus equivalencias.

## Archivos

- `create_career.go` - Verifica la carrera de Ingeniería de Sistemas existente (código ISIS)
- `seed_ing_sistemas.go` - Puebla el pensum completo con materias y equivalencias
- `run_seed.go` - Script principal que ejecuta ambos scripts en secuencia
- `run_complete_seed.sh` - **Script completo recomendado** (configura entorno + ejecuta todo)
- `setup_env.sh` - Configura variables de entorno desde el archivo .env
- `go.mod` - Dependencias de Go para los scripts

## Uso

### Opción 1: Ejecutar todo el proceso (Recomendado)

```bash
cd scripts
./run_complete_seed.sh
```

Este script:
1. Configura las variables de entorno desde el archivo `.env`
2. Verifica que existe la carrera ISIS (ID: 1)
3. Puebla el pensum completo
4. Limpia archivos temporales

### Opción 2: Ejecutar scripts individualmente

1. Configurar variables de entorno:
```bash
cd scripts
source ./setup_env.sh
```

2. Verificar la carrera:
```bash
go run create_career.go
```

3. Poblar el pensum:
```bash
go run seed_ing_sistemas.go
```

## Configuración de Base de Datos

Los scripts están configurados para usar **Supabase** con las credenciales del archivo `.env`:

- **Host**: aws-0-us-east-2.pooler.supabase.com
- **Usuario**: postgres.pptscezwliowljhgqocx
- **Base de datos**: postgres
- **SSL**: Requerido

El script `setup_env.sh` lee automáticamente el archivo `.env` y configura la variable `DATABASE_URL`.

## Carrera Existente

Los scripts están diseñados para trabajar con la carrera de Ingeniería de Sistemas que ya existe en Supabase:

- **ID**: 1
- **Código**: ISIS
- **Nombre**: Ingeniería de Sistemas

## Datos Incluidos

### Plan de Estudio
- **Versión**: 2023-1
- **Total de materias**: 65
- **Tipologías**:
  - Fundamentación Obligatoria: 7 materias (28 créditos)
  - Fundamentación Optativa: 15 materias (58 créditos)
  - Disciplinar Obligatoria: 18 materias (54 créditos)
  - Disciplinar Optativa: 25 materias (75 créditos)

### Equivalencias
Se incluyen 12 equivalencias entre materias antiguas y nuevas, específicas para Ingeniería de Sistemas:

```json
{
  "3006914": ["3010651"],
  "3007742": ["3010435"],
  "3007743": ["3010426"],
  "3007855": ["3010476"],
  "3007849": ["3010440"],
  "3007322": ["3010415"],
  "3007746": ["3010114"],
  "3009550": ["3007862"],
  "3007844": ["3010408"],
  "3007845": ["3010407"],
  "3007846": ["3010439"],
  "3008883": ["3011166"]
}
```

## Estructura de Datos

Los datos están basados en el prototipo Next.js y siguen la estructura:

```typescript
// Materia
{
  codigo: "1000004-M",
  nombre: "Cálculo Diferencial",
  creditos: 4,
  tipologia: "B - Fundamentación Obligatoria"
}
```

## Notas

- Los scripts son idempotentes: pueden ejecutarse múltiples veces sin duplicar datos
- Si una materia ya existe, se reutiliza en lugar de crear una nueva
- Las equivalencias se crean automáticamente basándose en los datos del prototipo
- Todos los créditos se calculan automáticamente por tipología
- El script verifica automáticamente la existencia de la carrera ISIS antes de proceder 