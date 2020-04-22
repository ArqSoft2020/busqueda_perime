# # busqueda-perime!

Microservicio de **Perime**. Se encarga de realizar la busquedad de las publicaciones con respecto a diferentes parametros


# Antes de iniciar

-instalar docker y docker-compose, configurar variables de entorno de go
-ejecutar docker-compose up 
-ejecutar comandos SQL en la base de datos antes de realizar cualquier operación con la API del microserivicio de busquedad. (ver siguiente item)


## Comandos SQL

- Luego de ejecutar el contenedor ejecutar:
	> docker exec -it perime-busqueda-db mysql -p	
- Introducir la contaseña, por decto del contenedor
	> password
- ejecutar las siguientes lineas de código:
use perime_search_db;
CREATE TABLE categorys ( id INT AUTO_INCREMENT PRIMARY KEY, Name_Category VARCHAR(50) NOT NULL, Type_Category VARCHAR(50) NOT NULL);


CREATE TABLE products ( id INT AUTO_INCREMENT PRIMARY KEY, Id_Category INT NOT NULL, Name_Product 
VARCHAR(50) NOT NULL, Description_Product VARCHAR(50) NOT NULL);

## Operaciones

Se listraran las operaciones:

|Categoria          |Ruta                        |Metodo                        |
|----------------|-------------------------------|-----------------------------|
|Obtener categorias|`/categorias`            |GET            |
|Crear una categoria          |`/categoria`            |POST          |
|Obtener una categoria     |`/categoria/{id:[0-9]}`|GET            |
|Modificar una categoria         |`/categoria/{id:[0-9]}`|PUT           |
|Eliminar una categoria       |`/categoria/{id:[0-9]}`|DELETE            |
  




|Producto         |Ruta                        |Metodo                        |
|----------------|-------------------------------|-----------------------------|
|Obtener productos|`/productos`            |GET            |
|Crear un producto          |`/producto`            |POST          |
|Obtener un producto   |`/producto/{id:[0-9]}`|GET            |
|Modificar un producto         |`/producto/{id:[0-9]}`|PUT           |
|Eliminar un producto      |`/producto/{id:[0-9]}`|DELETE            |


Repositorio en mi git personal: https://gitlab.com/kigama/busqueda_perime
