package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
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

type Purchase struct {
	User     User
	Ticket   ConcertTicket
	Quantity int
	TotalCost float64
}

type CustomEntry struct {
	widget.Entry
}

var purchases []Purchase

func NewCustomEntry() *CustomEntry {
	entry := &CustomEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *CustomEntry) MinSize() fyne.Size {
	return fyne.NewSize(200, e.Entry.MinSize().Height)
}

func formatCurrency(amount float64) string {
	tag := language.MustParse("id-ID")
	p := message.NewPrinter(tag)
	return p.Sprintf("Rp. %v", amount)
}

func buyTicket(user User, ticket ConcertTicket, quantity int, mainWindow fyne.Window, purchaseList *fyne.Container, totalPurchaseLabel *widget.Label) {
	totalCost := ticket.Price * float64(quantity)
	message := fmt.Sprintf("Pembelian tiket atas nama  %s (%s) untuk konser %s sejumlah %d dengan total biaya %s", user.Name, user.Email, ticket.ConcertName, quantity, formatCurrency(totalCost))

	dialog.ShowInformation("Pembelian Tiket", message, mainWindow)

	purchase := Purchase{User: user, Ticket: ticket, Quantity: quantity, TotalCost: totalCost}
	purchases = append(purchases, purchase)

	purchaseList.Add(widget.NewLabel(message))

	totalPurchase := 0.0
	for _, purchase := range purchases {
		totalPurchase += purchase.TotalCost
	}
	totalPurchaseLabel.SetText(fmt.Sprintf("Total Pembelian: %s", formatCurrency(totalPurchase)))
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Pembelian Tiket Konser")

	concertTicket := ConcertTicket{ConcertName: "Konser A", Price: 100000.0}
	user := User{Name: "John Doe", Email: "john@example.com"}

	userNameEntry := NewCustomEntry()
	userNameEntry.SetPlaceHolder("Nama Pengguna")

	emailEntry := NewCustomEntry()
	emailEntry.SetPlaceHolder("Email")

	quantityEntry := NewCustomEntry()
	quantityEntry.SetPlaceHolder("Jumlah Tiket")

	purchaseList := container.NewVBox()
	totalPurchaseLabel := widget.NewLabel("Total Pembelian: Rp. 0")

	buyButton := widget.NewButton("Beli Tiket", func() {
		user.Name = userNameEntry.Text
		user.Email = emailEntry.Text

		quantity := 0
		fmt.Sscanf(quantityEntry.Text, "%d", &quantity)

		if quantity > 0 {
			buyTicket(user, concertTicket, quantity, myWindow, purchaseList, totalPurchaseLabel)
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

	myWindow.SetContent(container.New(layout.NewCenterLayout(), container.NewVBox(
		form,
		buyButton,
		widget.NewLabel("Daftar Pembelian:"),
		purchaseList,
		totalPurchaseLabel,
	)))

	myWindow.ShowAndRun()
}