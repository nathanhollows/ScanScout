package models

import (
	"archive/zip"
	"context"
	"errors"
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

	ID                string    `bun:",pk,type:varchar(36)" json:"id"`
	Name              string    `bun:",type:varchar(255)" json:"name"`
	UserID            string    `bun:",type:varchar(36)" json:"user_id"`
	User              User      `bun:"rel:has-one,join:user_id=user_id" json:"user"`
	Teams             Teams     `bun:"rel:has-many,join:id=instance_id" json:"teams"`
	InstanceLocations Locations `bun:"rel:has-many,join:id=instance_id" json:"instance_locations"`
	Scans             Scans     `bun:"rel:has-many,join:id=instance_id" json:"scans"`
}

type Instances []Instance

func (i *Instance) Save(ctx context.Context) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	_, err := db.NewInsert().Model(i).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) Update(ctx context.Context) error {
	_, err := db.NewUpdate().Model(i).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Deleting an instance will cascade delete all teams, locations, and scans
func (i *Instance) Delete(ctx context.Context) error {
	// Delete teams
	for _, team := range i.Teams {
		err := team.Delete(ctx)
		if err != nil {
			return err
		}
	}

	// Delete locations
	for _, location := range i.InstanceLocations {
		err := location.Delete(ctx)
		if err != nil {
			return err
		}
	}

	// Delete scans
	for _, scan := range i.Scans {
		err := scan.Delete(ctx)
		if err != nil {
			return err
		}
	}

	_, err := db.NewDelete().Model(i).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetCurrentUserInstance gets the current instance from the context
func GetCurrentUserInstance(ctx context.Context) (*Instance, error) {
	user, ok := ctx.Value(UserIDKey).(*User)
	if !ok {
		return nil, errors.New("User not found in context")
	}

	if user.CurrentInstance == nil {
		return nil, errors.New("Current instance not found")
	}

	return user.CurrentInstance, nil
}

// FindAllInstances finds all instances
func FindAllInstances(ctx context.Context) (Instances, error) {
	instances := Instances{}
	user, ok := ctx.Value(UserIDKey).(*User)
	if !ok {
		return nil, errors.New("User not found in context")
	}
	err := db.NewSelect().Model(&instances).Where("user_id = ?", user.UserID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

// FindInstanceByID finds an instance by ID
func FindInstanceByID(ctx context.Context, id string) (*Instance, error) {
	instance := &Instance{}
	err := db.NewSelect().
		Model(instance).
		Where("id = ?", id).
		Relation("InstanceLocations").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func GenerateQRCodeArchive(ctx context.Context) (string, error) {
	instanceID := ctx.Value(UserIDKey).(*User).CurrentInstanceID
	locations, err := FindAllLocations(ctx)
	if err != nil {
		return "", err
	}
	for _, location := range locations {
		err = location.Coords.GenerateQRCode()
		if err != nil {
			return "", err
		}
	}

	// Create a zip file
	path := "./assets/codes/" + instanceID + ".zip"
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
		paths = append(paths, location.Coords.getQRFilename(true))
		// Commented out because we don't need to scan out
		// TODO: Implement scan out on a Locations level
		// if location.Coords.MustScanOut {
		// 	paths = append(paths, location.Coords.getQRFilename(false))
		// }
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

func (i *Instance) ZipPosters(ctx context.Context) (string, error) {
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
	instanceLocations, err := FindAllLocations(ctx)
	if err != nil {
		return "", err
	}
	for _, instanceLocation := range instanceLocations {
		paths = append(paths, instanceLocation.Coords.getQRFilename(true))
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

func GeneratePosters(ctx context.Context) (string, error) {
	instanceID := ctx.Value(UserIDKey).(*User).CurrentInstanceID
	instance, err := FindInstanceByID(ctx, instanceID)

	// Set up the document
	pdf := fpdf.New(fpdf.OrientationPortrait, fpdf.UnitMillimeter, fpdf.PageSizeA4, "")
	pdf.AddUTF8Font("ArchivoBlack", "", "./assets/fonts/ArchivoBlack-Regular.ttf")
	pdf.AddUTF8Font("OpenSans", "", "./assets/fonts/OpenSans.ttf")

	for _, location := range instance.InstanceLocations {
		location.Coords.GenerateQRCode()
		generatePosterPage(pdf, &location.Coords, instance, true)
		// if location.Coords.MustScanOut {
		// 	generatePosterPage(pdf, &location.Coords, i, false)
		// }
	}

	path := "./assets/posters/" + instance.ID + " posters.pdf"
	err = pdf.OutputFileAndClose(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func generatePosterPage(pdf *fpdf.Fpdf, coords *Coords, instance *Instance, scanIn bool) {
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
	locationName := coords.Name
	pdf.SetY(40)
	pdf.SetX((210 - pdf.GetStringWidth(locationName)) / 2)
	pdf.Cell(40, 70, locationName)

	// Add the image
	pdf.Image(coords.getQRPath(scanIn), 50, 90, 110, 0, false, "", 0, "")

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

	// Add the date to the bottom of the poster in small text
	pdf.SetFont("OpenSans", "", 10)
	pdf.SetY(270)
	pdf.SetX((210 - pdf.GetStringWidth("4th March")) / 2)
	pdf.Cell(0, 0, "4th March")
}
