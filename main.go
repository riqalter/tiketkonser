package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type ConcertTicket struct {
	ConcertName string
	Price       float64
}

type User struct {
	Name  string
	Email string
}

func formatCurrency(amount float64) string {
	tag := language.MustParse("id-ID")
	p := message.NewPrinter(tag)
	return p.Sprintf("Rp. %v", amount)
}

func buyTicket(user User, ticket ConcertTicket, quantity int, mainWindow fyne.Window) {
	totalCost := ticket.Price * float64(quantity)
	message := fmt.Sprintf("Pembelian tiket oleh %s (%s) untuk konser %s sejumlah %d dengan total biaya %s", user.Name, user.Email, ticket.ConcertName, quantity, formatCurrency(totalCost))

	dialog.ShowInformation("Pembelian Tiket", message, mainWindow)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Pembelian Tiket Konser")

	concertTicket := ConcertTicket{ConcertName: "Konser A", Price: 100000.0}
	user := User{Name: "John Doe", Email: "john@example.com"}

	userNameEntry := widget.NewEntry()
	userNameEntry.SetPlaceHolder("Nama Pengguna")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")

	quantityEntry := widget.NewEntry()
	quantityEntry.SetPlaceHolder("Jumlah Tiket")

	buyButton := widget.NewButton("Beli Tiket", func() {
		user.Name = userNameEntry.Text
		user.Email = emailEntry.Text

		quantity := 0
		fmt.Sscanf(quantityEntry.Text, "%d", &quantity)

		if quantity > 0 {
			buyTicket(user, concertTicket, quantity, myWindow)
		} else {
			dialog.ShowError(fmt.Errorf("Jumlah tiket harus lebih dari 0"), myWindow)
		}
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Nama Pengguna", Widget: userNameEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Jumlah Tiket", Widget: quantityEntry},
		},
	}

	myWindow.SetContent(container.NewVBox(
		form,
		buyButton,
	))

	myWindow.ShowAndRun()
}
