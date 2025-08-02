package data

import (
	"bytes"
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"io"
	"log"
	"net/http"
)

// Helper para enviar datos a un endpoint (POST) --> Mando un struct como JSON
func EnviarDatos(url string, data any) error {
	//Convierte el struct(data) a un JSON
	jsonData, err := json.Marshal(data)

	//Si no pudo serializar, devuelvo error
	if err != nil {
		return err
	}

	//log.Printf("JSON a enviar a: %s", string(jsonData))

	//POST a la url con el JSON
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	//Verifico error
	if err != nil {
		return err
	}
	//Cierro la rta, salio bien
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Error closing body")
		}
	}(resp.Body)

	return nil
}

// Helper para recibir datos desde un endpoint (GET) --> Pasa de JSON a struct
func RecibirDatos(url string, data any) error {
	//Llamo al endpoint y verifico error
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("error closing body")
		}
	}(resp.Body)

	//Leo el contenido
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//Deserializacion del JSON y lo paso a data
	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	return nil
}

// Leer Body JSON recibidos por POST o PUT
// Deserializa JSON en un struct de Go.
func LeerJson(w http.ResponseWriter, r *http.Request, mensaje any) error {
	err := json.NewDecoder(r.Body).Decode(mensaje)
	if err != nil {
		logger.Error("Error al decodificar el mensaje: %s", err.Error())
		http.Error(w, "Error al decodificar mensaje", http.StatusBadRequest)
		return err
	}
	// logger.Info("Me llegó un mensaje: %+v", mensaje)

	return nil
}

// Enviar datos por POST y obtener la respuesta (útil si la respuesta es un JSON)
func EnviarDatosConRespuesta(url string, data any) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	return resp, nil // importante: quien la use debe hacer defer resp.Body.Close()
}

func EnviarDatosYRecibirRespuesta(url string, dataEnviar any, dataRecibir any) error {
	jsonData, err := json.Marshal(dataEnviar)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Respuesta: %s", string(body)) // útil para debug

	err = json.Unmarshal(body, dataRecibir)
	if err != nil {
		return err
	}

	return nil
}
