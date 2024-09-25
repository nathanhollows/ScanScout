package models

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/pkg/db"
	"github.com/uptrace/bun"
)

// Instance represents a single planned activity belonging to a user
// Instance is used to match users, teams, locations, and scans
type Instance struct {
	baseModel

	ID        string       `bun:",pk,type:varchar(36)" json:"id"`
	Name      string       `bun:",type:varchar(255)" json:"name"`
	UserID    string       `bun:",type:varchar(36)" json:"user_id"`
	StartTime bun.NullTime `bun:",nullzero" json:"start_time"`
	EndTime   bun.NullTime `bun:",nullzero" json:"end_time"`
	Status    GameStatus   `bun:"-" json:"status"`

	Teams     Teams            `bun:"rel:has-many,join:id=instance_id" json:"teams"`
	Locations Locations        `bun:"rel:has-many,join:id=instance_id" json:"instance_locations"`
	Scans     Scans            `bun:"rel:has-many,join:id=instance_id" json:"scans"`
	Settings  InstanceSettings `bun:"rel:has-one,join:id=instance_id" json:"settings"`
}

type Instances []Instance

func (i *Instance) Save(ctx context.Context) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	_, err := db.DB.NewInsert().Model(i).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) Update(ctx context.Context) error {
	_, err := db.DB.NewUpdate().Model(i).WherePK().Exec(ctx)
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
	for _, location := range i.Locations {
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

	_, err := db.DB.NewDelete().Model(i).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// FindAllInstances finds all instances
func FindAllInstances(ctx context.Context, userID string) (Instances, error) {
	instances := Instances{}
	err := db.DB.NewSelect().Model(&instances).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

// FindInstanceByID finds an instance by ID
func FindInstanceByID(ctx context.Context, id string) (*Instance, error) {
	instance := &Instance{}
	err := db.DB.NewSelect().
		Model(instance).
		Where("id = ?", id).
		Relation("Locations").
		Relation("Settings").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// GetStatus returns the status of the instance
func (i *Instance) GetStatus() GameStatus {
	// If the start time is null, the game is closed
	if i.StartTime.Time.IsZero() {
		return Closed
	}

	// If the start time is in the future, the game is scheduled
	if i.StartTime.Time.UTC().After(time.Now().UTC()) {
		return Scheduled
	}

	// If the end time is in the past, the game is closed
	if !i.EndTime.Time.IsZero() && i.EndTime.Time.Before(time.Now().UTC()) {
		return Closed
	}

	// If the start time is in the past, the game is active
	return Active

}

// LoadSettings loads the settings for an instance
func (i *Instance) LoadSettings(ctx context.Context) error {
	if i.Settings.InstanceID == "" {
		i.Settings = InstanceSettings{}
		err := db.DB.NewSelect().Model(&i.Settings).Where("instance_id = ?", i.ID).Scan(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadLocations loads the locations for an instance
func (i *Instance) LoadLocations(ctx context.Context) error {
	if len(i.Locations) > 0 {
		return nil
	}

	var err error
	i.Locations, err = FindAllLocations(ctx, i.ID)
	if err != nil {
		return err
	}

	return nil
}

// LoadTeams loads the teams for an instance
func (i *Instance) LoadTeams(ctx context.Context) error {
	if len(i.Teams) > 0 {
		return nil
	}

	var err error
	i.Teams, err = FindAllTeams(ctx, i.ID)
	if err != nil {
		return err
	}

	return nil
}

func (i *Instance) ZipPosters(ctx context.Context) (string, error) {
	// Create a zip file
	path := "./assets/posters/" + i.ID + ".zip"
	archive, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// Collect the paths
	var paths []string
	// instanceLocations, err := FindAllLocations(ctx, i.ID)
	if err != nil {
		return "", err
	}
	// for _, instanceLocation := range instanceLocations {
	// 	paths = append(paths, instanceLocation.Marker.getQRFilename(true))
	// }

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

func GeneratePosters(ctx context.Context, instanceID string) (string, error) {
	instance, err := FindInstanceByID(ctx, instanceID)

	// Set up the document
	pdf := fpdf.New(fpdf.OrientationPortrait, fpdf.UnitMillimeter, fpdf.PageSizeA4, "")
	pdf.AddUTF8Font("ArchivoBlack", "", "./assets/fonts/ArchivoBlack-Regular.ttf")
	pdf.AddUTF8Font("OpenSans", "", "./assets/fonts/OpenSans.ttf")

	// for _, location := range instance.Locations {
	// 	location.Marker.GenerateQRCode()
	// 	generatePosterPage(pdf, &location.Marker, instance, true)
	// 	// if location.Coords.MustScanOut {
	// 	// 	generatePosterPage(pdf, &location.Coords, i, false)
	// 	// }
	// }

	path := "./assets/posters/" + instance.ID + " posters.pdf"
	err = pdf.OutputFileAndClose(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func generatePosterPage(pdf *fpdf.Fpdf, coords *Marker, instance *Instance, scanIn bool) {
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
	// pdf.Image(coords.getQRPath(scanIn), 50, 90, 110, 0, false, "", 0, "")

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
