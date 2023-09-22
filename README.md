
<p align="center">
<img src="material/diggity-black.png" style="display: block; margin-left: auto; margin-right: auto; width: 50%;">
</p>

# BOM Diggity
[![Github All Releases](https://img.shields.io/github/downloads/carbonetes/diggity/total.svg)]()
[![Go Report Card](https://goreportcard.com/badge/github.com/carbonetes/diggity)](https://goreportcard.com/report/github.com/carbonetes/diggity)
[![GitHub release](https://img.shields.io/github/release/carbonetes/diggity.svg)](https://github.com/carbonetes/diggity/releases/latest)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/carbonetes/diggity.svg)](https://github.com/carbonetes/diggity)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/carbonetes/diggity/blob/main/LICENSE)

BOM Diggity is an innovative open-source tool developed to streamline the critical process of generating comprehensive Software Bill of Materials (SBOMs) for software projects of all sizes. In an era where software supply chain security and transparency have taken center stage, having an accurate and up-to-date SBOM is essential for developers, organizations, and the entire software ecosystem.

## Integration with Jacked
Diggity seamlessly integrates with our complementary open-source project, [Jacked](https://github.com/carbonetes/jacked). Jacked specializes in open-source vulnerability analysis, providing robust security insights for your software components. Diggity and Jacked offer a comprehensive solution for not only generating SBOMs but also assessing and mitigating security risks in your software supply chain.

### The Significance of a Software Bill of Materials (SBOM)
A Software Bill of Materials (SBOM) serves as a holistic inventory, cataloging every component, dependency, and third-party library that makes up a software application. This transparency-rich resource offers numerous benefits:

- **Security Assurance**: Identifying and addressing vulnerabilities within your software stack is made more accessible with a well-maintained SBOM.
- **Compliance Confidence**: Ensure compliance with licensing and legal requirements by having a clear understanding of your software's composition.
- **Enhanced Trust** Disclose your software's building blocks to users and stakeholders, fostering trust and transparency.
- **Operational Efficiency**: Streamline maintenance, updates, and collaboration within your development teams.

### Key Features That Empower You
Diggity empowers developers, DevOps teams, and organizations with a range of features designed to make SBOM generation and management a seamless part of your software development process:
- **Automated Scanning**: Diggity automatically scans your project's source code and dependencies, intelligently piecing together an SBOM. This automation significantly reduces the manual effort required for creating and maintaining SBOMs.
- **Multiple SBOM Formats**: Flexibility is key. Diggity supports various industry-standard SBOM formats, including CycloneDX and SPDX, ensuring compatibility with your existing toolchain and compliance requirements.
- **Customization Options**: Tailor Diggity's SBOM generation process to the unique needs of your project. Customize component identification rules and output formats to match your project's specifics.
- **Seamless Integration**: Integrate Diggity into your CI/CD pipelines effortlessly. Continuously update SBOMs as your project evolves, ensuring that your SBOMs remain accurate and up-to-date.
- **Detailed Reporting**: Stay informed with detailed reports. Diggity provides insights into identified components, vulnerabilities, and licensing information, enabling proactive risk management and decision-making.

## Supported Ecosystems 

### Package Managers and Build Tools
- APK (/apk/db/installed)
- DPKG (/dpkg/status)
- RPM (Packages, Packages.db, rpmdb.sqlite)
- Pacman (/packman/local/*)
- Conan (conan.lock, conanfile.txt)
- Pub Package Manager (pubspec.yaml, pubspec.lock)
- NPM (package.json, package-lock.json)
- Yarn (yarn.json)
- PNPM (pnpm-lock.yaml)
- NuGet (*.deps.json)
- Go Modules (go.mod, /gobin/*)
- Cabal (stack.yaml, stack.yaml.lock, cabal.project.freeze)
- Hex (rebar.lock, mix.lock)
- Maven (pom.xml, pom.properties, MANIFEST.MF)
- Graddle (buildscript-gradle.lockfile, .build.gradle)
- Composer (composer.lock)
- Pip (wheel, *.egg-info, requirements.txt, METADATA)
- Poetry (poetry.lock)
- RubyGems (*.gemspec, Gemfile.lock)
- Cargo (cargo.lock)
- Cocoapods (Podfile.lock)
- Swift Package Manager (Package.resolved, .package.resolved)
- Nix (/nix/store/*)

### Languages
- Java
- Python
- PHP
- Javascript
- Rust
- Swift
- Objective-C
- Ruby
- C/C++
- Go
- Dart
- C#/F#/Visual Basic
- Haskell
- Erlang

### Plugins
- Jenkins Plugins (*.jpi, *.hpi)

## Installation

### Distributions
- Mac
  - darwin_amd64.tar.gz
  - darwin_arm64.tar.gz
- Linux
  - deb
    - linux_amd64.deb
    - linux_arm64.deb
    - linux_ppc64le.deb
  - rpm
    - linux_amd64.rpm
    - linux_arm64.rpm
    - linux_ppc64le.rpm
  - archive
    - linux_amd64.tar.gz
    - linux_arm64.tar.gz
    - linux_ppc64le.tar.gz
- Windows
  - windows_amd64.zip
### Recommended
```bash
curl -sSfL https://raw.githubusercontent.com/carbonetes/diggity/main/install.sh | sh -s -- -d /usr/local/bin
```
you can specify a release version and destination directory for the installation:

```
curl -sSfL https://raw.githubusercontent.com/carbonetes/diggity/main/install.sh | sh -s -- -d <DESTINATION_DIR> -v <RELEASE_VERSION>
```
### Homebrew
```sh
brew tap carbonetes/diggity
brew install diggity
```

### Scoop
```sh
scoop bucket add diggity https://github.com/carbonetes/diggity-bucket
scoop install diggity
```

## Getting Started

**Note**: Before you begin, make sure you have both Diggity and Docker installed on your system.

Start by pulling the container image for which you want to generate an SBOM. You can use the docker pull command to retrieve the image from a container registry. Replace `your-image:tag` with the actual image and tag you want to analyze.
```bash
docker pull your-image:tag
```

Diggity can now analyze your container image to identify its software components and generate an SBOM. Run the following command to perform the analysis:
```
diggity your-image:tag -o sbom.json
```
- `your-image:tag`: Replace this with the name and tag of the container image you pulled.
- `-o sbom.json`: This option specifies the output file for the generated SBOM. You can choose any file name and format you prefer.
<br />

Diggity will inspect the container image's filesystem and metadata to identify installed packages, libraries, and other dependencies.

<!-- <details>
<summary>Result</summary>

![Diggity](material/diggity.gif)

</details> -->

## Scanning Tarball and Directory
Use the following command to analyze the contents of the Tar file:
```bash
diggity /path/to/your/file.tar -o sbom.json
```
- `/path/to/your/file.tar`: Replace this with the actual path to your Tar file.
<br />

Diggity will inspect the contents of the Tar file and identify software components and dependencies.


And, to analyze the contents of the directory:
```bash
diggity /path/to/your/directory -o sbom.json
```
- `/path/to/your/directory`: Replace this with the actual path to your directory.
<br />

Diggity will scan the directory's files and identify software components, libraries, and dependencies.

## Secret detection
  - User-defined patterns
  - Efficient scanning of container images

<details>
<summary>Result</summary>

![Diggity](material/secret-result.png)

</details>

## Useful Commands and Flags 🚩
```
diggity [command] [flag]
```
### Available Commands and their flags with description:

```
diggity config [flag]
```
|     Flag      |               Description                |
| :------------ | :--------------------------------------- |
| `-d, --display` | Displays the contents of the configuration file. |
| `-h, --help` | Help for configuration.       |
| `-p, --path` | Displays the path of the configuration file.          |
| `-r, --reset` | Restores default configuration file.   |

### Output formats

The output format for Diggity is also configurable using the
`-o` (or `--output`) option:

Available `formats` include:
- `table`: A columnar summary (default).
- `json`: Use this to get as much information out of Diggity.
- `cyclonedx-xml`: An XML report conforming to the [CycloneDX 1.5 specification](https://github.com/CycloneDX/specification/blob/master/schema/bom-1.5.xsd).
- `cyclonedx-json`: A JSON report conforming to the [CycloneDX 1.5 specification](https://github.com/CycloneDX/specification/blob/master/schema/bom-1.5.schema.json).
- `spdx-tag-value`: A tag-value formatted report conforming to the [SPDX 2.3 specification](https://github.com/spdx/spdx-spec/blob/development/v2.3.1/examples/SPDXTagExample-v2.3.spdx).
- `spdx-json`: A JSON report conforming to the [SPDX 2.3 JSON Schema](https://github.com/spdx/spdx-spec/blob/development/v2.3.1/examples/SPDXJSONExample-v2.3.spdx.json) format.
- `spdx-yml`: A YAML report conforming to the [SPDX 2.3 YAML Schema](https://github.com/spdx/spdx-spec/blob/development/v2.3.1/examples/SPDXYAMLExample-2.3.spdx.yaml) format.
- `github-json`: A JSON report conforming to the [dependency snapshot](https://docs.github.com/en/rest/dependency-graph/dependency-submission?apiVersion=2022-11-28) format of Github.


## Configuration

Configuration search paths:
- `<HOME>/.diggity.yaml`

Configuration options (example values are the default):

```yaml
secret-config:
  # enables/disables parsing of secrets
  disabled: false
  # secret content regex are searched within files that match the provided regular expression
  secret-regex: API_KEY|SECRET_KEY|DOCKER_AUTH
  # excludes/includes secret searching for each specified filename
  excludes-filenames: []
  # exclude files exceeding the specified size
  max-file-size: 10485760
  # explicitly define file extensions to consider for secret search. 
  extensions: []  # default extensions are added upon config file generation.
# specify enabled parsers ([apk debian java npm composer python gem rpm dart nuget go rust conan hackage pod hex portage alpmdb]) (default all)
enabled-parsers: []
# disables file listing from package metadata
disable-file-listing: false
# disables the timeout when pulling an image from server
disable-pull-timeout: false
# disable all output except SBOM result
quiet: false
# save the sbom result to the output file instead of writing to standard output
output-file: ""
# supported output types: [json table cyclonedx-xml cyclonedx-json spdx-json spdx-tag-value spdx-yml github-json] (default [table])
output: []
registry: 
  # registry uri endpoint
  uri: ""
  # username credential for private registry access
  username: ""
  # password credential for private registry access
  password: ""
  # access token for private registry access
  token: ""
attestation:
  # path to generated cosign.key
  key: cosign.key
  # path to generated cosign.pub
  pub: cosign.pub
  # password associated with the generated cosign key-pair
  password: ""
```
## Private Registry Authentication
### Local Docker Credentials
When a container image runtime is not present in the local machine, Diggity can pull images from private registries using the provided credentials in your diggity config or as a flag. (--regisytryURI, --registryUsername, (--registryPassword or --registryToken))

An example `.diggity.yaml` looks something like this:
```yaml
registry:
  uri: "https://index.docker.io"
  username: "docker_username"
  password: "docker_password"
  token: ""
```

### AWS ECR Credentials
To pull images from AWS Elastic Container Registry (ECR), provide your account credentials in your diggity config. 
The URI follows the `<aws_account_id>.dkr.ecr.<region>.amazonaws.com` format, and the username would be  `AWS`. 
For the password, run the following command via AWS CLI to obtain your authentication token:

```
aws ecr get-login-password
```

Output:
```
<password>
```
Note that the authentication token is valid for 12 hours. 
For more information, check this [reference](https://docs.aws.amazon.com/cli/latest/reference/ecr/get-login-password.html).

Your `.diggity.yaml` should look something like this:
```yaml
registry:
  uri: "<aws_account_id>.dkr.ecr.<region>.amazonaws.com"
  username: "AWS"
  password: "<password>"
  token: ""
```
### Google Container Registry Credentials
To pull images from Google Container Registry, provide your account credentials in your diggity config. 
The URI follows the `gcr.io, us.gcr.io, eu.gcr.io, or asia.gcr.io` format depending on your service account, and the username would be  `oauth2accesstoken`. 
For the password, run the following command via Google CLI tool to obtain your authentication token:
```
gcloud auth print-access-token
```
Note that the authentication token is valid for about an hour only. 
For more information, check this [reference](https://cloud.google.com/container-registry/docs/advanced-authentication).

Your `.diggity.yaml` should look something like this:
```yaml
registry:
  uri: "gcr.io"
  username: "oauth2accesstoken"
  password: "<token>"
  token: ""
```
### JFrog Container Registry Credentials
To pull images from JFrog Container Registry, provide your account credentials in your diggity config. 
The URI follows the `<server-name>.jfrog.io` format. 
For the password, run the following command in your terminal `docker login -u[username] [server-name].jfrog.io`:

Note that the authentication token is valid for about an hour only. 
For more information, check this [reference](https://www.jfrog.com/confluence/display/JFROG/Getting+Started+with+Artifactory+as+a+Docker+Registry).

Your `.diggity.yaml` should look something like this:
```yaml
registry:
  uri: "diggity.jfrog.io"
  username: "diggity@carbonetes.com"
  password: "<token>"
  token: ""
 ```
## Attestation
Diggity is integrated with [Cosign](https://docs.sigstore.dev/cosign/overview/), which allows you to sign and verify SBOM attestations on images you own. To run attestations, make sure to install Cosign on your machine. Then, generate your cosign key-pair associated with a password using the following command:

```
cosign generate-key-pair
```

This should generate the **cosign.key** and **cosign.pub** files. Specify their respective paths and password in your `.diggity.yaml` config file:

```yaml
attestation:
  key: path/to/cosign.key
  pub: path/to/cosign.pub
  password: "<password>"
```

Alternatively, you could specify the information using flags.

|     Flag      |               Description                |
| :------------ | :--------------------------------------- |
| `-k, --key` | Path to cosign.key used for the SBOM Attestation. |
| `-p, --pub` | Path to cosign.pub used for the SBOM Attestation.       |
| `--password` | Password for the generated cosign key-pair.          |

To run an attestation, make sure that your registry is logged into your machine. Run the following command:

```
diggity attest <image>
```

The attestation metadata can be saved to a file using:

```
diggity attest <image> -f <filename>
```

You can also pass in an already generated SBOM file using the **predicate** flag:

```
diggity attest <image> --predicate <path/to/bom_file>
```

## SLSA Provenance
Include provenance metadata to your SBOMs to provide an additional level of assurance about the secure process used
to build the software. To reference your provenance file, run the following command: 

```
diggity <image> -o json --provenance <path/to/provenance_file>
```

You can also include your provenance metadata in SBOM attestations using the following command:

```
 diggity attest <image> --provenance <path/to/provenance_file>
```

## Contribute to the Project
We enthusiastically welcome contributions from the community! Whether you're interested in reporting issues, submitting pull requests, enhancing documentation, or just offering suggestions, your participation is invaluable. Find more details in our [Contribution Guidelines](https://github.com/carbonetes/diggity/blob/main/CONTRIBUTING.md).

## Get in Touch
Have questions, ideas, or need assistance? Don't hesitate to reach out to us at [eng@carbonetes.com](mailto:eng@carbonetes.com). We're here to support you.
<br />
Diggity is committed to simplifying SBOM generation, enhancing software security, and fostering transparency across the software development landscape. Join us in this mission!

## License
Diggity is released under the permissive [Apache 2.0](https://choosealicense.com/licenses/apache-2.0/), promoting openness and collaboration.