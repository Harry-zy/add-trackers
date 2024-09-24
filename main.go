package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hekmon/transmissionrpc/v2"
)

func main() {
	var serverAddress, username, password, trackerURLsInput string
	var port int
	var isHttps bool

	// 提示用户输入 Transmission 服务器地址 (例如: 192.168.xxx.xxx)
	fmt.Println("请输入 Transmission 服务器的完整地址 (例如: 192.168.xxx.xxx):")
	fmt.Scanln(&serverAddress)

	// 提示用户输入 Transmission 服务器端口号
	fmt.Println("请输入 Transmission 服务器服务器端口号:")
	fmt.Scanln(&port)

	// 提示用户输入 Transmission 服务器是否为HTTPS协议
	fmt.Println("请输入 Transmission 服务器是否为HTTPS协议:")
	fmt.Scanln(&isHttps)

	// 提示用户输入 Transmission 的用户名和密码
	fmt.Println("请输入 Transmission 服务器的用户名:")
	fmt.Scanln(&username)

	fmt.Println("请输入 Transmission 服务器的密码:")
	fmt.Scanln(&password)

	// 提示用户输入多个 Tracker URLs，并使用分号 ; 分隔
	fmt.Println("请输入多个 Tracker URLs (使用分号 ; 分隔):")
	fmt.Scanln(&trackerURLsInput)

	// 解析 tracker URLs 输入，按分号分隔，并替换为换行符
	trackerURLs := strings.Split(trackerURLsInput, ";")
	// 创建一个 Transmission 客户端
	client, err := transmissionrpc.New(serverAddress, username, password, &transmissionrpc.AdvancedConfig{
		Port:  uint16(port),
		HTTPS: isHttps,
	})
	if err != nil {
		log.Fatalf("无法连接到 Transmission 服务器: %v", err)
	}

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取所有 torrent 的 IDs
	torrents, err := client.TorrentGetAll(ctx)
	if err != nil {
		log.Fatalf("获取 torrent 列表失败: %v", err)
	}

	// 创建 Torrent ID 列表
	var allTorrentIDs []int64
	for _, torrent := range torrents {
		if torrent.ID != nil {
			allTorrentIDs = append(allTorrentIDs, *torrent.ID)
		}
	}

	// 分批次处理，每次 200 个 torrent
	err = addTrackerInBatches(client, allTorrentIDs, trackerURLs)
	if err != nil {
		log.Fatalf("批量添加 trackers 失败: %v", err)
	}

	fmt.Println("所有 trackers 已添加完成。")
}

// 分批添加 Tracker，每批处理 200 个种子
func addTrackerInBatches(client *transmissionrpc.Client, allTorrentIDs []int64, trackers []string) error {
	const batchSize = 200

	for i := 0; i < len(allTorrentIDs); i += batchSize {
		// 计算当前批次的结束位置
		end := i + batchSize
		if end > len(allTorrentIDs) {
			end = len(allTorrentIDs)
		}

		// 取当前批次的 torrent IDs
		batch := allTorrentIDs[i:end]

		// 处理当前批次
		err := addTrackersToBatch(client, batch, trackers)
		if err != nil {
			// 打印错误信息，但继续处理后续批次
			fmt.Printf("Error adding trackers to batch %d-%d: %v\n", i, end, err)
		}

		// 模拟重新建立连接（重新创建 Client 实例）
		time.Sleep(2 * time.Second) // 暂停一段时间来模拟重新建立连接
		fmt.Println("Re-establishing connection...")
	}

	return nil
}

// 为当前批次的种子添加 trackers
func addTrackersToBatch(client *transmissionrpc.Client, batch []int64, trackers []string) error {
	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 遍历每个 torrent，添加 trackers
	for _, torrentID := range batch {
		fmt.Printf("正在为 Torrent ID: %d 添加 Tracker...\n", torrentID)

		// 添加 trackers
		err := client.TorrentSet(ctx, transmissionrpc.TorrentSetPayload{
			IDs:        []int64{torrentID},
			TrackerAdd: trackers,
		})
		if err != nil {
			log.Printf("无法为 Torrent ID %d 添加 trackers: %v", torrentID, err)
			continue
		}
		fmt.Printf("成功为 Torrent ID: %d 添加 Trackers: %v\n", torrentID, trackers)
	}

	return nil
}
