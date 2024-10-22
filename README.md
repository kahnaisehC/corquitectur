# Corquitecture 
## Idioma
[Español](#descripción-general)/[English](#overview)

## Descripción general

Este proyecto fue creado para competir en **"Cumbre de Datos 2024: LA REVOLUCIÓN DE LA IA"**.  
Corquitecture permite a los usuarios subir un archivo `.csv` a través de un formulario web, y la aplicación lo analiza automáticamente e inserta su contenido en una base de datos PostgreSQL. La aplicación crea una tabla en la base de datos basada en los encabezados del archivo CSV, y cada fila del CSV se almacena como una fila en la base de datos.  
El objetivo del proyecto es facilitar el uso de archivos `.csv` creando una arquitectura que sirva como respaldo para los conjuntos de datos "colgantes", otorgándoles seguridad de tipos y vinculándolos para crear relaciones entre las tablas.

## Funcionalidades

- Interfaz web para subir archivos CSV.
- Crea automáticamente una tabla en PostgreSQL basada en los encabezados del CSV.
- Inserta los datos del CSV en la tabla correspondiente en PostgreSQL.
- Filtrado de datos incorrectos.
- API en formato JSON para consumir las tablas.
- Visualización completa y listado de las tablas.

## Requisitos previos

Antes de ejecutar la aplicación web, asegúrate de tener instalados los siguientes programas:

- [Go](https://golang.org/doc/install) (se recomienda v1.22 o superior)
- [PostgreSQL](https://www.postgresql.org/download/)

## Configuración

### 1. Instalar dependencias de Go

Asegúrate de tener instalados los módulos de Go necesarios para el proyecto. Ejecuta los siguientes comandos en el directorio raíz del proyecto:

```bash
cd corquitecture
go mod init 
go mod tidy
```

### 2. Configuración de PostgreSQL

Asegúrate de tener una instancia de PostgreSQL en ejecución y crea una base de datos para el proyecto:

```sql
CREATE DATABASE nombre_de_tu_base_de_datos;
```

También necesitas un usuario de PostgreSQL con los privilegios necesarios:

```sql
CREATE USER tu_usuario WITH PASSWORD 'tu_contraseña';
GRANT ALL PRIVILEGES ON DATABASE nombre_de_tu_base_de_datos TO tu_usuario;
```

### 3. Variables de entorno

La configuración de conexión a PostgreSQL se gestiona a través de variables de entorno. Estas se configuran en el archivo `.env`. Las variables requeridas son:

```bash
USERNAME="tu_usuario"
SECRET="tu_contraseña"
PORT="5432"
DB="nombre_de_tu_base_de_datos"
HOST="tu_host"
```

### 4. Ejecutar la aplicación web

Puedes ejecutar la aplicación web usando el siguiente comando:

```bash
go run main.go
```

El servidor web se iniciará en el puerto `3000` por defecto.

## Uso

1. Dirígete a la página (localhost:3000 si lo pruebas localmente) y carga un archivo CSV.  
2. Una vez que el archivo se haya cargado, puedes modificar el nombre y el tipo de cada columna con sus parámetros respectivos:

- **varchar**: este tipo tiene solo un parámetro:
    - longitud: la cantidad de caracteres del campo (por defecto 50).
- **integer**: este tipo tiene dos parámetros:
    - límite inferior: determina el valor mínimo posible para esa columna (por defecto 1).
    - límite superior: determina el valor máximo posible para esa columna (por defecto 100).
- **decimal/numeric**: basado en el tipo `Numeric` de PostgreSQL, tiene dos parámetros:
    - cantidad_de_dígitos: la cantidad máxima de dígitos que tendrá esta columna (por defecto 5).
    - posición_coma: la posición de la coma. Ten en cuenta que si se establece en un número negativo, la coma se moverá a la derecha (ejemplo: `decimal(5, -3)(12345)` == `12345000`) (por defecto 2).

3. ¡Dale a submit y busca el nombre de tu archivo entre las tablas!

## Problemas conocidos

1. **Detección automática de tipos**:
    - Las columnas se asignan como `varchar(50)` automáticamente en la interfaz, lo que puede causar errores de tipo si el usuario no tiene cuidado al usar la herramienta.
    - **Posible Solución**: Implementar una lógica para inferir o especificar explícitamente los tipos de columnas en función del contenido del CSV.

2. **Manejo de errores durante la subida**:
   - No hay retroalimentación detallada al usuario si ocurre un error al procesar el archivo CSV. Los usuarios solo ven un mensaje de error genérico.
   - **Posible solución**: Mejorar el sistema de reporte de errores para mostrar a los usuarios razones específicas del fallo (por ejemplo, problemas de conexión a la base de datos o errores de formato de archivo).

3. **Comportamiento errático con el tipo decimal**:
    - Se han encontrado errores al usar el tipo _decimal_ en la interfaz. Por ejemplo, si ingresas (0003.43) en un campo, se filtrará como un dato incorrecto.
    - **Posible solución**: Mejorar la lógica en la verificación del tipo decimal.

4. **Soporte limitado para tipos de datos**:
    - Solo hay tres tipos de datos disponibles en el cargador.
    - **Posible solución**: Añadir soporte para más tipos de datos en futuras actualizaciones.

## Contribuciones

Siéntete libre de hacer un fork de este repositorio y crear pull requests. ¡Las sugerencias para nuevas características u optimizaciones son bienvenidas!

---

## Overview

This project was created to compete in "Cumbre de Datos 2024: LA REVOLUCIÓN DE LA IA"
Corquitecture allows users to upload a `.csv` file through a web form, and the application automatically parses the file and inserts its contents into a PostgreSQL database. The app creates a database table based on the headers of the CSV file, and each row in the CSV is stored as a row in the database.
The objective of the project is to facilitate the usage of the `.csv` files while creating an architechture to backbone the "hanging" datasets by giving it type-safety and linking them in order to create relations between the tables.

## Features

- Web form interface to upload CSV files.
- Automatically creates PostgreSQL table based on CSV headers.
- Inserts CSV data into the corresponding PostgreSQL table.
- Broken datum filtering
- JSON API to consume the tables
- Full table display and listing of the tables

## Prerequisites

Before running the web application, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (v1.22+ recommended)
- [PostgreSQL](https://www.postgresql.org/download/)

## Setup

### 1. Install Go Dependencies

Ensure you have the required Go modules installed for the project. Run the following commands in the project root:

```bash
cd corquitecture
go mod init 
go mod tidy
```

### 2. PostgreSQL Setup

Ensure that you have a running PostgreSQL instance and create a database for the project:

```sql
CREATE DATABASE your_database_name;
```

You also need a PostgreSQL user with the necessary privileges:

```sql
CREATE USER your_username WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE your_database_name TO your_username;
```

### 3. Environment Variables

The PostgreSQL connection settings are managed through environment variables. These are configured in the `.env` file. The required variables are: 

```bash
USERNAME="your_username"
SECRET="your_password"
PORT="5432"
DB="your_database_name"
HOST="your_host"
```

### 4. Running the Web Application

You can run the web application using the following command:

```bash
go run main.go
```

The web server will start on port `3000` by default.

## Usage

1. Go to the page (localhost:3000 if trying it locally) and upload a csv File.
2. Once the file has been loaded, the name and the type of each column could be modified with its respective parameters: 

- **varchar**: this type has only one parameter:
    -  length: the amount of characters of the field (default 50).
- **integer**: this type has two parameters:
    - lower-bound: determines the lowest possible value of that specific column (default 1).
    - upper-bound: determines the highest possible value of that specific column (default 100).
- **decimal/numeric**: based on the Numeric type of Postgres, has two parameters:
    - amount-of-digits: the maximum amount of digits that this column will have (default 5)
    - comma-position: the position of the comma. Note that if it is set to a negative number will translate the comma to the right (decimal(5, -3)(12345) == 12345000) (default 2)

3. Click submit and search for the name of your file between the tables!

## Known Issues

1. **Automatic type detection**:
    - Columns are asigned varchar(50) automatically in the interface which may create some type errors along the way if the user is not carefull while using this tool
    - **Posible Solution**: Implement logic to infer or explicitly specify column data types based on CSV content.

2. **Error Handling During Upload**:
   - There is no detailed user feedback if an error occurs while processing the CSV file. Users only see a generic error message.
   - **Posible Solution**: Improve the error reporting system to show users specific reasons for failure (e.g., database connection issues, file format problems).

3. **Buggy behaviour when the type is decimal**:
    - Bugs where find when using the _decimal_ type in the interface. For example if we put (0003.43) to some field it will be filtered as broken data.
    - **Posible Solution**: Improve the logic in the decimal typecheck.

4. **Low support for types**:
    - There are only three types available in the uploader. 
    - **Posible Solution**: Add support to more data types in the future!!

## Contributing

Feel free to fork this repository and create pull requests. Suggestions for new features or optimizations are welcome!


