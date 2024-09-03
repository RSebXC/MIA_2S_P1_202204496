package structures

import (
	"encoding/binary"
	"fmt" // Paquete para imprimir mensajes de depuración
	"os"
)

type EBR struct {
	Ebr_mount [1]byte
	Ebr_fit   [1]byte
	Ebr_start int32
	Ebr_size  int32
	Ebr_next  int32
	Ebr_name  [16]byte
}

// Inicialización con valores nulos
func NewEBR() *EBR {
	return &EBR{
		Ebr_mount: [1]byte{0},
		Ebr_fit:   [1]byte{0},
		Ebr_start: -1, // Indica que no está en uso
		Ebr_size:  0,
		Ebr_next:  -1, // Indica que no hay siguiente EBR
		Ebr_name:  [16]byte{},
	}
}

// Serializa el EBR al archivo en la ubicación especificada
func (ebr *EBR) Serialize(path string, start int32) error {
	// Abre el archivo en modo escritura
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Moverse a la posición donde debe comenzar el EBR
	_, err = file.Seek(int64(start), 0)
	if err != nil {
		return err
	}

	fmt.Printf("Serializando EBR en la posición %d del archivo %s\n", start, path)

	// Escribir el EBR en esa posición
	err = binary.Write(file, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}

	fmt.Printf("EBR serializado con éxito en la posición %d. Detalles del EBR: %+v\n", start, ebr)
	return nil
}

// Deserializa el EBR desde el archivo en la ubicación especificada
func (ebr *EBR) Deserialize(path string, start int32) error {
	// Abre el archivo en modo lectura
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Moverse a la posición donde debe comenzar el EBR
	_, err = file.Seek(int64(start), 0)
	if err != nil {
		return err
	}

	fmt.Printf("Deserializando EBR de la posición %d del archivo %s\n", start, path)

	// Leer los datos del EBR desde esa posición
	err = binary.Read(file, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}

	fmt.Printf("EBR deserializado con éxito desde la posición %d. Detalles del EBR: %+v\n", start, ebr)
	return nil
}
