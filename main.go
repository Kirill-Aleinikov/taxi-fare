package main

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	baseCost              = 200.0
	costPerKm             = 30.0
	childSeatCost         = 100.0
	trunkCost             = 50.0
	differentDistrictCost = 200.0
)

type Trip struct {
	Distance            float64
	HasChildSeat        bool
	HasTrunk            bool
	IsDifferentDistrict bool
	Tariff              string
}

func calculateCost(trip Trip) float64 {
	tariffMultipliers := map[string]float64{
		"econom":       1.0,
		"standart":     1.2,
		"dorogoi":      1.5,
		"premium":      2.0,
		"o4en_dorogoi": 2.5,
	}

	totalCost := baseCost + (trip.Distance * costPerKm)

	if trip.HasChildSeat {
		totalCost += childSeatCost
	}
	if trip.HasTrunk {
		totalCost += trunkCost
	}
	if trip.IsDifferentDistrict {
		totalCost += differentDistrictCost
	}

	multiplier, ok := tariffMultipliers[trip.Tariff]
	if !ok {
		multiplier = 1.0
	}
	totalCost *= multiplier

	return totalCost
}

func handler(w http.ResponseWriter, r *http.Request) {
	var result string
	if r.Method == "POST" {
		r.ParseForm()
		distance, _ := strconv.ParseFloat(r.FormValue("distance"), 64)
		tariff := r.FormValue("tariff")
		childSeat := r.FormValue("childSeat") == "on"
		trunk := r.FormValue("trunk") == "on"
		district := r.FormValue("district") == "on"

		trip := Trip{
			Distance:            distance,
			Tariff:              tariff,
			HasChildSeat:        childSeat,
			HasTrunk:            trunk,
			IsDifferentDistrict: district,
		}

		cost := calculateCost(trip)
		result = fmt.Sprintf("Итоговая стоимость поездки: %.2f руб.", cost)
	}

	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Калькулятор такси</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					margin: 0;
					padding: 0;
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
				}
				.container {
					background-color: #fff;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
					width: 300px;
				}
				h1 {
					text-align: center;
					color: #333;
				}
				form {
					display: flex;
					flex-direction: column;
				}
				label {
					margin-bottom: 8px;
					font-weight: bold;
					color: #555;
				}
				input[type="text"], select {
					padding: 8px;
					margin-bottom: 16px;
					border: 1px solid #ccc;
					border-radius: 4px;
					font-size: 16px;
				}
				input[type="checkbox"] {
					margin-bottom: 16px;
				}
				input[type="submit"] {
					background-color: #28a745;
					color: white;
					padding: 10px;
					border: none;
					border-radius: 4px;
					font-size: 16px;
					cursor: pointer;
				}
				input[type="submit"]:hover {
					background-color: #218838;
				}
				.result {
					margin-top: 20px;
					text-align: center;
					font-size: 18px;
					color: #333;
				}
				.form-group {
					display: flex;
					justify-content: space-between;
				}
			</style>
			<script>
				function hideResult() {
					setTimeout(function() {
						document.querySelector('.result').style.display = 'none';
					}, 10000);
				}
			</script>
		</head>
		<body onload="hideResult()">
			<div class="container">
				<h1>Калькулятор такси</h1>
				<form method="POST">
					<label for="distance">Расстояние (км):</label>
					<input type="text" name="distance" required><br>

					<label for="tariff">Тариф:</label>
					<select name="tariff">
						<option value="econom">Эконом</option>
						<option value="standart">Стандарт</option>
						<option value="dorogoi">Дорогой</option>
						<option value="premium">Премиум</option>
						<option value="o4en_dorogoi">Очень дорогой</option>
					</select><br>

					<div class="form-group">
						<label for="childSeat">Детское кресло:</label>
						<input type="checkbox" name="childSeat">
					</div>

					<div class="form-group">
						<label for="trunk">Багажник:</label>
						<input type="checkbox" name="trunk">
					</div>

					<div class="form-group">
						<label for="district">Подача в другой район:</label>
						<input type="checkbox" name="district">
					</div>

					<input type="submit" value="Рассчитать">
				</form>
				<div class="result">%s</div>
			</div>
		</body>
		</html>
	`, result)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
