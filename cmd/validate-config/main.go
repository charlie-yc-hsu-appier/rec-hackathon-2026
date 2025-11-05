package main

import (
	"fmt"
	"os"
	"path"
	"rec-vendor-api/internal/config"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

// VendorsOnlyConfig represents a config with only vendors section
type VendorsOnlyConfig struct {
	Vendors []config.Vendor `mapstructure:"vendors"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: validate_config <vendors.yaml>")
		fmt.Println("Example: validate_config ./deploy/rec-vendor-api/secrets/vendors.yaml")
		os.Exit(1)
	}

	configPath := os.Args[1]

	cfg := &VendorsOnlyConfig{}

	// Use the existing config loader
	err := loadVendorConfig(configPath, cfg)
	if err != nil {
		fmt.Printf("❌ Configuration validation failed: %v\n", err)
		os.Exit(1)
	}

	// Validate the loaded configuration for supported macros only
	err = validateVendors(cfg.Vendors)
	if err != nil {
		fmt.Printf("❌ Macro validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Macro validation successful!\n")
	fmt.Printf("📊 Validated %d vendors for supported macros:\n", len(cfg.Vendors))
	for i, vendor := range cfg.Vendors {
		fmt.Printf("  %d. %s\n", i+1, vendor.Name)
	}
}

func loadVendorConfig(configPath string, cfg *VendorsOnlyConfig) error {
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", configPath)
	}

	// Load vendors-only configuration
	return loadVendorsOnly(configPath, cfg)
}

func loadVendorsOnly(configPath string, cfg *VendorsOnlyConfig) error {
	configName := path.Base(configPath)
	ext := path.Ext(configPath)
	dir := path.Dir(configPath)

	if ext != ".yaml" {
		return fmt.Errorf("only accept .yaml file")
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(configName)
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Try to unmarshal vendors directly
	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func validateVendors(vendors []config.Vendor) error {
	if len(vendors) == 0 {
		return fmt.Errorf("no vendors found in configuration")
	}

	var errors []string

	for _, vendor := range vendors {
		// Only validate URL macros - skip other validations
		if err := validateMacros(vendor.Request.URL, vendor.Name, "request.url"); err != nil {
			errors = append(errors, err.Error())
		}

		if err := validateMacros(vendor.Tracking.URL, vendor.Name, "tracking.url"); err != nil {
			errors = append(errors, err.Error())
		}

		// Validate query macros
		for _, query := range vendor.Request.Queries {
			if err := validateMacros(query.Value, vendor.Name, fmt.Sprintf("request.queries.%s", query.Key)); err != nil {
				errors = append(errors, err.Error())
			}
		}

		for _, query := range vendor.Tracking.Queries {
			if err := validateMacros(query.Value, vendor.Name, fmt.Sprintf("tracking.queries.%s", query.Key)); err != nil {
				errors = append(errors, err.Error())
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("found %d validation errors:\n- %s", len(errors), strings.Join(errors, "\n- "))
	}

	return nil
}

func validateMacros(text, vendorName, field string) error {
	// Extract macros from text
	macroRegex := regexp.MustCompile(`\{([^}]+)\}`)
	matches := macroRegex.FindAllStringSubmatch(text, -1)

	supportedMacros := map[string]bool{
		"width":              true,
		"height":             true,
		"user_id_lower":      true,
		"user_id_case_by_os": true,
		"click_id_base64":    true,
		"click_id":           true,
		"web_host":           true,
		"bundle_id":          true,
		"adtype":             true,
		"partner_id":         true,
		"subid":              true,
		"keeta_campaign_id":  true,
		"client_ip":          true,
		"latitude":           true,
		"longitude":          true,
		"product_url":        true,
	}

	for _, match := range matches {
		macro := match[1]
		if !supportedMacros[macro] {
			return fmt.Errorf("vendor %s: unsupported macro {%s} in %s", vendorName, macro, field)
		}
	}

	return nil
}
