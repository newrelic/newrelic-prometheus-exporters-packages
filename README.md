[![Community Project header](https://github.com/newrelic/opensource-website/raw/master/src/images/categories/Community_Project.png)](https://opensource.newrelic.com/oss-category/#community-project)

# Prometheus Exporter Packages

This project packages several Prometheus exporters as native operating system packages, with the goal of providing a better installation experience for Prometheus exporters.

All native packages are available for installation in New Relic's public repositories.

## Installation

To use the packages generated by this project you need to:

- Determine your OS version.
  - Debian, Red Hat, CentOS, Amazon Linux:

    ```bash
    cat /etc/os-release
    ```

  - Ubuntu

    ```bash
    cat /etc/lsb-release
    ```

  - SuSE Linux Enterprise Server

    ```bash
    cat /etc/os-release | grep VERSION_ID
    ```

- Enable New Relic's GPG key (APT and YUM packages).
  
```bash
    curl -s https://download.newrelic.com/infrastructure_agent/gpg/newrelic-infra.gpg | sudo apt-key add -
```

- Add the New Relic's repository to the operating system package manager.
  - Debian based:

    ```bash
    printf "deb [arch=amd64] https://download.newrelic.com/infrastructure_agent/linux/apt <VERSION> main" | sudo tee -a /etc/apt/sources.list.d/newrelic-infra.list
    ```

    Note: replace \<VERSION> with either **jessie**, **stretch** or **buster**, depending on your Debian version

  - YUM based (Amazon Linux, Amazon Linux 2, RHEL, CentOS):

    ```bash
    sudo curl -o /etc/yum.repos.d/newrelic-infra.repo https://download.newrelic.com/infrastructure_agent/linux/yum/el/<VERSION>/x86_64/newrelic-infra.repo
    ```

     Note:
    - for Amazon Linux replace \<VERSION> with **6** and for Amazon Linux 2 replace with **7**
    - for CentOS and RHEL, replace \<VERSION> with the version you are using, (**5**, **6**, **7** or **8**)

- Refresh the repositories.
  - Debian, Ubuntu

    ```bash
    sudo apt-get update
    ```

  - Amazon Linux, Amazon Linux 2, RHEL, CentOS

    ```bash
    sudo yum -q makecache -y --disablerepo='*' --enablerepo='newrelic-infra'
    ```
  
  - SLES

    ```bash
    sudo zypper -n ref -r newrelic-infra
    ```

## Getting Started

Before installing any of the exporters, be sure to read the documentation for the specific exporter you want to install.
Each exporter may have specific configuration options that will require you to modify to make it work for your environment.

Also make sure you read the documentation about the [Infrastructure agent](https://github.com/newrelic/infrastructure-agent). The agent, in conjunction with our [Prometheus Open Metrics integration](https://github.com/newrelic/nri-prometheus), is the component responsible for sending the metrics, provided by the exporter, to New Relic.

If you followed the instructions in the [installation section](#installation) you can now install the exporter using a single command.

- Ubuntu, Debian

```bash
sudo apt-get install <exporter package name>
```

- Amazon Linux, Amazon Linux 2, RHEL, CentOS

```bash
sudo yum install <exporter package name>
```

- SLES

```bash
sudo zypper install <exporter package name>
```

## Adding a new exporter

In order to add a new exporter a new folder in the path `exporters/{exportername}` should be created. You can refer to `githubactions` exporter example in order to doublecheck parameters and fields of scripts and definitions

In each folder we expect to find:
  - `LICENSE`: license of the exporter taken from its repository
  - `exporter.yml`: definition of the exporter
  - `Makefile` a makefile having the target `all` that builds the exporter and the installation packages
  - `win_build.ps1` a powershell script building the exporter and creating the `msi` package.

The definition file requieres the following fields:
``` yaml
# name of the exporter, should mach with the folder name
name: githubactions
# version of the package created
version: 1.2.2
# URL to the git project hosting the exporter
exporter_repo_url: https://github.com/Spendesk/github-actions-exporter
# Tag of the exporter to checkout
exporter_tag: v1.2
# Commit of the exporter to checkout (used if tag property is empty)
exporter_commit: ifTagIsSetThisIsNotUsed
# Changelog to add to the new release
exporter_changelog: "Changelog for the current version, nothing relly changed, just testing pipeline"
# Enable packages for Linux
package_linux: true
# Enable packages for Windows
package_windows: true
# Exporter GUID used in the msi package Requiered if package_windows is set to true
exporter_guid: 7B629E90-530F-4FAA-B7FE-1F1B30A95714
# Lincense GUID used in the msi package Requiered if package_windows is set to true
license_guid: 95E897AC-895A-43BE-A5EF-D72AD58E4ED1
```

When added open a PR and once merged to master a github action workflow will start building and uploading packages to Github. 

 - In case one exporter definition has been modified or added the exporter will be released for the os requested and a Github release will be created
 - In case two exporter definitions have been modified the pipeline fail
 - In case no exporters definition have been modified the pipeline terminates

Please notice that exporters have their own `build` script but they share the packaging scripts, located under `./scripts`

## Support

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices. Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub. You can find this project's topic/threads here:

https://discuss.newrelic.com/t/prometheus-exporters-packages/116524

## Contributing
We encourage your contributions to improve Prometheus Exporter Packages! Keep in mind when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.
If you have any questions, or to execute our corporate CLA, required if your contribution is on behalf of a company,  please drop us an email at opensource@newrelic.com.

## License
Prometheus Exporter Packages is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.
