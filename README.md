# CEP Weather API

A Go microservice that receives a Brazilian zipcode (CEP), identifies the city, and returns the current temperature in Celsius, Fahrenheit, and Kelvin. This service is designed to be deployed on Google Cloud Run.

## Google Cloud Link:

```
https://cep-weather-rkjaprhgxq-uc.a.run.app/
```
You can try diferent cep values using /weather/[cep number] endpoint:
 
```
https://cep-weather-rkjaprhgxq-uc.a.run.app/weather/01001000
```

## Project Structure

```
/cep-weather/
├── cmd/
│   ├── main.go
│   └── main_test.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── .env
```

## Features

- Accepts a valid 8-digit Brazilian zipcode (CEP)
- Validates the zipcode format
- Fetches location information using ViaCEP API
- Retrieves current temperature data using WeatherAPI
- Converts temperature to Celsius, Fahrenheit, and Kelvin
- Returns appropriate HTTP status codes and messages for different scenarios

## API Endpoints

### GET /weather/{zipcode}

**Success Response (200 OK)**
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

**Invalid Zipcode Format (422 Unprocessable Entity)**
```json
{
  "message": "invalid zipcode"
}
```

**Zipcode Not Found (404 Not Found)**
```json
{
  "message": "can not find zipcode"
}
```

## Prerequisites

- Go 1.20+
- Docker and Docker Compose
- WeatherAPI key (get one at [WeatherAPI](https://www.weatherapi.com/))

## Running Locally

### 1. Setup Environment Variables

Create an `.env` file in the root of the project:

```
WEATHER_API_KEY=your_weatherapi_key_here
```

### 2. Build and Run with Docker Compose

```bash
docker-compose up --build
```

This will build the Docker image and start the container. The service will be available at http://localhost:8080.

### 3. Testing the API with curl

**Testing a valid zipcode:**
```bash
curl http://localhost:8080/weather/01001000
```

**Testing an invalid zipcode format:**
```bash
curl http://localhost:8080/weather/123
```

**Testing a non-existent zipcode:**
```bash
curl http://localhost:8080/weather/99999999
```

## Running Tests

To run the automated tests in the project, navigate to the `cmd` directory and use the `go test` command:

```bash
cd cmd
go test -v
```

You can also run the tests with coverage information:

```bash
cd cmd
go test -cover
```

## External APIs Used

- [ViaCEP](https://viacep.com.br/) - For translating Brazilian zipcodes to city names
- [WeatherAPI](https://www.weatherapi.com/) - For retrieving current weather data

## Temperature Conversion Formulas

- Celsius to Fahrenheit: `F = C * 1.8 + 32`
- Celsius to Kelvin: `K = C + 273.15`

## License

This project is available under the MIT License.
