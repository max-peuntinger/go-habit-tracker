package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./habits.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	setupDatabase()

	for {
		showMenu()
		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			addHabit()
		case 2:
			trackHabit()
		case 3:
			listHabits()
		case 4:
			viewHabitDetails()
		case 5:
			return
		default:
			fmt.Println("Invalid choice!")
		}
	}
}

func setupDatabase() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS habit (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        start_date DATE NOT NULL
    )`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tracking (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        habit_id INTEGER,
        date DATE NOT NULL,
        FOREIGN KEY (habit_id) REFERENCES habit(id)
    )`)
	if err != nil {
		panic(err)
	}
}

func showMenu() {
	fmt.Println("1. Add Habit")
	fmt.Println("2. Track Habit")
	fmt.Println("3. List Habits")
	fmt.Println("4. View Habit Details")
	fmt.Println("5. Exit")
	fmt.Print("Enter your choice: ")
}

func addHabit() {
	fmt.Print("Enter habit name: ")
	var name string
	fmt.Scan(&name)

	_, err := db.Exec("INSERT INTO habit(name, start_date) VALUES(?, ?)", name, time.Now().Format("2006-01-02"))
	if err != nil {
		fmt.Println("Error adding habit:", err)
		return
	}
	fmt.Println("Habit added successfully!")
}

func trackHabit() {
	listHabits()

	fmt.Print("Enter habit ID to track: ")
	var id int
	fmt.Scan(&id)

	_, err := db.Exec("INSERT INTO tracking(habit_id, date) VALUES(?, ?)", id, time.Now().Format("2006-01-02"))
	if err != nil {
		fmt.Println("Error tracking habit:", err)
		return
	}
	fmt.Println("Habit tracked successfully!")
}

func listHabits() {
	rows, err := db.Query("SELECT id, name, start_date FROM habit")
	if err != nil {
		fmt.Println("Error fetching habits:", err)
		return
	}
	defer rows.Close()

	fmt.Println("ID | Habit Name | Start Date")
	for rows.Next() {
		var id int
		var name, startDate string
		rows.Scan(&id, &name, &startDate)
		fmt.Printf("%d | %s | %s\n", id, name, startDate)
	}
}

func viewHabitDetails() {
	// First, list all habits to let the user choose which one they want to see details for
	listHabits()

	fmt.Print("Enter habit ID to view details: ")
	var id int
	fmt.Scan(&id)

	rows, err := db.Query("SELECT date FROM tracking WHERE habit_id=?", id)
	if err != nil {
		fmt.Println("Error fetching habit details:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Tracked Dates for Habit ID:", id)
	for rows.Next() {
		var date string
		rows.Scan(&date)
		fmt.Println(date)
	}
}
