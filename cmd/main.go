package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
	Erro        bool   `json:"erro"`
}

type WeatherAPIResponse struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Println("AVISO: Variável de ambiente WEATHER_API_KEY não está configurada")
	} else {
		log.Println("WEATHER_API_KEY configurada com sucesso")
	}

	http.HandleFunc("/weather/", weatherHandler)

	log.Printf("Servidor iniciado na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	path := r.URL.Path
	re := regexp.MustCompile(`/weather/(.+)`)
	matches := re.FindStringSubmatch(path)

	if len(matches) < 2 {
		respondWithError(w, http.StatusBadRequest, "CEP não fornecido")
		return
	}

	cep := matches[1]
	log.Printf("CEP recebido: %s", cep)

	reValidCep := regexp.MustCompile(`^\d{8}$`)
	if !reValidCep.MatchString(cep) {
		log.Printf("CEP inválido: %s", cep)
		respondWithError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	location, err := getCEPInfo(cep)
	if err != nil {
		log.Printf("Erro ao buscar informações do CEP: %v", err)
		respondWithError(w, http.StatusNotFound, "can not find zipcode")
		return
	}
	log.Printf("Localidade encontrada: %s/%s", location.Localidade, location.UF)

	tempC, err := getWeather(location.Localidade)
	if err != nil {
		log.Printf("Erro ao obter dados do clima: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro ao obter dados do clima")
		return
	}
	log.Printf("Temperatura obtida: %.2f°C", tempC)

	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15

	// Preparar resposta
	response := TemperatureResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

var getCEPInfo = func(cep string) (*ViaCEPResponse, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	log.Printf("Consultando ViaCEP: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status de resposta inválido: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
	}
	log.Printf("Resposta ViaCEP: %s", string(body))

	var cepInfo ViaCEPResponse
	if err := json.Unmarshal(body, &cepInfo); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	if cepInfo.Erro || cepInfo.Localidade == "" {
		return nil, fmt.Errorf("CEP não encontrado")
	}

	return &cepInfo, nil
}

var getWeather = func(city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("chave da API de clima não configurada")
	}

	query := city
	encodedQuery := url.QueryEscape(query)

	apiURL := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, encodedQuery)
	log.Printf("Consultando WeatherAPI: %s", apiURL)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, fmt.Errorf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("erro ao consultar API do clima (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("erro ao ler resposta: %v", err)
	}
	log.Printf("Resposta WeatherAPI: %s", string(body))

	var weatherInfo WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherInfo); err != nil {
		return 0, fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	return weatherInfo.Current.TempC, nil
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}
