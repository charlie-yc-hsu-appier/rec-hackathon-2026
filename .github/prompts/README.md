# Vendor onboarding prompts

Used by the **Update Vendor** workflow (`.github/workflows/update-vendor.yml`).

- **vendor-onboard-prompt.md** – Template for adding/updating a vendor in `config-template/vendors.yaml`. The placeholder `{{VENDOR_JSON}}` is replaced with the vendor JSON parsed from the GitHub issue.

The issue form (AI provider: Gemini or Claude) chooses which AI runs the prompt.
