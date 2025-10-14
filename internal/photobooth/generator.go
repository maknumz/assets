package photobooth

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	xdraw "golang.org/x/image/draw"
)

// Generate creates a collage image based on the provided template name.
func Generate(templateName string, imagePaths []string, outputPath string) error {
	tmpl, err := findTemplate(templateName)
	if err != nil {
		return err
	}

	if len(imagePaths) != len(tmpl.Slots) {
		return fmt.Errorf("template %q ต้องการ %d รูป แต่ส่งมา %d รูป", tmpl.Name, len(tmpl.Slots), len(imagePaths))
	}

	if err := ensureDir(outputPath); err != nil {
		return err
	}

	canvas := image.NewRGBA(image.Rect(0, 0, CanvasWidth, CanvasHeight))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	for i, slot := range tmpl.Slots {
		src, err := loadImage(imagePaths[i])
		if err != nil {
			return fmt.Errorf("ไม่สามารถอ่านไฟล์ %s: %w", imagePaths[i], err)
		}

		cropped := cropToAspect(src, slot.Dx(), slot.Dy())
		resized := image.NewRGBA(image.Rect(0, 0, slot.Dx(), slot.Dy()))
		xdraw.CatmullRom.Scale(resized, resized.Bounds(), cropped, cropped.Bounds(), draw.Over, nil)

		draw.Draw(canvas, slot, resized, image.Point{}, draw.Over)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("ไม่สามารถสร้างไฟล์เอาต์พุต: %w", err)
	}
	defer outFile.Close()

	options := &jpeg.Options{Quality: 95}
	if strings.EqualFold(filepath.Ext(outputPath), ".png") {
		return png.Encode(outFile, canvas)
	}

	return jpeg.Encode(outFile, canvas, options)
}

// findTemplate retrieves a template by name.
func findTemplate(name string) (*Template, error) {
	for _, tmpl := range templates {
		if tmpl.Name == name {
			return &tmpl, nil
		}
	}
	return nil, fmt.Errorf("ไม่พบ template ชื่อ %s", name)
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func cropToAspect(src image.Image, targetWidth, targetHeight int) image.Image {
	bounds := src.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	targetRatio := float64(targetWidth) / float64(targetHeight)
	srcRatio := float64(srcWidth) / float64(srcHeight)

	var cropRect image.Rectangle
	if srcRatio > targetRatio {
		// Source is wider than target, crop horizontally.
		newWidth := int(float64(srcHeight) * targetRatio)
		if newWidth > srcWidth {
			newWidth = srcWidth
		}
		left := bounds.Min.X + (srcWidth-newWidth)/2
		cropRect = image.Rect(left, bounds.Min.Y, left+newWidth, bounds.Min.Y+srcHeight)
	} else {
		// Source is taller than target, crop vertically.
		newHeight := int(float64(srcWidth) / targetRatio)
		if newHeight > srcHeight {
			newHeight = srcHeight
		}
		top := bounds.Min.Y + (srcHeight-newHeight)/2
		cropRect = image.Rect(bounds.Min.X, top, bounds.Min.X+srcWidth, top+newHeight)
	}

	cropped := image.NewRGBA(image.Rect(0, 0, cropRect.Dx(), cropRect.Dy()))
	draw.Draw(cropped, cropped.Bounds(), src, cropRect.Min, draw.Src)
	return cropped
}

func ensureDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("ไม่สามารถสร้างโฟลเดอร์ %s: %w", dir, err)
	}
	return nil
}

// ValidateTemplateName returns an error if the template name is unknown.
func ValidateTemplateName(name string) error {
	_, err := findTemplate(name)
	return err
}

// TemplateSummary provides a short description for listing purposes.
type TemplateSummary struct {
	Name        string
	Description string
	Slots       int
}

// Summaries returns descriptions for all available templates.
func Summaries() []TemplateSummary {
	list := Templates()
	result := make([]TemplateSummary, len(list))
	for i, tmpl := range list {
		result[i] = TemplateSummary{
			Name:        tmpl.Name,
			Description: tmpl.Description,
			Slots:       len(tmpl.Slots),
		}
	}
	return result
}

// SlotCountRange returns the minimum and maximum number of images supported.
func SlotCountRange() (int, int) {
	min := len(templates[0].Slots)
	max := len(templates[0].Slots)
	for _, tmpl := range templates {
		if len(tmpl.Slots) < min {
			min = len(tmpl.Slots)
		}
		if len(tmpl.Slots) > max {
			max = len(tmpl.Slots)
		}
	}
	return min, max
}

// FilterTemplatesBySlots returns templates that match the provided slot count.
func FilterTemplatesBySlots(count int) []TemplateSummary {
	matches := []TemplateSummary{}
	for _, tmpl := range Templates() {
		if len(tmpl.Slots) == count {
			matches = append(matches, TemplateSummary{
				Name:        tmpl.Name,
				Description: tmpl.Description,
				Slots:       len(tmpl.Slots),
			})
		}
	}
	return matches
}

// ValidateSlotCount ensures that the template slot counts cover the requested number of images.
func ValidateSlotCount(count int) error {
	if count <= 0 {
		return errors.New("จำนวนรูปต้องมากกว่า 0")
	}
	if _, ok := templateCountIndex()[count]; ok {
		return nil
	}
	return fmt.Errorf("ยังไม่มี template สำหรับ %d รูป", count)
}

func templateCountIndex() map[int]struct{} {
	index := make(map[int]struct{})
	for _, tmpl := range templates {
		index[len(tmpl.Slots)] = struct{}{}
	}
	return index
}
