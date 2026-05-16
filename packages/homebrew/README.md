# 🍺 Homebrew Formula for Obelisk CLI

This directory contains the Homebrew formula for Obelisk CLI (macOS and Linux).

## 🎯 What is Homebrew?

Homebrew is the most popular package manager for macOS and Linux. It allows users to install software with:

```bash
brew install obelisk-cli
```

## 📋 Formula File

- **obelisk-cli.rb** - Homebrew formula definition

## 🚀 Testing Locally

### Step 1: Calculate SHA256

```bash
# Download the source tarball
curl -L https://github.com/Swif7ify/Obelisk-CLI/archive/refs/tags/v0.1.0.tar.gz -o obelisk-0.1.0.tar.gz

# Calculate checksum
shasum -a 256 obelisk-0.1.0.tar.gz
```

### Step 2: Update Formula

Edit `obelisk-cli.rb` and update the SHA256:

```ruby
sha256 "YOUR_SHA256_HERE"
```

### Step 3: Test Installation

```bash
# Install from local formula
brew install --build-from-source ./obelisk-cli.rb

# Test the installation
obelisk version

# Run tests
brew test obelisk-cli

# Uninstall
brew uninstall obelisk-cli
```

### Step 4: Audit Formula

```bash
# Check for issues
brew audit --strict --online obelisk-cli.rb

# Fix style issues
brew style obelisk-cli.rb
```

## 📤 Publishing to Homebrew

### Option 1: Homebrew Core (Official)

**Requirements:**

- 30+ GitHub stars
- 75+ forks or 30+ watchers
- Stable project (not pre-release)
- Notable user base

**Process:**

1. **Fork homebrew-core:**

    ```bash
    # Go to: https://github.com/Homebrew/homebrew-core
    # Click "Fork"
    ```

2. **Clone and create branch:**

    ```bash
    git clone https://github.com/YOUR_USERNAME/homebrew-core.git
    cd homebrew-core
    git checkout -b obelisk-cli
    ```

3. **Add formula:**

    ```bash
    cp /path/to/obelisk-cli.rb Formula/obelisk-cli.rb
    ```

4. **Test thoroughly:**

    ```bash
    brew install --build-from-source Formula/obelisk-cli.rb
    brew test obelisk-cli
    brew audit --strict --online Formula/obelisk-cli.rb
    ```

5. **Submit PR:**

    ```bash
    git add Formula/obelisk-cli.rb
    git commit -m "obelisk-cli 0.1.0 (new formula)"
    git push origin obelisk-cli

    # Create PR on GitHub
    ```

### Option 2: Homebrew Tap (Easier, Recommended)

Create your own tap for easier distribution:

1. **Create tap repository:**

    ```bash
    # Create repo: homebrew-obelisk
    # URL: https://github.com/Swif7ify/homebrew-obelisk
    ```

2. **Add formula:**

    ```bash
    git clone https://github.com/Swif7ify/homebrew-obelisk.git
    cd homebrew-obelisk
    mkdir Formula
    cp /path/to/obelisk-cli.rb Formula/
    git add Formula/obelisk-cli.rb
    git commit -m "Add obelisk-cli formula"
    git push
    ```

3. **Users install with:**
    ```bash
    brew tap swif7ify/obelisk
    brew install obelisk-cli
    ```

## 🔄 Updating for New Versions

### Update Formula

```ruby
class ObeliskCli < Formula
  desc "AI-Powered Automated Tech Lead for Modern Codebases"
  homepage "https://github.com/Swif7ify/Obelisk-CLI"
  url "https://github.com/Swif7ify/Obelisk-CLI/archive/refs/tags/v0.2.0.tar.gz"  # Update version
  sha256 "NEW_SHA256_HERE"  # Update checksum
  license "MIT"

  # ... rest of formula
end
```

### Submit Update

```bash
# For homebrew-core
git checkout -b obelisk-cli-0.2.0
git add Formula/obelisk-cli.rb
git commit -m "obelisk-cli 0.2.0"
git push origin obelisk-cli-0.2.0

# For your tap
git add Formula/obelisk-cli.rb
git commit -m "Update obelisk-cli to 0.2.0"
git push
```

## 🧪 Testing Checklist

Before submitting:

- [ ] Formula installs successfully
- [ ] Binary runs: `obelisk version`
- [ ] Tests pass: `brew test obelisk-cli`
- [ ] Audit passes: `brew audit --strict --online obelisk-cli.rb`
- [ ] Style is correct: `brew style obelisk-cli.rb`
- [ ] SHA256 is correct
- [ ] Version matches GitHub release
- [ ] License is specified
- [ ] Description is clear

## 📝 Formula Best Practices

### Naming

- Use lowercase with hyphens: `obelisk-cli`
- Match GitHub repo name when possible
- Avoid redundant suffixes

### Dependencies

```ruby
depends_on "go" => :build  # Build-time only
depends_on "openssl"       # Runtime dependency
```

### Installation

```ruby
def install
  # Build from source
  system "go", "build", *std_go_args(ldflags: "-s -w")

  # Install additional files
  bin.install "obelisk"
  man1.install "docs/obelisk.1"
  bash_completion.install "completions/obelisk.bash"
end
```

### Testing

```ruby
def test
  # Test version
  assert_match version.to_s, shell_output("#{bin}/obelisk version")

  # Test functionality
  system bin/"obelisk", "--help"
end
```

## 🆘 Troubleshooting

### "Formula doesn't install"

```bash
# Check build logs
brew install --verbose --debug obelisk-cli
```

### "SHA256 mismatch"

```bash
# Recalculate checksum
curl -L https://github.com/Swif7ify/Obelisk-CLI/archive/refs/tags/v0.1.0.tar.gz | shasum -a 256
```

### "Audit fails"

```bash
# See specific issues
brew audit --strict --online obelisk-cli.rb

# Common issues:
# - Missing license
# - Incorrect URL
# - Style violations
```

### "Tests fail"

```bash
# Run tests manually
brew install obelisk-cli
brew test obelisk-cli

# Debug
obelisk version
obelisk --help
```

## 📊 Distribution Comparison

| Method         | Difficulty | Reach  | Updates   |
| -------------- | ---------- | ------ | --------- |
| Homebrew Core  | Hard       | High   | Automatic |
| Homebrew Tap   | Easy       | Medium | Manual    |
| Direct Install | Easiest    | Low    | Manual    |

**Recommendation:** Start with a Homebrew Tap, then submit to Core once you have traction.

## 🎯 Homebrew Tap Setup

### Create Tap Repository

1. **Create repo:** `homebrew-obelisk`
2. **Add formula:** `Formula/obelisk-cli.rb`
3. **Users install:**
    ```bash
    brew tap swif7ify/obelisk
    brew install obelisk-cli
    ```

### Tap Structure

```
homebrew-obelisk/
├── Formula/
│   └── obelisk-cli.rb
└── README.md
```

### Tap README

````markdown
# Swif7ify Homebrew Tap

## Installation

```bash
brew tap swif7ify/obelisk
brew install obelisk-cli
```
````

## Formulae

- **obelisk-cli** - AI-Powered Automated Tech Lead

```

## 📚 Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)
- [Homebrew Tap Documentation](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
- [Formula Style Guide](https://docs.brew.sh/Formula-Cookbook#style-guide)

## 🤝 Contributing

To improve the formula:
1. Test on macOS and Linux
2. Ensure all tests pass
3. Follow Homebrew style guide
4. Submit PR with clear description

---

**Note:** Homebrew Tap is the easiest way to distribute. Submit to Homebrew Core later once you have an established user base!
```
