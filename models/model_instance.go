package models

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
)

// Instance represents a single planned activity belonging to a user
// Instance is used to match users, teams, locations, and scans
type Instance struct {
	baseModel

	ID        string    `bun:",pk,type:varchar(36)" json:"id"`
	Name      string    `bun:",type:varchar(255)" json:"name"`
	UserId    string    `bun:",type:varchar(36)" json:"user_id"`
	User      User      `bun:"rel:has-one,join:user_id=user_id" json:"user"`
	Teams     Teams     `bun:"rel:has-many,join:id=instance_id" json:"teams"`
	Locations Locations `bun:"rel:has-many,join:id=instance_id" json:"locations"`
	Scans     Scans     `bun:"rel:has-many,join:id=instance_id" json:"scans"`
}

type Instances []Instance

func CreateDummyInstance() (*Instance, error) {
	instance := &Instance{
		ID:   uuid.New().String(),
		Name: "LAWS200 Wānanga",
	}
	err := instance.Save()
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (i *Instance) Save() error {
	ctx := context.Background()
	_, err := db.NewInsert().Model(i).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) Update() error {
	ctx := context.Background()
	_, err := db.NewUpdate().Model(i).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) Delete() error {
	ctx := context.Background()
	_, err := db.NewDelete().Model(i).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// FindAllInstances finds all instances
func FindAllInstances(userId string) (Instances, error) {
	ctx := context.Background()
	instances := Instances{}
	err := db.NewSelect().Model(&instances).Where("user_id = ?", userId).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

// FindInstanceByID finds an instance by ID
func FindInstanceByID(id string) (*Instance, error) {
	CreateDummyInstance()
	ctx := context.Background()
	instance := &Instance{}
	err := db.NewSelect().
		Model(instance).
		Where("id = ?", id).
		Relation("Locations").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (i *Instance) ZipQRCodes() (string, error) {
	locations, err := FindAllLocations()
	if err != nil {
		return "", err
	}
	for _, location := range locations {
		err = location.GenerateQRCode()
		if err != nil {
			return "", err
		}
	}

	// Create a zip file
	path := "./assets/codes/" + i.ID + ".zip"
	archive, err := os.Create(path)
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// Collect the paths
	var paths []string
	for _, location := range locations {
		paths = append(paths, location.getQRFilename(true))
		if location.MustScanOut {
			paths = append(paths, location.getQRFilename(false))
		}
	}

	// Add each file to the zip
	adder := func(path string) error {
		file, err := os.Open("./assets/codes/" + path)
		if err != nil {
			return err
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = path
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	}

	for _, path := range paths {
		err = adder(path)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func (i *Instance) ZipPosters() (string, error) {
	// Create a zip file
	path := "./assets/posters/" + i.ID + ".zip"
	archive, err := os.Create(path)
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// Collect the paths
	var paths []string
	locations, err := FindAllLocations()
	if err != nil {
		return "", err
	}
	for _, location := range locations {
		paths = append(paths, location.getQRFilename(true))
	}

	// Add each file to the zip
	adder := func(path string) error {
		file, err := os.Open("./assets/posters/" + path)
		if err != nil {
			return err
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = path
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	}

	for _, path := range paths {
		err = adder(path)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func (i *Instance) GeneratePosters() (string, error) {
	// Set up the document
	pdf := fpdf.New(fpdf.OrientationPortrait, fpdf.UnitMillimeter, fpdf.PageSizeA4, "")
	pdf.AddUTF8Font("ArchivoBlack", "", "./assets/fonts/ArchivoBlack-Regular.ttf")
	pdf.AddUTF8Font("OpenSans", "", "./assets/fonts/OpenSans.ttf")

	for _, location := range i.Locations {
		location.GenerateQRCode()
		generatePosterPage(pdf, location, i, true)
		if location.MustScanOut {
			generatePosterPage(pdf, location, i, false)
		}
	}

	path := "./assets/posters/" + i.ID + ".pdf"
	err := pdf.OutputFileAndClose(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func generatePosterPage(pdf *fpdf.Fpdf, location *Location, instance *Instance, scanIn bool) {
	pdf.AddPage()

	if !scanIn {
		// Add a background color
		pdf.SetFillColor(256, 216, 216)
		pdf.Rect(0, 0, 210, 297, "F")
	}

	// Add the instance title
	pdf.SetFont("ArchivoBlack", "", 28)
	title := strings.ToUpper(instance.Name)
	pdf.SetY(32)
	pdf.SetX((210 - pdf.GetStringWidth(title)) / 2)
	pdf.Cell(130, 32, title)

	// Add the location name
	pdf.SetFont("OpenSans", "", 24)
	locationName := location.Name
	pdf.SetY(40)
	pdf.SetX((210 - pdf.GetStringWidth(locationName)) / 2)
	pdf.Cell(40, 70, locationName)

	// Add the image
	pdf.Image(location.getQRPath(scanIn), 50, 90, 110, 0, false, "", 0, "")

	// Add Scan In/Out
	scanText := "Scan In"
	if !scanIn {
		scanText = "Scan Out"
	}
	pdf.SetY(180)
	pdf.SetX((210 - pdf.GetStringWidth(scanText)) / 2)
	pdf.Cell(40, 70, scanText)

	// Add reminder to scan out
	reminderText := "Remember to scan out before\nmoving to the next location!"
	if !scanIn {
		reminderText = "You must scan out before\nmoving to the next location!"
	}
	pdf.SetY(240)
	pdf.SetX((210 - 130) / 2)
	pdf.MultiCell(130, 10, reminderText, fpdf.BorderNone, fpdf.AlignCenter, false)
}
