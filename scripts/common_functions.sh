#!/usr/bin/env bash

# loadVariables reads with yq the exporter definition in EXPORTER_PATH and exports the variables
loadVariables(){
    exporter_path=$1
    export NAME=$(cat $exporter_path | yq e .name -)
    export VERSION=$(cat $exporter_path | yq e .version -)
    export EXPORTER_REPO_URL=$(cat $exporter_path | yq e .exporter_repo_url -)
    export EXPORTER_LICENSE_PATH=$(cat $exporter_path | yq e .exporter_license_path -)
    export EXPORTER_TAG=$(cat $exporter_path | yq e .exporter_tag -)
    export EXPORTER_COMMIT=$(cat $exporter_path | yq e .exporter_commit -)
    export EXPORTER_CHANGELOG=$(cat $exporter_path | yq e .exporter_changelog -)
    export EXPORTER_CONFIG_FILES=$(cat $exporter_path | yq e '.exporter_config_files // ""' -)
    export UPGRADE_GUID=$(cat $exporter_path | yq e .upgrade_guid -)
    export NRI_GUID=$(cat $exporter_path | yq e .nri_guid -)
    export EXPORTER_GUID=$(cat $exporter_path | yq e .exporter_guid -)
    export CONFIG_GUID=$(cat $exporter_path | yq e .config_guid -)
    export LICENSE_GUID=$(cat $exporter_path | yq e .license_guid -)
    export PACKAGE_LINUX=$(cat $exporter_path | yq e .package_linux -)
    export PACKAGE_WINDOWS=$(cat $exporter_path | yq e .package_windows -)
    if [[ -z $EXPORTER_TAG ]]
    then
        export EXPORTER_HEAD=$EXPORTER_COMMIT
    else
        export  EXPORTER_HEAD=$EXPORTER_TAG 
    fi

    export PACKAGE_NAME=nri-${NAME}

    PACKAGE_LINUX_GOARCHS=$(cat $exporter_path | yq e .package_linux_goarchs -)
    export PACKAGE_LINUX_GOARCHS="${PACKAGE_LINUX_GOARCHS:-amd64}"
}

# setStepOutput exposes the environment variables needed by next github actions steps steps
setStepOutput(){
    echo "NAME=${NAME}" >> $GITHUB_OUTPUT
    echo "PACKAGE_NAME=${PACKAGE_NAME}" >> $GITHUB_OUTPUT
    echo "EXPORTER_HEAD=${EXPORTER_HEAD}" >> $GITHUB_OUTPUT
    echo "EXPORTER_REPO_URL=${EXPORTER_REPO_URL}" >> $GITHUB_OUTPUT
    echo "EXPORTER_LICENSE_PATH=${EXPORTER_LICENSE_PATH}" >> $GITHUB_OUTPUT
    echo "VERSION=${VERSION}" >> $GITHUB_OUTPUT
    echo "EXPORTER_CHANGELOG=${EXPORTER_CHANGELOG}" >> $GITHUB_OUTPUT
    echo "CREATE_RELEASE=${CREATE_RELEASE}" >> $GITHUB_OUTPUT
    echo "EXPORTER_PATH=${EXPORTER_PATH}" >> $GITHUB_OUTPUT
    echo "PACKAGE_LINUX=${PACKAGE_LINUX}" >> $GITHUB_OUTPUT
    echo "PACKAGE_WINDOWS=${PACKAGE_WINDOWS}" >> $GITHUB_OUTPUT
    echo "UPGRADE_GUID=${UPGRADE_GUID}" >> $GITHUB_OUTPUT
    echo "EXPORTER_GUID=${EXPORTER_GUID}" >> $GITHUB_OUTPUT
    echo "NRI_GUID=${NRI_GUID}" >> $GITHUB_OUTPUT
    echo "LICENSE_GUID=${LICENSE_GUID}" >> $GITHUB_OUTPUT
    echo "CONFIG_GUID=${CONFIG_GUID}" >> $GITHUB_OUTPUT
}

# shouldDoRelease checks if any exporter.yml has been modified, if so we set CREATE_RELEASE to true setting the variable EXPORTER_PATH
shouldDoRelease(){
    old=$(git describe --tags --abbrev=0)  # ERROR PRONE IF THERE IS NO PREVIOUS TAG!
    export EXPORTER_PATH=$(git --no-pager diff  --name-only $old "exporters/**/exporter.yml")
    CREATE_RELEASE=false

    if [ -z "$EXPORTER_PATH" ]
    then
        echo "No definition has been modified"
        echo "CREATE_RELEASE=${CREATE_RELEASE}" >> $GITHUB_OUTPUT
        exit 0
    fi

    if (( $(git --no-pager diff  --name-only $old "exporters/**/exporter.yml"| wc -l) > 1 ))
    then
        echo "Only one definition should be modified at the same time"
        git --no-pager diff  --name-only $old "exporters/**/exporter.yml"
        echo "CREATE_RELEASE=${CREATE_RELEASE}" >> $GITHUB_OUTPUT
        exit 1
    fi
    CREATE_RELEASE=true
}


# checkExporter runs a series of tests to find common issues
checkExporter(){
    ERRORS=""
    # checking variables in the yaml file
    if [ -z "$NAME" ];then
        ERRORS=$ERRORS" - name is missing from exporter.yml"
    fi
    if [ -z "$EXPORTER_HEAD" ];then
        ERRORS=$ERRORS" - exporter_tag and exporter_commit are missing from exporter.yml"
    fi
    if [ -z "$EXPORTER_REPO_URL" ];then
        ERRORS=$ERRORS" - exporter_repo_url is missing from exporter.yml"
    fi
    if [ -z "$EXPORTER_LICENSE_PATH" ];then
        ERRORS=$ERRORS" - exporter_license_path is missing from exporter.yml"
    fi
    if [ -z "$VERSION" ];then
        ERRORS=$ERRORS" - version is missing from exporter.yml"
    fi
    if [ -z "$PACKAGE_LINUX" ];then
        ERRORS=$ERRORS" - package_linux is missing from exporter.yml"
    fi
    if [ -z "$PACKAGE_LINUX_GOARCHS" ];then
        ERRORS=$ERRORS" - package_linux_goarchs is missing from exporter.yml"
    fi
    if [ -z "$PACKAGE_WINDOWS" ];then
        ERRORS=$ERRORS" - package_windows is missing from exporter.yml"
    fi

    # checking if the linux packaging is required if the file are present
    if [ "$PACKAGE_LINUX" = "true" ];then
        if [ ! -f "./exporters/$NAME/$NAME-config.yml.sample" ]; then
            ERRORS=$ERRORS" - the file ./exporters/$NAME/$NAME-config.yml.sample should exist"
        fi
        if [ ! -f "./exporters/$NAME/build-exporter-linux.sh" ]; then
            ERRORS=$ERRORS" - the file ./exporters/$NAME/build-exporter-linux.sh should exist"
        fi
    fi

    # checking if the windows packaging is required if the file and GUIID are present and not reused
    if [ "$PACKAGE_WINDOWS" = "true" ];then
        if [ -z "$UPGRADE_GUID" ];then
            ERRORS=$ERRORS" - upgrade_guid is missing from exporter.yml"
        else
            if [ $(grep $UPGRADE_GUID exporters/*/exporter.yml | wc -l) != 1 ];then
                ERRORS=$ERRORS" - upgrade_guid was already used in a different exporter"
            fi
            if [[ ! $UPGRADE_GUID =~ ^\{?[A-F0-9a-f]{8}-[A-F0-9a-f]{4}-[A-F0-9a-f]{4}-[A-F0-9a-f]{4}-[A-F0-9a-f]{12}\}?$ ]]; then
                ERRORS=$ERRORS" - upgrade_guid is not a GUID"
            fi
        fi
        
        if [ -z "$EXPORTER_GUID" ];then
            ERRORS=$ERRORS" - exporter_guid is missing from exporter.yml"
        else
            if [ $(grep $EXPORTER_GUID exporters/*/exporter.yml | wc -l) != 1 ];then
                ERRORS=$ERRORS" - exporter_guid was already used in a different exporter"
            fi
            if [[ ! $EXPORTER_GUID =~ ^\{?[A-F0-9a-f]{8}-[A-F0-9a-f]{4}-[A-F0-9a-f]{4}-[A-F0-9a-f]{4}-[A-F0-9a-f]{12}\}?$ ]]; then
                ERRORS=$ERRORS" - exporter_guid is not a GUID"
            fi
        fi

        if [ -z "$NRI_GUID" ];then
            ERRORS=$ERRORS" - nri_guid is missing from exporter.yml"
        else
            if [ $(grep $NRI_GUID exporters/*/exporter.yml | wc -l) != 1 ];then
                ERRORS=$ERRORS" - nri_guid was already used in a different exporter"
            fi
            if [[ ! $NRI_GUID =~ ^\{?[A-F0-9a-f]{8}-[A-F0-9a-f]{4}-[A-F0-9a-f]{4}-[A-F0-9a-f]{4}-[A-F0-9a-f]{12}\}?$ ]]; then
                ERRORS=$ERRORS" - nri_guid is not a GUID"
            fi
        fi

        if [ -z "$LICENSE_GUID" ];then
            ERRORS=$ERRORS" - license_guid is missing from exporter.yml"
        else
            if [ $(grep $LICENSE_GUID exporters/*/exporter.yml | wc -l) != 1 ];then
                ERRORS=$ERRORS" - license_guid was already used in a different exporter"
            fi
            if [[ ! $LICENSE_GUID =~ ^\{?[A-F0-9a-f]{8}-[A-F0-9a-f]{4}-[A-F0-9a-f]{4}-[A-F0-9a-f]{4}-[A-F0-9a-f]{12}\}?$ ]]; then
                ERRORS=$ERRORS" - license_guid is not a GUID"
            fi
        fi

        if [ -z "$CONFIG_GUID" ];then
            ERRORS=$ERRORS" - config_guid is missing from exporter.yml"
        else
            if [ $(grep $CONFIG_GUID exporters/*/exporter.yml | wc -l) != 1 ];then
                ERRORS=$ERRORS" - config_guid was already used in a different exporter"
            fi
            if [[ ! $CONFIG_GUID =~ ^\{?[A-F0-9a-f]{8}-[A-F0-9a-f]{4}-[A-F0-9a-f]{4}-[A-F0-9a-f]{4}-[A-F0-9a-f]{12}\}?$ ]]; then
                ERRORS=$ERRORS" - config_guid is not a GUID"
            fi
        fi

        if [ ! -f "./exporters/$NAME/build-exporter-windows.sh" ]; then
            ERRORS=$ERRORS" - the file ./exporters/$NAME/build-exporter-windows.sh should exist"
        fi
    fi

    # checking license file and if the name of the folder is the same in the definition
    IFS="/" read tmp exporter_name exporter_yaml <<< "$EXPORTER_PATH"
    if [ "$exporter_name" != "$NAME" ]; then
        ERRORS=$ERRORS" - The exporter.yml is in a wrong folder"
    fi
}
