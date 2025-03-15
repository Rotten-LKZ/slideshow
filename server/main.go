package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type FolderTree map[string]interface{}

type Config struct {
	Basic Basic `yaml:"basic"`
	Anime Anime `yaml:"anime"`
}

type Basic struct {
	ImagesPath string `yaml:"images_path"`
	Url        string `yaml:"url"`
	Port       int    `yaml:"port"`
}

type Anime struct {
	Duration        int    `yaml:"duration"`
	AnimeDuration   int    `yaml:"anime_duration"`
	BackgroundImage string `yaml:"background_image"`
}

func main() {
	dataBytes, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("读取配置文件失败：", err)
		return
	}
	config := Config{}
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		fmt.Println("解析配置文件失败：", err)
		return
	}

	mp := make(map[string]any, 2)
	err = yaml.Unmarshal(dataBytes, mp)
	if err != nil {
		fmt.Println("解析配置文件失败：", err)
		return
	}

	fs := http.FileServer(http.Dir(config.Basic.ImagesPath))
	http.Handle("/g/", http.StripPrefix("/g/", fs))

	http.Handle("/p/", http.StripPrefix("/p/", http.FileServer(http.Dir("."))))

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		prefix := fmt.Sprintf("%s:%d/g/", config.Basic.Url, config.Basic.Port)
		j, err := generateJSON(config.Basic.ImagesPath, prefix, config.Anime)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "HELLO WORLD")
	})

	fmt.Printf("服务器已在 %s:%d 启动\n", config.Basic.Url, config.Basic.Port)
	port := fmt.Sprintf(":%d", config.Basic.Port)
	http.ListenAndServe(port, nil)
}

// 递归遍历目录并构建文件夹树
func buildFolderTree(rootPath string) (FolderTree, error) {
	tree := make(FolderTree)

	// 读取目录下的所有文件和文件夹
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// 如果是文件夹，递归处理
			subTree, err := buildFolderTree(filepath.Join(rootPath, entry.Name()))
			if err != nil {
				return nil, err
			}
			tree[entry.Name()] = subTree
		} else {
			// 如果是文件，直接添加到当前文件夹的文件列表中
			if _, ok := tree["files/"]; !ok {
				tree["files/"] = []string{}
			}
			tree["files/"] = append(tree["files/"].([]string), url.PathEscape(entry.Name()))
		}
	}

	return tree, nil
}

// 生成JSON字符串
func generateJSON(rootPath string, prefix string, config any) ([]byte, error) {
	tree, err := buildFolderTree(rootPath)
	if err != nil {
		return nil, err
	}
	tree["BaseUrl/"] = prefix
	tree["AnimeConfig/"] = config

	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
