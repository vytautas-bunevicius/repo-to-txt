# repo-to-txt

**repo-to-txt** is a versatile Command-Line Interface (CLI) tool written in Go that consolidates all contents of a GitHub repository into a single `.txt` file. The output file is automatically named after the repository and can be saved to a specified directory, ensuring organized and easily identifiable documentation of repository contents.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
  - [Building from Source](#building-from-source)
  - [Making repo-to-txt Globally Accessible](#making-repo-to-txt-globally-accessible)
- [Usage](#usage)
  - [Interactive Mode](#interactive-mode)
  - [Command-Line Flags](#command-line-flags)
- [Excluding Specific Folders](#excluding-specific-folders)
  - [Interactive Exclusions](#interactive-exclusions)
  - [Command-Line Exclusions](#command-line-exclusions)
- [Including Specific File Extensions](#including-specific-file-extensions)
  - [Command-Line Inclusion](#command-line-inclusion)
- [Including Specific Files](#including-specific-files)
  - [Command-Line File Inclusion](#command-line-file-inclusion)
- [Clipboard Copying](#clipboard-copying)
  - [Installing Clipboard Utilities](#installing-clipboard-utilities)
  - [Using Clipboard Copying](#using-clipboard-copying)
- [Authentication Methods](#authentication-methods)
  - [No Authentication](#no-authentication)
  - [HTTPS Authentication with PAT](#https-authentication-with-pat)
  - [SSH Authentication](#ssh-authentication)
- [Examples](#examples)
- [Error Handling](#error-handling)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Automatic Output Naming**: Generates a `.txt` file named after the repository.
- **Customizable Output Directory**: Allows specifying the directory where the output file should be saved.
- **Single Consolidated File**: Merges all repository contents into one `.txt` file with clear file path separators.
- **Support for Public and Private Repositories**: Clone public repositories without authentication or private repositories using HTTPS or SSH.
- **Excluding Specific Folders**: Specify folders to exclude from the output using command-line flags or interactive prompts.
- **Including Specific File Extensions**: Optionally include only specified file extensions to focus on relevant files.
- **Including Specific Files**: Select exact file names to include in the consolidated `.txt` output, with their paths clearly indicated.
- **Flexible Input Methods**: Supports both interactive prompts and command-line flags for providing inputs.
- **Cross-Platform Compatibility**: Works seamlessly on Windows, macOS, and Linux.
- **Security Enhancements**:
  - Handles sensitive information like Personal Access Tokens (PATs) securely.
  - Supports SSH keys with passphrases, enhancing SSH authentication security.
- **Clipboard Copying**: Optionally copy the generated `.txt` file content directly to the clipboard for quick access.
- **Improved Error Handling and Logging**: Provides more descriptive error messages to aid in troubleshooting.

## Prerequisites

- **Go**: Ensure that [Go](https://golang.org/dl/) is installed on your system. The tool is compatible with Go version 1.16 and above.
- **Git**: Required for cloning repositories. Ensure that [Git](https://git-scm.com/downloads) is installed and available in your system's PATH.
- **Clipboard Utilities** (Optional): To enable the clipboard copying feature, install one of the supported clipboard utilities as detailed below.

## Installation

### Building from Source

1. **Clone the Repository**

   ```sh
   git clone https://github.com/vytautas-bunevicius/repo-to-txt.git
   cd repo-to-txt
   ```

2. **Build the Executable**

   ```sh
   go build -o repo-to-txt ./cmd/repo-to-txt
   ```

   This command compiles the Go program and generates an executable named `repo-to-txt` in the current directory.

### Making repo-to-txt Globally Accessible

To run `repo-to-txt` from any directory, you need to add the directory containing the executable to your system's `PATH` environment variable.

**Unix-like systems (Linux, macOS):**

1. **Move the executable to a directory in your PATH:**
   Common locations include `/usr/local/bin` or `/usr/bin`.

   ```sh
   sudo mv repo-to-txt /usr/local/bin/
   ```

2. **Alternatively, add the executable's current directory to your PATH:**
   Open your shell's configuration file (e.g., `~/.bashrc` or `~/.zshrc`) and add the following line, replacing `/path/to` with the actual path to the executable:

   ```bash
   export PATH="$PATH:/path/to"
   ```

   Then, reload the shell configuration:

   ```bash
   source ~/.bashrc # or source ~/.zshrc
   ```

**Windows:**

1. **Move the executable to a directory in your PATH:**
   You can find the directories in your PATH by searching for "environment variables" in the Start menu and selecting "Edit the system environment variables".

2. **Alternatively, add the executable's current directory to your PATH:**
   - Search for "Environment Variables" in the Start menu.
   - Click on "Environment Variables" in the System Properties window.
   - Under "User variables" or "System variables," find and select the `Path` variable, then click "Edit."
   - Click "New" and add the path to the directory containing `repo-to-txt.exe`.
   - Click "OK" on all open dialog boxes to apply the changes.

## Usage

`repo-to-txt` can be used in two modes:

1. **Interactive Mode**: Prompts the user step-by-step for necessary inputs.
2. **Command-Line Flags**: Allows users to provide all inputs upfront for automation and scripting.

### Interactive Mode

Simply run the executable without any flags, and the tool will guide you through the process.

```sh
repo-to-txt
```

**Sample Interaction:**

```
Enter the GitHub repository URL (HTTPS or SSH): https://github.com/vytautas-bunevicius/repo-to-txt.git
? Select authentication method:
❯ No Authentication
  HTTPS with PAT
  SSH
Enter the output directory (default "/home/user/Downloads"): /path/to/output
Enter folders to exclude (comma-separated, leave empty to include all): vendor, tests
Enter file extensions to include (comma-separated, leave empty to include all): .go,.md
Enter exact file names to copy (comma-separated, leave empty to copy all files): prompt.go, main.go
Do you want to copy the output to the clipboard?
❯ Yes
  No
2024/09/15 10:35:09 Welcome to repo-to-txt!
Cloning repository: https://github.com/vytautas-bunevicius/repo-to-txt.git
Enumerating objects: 132, done.
Total 132 (delta 0), reused 0 (delta 0), pack-reused 132 (from 1)
2024/09/15 10:35:17 Added prompt.go to /path/to/output/repo-to-txt.txt
2024/09/15 10:35:17 Added main.go to /path/to/output/repo-to-txt.txt
2024/09/15 10:35:17 Specified files' contents written to /path/to/output/repo-to-txt.txt
2024/09/15 10:35:17 Specified files' contents have been copied to the clipboard.
```

### Command-Line Flags

Provide necessary inputs via flags to run the tool non-interactively.

```sh
repo-to-txt -repo=<repository_url> -auth=<authentication_method> [additional_flags]
```

**Available Flags:**

- `-repo`: **(Required)** GitHub repository URL (HTTPS or SSH).
- `-auth`: Authentication method. Options: `none`, `https`, `ssh`.
- `-username`: GitHub username (required for HTTPS).
- `-pat`: GitHub Personal Access Token (required for HTTPS).
- `-ssh-key`: Path to SSH private key (required for SSH).
- `-ssh-passphrase`: Passphrase for SSH private key (if protected).
- `-output-dir`: The directory where the output file should be saved. Defaults to the user's Downloads directory.
- `-exclude`: Comma-separated list of folders to exclude from the output.
- `-include-ext`: Comma-separated list of file extensions to include (e.g., `.go,.md`). If not set, defaults to excluding certain non-code files like `.ipynb`.
- `-files`: Comma-separated list of exact file names to copy from the repository.
- `-copy-clipboard`: Copy the output to the clipboard after creation. Options: `true`, `false`.
- `-version`: Print the version number and exit.

**Note**: The output file is automatically named after the repository (e.g., `repository-name.txt`).

**Example Command:**

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/repo-to-txt.git -auth=none -output-dir=/path/to/output -exclude="vendor,tests" -include-ext=".go,.md" -files="prompt.go,main.go" -copy-clipboard=true
```

## Excluding Specific Folders

You can specify folders that you want to exclude from the `.txt` output. This can be done either interactively or via command-line flags.

### Interactive Exclusions

When running the tool in interactive mode, after selecting the authentication method, you will be prompted to enter folders to exclude.

**Sample Interaction:**

```
Enter folders to exclude (comma-separated, leave empty to include all): vendor, tests
```

### Command-Line Exclusions

Use the `-exclude` flag followed by a comma-separated list of folder names to exclude from the output.

**Usage Example:**

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/repo-to-txt.git -auth=none -output-dir=/path/to/output -exclude="vendor,tests"
```

## Including Specific File Extensions

By default, the tool excludes non-code files like `.ipynb`. You can specify which file extensions to include using the `-include-ext` flag.

### Command-Line Inclusion

Use the `-include-ext` flag followed by a comma-separated list of file extensions to include only those files in the output.

**Usage Example:**

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/repo-to-txt.git -auth=none -output-dir=/path/to/output -include-ext=".go,.md"
```

**Note**: If `-include-ext` is not set, the tool defaults to excluding non-code files such as `.ipynb`.

## Including Specific Files

In addition to excluding folders and including specific file extensions, you can select exact file names to include in the consolidated `.txt` output. This allows you to focus on particular files within the repository.

### Command-Line File Inclusion

Use the `-files` flag followed by a comma-separated list of exact file names you wish to include in the output.

**Usage Example:**

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/repo-to-txt.git -auth=none -output-dir=/path/to/output -files="prompt.go,main.go"
```

**Note**:

- If multiple files with the same name exist in different directories, the tool will prompt you to select which one to include.
- The output `.txt` file will include the path of each specified file as a separator before its contents.

**Sample Output (`repo-to-txt.txt`):**

```
=== prompt.go ===
<Contents of prompt.go>

=== main.go ===
<Contents of main.go>
```

## Clipboard Copying

`repo-to-txt` offers an optional feature to copy the generated `.txt` file content directly to the clipboard for quick access.

### Installing Clipboard Utilities

To enable clipboard copying, you need to have one of the supported clipboard utilities installed on your system. The tool relies on external utilities to interact with the system clipboard.

#### **Supported Clipboard Utilities:**

- **For X11 (Most Linux Distributions):**
  - `xclip`
  - `xsel`

- **For Wayland (Modern Linux Display Server):**
  - `wl-clipboard`

- **For Termux (Android Terminal Emulator):**
  - `Termux:API` add-on

#### **Installing `xclip` (Recommended for X11 Environments):**

**Ubuntu/Debian:**

```bash
sudo apt update
sudo apt install xclip
```

**Fedora:**

```bash
sudo dnf install xclip
```

**Arch Linux:**

```bash
sudo pacman -S xclip
```

#### **Installing `xsel`:**

**Ubuntu/Debian:**

```bash
sudo apt update
sudo apt install xsel
```

**Fedora:**

```bash
sudo dnf install xsel
```

**Arch Linux:**

```bash
sudo pacman -S xsel
```

#### **Installing `wl-clipboard` (For Wayland Environments):**

**Ubuntu/Debian:**

```bash
sudo apt update
sudo apt install wl-clipboard
```

**Fedora:**

```bash
sudo dnf install wl-clipboard
```

**Arch Linux:**

```bash
sudo pacman -S wl-clipboard
```

#### **Installing Termux:API (For Termux on Android):**

1. **Install Termux:API Add-on:**
   - Download and install the [Termux:API](https://play.google.com/store/apps/details?id=com.termux.api) app from the Google Play Store or F-Droid.

2. **Install `termux-api` Package:**

   ```bash
   pkg install termux-api
   ```

#### **Verify Installation**

After installing a clipboard utility, verify that it's accessible:

- **For `xclip`:**

  ```bash
  xclip -version
  ```

- **For `xsel`:**

  ```bash
  xsel --version
  ```

- **For `wl-clipboard`:**

  ```bash
  wl-copy --version
  ```

- **For Termux:API:**

  ```bash
  termux-clipboard-get --help
  ```

If the installation was successful, these commands should display version information or help text.

### Using Clipboard Copying

When running `repo-to-txt`, use the `-copy-clipboard` flag to enable clipboard copying.

**Command-Line Example:**

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/repo-to-txt.git -auth=none -output-dir=/path/to/output -copy-clipboard=true
```

**Interactive Mode:**

During interactive prompts, you will be asked whether you want to copy the output to the clipboard. Select **"Yes"** to enable this feature.

**Sample Interaction:**

```
Do you want to copy the output to the clipboard?
❯ Yes
  No
```

**Error Handling:**

If no supported clipboard utility is found, the tool will log an error message but continue execution:

```
2024/09/15 10:35:17 Repository contents written to /path/to/output/repo-to-txt.txt
2024/09/15 10:35:17 Error: failed to copy content to clipboard: No clipboard utilities available. Please install xsel, xclip, wl-clipboard or Termux:API add-on for termux-clipboard-get/set.
```

**Resolution:**

- Install one of the supported clipboard utilities as detailed in the [Clipboard Copying](#clipboard-copying) section.
- Re-run the tool to enable clipboard copying.

**Benefits of Clipboard Copying:**

- **Quick Access:** Easily paste the consolidated repository contents into documents, emails, or other applications without manually opening the `.txt` file.
- **Automation:** Facilitates workflows where repository content needs to be shared or processed further immediately after generation.

## Authentication Methods

`repo-to-txt` supports multiple authentication methods to accommodate both public and private repositories.

### No Authentication

Use this method to clone **public** repositories without providing any authentication details.

**Usage Example:**

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/repo-to-txt.git -auth=none -output-dir=/path/to/output
```

### HTTPS Authentication with PAT

Use this method to clone **private** repositories using your GitHub username and Personal Access Token (PAT).

**Usage Example:**

```sh
repo-to-txt -repo=https://github.com/your-username/private-repo.git -auth=https -username=your_username -pat=your_PAT -output-dir=/path/to/output
```

**Security Note:** Be cautious when using command-line flags to provide sensitive information like PATs, as they can be exposed in process listings. Consider using interactive prompts for improved security.

### SSH Authentication

Use this method to clone **private** repositories using SSH keys.

**Usage Example:**

```sh
repo-to-txt -repo=git@github.com:your-username/private-repo.git -auth=ssh -ssh-key=/path/to/id_rsa -ssh-passphrase="your_passphrase" -output-dir=/path/to/output -exclude="vendor,tests"
```

**Prerequisites:**

- Ensure that your SSH key is added to your GitHub account.
- The default SSH key path is `~/.ssh/id_rsa`. If your key is located elsewhere, specify the path using the `-ssh-key` flag.
- If your SSH key is protected with a passphrase, provide it using the `-ssh-passphrase` flag. If your key does not have a passphrase, you can omit this flag.

## Examples

### 1. Cloning a Public Repository Without Authentication and Excluding Folders to a Specific Directory

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/repo-to-txt.git -auth=none -output-dir=/path/to/output -exclude="vendor,tests"
```

**Output:**

```
2024/09/15 10:35:09 Welcome to repo-to-txt!
Cloning repository: https://github.com/vytautas-bunevicius/repo-to-txt.git
Enumerating objects: 132, done.
Total 132 (delta 0), reused 0 (delta 0), pack-reused 132 (from 1)
2024/09/15 10:35:17 Skipping file .git/index: binary file
2024/09/15 10:35:17 Skipping file .git/objects/pack/pack-df8004320515f54f375aa24042fb596aa7f3a2b3.idx: binary file
2024/09/15 10:35:17 Skipping file .git/objects/pack/pack-df8004320515f54f375aa24042fb596aa7f3a2b3.pack: binary file
2024/09/15 10:35:17 Skipping file repo-to-txt: binary file
2024/09/15 10:35:17 Repository contents written to /path/to/output/repo-to-txt.txt
2024/09/15 10:35:17 Repository contents have been copied to the clipboard.
```

### 2. Cloning a Private Repository Using HTTPS Authentication and Excluding Folders to a Specific Directory

```sh
repo-to-txt -repo=https://github.com/your-username/private-repo.git -auth=https -username=your_username -pat=your_PAT -output-dir=/path/to/output -exclude="vendor,tests"
```

**Output:**

```
2024/09/15 10:35:09 Welcome to repo-to-txt!
Cloning repository: https://github.com/your-username/private-repo.git
Enumerating objects: 132, done.
Total 132 (delta 0), reused 0 (delta 0), pack-reused 132 (from 1)
2024/09/15 10:35:17 Skipping file .git/index: binary file
2024/09/15 10:35:17 Skipping file .git/objects/pack/pack-df8004320515f54f375aa24042fb596aa7f3a2b3.idx: binary file
2024/09/15 10:35:17 Skipping file .git/objects/pack/pack-df8004320515f54f375aa24042fb596aa7f3a2b3.pack: binary file
2024/09/15 10:35:17 Skipping file repo-to-txt: binary file
2024/09/15 10:35:17 Repository contents written to /path/to/output/private-repo.txt
2024/09/15 10:35:17 Repository contents have been copied to the clipboard.
```

### 3. Cloning a Private Repository Using SSH Authentication Without Excluding Any Folders to a Specific Directory

```sh
repo-to-txt -repo=git@github.com:your-username/private-repo.git -auth=ssh -ssh-key=/path/to/id_rsa -output-dir=/path/to/output
```

**Output:**

```
2024/09/15 10:35:09 Welcome to repo-to-txt!
Cloning repository: git@github.com:your-username/private-repo.git
Enumerating objects: 132, done.
Total 132 (delta 0), reused 0 (delta 0), pack-reused 132 (from 1)
2024/09/15 10:35:17 Skipping file .git/index: binary file
2024/09/15 10:35:17 Skipping file .git/objects/pack/pack-df8004320515f54f375aa24042fb596aa7f3a2b3.idx: binary file
2024/09/15 10:35:17 Skipping file .git/objects/pack/pack-df8004320515f54f375aa24042fb596aa7f3a2b3.pack: binary file
2024/09/15 10:35:17 Skipping file repo-to-txt: binary file
2024/09/15 10:35:17 Repository contents written to /path/to/output/private-repo.txt
2024/09/15 10:35:17 Repository contents have been copied to the clipboard.
```

### 4. Using Interactive Mode with Exclusions and Specifying Output Directory

Run the tool without any flags and follow the interactive prompts.

```sh
repo-to-txt
```

**Sample Interaction:**

```
Enter the GitHub repository URL (HTTPS or SSH): https://github.com/vytautas-bunevicius/repo-to-txt.git
? Select authentication method:
❯ No Authentication
  HTTPS with PAT
  SSH
Enter the output directory (default "/home/user/Downloads"): /path/to/output
Enter folders to exclude (comma-separated, leave empty to include all): vendor, tests
Enter file extensions to include (comma-separated, leave empty to include all): .go,.md
Enter exact file names to copy (comma-separated, leave empty to copy all files): prompt.go, main.go
Do you want to copy the output to the clipboard?
❯ Yes
  No
2024/09/15 10:35:09 Welcome to repo-to-txt!
Cloning repository: https://github.com/vytautas-bunevicius/repo-to-txt.git
Enumerating objects: 132, done.
Total 132 (delta 0), reused 0 (delta 0), pack-reused 132 (from 1)
2024/09/15 10:35:17 Added prompt.go to /path/to/output/repo-to-txt.txt
2024/09/15 10:35:17 Added main.go to /path/to/output/repo-to-txt.txt
2024/09/15 10:35:17 Specified files' contents written to /path/to/output/repo-to-txt.txt
2024/09/15 10:35:17 Specified files' contents have been copied to the clipboard.
```

### 5. Including Specific File Extensions and Specifying Output Directory

Clone a repository and include only `.go` and `.md` files in the output, saving it to a specific directory.

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/repo-to-txt.git -auth=none -output-dir=/path/to/output -include-ext=".go,.md"
```

**Output:**

```
2024/09/15 10:35:09 Welcome to repo-to-txt!
Cloning repository: https://github.com/vytautas-bunevicius/repo-to-txt.git
Enumerating objects: 132, done.
Total 132 (delta 0), reused 0 (delta 0), pack-reused 132 (from 1)
2024/09/15 10:35:17 Skipping file .git/index: binary file
2024/09/15 10:35:17 Skipping file .git/objects/pack/pack-df8004320515f54f375aa24042fb596aa7f3a2b3.idx: binary file
2024/09/15 10:35:17 Skipping file .git/objects/pack/pack-df8004320515f54f375aa24042fb596aa7f3a2b3.pack: binary file
2024/09/15 10:35:17 Skipping file repo-to-txt: binary file
2024/09/15 10:35:17 Repository contents written to /path/to/output/repo-to-txt.txt
2024/09/15 10:35:17 Repository contents have been copied to the clipboard.
```

### 6. Cloning a Repository with an SSH Key Passphrase and Specifying Output Directory

If your SSH key is protected with a passphrase, provide it using the `-ssh-passphrase` flag and specify the output directory.

```sh
repo-to-txt -repo=git@github.com:your-username/private-repo.git -auth=ssh -ssh-key=/path/to/id_rsa -ssh-passphrase="your_passphrase" -output-dir=/path/to/output
```

**Output:**

```
2024/09/15 10:35:09 Welcome to repo-to-txt!
Cloning repository: git@github.com:your-username/private-repo.git
Enumerating objects: 132, done.
Total 132 (delta 0), reused 0 (delta 0), pack-reused 132 (from 1)
2024/09/15 10:35:17 Skipping file .git/index: binary file
2024/09/15 10:35:17 Skipping file .git/objects/pack/pack-df8004320515f54f375aa24042fb596aa7f3a2b3.idx: binary file
2024/09/15 10:35:17 Skipping file .git/objects/pack/pack-df8004320515f54f375aa24042fb596aa7f3a2b3.pack: binary file
2024/09/15 10:35:17 Skipping file repo-to-txt: binary file
2024/09/15 10:35:17 Repository contents written to /path/to/output/private-repo.txt
2024/09/15 10:35:17 Repository contents have been copied to the clipboard.
```

## Error Handling

The tool provides descriptive error messages to help you troubleshoot issues. Common errors include:

- **Invalid Repository URL**: Ensure that the repository URL is correct and accessible.
- **Authentication Failures**: Verify your authentication credentials or SSH key.
- **Network Issues**: Check your internet connection and firewall settings.
- **Permission Issues**: Ensure you have the necessary permissions to clone the repository and write to the output directory.
- **SSH Passphrase Errors**: If using an SSH key with a passphrase, ensure that the passphrase is correct.
- **Clipboard Utility Not Found**: If clipboard copying is enabled but no supported clipboard utility is installed, you'll receive an error prompting you to install one.

**Example Error Message:**

```
2024/09/15 10:35:17 Error: failed to copy content to clipboard: No clipboard utilities available. Please install xsel, xclip, wl-clipboard or Termux:API add-on for termux-clipboard-get/set.
```

**Resolution:**

- **For Clipboard Utility Errors:**
  - Install one of the supported clipboard utilities as detailed in the [Clipboard Copying](#clipboard-copying) section.
  - Re-run the tool to enable clipboard copying.

- **For Authentication Errors:**
  - Double-check your GitHub credentials or SSH key setup.
  - Ensure that your PAT has the necessary scopes for repository access.
  - Verify that your SSH key is correctly added to your GitHub account.

- **For Repository Cloning Errors:**
  - Ensure that the repository URL is correct.
  - Verify that you have access permissions for private repositories.
  - Check your network connection and proxy settings if applicable.
