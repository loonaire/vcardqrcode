package main

import (
	"bytes"
	"errors"
	"image"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skip2/go-qrcode"
)

const (
	applicationId = "ContactQRCode"
)

var (
	//appPreferences = []string{"fullname", "firstname", "lastname", "email", "phone", "mobile", "address", "city", "postalCode", "country", "organization", "title", "website"}

	contactQRCodeApp fyne.App
	mainWindow       fyne.Window
	qrCodeLabel      widget.Label
	qrCodeCanvas     *canvas.Image
	qrCodeTabContent *fyne.Container
)

func main() {

	contactQRCodeApp = app.NewWithID(applicationId)
	mainWindow = contactQRCodeApp.NewWindow("ContactQRCode")

	mainWindow.SetContent(container.NewBorder(nil, nil, nil, nil, mainScreen(contactQRCodeApp)))
	mainWindow.Resize(fyne.NewSize(600, 400))

	mainWindow.ShowAndRun()

	/*
		myApp := app.NewWithID("com.example.loonaire.test")

		myWindow := myApp.NewWindow("Contact QRCode")

		//myApp.Preferences().SetString("hello", "world")

		lvl := widget.NewLabel(myApp.Preferences().String("hello"))

		// Définir le contenu de la fenêtre
		myWindow.SetContent(container.NewBorder(lvl, nil, nil, nil, mainScreen()))

		// Afficher et exécuter l'application
		myWindow.Resize(fyne.NewSize(400, 300))
		myWindow.ShowAndRun()
	*/

}

func mainScreen(a fyne.App) fyne.CanvasObject {
	tab1 := container.NewTabItemWithIcon("QRCode", theme.FileImageIcon(), QRCodeTab(a))
	tab2 := container.NewTabItemWithIcon("Informations de contact", theme.AccountIcon(), contactInfoTab(a))
	tabs := container.NewAppTabs(
		tab1,
		tab2,
	)

	return container.NewBorder(nil, nil, nil, nil, tabs)
}

func QRCodeTab(a fyne.App) fyne.CanvasObject {
	qrCodeLabel.SetText(a.Preferences().String("vCard"))
	qrCodeCanvas = generateQrCodeCanvas(a)
	qrCodeTabContent = container.NewBorder(&qrCodeLabel, nil, nil, nil, qrCodeCanvas)
	return qrCodeTabContent
}

func updateQRCodeTab(a fyne.App) {
	qrCodeLabel.SetText(a.Preferences().String("vCard"))

	log.Println(qrCodeTabContent.Objects)

	qrCodeCanvas.Image = generateQrCodeImage(a)
	qrCodeCanvas.Refresh()
}

func contactInfoTab(a fyne.App) fyne.CanvasObject {
	appPreferences := a.Preferences()
	fullname := widget.NewEntry()
	fullname.SetText(appPreferences.String("fullname"))
	fullname.Validator = validation.NewRegexp(`.+`, "Le champ ne peux pas être vide")

	firstname := widget.NewEntry()
	firstname.SetText(appPreferences.String("firstname"))

	lastname := widget.NewEntry()
	lastname.SetText(appPreferences.String("lastname"))

	email := widget.NewEntry()
	email.SetText(appPreferences.String("email"))
	email.Validator = validation.NewRegexp(`^([^\s.][\w-_.]*[^.])(@\w+)(\.\w+(\.\w+)?[^.\W])$`, "Adresse mail invalide")

	phone := widget.NewEntry()
	phone.SetText(appPreferences.String("phone"))
	phone.Validator = validation.NewRegexp(`^$|^\+?[0-9]{10,}$`, "Numero invalide")

	mobile := widget.NewEntry()
	mobile.SetText(appPreferences.String("mobile"))
	mobile.Validator = validation.NewRegexp(`^$|^\+?[0-9]{10,}$`, "Numero invalide")

	address := widget.NewEntry()
	address.SetText(appPreferences.String("address"))

	city := widget.NewEntry()
	city.SetText(appPreferences.String("city"))

	postalCode := widget.NewEntry()
	postalCode.SetText(appPreferences.String("postalCode"))
	postalCode.Validator = validation.NewRegexp(`^$|^[0-9]*$`, "Code postal invalide")

	country := widget.NewEntry()
	country.SetText(appPreferences.String("country"))

	organization := widget.NewEntry()
	organization.SetText(appPreferences.String("organization"))

	title := widget.NewEntry()
	title.SetText(appPreferences.String("title"))

	website := widget.NewEntry()
	website.SetText(appPreferences.String("url"))

	vCardForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Nom complet", Widget: fullname},
			{Text: "Prénom", Widget: firstname},
			{Text: "Nom", Widget: lastname},
			{Text: "Email", Widget: email, HintText: "A valid email address"},
			{Text: "Téléphone", Widget: phone},
			{Text: "Mobile", Widget: mobile},
			{Text: "Adresse", Widget: address},
			{Text: "Code postal", Widget: postalCode},
			{Text: "Ville", Widget: city},
			{Text: "Pays", Widget: country},
			{Text: "Entreprise", Widget: organization},
			{Text: "Titre", Widget: title},
			{Text: "Site web", Widget: website},
		},
		OnSubmit: func() {
			appPreferences.SetString("fullname", fullname.Text)
			appPreferences.SetString("firstname", firstname.Text)
			appPreferences.SetString("lastname", lastname.Text)
			appPreferences.SetString("email", email.Text)
			appPreferences.SetString("phone", phone.Text)
			appPreferences.SetString("mobile", mobile.Text)
			appPreferences.SetString("address", address.Text)
			appPreferences.SetString("postalCode", postalCode.Text)
			appPreferences.SetString("city", city.Text)
			appPreferences.SetString("country", country.Text)
			appPreferences.SetString("organization", organization.Text)
			appPreferences.SetString("title", title.Text)
			appPreferences.SetString("url", website.Text)

			appPreferences.SetString("vCard", generatevCardCode(a))

			//qrcodeImage = generateQrCodeCanvas(a)
			updateQRCodeTab(a)
		},
	}
	vCardForm.SubmitText = "Sauvegarder"

	return vCardForm

	//return container.NewVScroll(vCardForm)
}

func generatevCardCode(a fyne.App) string {

	vCardArray := []string{}

	vCardArray = append(vCardArray, "BEGIN:VCARD")
	vCardArray = append(vCardArray, "VERSION:4.0")

	vCardArray = append(vCardArray, "FN:"+a.Preferences().String("fullname"))

	vCardName := "N:"
	if a.Preferences().String("lastname") != "" {
		vCardName += a.Preferences().String("lastname")
	}
	vCardName += ";"
	if a.Preferences().String("firstname") != "" {
		vCardName += a.Preferences().String("firstname")
	}
	vCardArray = append(vCardArray, vCardName+";;;")

	vCardArray = append(vCardArray, "EMAIL;TYPE=WORK:"+a.Preferences().String("email"))

	if a.Preferences().String("phone") != "" {
		vCardArray = append(vCardArray, "TEL;TYPE=LANDLINE:"+a.Preferences().String("phone"))
	}
	if a.Preferences().String("mobile") != "" {
		vCardArray = append(vCardArray, "TEL;TYPE=CELL:"+a.Preferences().String("mobile"))
	}

	vCardAddress := "ADR;TYPE=WORK:;;"
	if a.Preferences().String("address") != "" {
		vCardAddress += a.Preferences().String("address")
	}
	vCardAddress += ";"
	if a.Preferences().String("city") != "" {
		vCardAddress += a.Preferences().String("city")
	}
	vCardAddress += ";;"
	if a.Preferences().String("postalCode") != "" {
		vCardAddress += a.Preferences().String("postalCode")
	}
	vCardAddress += ";"
	if a.Preferences().String("country") != "" {
		vCardAddress += a.Preferences().String("country")
	}

	vCardArray = append(vCardArray, vCardAddress)

	if a.Preferences().String("organization") != "" {
		vCardArray = append(vCardArray, "ORG:"+a.Preferences().String("organization"))
	}

	if a.Preferences().String("title") != "" {
		vCardArray = append(vCardArray, "TITLE:"+a.Preferences().String("title"))
	}

	if a.Preferences().String("url") != "" {
		vCardArray = append(vCardArray, "URL:"+a.Preferences().String("url"))
	}

	vCardArray = append(vCardArray, "END:VCARD")
	// NOTE : \r\n is used because on windows, google lens read a url instead of vcard
	return strings.Join(vCardArray, "\r\n") + "\r\n"

}

func generateQrCodeImage(a fyne.App) image.Image {
	data, err := qrcode.Encode(a.Preferences().String("vCard"), qrcode.Highest, 1024)
	//data, err := qrcode.Encode("hello", qrcode.Low, 512)
	if err != nil {
		dialog.ShowError(errors.New("erreur lors de la création du QRcode"), mainWindow)
		return nil
	}
	img, _, _ := image.Decode(bytes.NewReader(data))
	return img
}

func generateQrCodeCanvas(a fyne.App) *canvas.Image {

	qrcodeImage := canvas.NewImageFromImage(generateQrCodeImage(a))
	qrcodeImage.FillMode = canvas.ImageFillContain

	return qrcodeImage
}
