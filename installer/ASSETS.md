# 🎨 Installer Assets Guide

This guide explains how to create professional installer assets for Obelisk CLI.

## 📋 Required Assets

### 1. Application Icon (icon.ico)

**Specifications:**

- Format: ICO (Windows Icon)
- Sizes: 16x16, 32x32, 48x48, 256x256 (multi-resolution)
- Color: 32-bit with alpha channel
- Theme: Obelisk/pillar/monument theme

**How to Create:**

#### Option A: Using Online Tools (Easiest)

1. Create a 256x256 PNG with your design
2. Go to https://convertio.co/png-ico/
3. Upload your PNG
4. Download the ICO file
5. Save as `installer/icon.ico`

#### Option B: Using GIMP (Free)

1. Create 256x256 image with transparent background
2. Design your icon (obelisk pillar, monument, etc.)
3. File → Export As → icon.ico
4. Check "Compressed (PNG)" option
5. Save as `installer/icon.ico`

#### Option C: Using Photoshop

1. Create 256x256 canvas
2. Design your icon
3. Install ICO plugin
4. Save as ICO format with multiple sizes

**Design Tips:**

- Use the 🏛️ obelisk/pillar theme
- Keep it simple and recognizable at small sizes
- Use contrasting colors
- Include transparency for rounded corners

---

### 2. Installer Banner (banner.bmp)

**Specifications:**

- Format: BMP (24-bit)
- Size: 493 x 58 pixels
- Location: Top of installer window
- Content: Logo + "Obelisk CLI" text

**How to Create:**

1. **Create Canvas:**
    - Open image editor (GIMP, Photoshop, Paint.NET)
    - New image: 493 x 58 pixels
    - Background: White or gradient

2. **Design:**

    ```
    [🏛️ Icon]  Obelisk CLI - AI-Powered Automated Tech Lead
    ```

    - Left: Small icon (32x32)
    - Center: Product name
    - Right: Tagline or version

3. **Export:**
    - File → Export As → banner.bmp
    - Format: 24-bit BMP (no alpha channel)
    - Save as `installer/banner.bmp`

**Design Tips:**

- Keep text readable at small size
- Use professional fonts (Arial, Segoe UI)
- Match your brand colors
- Test on different Windows themes

---

### 3. Installer Dialog (dialog.bmp)

**Specifications:**

- Format: BMP (24-bit)
- Size: 493 x 312 pixels
- Location: Left side of installer window
- Content: Branding image

**How to Create:**

1. **Create Canvas:**
    - New image: 493 x 312 pixels
    - Vertical orientation

2. **Design Options:**

    **Option A: Gradient + Logo**

    ```
    ┌─────────────┐
    │             │
    │   Gradient  │
    │             │
    │   🏛️ Logo   │
    │             │
    │  Obelisk    │
    │    CLI      │
    │             │
    └─────────────┘
    ```

    **Option B: Feature Showcase**

    ```
    ┌─────────────┐
    │ 🏛️ Obelisk  │
    │             │
    │ ✓ Security  │
    │ ✓ Quality   │
    │ ✓ AI-Powered│
    │             │
    │ Version 0.1 │
    └─────────────┘
    ```

    **Option C: Abstract Design**
    - Geometric patterns
    - Code snippets background
    - Tech-themed graphics

3. **Export:**
    - File → Export As → dialog.bmp
    - Format: 24-bit BMP
    - Save as `installer/dialog.bmp`

**Design Tips:**

- Vertical composition
- Professional appearance
- Not too busy or distracting
- Matches banner style

---

## 🎨 Quick Start Templates

### Using PowerShell to Generate Placeholder Assets

```powershell
# This creates basic placeholder images
# Replace with professional designs later

# Create 256x256 icon placeholder
Add-Type -AssemblyName System.Drawing
$bmp = New-Object System.Drawing.Bitmap(256, 256)
$graphics = [System.Drawing.Graphics]::FromImage($bmp)
$graphics.Clear([System.Drawing.Color]::Blue)
$font = New-Object System.Drawing.Font("Arial", 48, [System.Drawing.FontStyle]::Bold)
$brush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::White)
$graphics.DrawString("O", $font, $brush, 80, 80)
$bmp.Save("installer\icon-temp.png")
# Convert PNG to ICO using online tool

# Create banner placeholder
$banner = New-Object System.Drawing.Bitmap(493, 58)
$g = [System.Drawing.Graphics]::FromImage($banner)
$g.Clear([System.Drawing.Color]::White)
$font = New-Object System.Drawing.Font("Arial", 16, [System.Drawing.FontStyle]::Bold)
$brush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::Black)
$g.DrawString("Obelisk CLI", $font, $brush, 10, 15)
$banner.Save("installer\banner.bmp")

# Create dialog placeholder
$dialog = New-Object System.Drawing.Bitmap(493, 312)
$g = [System.Drawing.Graphics]::FromImage($dialog)
$g.Clear([System.Drawing.Color]::LightBlue)
$font = New-Object System.Drawing.Font("Arial", 24, [System.Drawing.FontStyle]::Bold)
$brush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::White)
$g.DrawString("Obelisk CLI", $font, $brush, 150, 130)
$dialog.Save("installer\dialog.bmp")
```

---

## 🎯 Professional Design Services

If you want professional assets:

### Freelance Platforms

- **Fiverr** - $5-50 for icon + banners
- **Upwork** - Professional designers
- **99designs** - Design contests

### DIY Tools

- **Canva** - Easy templates
- **Figma** - Professional design
- **GIMP** - Free Photoshop alternative

### AI Generation

- **DALL-E** - Generate icon concepts
- **Midjourney** - Create unique designs
- **Stable Diffusion** - Free AI art

---

## ✅ Validation Checklist

Before using your assets:

### Icon (icon.ico)

- [ ] 256x256 resolution included
- [ ] ICO format (not PNG/JPG)
- [ ] Transparent background
- [ ] Looks good at 16x16 (taskbar size)
- [ ] File size < 100KB

### Banner (banner.bmp)

- [ ] Exactly 493 x 58 pixels
- [ ] 24-bit BMP format
- [ ] Text is readable
- [ ] Matches brand colors
- [ ] File size < 500KB

### Dialog (dialog.bmp)

- [ ] Exactly 493 x 312 pixels
- [ ] 24-bit BMP format
- [ ] Professional appearance
- [ ] Not too busy
- [ ] File size < 1MB

---

## 🔧 Testing Your Assets

```powershell
# Build installer with your assets
.\installer\build-installer.ps1

# Run the installer to preview
.\installer\ObeliskCLI-0.1.0-x64.msi

# Check:
# - Icon appears in installer window
# - Banner looks good at top
# - Dialog image displays on left
# - All text is readable
```

---

## 📚 Resources

- [Windows Icon Guidelines](https://docs.microsoft.com/en-us/windows/apps/design/style/iconography)
- [WiX UI Customization](https://wixtoolset.org/documentation/manual/v3/wixui/wixui_customizations.html)
- [Icon Design Best Practices](https://material.io/design/iconography)

---

**Note:** The installer will work without these assets (using defaults), but professional assets make your software look more trustworthy and polished!
