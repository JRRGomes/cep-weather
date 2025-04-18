package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock para a função getCEPInfo
func mockGetCEPInfo(cep string) (*ViaCEPResponse, error) {
	if cep == "01001000" {
		return &ViaCEPResponse{
			CEP:        "01001-000",
			Logradouro: "Praça da Sé",
			Bairro:     "Sé",
			Localidade: "São Paulo",
			UF:         "SP",
			IBGE:       "3550308",
			Erro:       false,
		}, nil
	}
	return nil, fmt.Errorf("CEP não encontrado")
}

// Mock para a função getWeather
func mockGetWeather(city string) (float64, error) {
	if city == "São Paulo" {
		return 25.5, nil
	}
	return 0, fmt.Errorf("Cidade não encontrada")
}

func TestWeatherHandlerValidCEP(t *testing.T) {
	originalGetCEPInfo := getCEPInfo
	originalGetWeather := getWeather
	defer func() {
		getCEPInfo = originalGetCEPInfo
		getWeather = originalGetWeather
	}()

	getCEPInfo = mockGetCEPInfo
	getWeather = mockGetWeather

	// Criar request para teste
	req, err := http.NewRequest("GET", "/weather/01001000", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Criar ResponseRecorder para registrar a resposta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(weatherHandler)

	// Executar o handler
	handler.ServeHTTP(rr, req)

	// Verificar o status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler retornou status code incorreto: recebido %v esperado %v",
			status, http.StatusOK)
	}

	// Verificar o tipo de conteúdo
	expected := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expected {
		t.Errorf("tipo de conteúdo incorreto: recebido %v esperado %v",
			contentType, expected)
	}

	// Verificar o corpo da resposta
	var response TemperatureResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Erro ao decodificar resposta JSON: %v", err)
	}

	// Verificar valores de temperatura
	if response.TempC != 25.5 {
		t.Errorf("Temperatura em Celsius incorreta: recebida %.2f esperada %.2f", response.TempC, 25.5)
	}

	expectedTempF := 25.5*1.8 + 32
	if response.TempF != expectedTempF {
		t.Errorf("Temperatura em Fahrenheit incorreta: recebida %.2f esperada %.2f", response.TempF, expectedTempF)
	}

	expectedTempK := 25.5 + 273.15
	if response.TempK != expectedTempK {
		t.Errorf("Temperatura em Kelvin incorreta: recebida %.2f esperada %.2f", response.TempK, expectedTempK)
	}
}

func TestWeatherHandlerInvalidCEP(t *testing.T) {
	// Criar request para teste com CEP inválido
	req, err := http.NewRequest("GET", "/weather/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Criar ResponseRecorder para registrar a resposta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(weatherHandler)

	// Executar o handler
	handler.ServeHTTP(rr, req)

	// Verificar o status code
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler retornou status code incorreto: recebido %v esperado %v",
			status, http.StatusUnprocessableEntity)
	}

	// Verificar o corpo da resposta
	var response ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Erro ao decodificar resposta JSON: %v", err)
	}

	expected := "invalid zipcode"
	if response.Message != expected {
		t.Errorf("Mensagem de erro incorreta: recebida %v esperada %v",
			response.Message, expected)
	}
}

func TestWeatherHandlerNotFoundCEP(t *testing.T) {
	// Substituir as funções reais pelos mocks
	originalGetCEPInfo := getCEPInfo
	defer func() {
		getCEPInfo = originalGetCEPInfo
	}()

	getCEPInfo = mockGetCEPInfo

	// Criar request para teste com CEP não encontrado
	req, err := http.NewRequest("GET", "/weather/99999999", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Criar ResponseRecorder para registrar a resposta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(weatherHandler)

	// Executar o handler
	handler.ServeHTTP(rr, req)

	// Verificar o status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler retornou status code incorreto: recebido %v esperado %v",
			status, http.StatusNotFound)
	}

	// Verificar o corpo da resposta
	var response ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Erro ao decodificar resposta JSON: %v", err)
	}

	expected := "can not find zipcode"
	if response.Message != expected {
		t.Errorf("Mensagem de erro incorreta: recebida %v esperada %v",
			response.Message, expected)
	}
}
