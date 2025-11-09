package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("длина слайса не равна 3")
	}

	steps, err1 := strconv.Atoi(parts[0])
	if err1 != nil {
		return 0, "", 0, err1
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	duration, err2 := time.ParseDuration(parts[2])
	if err2 != nil {
		return 0, "", 0, err2
	}
	if duration < 1 {
		return 0, "", 0, fmt.Errorf("некорректная продолжительность")
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {

	stepLength := height * stepLengthCoefficient

	distanceKm := (float64(steps) * stepLength) / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)

	durationInHours := duration.Hours()

	meanSpeed := distance / durationInHours

	return meanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Проверка вида тренировки
	var info string
	switch activity {
	case "Бег":
		// Рассчет дистанции, средней скорости и калории для бега
		distance := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, duration)
		calories, _ := RunningSpentCalories(steps, weight, height, duration)

		// Строка с информацией о тренировке
		info = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
			activity, duration.Hours(), distance, meanSpeed, calories)
	case "Ходьба":
		// Рассчет дистанции, средней скорости и калории для ходьбы
		distance := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, duration)
		calories, _ := WalkingSpentCalories(steps, weight, height, duration)

		// Строка с информацией о тренировке
		info = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
			activity, duration.Hours(), distance, meanSpeed, calories)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	return info, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	//проверка параметров на корректность
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные данные")
	}
	//рассчет средней скорости через meanSpeed()
	meanSpeed := meanSpeed(steps, height, duration)

	//перевод продолжительности в минуты
	durationInMinutes := duration.Minutes()

	return (weight * meanSpeed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	//проверка параметров на корректность
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные данные")
	}

	//рассчет средней скорости через meanSpeed()
	meanSpeed := meanSpeed(steps, height, duration)

	//перевод продолжительности в минуты
	durationInMinutes := duration.Minutes()

	calories := (weight * meanSpeed * durationInMinutes) / minInH

	return walkingCaloriesCoefficient * calories, nil
}
