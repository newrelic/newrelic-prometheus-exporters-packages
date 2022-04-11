#!/usr/bin/env bash

# loadVariables reads with yq the exporter definition in EXPORTER_PATH and exports the variables
loadVariables(){
    export NAME=$(cat $EXPORTER_PATH | yq e .name -)
    export VERSION=$(cat $EXPORTER_PATH | yq e .version -)
    export EXPORTER_REPO_URL=$(cat $EXPORTER_PATH | yq e .exporter_repo_url -)
    export EXPORTER_LICENSE_PATH=$(cat $EXPORTER_PATH | yq e .exporter_license_path -)
    export EXPORTER_TAG=$(cat $EXPORTER_PATH | yq e .exporter_tag -)
    export EXPORTER_COMMIT=$(cat $EXPORTER_PATH | yq e .exporter_commit -)
    export EXPORTER_CHANGELOG=$(cat $EXPORTER_PATH | yq e .exporter_changelog -)
    export UPGRADE_GUID=$(cat $EXPORTER_PATH | yq e .upgrade_guid -)
    export NRI_GUID=$(cat $EXPORTER_PATH | yq e .nri_guid -)
    export EXPORTER_GUID=$(cat $EXPORTER_PATH | yq e .exporter_guid -)
    export CONFIG_GUID=$(cat $EXPORTER_PATH | yq e .config_guid -)
    export LICENSE_GUID=$(cat $EXPORTER_PATH | yq e .license_guid -)
    export PACKAGE_LINUX=$(cat $EXPORTER_PATH | yq e .package_linux -)
    export PACKAGE_WINDOWS=$(cat $EXPORTER_PATH | yq e .package_windows -)
    if [[ -z $EXPORTER_TAG ]]
    then
        export EXPORTER_HEAD=$EXPORTER_COMMIT
    else
        export  EXPORTER_HEAD=$EXPORTER_TAG 
    fi

    export PACKAGE_NAME=nri-${NAME}
}

# setStepOutput exposes the environment variables needed by next github actions steps steps
setStepOutput(){
    echo "::set-output name=NAME::${NAME}"
    echo "::set-output name=PACKAGE_NAME::${PACKAGE_NAME}"
    echo "::set-output name=EXPORTER_HEAD::${EXPORTER_HEAD}"
    echo "::set-output name=EXPORTER_REPO_URL::${EXPORTER_REPO_URL}"
    echo "::set-output name=EXPORTER_LICENSE_PATH::${EXPORTER_LICENSE_PATH}"
    echo "::set-output name=VERSION::${VERSION}"
    echo "::set-output name=EXPORTER_CHANGELOG::${EXPORTER_CHANGELOG}"
    echo "::set-output name=CREATE_RELEASE::${CREATE_RELEASE}"
    echo "::set-output name=EXPORTER_PATH::${EXPORTER_PATH}"
    echo "::set-output name=PACKAGE_LINUX::${PACKAGE_LINUX}"
    echo "::set-output name=PACKAGE_WINDOWS::${PACKAGE_WINDOWS}"
    echo "::set-output name=UPGRADE_GUID::${UPGRADE_GUID}"
    echo "::set-output name=EXPORTER_GUID::${EXPORTER_GUID}"
    echo "::set-output name=NRI_GUID::${NRI_GUID}"
    echo "::set-output name=LICENSE_GUID::${LICENSE_GUID}"
    echo "::set-output name=CONFIG_GUID::${CONFIG_GUID}"

}



# packageLinux runs the makefile with target all int the EXPORTER_PATH repo
packageLinux(){
    IFS="/" read tmp exporter_name exporter_yaml <<< "$EXPORTER_PATH"

    if [ $exporter_name != $NAME ]; then
        echo "The exporter.yml is in a wrong folder. The name in the definition '$NAME' does not match with the foldername '$exporter_name'" 
        exit 1
    fi

    current_pwd=$(pwd)
    make build-$exporter_name
    cd $current_pwd
}

# shouldDoRelease checks if any exporter.yml has been modified, if so we set CREATE_RELEASE to true setting the variable EXPORTER_PATH
shouldDoRelease(){
    old=$(git describe --tags --abbrev=0)
    export EXPORTER_PATH=$(git --no-pager diff  --name-only $old "exporters/**/exporter.yml")
    CREATE_RELEASE=false

    if [ -z "$EXPORTER_PATH" ]
    then
        echo "No definition has been modified"
        echo "::set-output name=CREATE_RELEASE::${CREATE_RELEASE}"
        exit 0
    fi

    if (( $(git --no-pager diff  --name-only $old "exporters/**/exporter.yml"| wc -l) > 1 ))
    then
        echo "Only one definition should be modified at the same time"
        git --no-pager diff  --name-only $old "exporters/**/exporter.yml"
        echo "::set-output name=CREATE_RELEASE::${CREATE_RELEASE}"
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
    if [ -z "$PACKAGE_WINDOWS" ];then
        ERRORS=$ERRORS" - package_windows is missing from exporter.yml"
    fi

    # checking if the linux packaging is required if the file are present
    if [ "$PACKAGE_LINUX" = "true" ];then
        if [ ! -f "./exporters/$NAME/$NAME-config.yml.sample" ]; then
            ERRORS=$ERRORS" - the file ./exporters/$NAME/$NAME-config.yml.sample should exist"
        fi
        if [ ! -f "./exporters/$NAME/build.sh" ]; then
            ERRORS=$ERRORS" - the file ./exporters/$NAME/build.sh should exist"
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

        if [ ! -f "./exporters/$NAME/win_build.ps1" ]; then
            ERRORS=$ERRORS" - the file ./exporters/$NAME/win_build.ps1 should exist"
        fi
    fi

    # checking license file and if the name of the folder is the same in the definition
    IFS="/" read tmp exporter_name exporter_yaml <<< "$EXPORTER_PATH"
    if [ "$exporter_name" != "$NAME" ]; then
        ERRORS=$ERRORS" - The exporter.yml is in a wrong folder"
    fi


}