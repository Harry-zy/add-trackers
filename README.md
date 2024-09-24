# Transmission Tracker Adder

`transmission-tracker-adder` 是一个基于 Go 的工具，旨在将多个 Tracker URL 添加到 Transmission 服务器上的种子中。该工具允许您与 Transmission 的 API 交互，连接到服务器，获取所有种子 ID，并批量添加 Tracker，以防止超时。

## Features / 功能

- 使用凭据（HTTP 或 HTTPS）连接到 Transmission 服务器。
- 从服务器获取所有种子 ID。
- 将多个 Tracker URL 批量添加到种子中（默认每批 200 个种子）。
- 在批次之间模拟重新连接，以避免连接超时。

## Requirements / 需求

- Go 1.18 或更高版本。
- [hekmon/transmissionrpc](https://github.com/hekmon/transmissionrpc/v2) 库，用于与 Transmission 服务器进行交互。

## Installation / 安装

1. 克隆仓库：

    ```bash
    git clone https://github.com/yourusername/transmission-tracker-adder.git
    ```

2. 导航到项目目录：

    ```bash
    cd transmission-tracker-adder
    ```

3. 安装依赖：

    ```bash
    go get -u github.com/hekmon/transmissionrpc/v2
    ```

## Usage / 使用

1. 构建项目：

    ```bash
    go build -o transmission-tracker-adder
    ```

2. 运行工具：

    ```bash
    ./transmission-tracker-adder
    ```

3. 工具会提示您输入以下内容：

    - **Transmission 服务器地址**（例如，`192.168.1.10`）
    - **Transmission 服务器端口**（例如，`9091`）
    - **是否为 HTTPS 协议（true/false）** — 指定服务器是否使用 HTTPS。
    - **用户名** — Transmission 服务器用户名。
    - **密码** — Transmission 服务器密码。
    - **Tracker URLs** — 输入多个 Tracker URLs，用分号（`;`）分隔。

4. 工具将从服务器获取所有种子 ID，并以 200 个种子的批次将提供的 Tracker 添加到种子中。

### Example Input / 示例输入

```bash
请输入 Transmission 服务器的完整地址 (例如: 192.168.xxx.xxx): 192.168.1.10
请输入 Transmission 服务器端口号: 9091
请输入 Transmission 服务器是否为HTTPS协议 (true/false): false
请输入 Transmission 服务器的用户名: admin
请输入 Transmission 服务器的密码: password123
请输入多个 Tracker URLs (使用分号 ; 分隔): https://tracker1.com;https://tracker2.com