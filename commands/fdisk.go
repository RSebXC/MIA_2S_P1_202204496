package commands

import (
	structures "MIA_202204496/structures" // Importa el paquete "structures" desde el directorio "EDD2021/structures"
	utils "MIA_202204496/utils"
	"errors"  // Paquete para manejar errores y crear nuevos errores con mensajes personalizados
	"fmt"     // Paquete para formatear cadenas y realizar operaciones de entrada/salida
	"regexp"  // Paquete para trabajar con expresiones regulares, útil para encontrar y manipular patrones en cadenas
	"strconv" // Paquete para convertir cadenas a otros tipos de datos, como enteros
	"strings" // Paquete para manipular cadenas, como unir, dividir, y modificar contenido de cadenas
)

// FDISK estructura que representa el comando fdisk con sus parámetros
type FDISK struct {
	size int    // Tamaño de la partición
	unit string // Unidad de medida del tamaño (K o M)
	fit  string // Tipo de ajuste (BF, FF, WF)
	path string // Ruta del archivo del disco
	typ  string // Tipo de partición (P, E, L)
	name string // Nombre de la partición
}

/*
	fdisk -size=1 -type=L -unit=M -fit=BF -name="Particion3" -path="/home/keviin/University/PRACTICAS/MIA_LAB_S2_2024/CLASEEXTRA/disks/Disco1.mia"
	fdisk -size=300 -path=/home/Disco1.mia -name=Particion1
	fdisk -type=E -path=/home/Disco2.mia -Unit=K -name=Particion2 -size=300
*/

// CommandFdisk parsea el comando fdisk y devuelve una instancia de FDISK
func ParserFdisk(tokens []string) (*FDISK, error) {
	cmd := &FDISK{} // Crea una nueva instancia de FDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando fdisk
	re := regexp.MustCompile(`-size=\d+|-unit=[kKmM]|-fit=[bBfF]{2}|-path="[^"]+"|-path=[^\s]+|-type=[pPeElL]|-name="[^"]+"|-name=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-size":
			// Convierte el valor del tamaño a un entero
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return nil, errors.New("el tamaño debe ser un número entero positivo")
			}
			cmd.size = size
		case "-unit":
			// Verifica que la unidad sea "K" o "M"
			if value != "K" && value != "M" {
				return nil, errors.New("la unidad debe ser K o M")
			}
			cmd.unit = strings.ToUpper(value)
		case "-fit":
			// Verifica que el ajuste sea "BF", "FF" o "WF"
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return nil, errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.fit = value
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-type":
			// Verifica que el tipo sea "P", "E" o "L"
			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return nil, errors.New("el tipo debe ser P, E o L")
			}
			cmd.typ = value
		case "-name":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return nil, errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size, -path y -name hayan sido proporcionados
	if cmd.size == 0 {
		return nil, errors.New("faltan parámetros requeridos: -size")
	}
	if cmd.path == "" {
		return nil, errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return nil, errors.New("faltan parámetros requeridos: -name")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if cmd.unit == "" {
		cmd.unit = "M"
	}

	// Si no se proporcionó el ajuste, se establece por defecto a "FF"
	if cmd.fit == "" {
		cmd.fit = "WF"
	}

	// Si no se proporcionó el tipo, se establece por defecto a "P"
	if cmd.typ == "" {
		cmd.typ = "P"
	}

	// Crear la partición con los parámetros proporcionados
	err := commandFdisk(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return cmd, nil // Devuelve el comando FDISK creado
}

func commandFdisk(fdisk *FDISK) error {
	// Convertir el tamaño a bytes
	sizeBytes, err := utils.ConvertToBytes(fdisk.size, fdisk.unit)
	if err != nil {
		fmt.Println("Error al convertir el tamaño:", err)
		return err
	}

	// Deserializar el MBR para obtener la estructura actual
	var mbr structures.MBR
	err = mbr.Deserialize(fdisk.path)
	if err != nil {
		fmt.Println("Error al deserializar el MBR:", err)
		return err
	}

	switch fdisk.typ {
	case "P":
		// Crear partición primaria
		err = createPrimaryPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error al crear la partición primaria:", err)
			return err
		}
	case "E":
		// Crear partición extendida
		err = createExtendedPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error al crear la partición extendida:", err)
			return err
		}
	case "L":
		// Obtener la partición extendida
		extendedPartition, _, _ := mbr.GetPartitionByType('E')
		if extendedPartition == nil {
			fmt.Println("Error: No se encontró una partición extendida.")
			return fmt.Errorf("no se encontró una partición extendida")
		}

		// Crear partición lógica
		err = createLogicalPartition(extendedPartition, int32(sizeBytes), fdisk.path, fdisk.name)
		if err != nil {
			fmt.Println("Error al crear la partición lógica:", err)
			return err
		}
	default:
		// Tipo de partición no válido
		fmt.Println("Error: Tipo de partición no válido. Use 'P' para primaria, 'E' para extendida, o 'L' para lógica.")
		return fmt.Errorf("tipo de partición no válido: %s", fdisk.typ)
	}

	fmt.Println("Partición creada exitosamente.")
	return nil
}
func createPrimaryPartition(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Obtener la primera partición disponible
	availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
	if availablePartition == nil {
		fmt.Println("No hay particiones disponibles.")
	}

	// Crear la partición con los parámetros proporcionados
	availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

	// Colocar la partición en el MBR
	if availablePartition != nil {
		mbr.Mbr_partitions[indexPartition] = *availablePartition
	}

	// Serializar el MBR en el archivo binario
	err = mbr.Serialize(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Partición primaria creada exitosamente.")
	return nil
}

func createExtendedPartition(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Verificar si ya existe una partición extendida
	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_type[0] == 'E' {
			fmt.Println("Ya existe una partición extendida en el disco.")
			return errors.New("partición extendida ya existe")
		}
	}

	// Obtener la primera partición disponible para la partición extendida
	availablePartition, startPartition, _ := mbr.GetFirstAvailablePartition()
	if availablePartition == nil {
		fmt.Println("No hay particiones disponibles.")
		return nil
	}

	// Crear la partición extendida en la tabla de particiones del MBR
	availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

	// Serializar el MBR en el archivo binario
	err = mbr.Serialize(fdisk.path)
	if err != nil {
		fmt.Println("Error serializando el MBR:", err)
		return err
	}

	// Crear una estructura EBR para la partición extendida
	var ebr structures.EBR
	ebr.Ebr_start = int32(startPartition)
	ebr.Ebr_size = 0  // Inicialmente sin tamaño para partición lógica
	ebr.Ebr_next = -1 // No hay siguiente EBR
	copy(ebr.Ebr_fit[:], fdisk.fit)
	copy(ebr.Ebr_name[:], fdisk.name)

	// Serializar el EBR en el archivo
	err = ebr.Serialize(fdisk.path, int32(startPartition))
	if err != nil {
		fmt.Println("Error serializando el EBR:", err)
		return err
	}

	fmt.Println("Partición extendida creada exitosamente.")
	return nil
}

func createLogicalPartition(partition *structures.Partition, size int32, path string, name string) error {
	// Verificar que la partición sea extendida
	if partition.Part_type[0] != 'E' {
		return errors.New("la partición no es extendida")
	}

	// Leer el primer EBR de la partición extendida
	var currentEBR structures.EBR
	err := currentEBR.Deserialize(path, partition.Part_start)
	if err != nil {
		return err
	}
	fmt.Printf("EBR inicial deserializado: %+v\n", currentEBR)

	// Buscar el último EBR en la lista
	for currentEBR.Ebr_next != -1 {
		fmt.Printf("Leyendo EBR siguiente en la posición: %d\n", currentEBR.Ebr_next)
		err = currentEBR.Deserialize(path, currentEBR.Ebr_next)
		if err != nil {
			return err
		}
	}

	// Calcular el inicio de la nueva partición lógica
	newEBRStart := currentEBR.Ebr_start + currentEBR.Ebr_size
	fmt.Printf("Calculado inicio del nuevo EBR: %d\n", newEBRStart)

	// Verificar que haya espacio suficiente para la nueva partición
	if newEBRStart+size > partition.Part_start+partition.Part_size {
		return errors.New("no hay suficiente espacio en la partición extendida")
	}

	// Crear el nuevo EBR
	newEBR := &structures.EBR{
		Ebr_start: newEBRStart,
		Ebr_size:  size,
		Ebr_next:  -1,
	}
	copy(newEBR.Ebr_fit[:], currentEBR.Ebr_fit[:])

	// Asignar el nombre proporcionado al nuevo EBR
	copy(newEBR.Ebr_name[:], name)
	fmt.Printf("Nombre asignado al nuevo EBR: %s\n", name)
	fmt.Printf("Nuevo EBR creado: %+v\n", newEBR)

	// Actualizar el EBR actual para que apunte al nuevo EBR
	currentEBR.Ebr_next = newEBRStart
	fmt.Printf("Actualizando EBR actual para apuntar al nuevo EBR en: %d\n", currentEBR.Ebr_next)

	// Serializar el EBR actualizado en su lugar
	err = currentEBR.Serialize(path, currentEBR.Ebr_start)
	if err != nil {
		return err
	}

	// Serializar el nuevo EBR
	err = newEBR.Serialize(path, newEBRStart)
	if err != nil {
		return err
	}

	fmt.Println("Partición lógica creada exitosamente.")
	return nil
}
