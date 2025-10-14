package manager

import (
	"fmt"
	"os"

	"github.com/trustwallet/assets-go-libs/file"
	"github.com/trustwallet/assets-go-libs/path"
	"github.com/trustwallet/assets/internal/config"
	"github.com/trustwallet/assets/internal/photobooth"
	"github.com/trustwallet/assets/internal/processor"
	"github.com/trustwallet/assets/internal/report"
	"github.com/trustwallet/assets/internal/service"
	"github.com/trustwallet/go-primitives/asset"
	"github.com/trustwallet/go-primitives/coin"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configPath, root string

func InitCommands() {
	rootCmd.Flags().StringVar(&configPath, "config", ".github/assets.config.yaml",
		"config file (default is $HOME/.github/assets.config.yaml)")
	rootCmd.Flags().StringVar(&root, "root", ".", "root path to files")

	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(fixCmd)
	rootCmd.AddCommand(updateAutoCmd)
	rootCmd.AddCommand(addTokenCmd)
	rootCmd.AddCommand(addTokenlistCmd)
	rootCmd.AddCommand(addTokenlistExtendedCmd)
	initPhotoboothCommand()
	rootCmd.AddCommand(photoboothCmd)
}

var (
	rootCmd = &cobra.Command{
		Use:   "assets",
		Short: "",
		Long:  "",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: " Execute validation checks",
		Run: func(cmd *cobra.Command, args []string) {
			assetsService := InitAssetsService()
			assetsService.RunJob(assetsService.Check)
		},
	}
	fixCmd = &cobra.Command{
		Use:   "fix",
		Short: "Perform automatic fixes where possible",
		Run: func(cmd *cobra.Command, args []string) {
			assetsService := InitAssetsService()
			assetsService.RunJob(assetsService.Fix)
		},
	}
	updateAutoCmd = &cobra.Command{
		Use:   "update-auto",
		Short: "Run automatic updates from external sources",
		Run: func(cmd *cobra.Command, args []string) {
			assetsService := InitAssetsService()
			assetsService.RunUpdateAuto()
		},
	}

	addTokenCmd = &cobra.Command{
		Use:   "add-token",
		Short: "Creates info.json template for the asset",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal("1 argument was expected")
			}

			err := CreateAssetInfoJSONTemplate(args[0])
			if err != nil {
				log.Fatalf("Can't create asset info json template: %v", err)
			}
		},
	}

	addTokenlistCmd = &cobra.Command{
		Use:   "add-tokenlist",
		Short: "Adds token to tokenlist.json",
		Run: func(cmd *cobra.Command, args []string) {
			handleAddTokenList(args, path.TokenlistDefault)
		},
	}

	addTokenlistExtendedCmd = &cobra.Command{
		Use:   "add-tokenlist-extended",
		Short: "Adds token to tokenlist-extended.json",
		Run: func(cmd *cobra.Command, args []string) {
			handleAddTokenList(args, path.TokenlistExtended)
		},
	}
	photoboothCmd = &cobra.Command{
		Use:   "photobooth [รูปภาพ...]",
		Short: "สร้างภาพ Photo Booth ขนาด 4x6 นิ้ว",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if photoboothList {
				return listPhotoboothTemplates()
			}

			if photoboothTemplate == "" {
				return fmt.Errorf("กรุณาระบุชื่อ template ด้วย --template")
			}

			if err := photobooth.ValidateTemplateName(photoboothTemplate); err != nil {
				return err
			}

			if len(args) == 0 {
				min, max := photobooth.SlotCountRange()
				return fmt.Errorf("กรุณาระบุไฟล์รูปภาพ (%d-%d รูปตาม template)", min, max)
			}

			output := photoboothOutput
			if output == "" {
				output = "photobooth.jpg"
			}

			if err := photobooth.Generate(photoboothTemplate, args, output); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "บันทึกไฟล์เรียบร้อย: %s\n", output)
			return nil
		},
	}
)

var (
	photoboothTemplate string
	photoboothOutput   string
	photoboothList     bool
	photoboothCount    int
)

func initPhotoboothCommand() {
	photoboothCmd.Flags().BoolVar(&photoboothList, "list", false, "แสดง template ที่มีให้เลือก")
	photoboothCmd.Flags().StringVar(&photoboothTemplate, "template", "", "ระบุชื่อ template ที่ต้องการใช้")
	photoboothCmd.Flags().StringVar(&photoboothOutput, "output", "photobooth.jpg", "ไฟล์เอาต์พุตสำหรับบันทึกภาพ")
	photoboothCmd.Flags().IntVar(&photoboothCount, "count", 0, "กรอง template ตามจำนวนรูป")
}

func listPhotoboothTemplates() error {
	summaries := photobooth.Summaries()

	if photoboothCount > 0 {
		if err := photobooth.ValidateSlotCount(photoboothCount); err != nil {
			return err
		}
		summaries = photobooth.FilterTemplatesBySlots(photoboothCount)
	}

	for _, summary := range summaries {
		fmt.Printf("%-18s (%d รูป) - %s\n", summary.Name, summary.Slots, summary.Description)
	}

	return nil
}

func handleAddTokenList(args []string, tokenlistType path.TokenListType) {
	if len(args) != 1 {
		log.Fatal("1 argument was expected")
	}

	c, tokenID, err := asset.ParseID(args[0])
	if err != nil {
		log.Fatalf("Can't parse token: %v", err)
	}

	chain, ok := coin.Coins[c]
	if !ok {
		log.Fatal("Invalid token")
	}

	err = AddTokenToTokenListJSON(chain, args[0], tokenID, tokenlistType)
	if err != nil {
		log.Fatalf("Can't add token: %v", err)
	}
}

func InitAssetsService() *service.Service {
	setup()

	paths, err := file.ReadLocalFileStructure(root, config.Default.ValidatorsSettings.RootFolder.SkipFiles)
	if err != nil {
		log.WithError(err).Fatal("Failed to load file structure.")
	}

	fileService := file.NewService(paths...)
	validatorsService := processor.NewService(fileService)
	reportService := report.NewService()

	return service.NewService(fileService, validatorsService, reportService, paths)
}

func setup() {
	if err := config.SetConfig(configPath); err != nil {
		log.WithError(err).Fatal("Failed to set config.")
	}

	logLevel, err := log.ParseLevel(config.Default.App.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("Failed to parse log level.")
	}

	log.SetLevel(logLevel)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
