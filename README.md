# repo-to-txt

**repo-to-txt** is a versatile Command-Line Interface (CLI) tool written in Go that consolidates all contents of a GitHub repository into a single `.txt` file. The output file is automatically named after the repository, ensuring organized and easily identifiable documentation of repository contents.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
  - [Building from Source](#building-from-source)
  - [Using Pre-built Binaries](#using-pre-built-binaries)
- [Usage](#usage)
  - [Interactive Mode](#interactive-mode)
  - [Command-Line Flags](#command-line-flags)
- [Excluding Specific Folders](#excluding-specific-folders)
  - [Interactive Exclusions](#interactive-exclusions)
  - [Command-Line Exclusions](#command-line-exclusions)
- [Including Specific File Extensions](#including-specific-file-extensions)
  - [Command-Line Inclusion](#command-line-inclusion)
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
- **Single Consolidated File**: Merges all repository contents into one `.txt` file with clear file path separators.
- **Support for Public and Private Repositories**: Clone public repositories without authentication or private repositories using HTTPS or SSH.
- **Excluding Specific Folders**: Specify folders to exclude from the output using command-line flags or interactive prompts.
- **Including Specific File Extensions**: Optionally include only specified file extensions to focus on relevant files.
- **Flexible Input Methods**: Supports both interactive prompts and command-line flags for providing inputs.
- **Cross-Platform Compatibility**: Works seamlessly on Windows, macOS, and Linux.
- **Security**: Handles sensitive information like Personal Access Tokens (PATs) securely.

## Prerequisites

- **Go**: Ensure that [Go](https://golang.org/dl/) is installed on your system. The tool is compatible with Go version 1.16 and above.
- **Git**: Required for cloning repositories. Ensure that [Git](https://git-scm.com/downloads) is installed and available in your system's PATH.

## Installation

### Building from Source

1. **Clone the Repository**

   ```sh
   git clone https://github.com/vytautas-bunevicius/repo-to-txt.git
   cd repo-to-txt
   ```

2. **Build the Executable**

   ```sh
   go build -o repo-to-txt main.go
   ```

   This command compiles the Go program and generates an executable named `repo-to-txt` in the current directory.

3. **(Optional) Move Executable to PATH**

   To run `repo-to-txt` from anywhere, move it to a directory that's in your system's PATH, such as `/usr/local/bin` on Unix-based systems.

   ```sh
   sudo mv repo-to-txt /usr/local/bin/
   ```

### Using Pre-built Binaries

Pre-built binaries may be available for different operating systems. Check the [Releases](https://github.com/vytautas-bunevicius/repo-to-txt/releases) section of the repository to download the appropriate binary for your system.

1. **Download the Binary**

   Navigate to the [Releases](https://github.com/vytautas-bunevicius/repo-to-txt/releases) page and download the binary corresponding to your operating system.

2. **Make the Binary Executable**

   ```sh
   chmod +x repo-to-txt
   ```

3. **Move Executable to PATH**

   ```sh
   sudo mv repo-to-txt /usr/local/bin/
   ```

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
Welcome to repo-to-txt!
Enter the GitHub repository URL (HTTPS or SSH):: https://github.com/vytautas-bunevicius/web-comment-monitor.git
? Select authentication method:  [Use arrow keys]
❯ No Authentication
  HTTPS with PAT
  SSH
Enter folders to exclude (comma-separated, leave empty to include all): vendor, tests
Repository contents written to web-comment-monitor.txt
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
- `-exclude`: Comma-separated list of folders to exclude from the output.
- `-include-ext`: Comma-separated list of file extensions to include (e.g., `.go,.md`). If not set, defaults to excluding certain non-code files like `.ipynb`.

**Note**: The output file is automatically named after the repository (e.g., `repository-name.txt`).

## Excluding Specific Folders

You can specify folders that you want to exclude from the `.txt` output. This can be done either interactively or via command-line flags.

### Interactive Exclusions

When running the tool in interactive mode, after selecting the authentication method, you will be prompted to enter folders to exclude.

**Sample Interaction:**

```
Enter folders to exclude (comma-separated, leave empty to include all): vendor, tests
```

### Command-Line Exclusions

Use the `-exclude` flag followed by a comma-separated list of folder names to exclude.

**Usage Example:**

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/web-comment-monitor.git -auth=none -exclude="vendor,tests"
```

## Including Specific File Extensions

By default, the tool excludes non-code files like `.ipynb`. You can specify which file extensions to include using the `-include-ext` flag.

### Command-Line Inclusion

Use the `-include-ext` flag followed by a comma-separated list of file extensions to include only those files in the output.

**Usage Example:**

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/web-comment-monitor.git -auth=none -include-ext=".go,.md"
```

**Note**: If `-include-ext` is not set, the tool defaults to excluding non-code files such as `.ipynb`.

## Authentication Methods

### No Authentication

Use this method to clone **public** repositories without providing any authentication details.

**Usage Example:**

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/web-comment-monitor.git -auth=none
```

### HTTPS Authentication with PAT

Use this method to clone **private** repositories using your GitHub username and Personal Access Token (PAT).

**Usage Example:**

```sh
repo-to-txt -repo=https://github.com/your-username/private-repo.git -auth=https -username=your_username -pat=your_PAT
```

**Security Note**: Be cautious when using command-line flags to provide sensitive information like PATs, as they can be exposed in process listings. Consider using environment variables or interactive prompts for improved security.

### SSH Authentication

Use this method to clone **private** repositories using SSH keys.

**Usage Example:**

```sh
repo-to-txt -repo=git@github.com:your-username/private-repo.git -auth=ssh -ssh-key=/path/to/id_rsa -exclude="vendor,tests"
```

**Prerequisites**:

- Ensure that your SSH key is added to your GitHub account.
- The default SSH key path is `~/.ssh/id_rsa`. If your key is located elsewhere, specify the path using the `-ssh-key` flag.

## Examples

### 1. Cloning a Public Repository Without Authentication and Excluding Folders

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/web-comment-monitor.git -auth=none -exclude="vendor,tests"
```

**Output**:

```
Welcome to repo-to-txt!
Cloning repository https://github.com/vytautas-bunevicius/web-comment-monitor.git...
...
Repository contents written to web-comment-monitor.txt
```

### 2. Cloning a Private Repository Using HTTPS Authentication and Excluding Folders

```sh
repo-to-txt -repo=https://github.com/your-username/private-repo.git -auth=https -username=your_username -pat=your_PAT -exclude="vendor,tests"
```

**Output**:

```
Welcome to repo-to-txt!
Cloning repository https://github.com/your-username/private-repo.git...
...
Repository contents written to private-repo.txt
```

### 3. Cloning a Private Repository Using SSH Authentication Without Excluding Any Folders

```sh
repo-to-txt -repo=git@github.com:your-username/private-repo.git -auth=ssh -ssh-key=/path/to/id_rsa
```

**Output**:

```
Welcome to repo-to-txt!
Cloning repository git@github.com:your-username/private-repo.git...
...
Repository contents written to private-repo.txt
```

### 4. Using Interactive Mode with Exclusions

Run the tool without any flags and follow the interactive prompts.

```sh
repo-to-txt
```

**Sample Interaction:**

```
Welcome to repo-to-txt!
Enter the GitHub repository URL (HTTPS or SSH):: https://github.com/vytautas-bunevicius/web-comment-monitor.git
? Select authentication method:  [Use arrow keys]
❯ No Authentication
  HTTPS with PAT
  SSH
Enter folders to exclude (comma-separated, leave empty to include all): vendor, tests
Repository contents written to web-comment-monitor.txt
```

### 5. Including Specific File Extensions

Clone a repository and include only `.go` and `.md` files in the output.

```sh
repo-to-txt -repo=https://github.com/vytautas-bunevicius/web-comment-monitor.git -auth=none -include-ext=".go,.md"
```

**Output**:

```
Welcome to repo-to-txt!
Cloning repository https://github.com/vytautas-bunevicius/web-comment-monitor.git...
...
Repository contents written to web-comment-monitor.txt
```

## Error Handling

The tool provides descriptive error messages to help you troubleshoot issues. Common errors include:

- **Invalid Repository URL**: Ensure that the repository URL is correct and accessible.
- **Authentication Failures**: Verify your authentication credentials or SSH key.
- **Network Issues**: Check your internet connection and firewall settings.
- **Permission Issues**: Ensure you have the necessary permissions to clone the repository and write to the output directory.

**Example Error Message:**

```
2024/09/13 10:43:38 Error cloning/pulling repository: failed to pull repository: invalid auth method
```

**Resolution**:

- If cloning a **public** repository, select "No Authentication" to avoid unnecessary authentication errors.
- Ensure that your PAT has the necessary scopes if using HTTPS authentication.
- Verify that your SSH key is correctly added to your GitHub account if using SSH authentication.