#!/bin/bash

# 获取 Go 模块名称
get_module_name() {
    local module_name=$(grep 'module' go.mod | cut -d ' ' -f 2)
    echo $module_name
}

# 获取项目的根目录
get_project_root() {
    local current_dir=$(pwd)
    local root_dir=$current_dir

    # 查找上一级目录，直到找到一个包含 go.mod 文件的目录
    while [ "$root_dir" != "/" ]; do
        if [ -f "$root_dir/go.mod" ]; then
            break
        fi
        root_dir=$(dirname "$root_dir")
    done

    if [ ! -f "$root_dir/go.mod" ]; then
        echo >&2 "error: unable to determine go mod project root directory"
        exit 1
    fi
    echo "$root_dir"
}

# 搜索 version.go 文件的位置
find_version_go_path() {
    local root_dir=$(get_project_root)
    local version_go_path
    # 在项目目录中搜索 version.go 文件
    version_go_path=$(find "$root_dir" -name 'version.go' | head -n 1)
    if [ -z "$version_go_path" ]; then
        echo >&2 "error: version.go file not found"
        exit 1
    fi
    echo "$version_go_path"
}

# 获取 version.go 的包路径
get_version_package_path() {
    local root_dir=$(get_project_root)
    local version_go_path=$(find_version_go_path)
    local module_name=$(get_module_name)

    # 计算 version.go 文件相对于项目根目录的路径
    local package_path=${version_go_path#"$root_dir/"}
    # 移除文件名，保留目录路径
    package_path=${package_path%/*}

    # 连接模块名和包路径
    echo "$module_name/$package_path"
}

# 获取版本信息
get_version() {
    VERSION="dev"
    if git rev-parse --git-dir > /dev/null 2>&1; then
        # 获取最近的标签名
        VERSION=$(git describe --tags --exact-match 2>/dev/null)
        if [ $? -ne 0 ]; then
            VERSION=$(git rev-parse --abbrev-ref HEAD)
        fi
    fi
    echo $VERSION
}

# 获取构建时间
get_build_time() {
    BUILDTIME=$(date -d "@${SOURCE_DATE_EPOCH:-$(date +%s)}" --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')
    echo $BUILDTIME
}

# 获取 Git 分支信息
get_git_branch() {
    # 检查是否存在 .git 目录
    if command -v git &> /dev/null && [ -d ".git" ]; then
        GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
    else
        GIT_BRANCH="Unknown"
    fi
    echo $GIT_BRANCH
}

# 获取 Git 提交哈希信息
get_git_revision() {
    # 检查是否存在 .git 目录
    if command -v git &> /dev/null && [ -d ".git" ]; then
        GIT_REVISION=$(git rev-parse HEAD 2>/dev/null)
        if [ -n "$(git status --porcelain --untracked-files=no)" ]; then
            GIT_REVISION="${GIT_REVISION}-unsupported"
        fi
    else
        GIT_REVISION="Unknown"
    fi
    echo $GIT_REVISION
}

# 主函数，组合所有信息
generate_ldflags() {
    local module_name=$(get_module_name)
    local version=$(get_version)
    local version_package_path=$(get_version_package_path)
    local build_time=$(get_build_time)
    local git_branch=$(get_git_branch)
    local git_revision=$(get_git_revision)

    # 格式化为 ldflags 选项
    LDFLAGS="-X '${version_package_path}.Version=${version}' \
             -X '${version_package_path}.BuildTime=${build_time}' \
             -X '${version_package_path}.GitBranch=${git_branch}' \
             -X '${version_package_path}.GitRevision=${git_revision}'"

    echo $LDFLAGS
}