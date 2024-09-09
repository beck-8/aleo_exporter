package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// 实时读取日志文件并寻找关键词
func tailLogFile(logFilePath string) {
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer file.Close()

	// 创建文件监视器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("无法创建文件监视器: %v", err)
	}
	defer watcher.Close()

	// 将日志文件添加到监视器
	err = watcher.Add(logFilePath)
	if err != nil {
		log.Fatalf("无法监视日志文件: %v", err)
	}

	// 移动文件指针到文件末尾
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatalf("无法移动到文件末尾: %v", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// 当文件被修改时读取新增的日志行
			if event.Op&fsnotify.Write == fsnotify.Write {
				// log.Println("Receive the write operation")

				scanner := bufio.NewScanner(file)

				for scanner.Scan() {
					line := scanner.Text()

					/*
						+------------------------------------------------------------------------------------------+
						|                                           OULA                                           |
						+==========================================================================================+
						|  GPU   | ProofRate | Valid | Power       | Memory             | Device                   |
						| -------+-----------+-------+-------------+--------------------+------------------------- |
						|  0     | 117028    | 0     | 312 W/320 W | 1008 MiB/10240 MiB | NVIDIA GeForce RTX 3080  |
						|  1     | 116583    | 0     | 310 W/320 W | 1008 MiB/10240 MiB | NVIDIA GeForce RTX 3080  |
						|  2     | 115021    | 0     | 311 W/320 W | 1008 MiB/10240 MiB | NVIDIA GeForce RTX 3080  |
						|  3     | 118355    | 0     | 296 W/320 W | 1008 MiB/10240 MiB | NVIDIA GeForce RTX 3080  |
						|  4     | 115778    | 0     | 294 W/320 W | 1008 MiB/10240 MiB | NVIDIA GeForce RTX 3080  |
						|  5     | 117695    | 0     | 296 W/320 W | 1008 MiB/10240 MiB | NVIDIA GeForce RTX 3080  |
						|  6     | 117500    | 0     | 277 W/320 W | 1008 MiB/10240 MiB | NVIDIA GeForce RTX 3080  |
						|  7     | 117795    | 0     | 284 W/320 W | 1008 MiB/10240 MiB | NVIDIA GeForce RTX 3080  |
						|  Total | 935755    | 0     |             |                    | Uptime  68415s           |
						+------------------------------------------------------------------------------------------+
												6bock
						+-------------------------------------------------------------------------------------------+
						| 2024-08-21T07:56:47                                                                       |
						|                                                                                           |
						| gpu[0]: (1m - 138750    5m - 138916    15m - 133666    30m - 162111    60m - 156944  )    |
						| gpu[1]: (1m - 138750    5m - 138916    15m - 133750    30m - 162347    60m - 157298  )    |
						| gpu[2]: (1m - 137500    5m - 137250    15m - 132111    30m - 160416    60m - 155423  )    |
						| gpu[3]: (1m - 140000    5m - 139833    15m - 134555    30m - 163291    60m - 158048  )    |
						| gpu[4]: (1m - 136666    5m - 136750    15m - 131583    30m - 159875    60m - 154715  )    |
						| gpu[5]: (1m - 137916    5m - 137750    15m - 132666    30m - 161486    60m - 156472  )    |
						| gpu[6]: (1m - 138750    5m - 138583    15m - 133388    30m - 162083    60m - 157111  )    |
						| gpu[7]: (1m - 139166    5m - 139333    15m - 134166    30m - 162763    60m - 157597  )    |
						| gpu[*]: (1m - 1107500   5m - 1107333   15m - 1065888   30m - 1294375   60m - 1253611 )    |
						|                                                                                           |
						+-------------------------------------------------------------------------------------------+
					*/

					var rateString string
					if strings.Contains(line, "Total") && strings.Contains(line, "Uptime") {
						// oula
						rateString = strings.TrimSpace(strings.Split(line, "|")[2])

					} else if strings.Contains(line, "gpu[*]:") {
						// 6block
						rateString = strings.Split(line, " ")[4]
					}

					if rateString == "" {
						continue
					} else if rateString == "N/A" {
						ProofRate.Set(0)
						continue
					}

					rateNum, err := strconv.ParseUint(rateString, 10, 64)
					if err != nil {
						log.Println(err)
						log.Println(line)
						totalLogError.Inc()
					} else {
						ProofRate.Set(float64(rateNum))
					}

				}

				if err := scanner.Err(); err != nil {
					log.Printf("读取日志文件时出错: %v", err)
					totalLogError.Inc()
				}

				// 移动文件指针到文件末尾，处理日志被清空的情况
				_, err = file.Seek(0, io.SeekEnd)
				if err != nil {
					log.Printf("无法移动到文件末尾: %v", err)
				}
			}

			// 处理日志被轮转的情况
			if event.Has(fsnotify.Remove | fsnotify.Rename) {
				log.Printf("Receive the %v operation", event.String())
				file.Close()

				// 等待新文件产生
				time.Sleep(1 * time.Second)

				file, err = os.Open(logFilePath)
				if err != nil {
					log.Fatalf("无法打开日志文件: %v", err)
				}
				defer file.Close()

			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("监视器错误: %v", err)
			totalLogError.Inc()
		}
	}
}
