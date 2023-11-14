package main

import (
	"fmt"
	"strings"
	"strconv"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const conferenceTickets = 50

var conferenceName = "concert"
var remainingTickets uint = 50

type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOfTickets uint
}

var bookings = make([]UserData, 0)

var myApp fyne.App
var myWindow fyne.Window

func validateUser(firstName string, lastName string, email string, userTickets uint) (bool, bool, bool) {
	isValidName := len(firstName) >= 2 && len(lastName) >= 2
	isValidEmail := strings.Contains(email, "@") && len(email) >= 2
	isValidTicketAmount := userTickets > 0 && userTickets <= remainingTickets

	return isValidName, isValidEmail, isValidTicketAmount
}

func main() {
	myApp = app.New()
	myWindow = myApp.NewWindow("Conference Booking")

	tabs := container.NewAppTabs(
		tabItem("Book Tickets", bookingTab()),
		tabItem("Booking List", bookingListTab()),
	)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

func bookingTab() fyne.CanvasObject {
	firstNameEntry := widget.NewEntry()
	lastNameEntry := widget.NewEntry()
	emailEntry := widget.NewEntry()
	ticketsEntry := widget.NewEntry()

	form := &widget.Form{
		OnSubmit: func() {
			firstName := firstNameEntry.Text
			lastName := lastNameEntry.Text
			email := emailEntry.Text
			ticketsStr := ticketsEntry.Text

			userTickets, err := parseTicketInput(ticketsStr)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			isValidName, isValidEmail, isValidTicketAmount := validateUser(firstName, lastName, email, userTickets)

			if isValidName && isValidEmail && isValidTicketAmount {
				bookTickets(userTickets, firstName, lastName, email)
				sendTicketsAsync(firstName, lastName, userTickets, email)
				dialog.ShowInformation("Booking Success", "Your tickets have been booked successfully.", myWindow)
			} else {
				var errorMsg string
				if !isValidName {
					errorMsg += "First name or last name is too short\n"
				}
				if !isValidEmail {
					errorMsg += "Email address is not valid\n"
				}
				if !isValidTicketAmount {
					errorMsg += "Number of tickets that you entered is not valid\n"
				}
				dialog.ShowError(fmt.Errorf(errorMsg), myWindow)
			}
		},
	}

	form.Append("First Name", firstNameEntry)
	form.Append("Last Name", lastNameEntry)
	form.Append("Email", emailEntry)
	form.Append("Number of Tickets", ticketsEntry)

	return container.New(layout.NewCenterLayout(), form)
}

func bookingListTab() fyne.CanvasObject {
	list := widget.NewList(
		func() int {
			return len(bookings)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(fmt.Sprintf("%s %s: %d tickets", bookings[id].firstName, bookings[id].lastName, bookings[id].numberOfTickets))
		},
	)

	return container.New(layout.NewCenterLayout(), list)
}

func parseTicketInput(input string) (uint, error) {
	tickets, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse tickets: %v", err)
	}
	return uint(tickets), nil
}

func bookTickets(userTickets uint, firstName string, lastName string, email string) {
	remainingTickets = remainingTickets - userTickets

	userData := UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}

	bookings = append(bookings, userData)

	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v\n", firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v\n", remainingTickets, conferenceName)
}

func sendTicketsAsync(firstName string, lastName string, userTickets uint, email string) {
	go func() {
		time.Sleep(10 * time.Second)
		tickets := fmt.Sprintf("%v tickets for %v %v", userTickets, firstName, lastName)
		// Consider adding error handling here
		fmt.Printf("Sending ticket: \n %v \n to email address %v\n", tickets, email)
	}()
}

func tabItem(name string, content fyne.CanvasObject) *container.TabItem {
	return container.NewTabItem(name, content)
}
