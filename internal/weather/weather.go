package weather

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const GetWeatherURL = "https://wttr.in/%s?format=j2"

type WeatherInfo struct {
	City         string
	TemperatureC string
	TemperatureF string
	Description  string
}

type apiResponse struct {
	CurrentCondition []struct {
		TempC       string `json:"temp_C"`
		TempF       string `json:"temp_F"`
		WeatherDesc []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
	} `json:"current_condition"`
	NearestArea []struct {
		AreaName []struct {
			Value string `json:"value"`
		} `json:"areaName"`
	} `json:"nearest_area"`
}

func GetTLSClient() (*http.Client, error) {
	conn, err := tls.Dial("tcp", "wttr.in:443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("certificate not found")
	}

	roots := x509.NewCertPool()
	for _, cert := range certs {
		roots.AddCert(cert)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: roots,
			},
		},
	}
	return client, nil
}

func GetWeather(city string, client *http.Client) (WeatherInfo, error) {
	url := fmt.Sprintf(GetWeatherURL, city)

	resp, err := client.Get(url)
	if err != nil {
		return WeatherInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherInfo{}, err
	}

	if resp.StatusCode != 200 {
		return WeatherInfo{}, fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	var data apiResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return WeatherInfo{}, err
	}

	if len(data.CurrentCondition) == 0 || len(data.NearestArea) == 0 {
		return WeatherInfo{}, fmt.Errorf("wrong data")
	}

	info := WeatherInfo{
		City:         data.NearestArea[0].AreaName[0].Value,
		TemperatureC: data.CurrentCondition[0].TempC,
		TemperatureF: data.CurrentCondition[0].TempF,
		Description:  data.CurrentCondition[0].WeatherDesc[0].Value,
	}

	return info, nil
}

func Show(weatherInfo WeatherInfo) {
	fmt.Print("\033[H\033[2J")
	fmt.Println()
	fmt.Println()
	if weatherInfo.City == "" {
		fmt.Println("\tWifi connection failed >_<")
	} else {
		fmt.Println("\t\tCity:\t\t", weatherInfo.City)
		fmt.Println("\t\tTemperature:\t", weatherInfo.TemperatureC, "(C)"+" / "+weatherInfo.TemperatureF+" (F)")
		fmt.Println("\t\tDescription:\t", weatherInfo.Description)
	}
	fmt.Println()
	fmt.Println()
}
