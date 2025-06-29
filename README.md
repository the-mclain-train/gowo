# gowo

<p align="center">
  <img src="./assets/gowo-logo.png" alt="gowo logo" width="200"/>
</p>

<p align="center">
  <strong>A CLI tool to easily manage Go workspaces by project definition(s).</strong>
</p>

---

`gowo` is a command-line tool designed to streamline the management of complex Go development environments. It allows you to define "projects" as named collections of Git repositories. With a single command, `gowo` can then create a fully configured Go workspace on your filesystem, cloning all necessary repositories and initializing a `go.work` file for you.

## Features

-   **Project Definitions:** Define simple, reusable projects as a list of repositories.
-   **Automated Workspace Creation:** Clones all repositories for a project into a dedicated directory.
-   **Automatic `go.work` Management:** Initializes a Go workspace (`go work init`) and automatically discovers and links all Go modules (`go work use`).
-   **Simple Configuration:** Uses a single YAML configuration file (viper) to manage all settings and project definitions.

## Installation

Make sure you have Go installed (v1.18 minimum) and that your `$GOPATH/bin` directory is in your system's `PATH`.

```sh
go install github.com/the-mclain-train/gowo@latest
```

`gowo` also requires `git` to be installed and available in your `PATH` to clone repositories.

## Quick Start

1.  **Set the root directory for your workspaces.** This is where `gowo` will create the workspace folders. If not set, will default to `~/gowo`.

    ```sh
    # You can use a tilde for your home directory
    gowo config --root "~/Developer/Workspaces"
    ```

2.  **Define a new project.** A project is just a name and a list of Git repositories (e.g., `github.com/user/repo`).

    ```sh
    gowo project add my-cool-project github.com/org/api github.com/org/webapp
    ```

3.  **Create the workspace.** This command will create a new folder, clone the repositories, and set up the `go.work` file.

    ```sh
    gowo workspace create --project my-cool-project --name my-cool-project-ws
    ```

4.  **You're ready!** Navigate to your new workspace directory.

    ```sh
    cd ~/Developer/Workspaces/my-cool-project-ws
    ```

## Usage

`gowo` is structured around three main commands: `config`, `project`, and `workspace`.

---

### `gowo config`

Manage the tool's configuration.

#### `gowo config --root <path>`

Sets the default parent directory where new workspaces will be created. This is a optional first step (will default to `~/gowo`).

-   **`--root`**, **`-r`** (required): The absolute or tilde-prefixed path to your workspaces directory.

**Example:**

```sh
gowo config --root "~/dev"
```

---

### `gowo project`

Manage project definitions. A project is a named collection of git repositories.
Alias: `p`

#### `gowo project add <projectName> <repo1> [repo2...]`

Adds a new project definition with a name and a list of repositories. The repositories should be listed in `host/user/repo` format.

**Example:**

```sh
gowo project add my-api-service github.com/corp/auth-svc github.com/corp/user-svc
```

#### `gowo project ls`

Lists all projects saved in your configuration file.
Alias: `list`

**Example:**

```sh
gowo project ls
```

#### `gowo project show <projectName>`

Displays the list of repositories associated with a given project.

**Example:**

```sh
gowo project show my-api-service
```

#### `gowo project rm <projectName>`

Removes a project definition from the configuration.
Alias: `remove`

**Example:**

```sh
gowo project rm my-api-service
```

#### `gowo project modify`

Modifies a project by adding or removing a single repository.

-   **`--project`**, **`-p`** (required): The name of the project to modify.
-   **`--repo`**, **`-R`** (required): The repository to add or remove.
-   **`--add`**, **`-a`**: The flag to add the repository.
-   **`--remove`**, **`-r`**: The flag to remove the repository.

You must specify exactly one of `--add` or `--remove`.

**Examples:**

```sh
# Add a new repo to the project
gowo project modify -p my-api-service -a -R github.com/corp/billing-svc

# Remove a repo from the project
gowo project modify -p my-api-service -r -R github.com/corp/user-svc
```

---

### `gowo workspace`

Manage Go workspaces. A workspace is a directory containing the cloned repositories of a project and a `go.work` file.
Alias: `ws`

#### `gowo workspace create`

Creates a new Go workspace from a project definition. This is the core command of `gowo`. It performs the following steps:
1.  Creates a new directory for the workspace.
2.  Clones all repositories defined in the project into that directory.
3.  Runs `go work init`.
4.  Scans for all `go.mod` files and links them using `go work use`.

-   **`--project`**, **`-p`** (required): The name of the project to use as a template.
-   **`--name`**, **`-n`** (required): The name for the new workspace directory.
-   **`--directory`**, **`-d`** (optional): A parent directory to create the workspace in, overriding the one set in the config.

**Example:**

```sh
# Create a workspace named "my-feature-branch" from the "my-api-service" project
gowo workspace create --project my-api-service --name my-feature-branch
```

## Configuration

`gowo` stores all its configuration in a single YAML file located at:

`~/.config/gowo/config.yaml`

This file is created automatically when you run a command for the first time (like `gowo config`). It stores the workspaces root directory and all your project definitions.

**Example `config.yaml`:**

```yaml
config:
  workspaces_directory: /Users/tom/dev
projects:
  my-api-service:
    repositories:
      - github.com/corp/auth-svc
      - github.com/corp/user-svc
      - github.com/corp/billing-svc
  another-project:
    repositories:
      - github.com/another-org/tool
```

## License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.
